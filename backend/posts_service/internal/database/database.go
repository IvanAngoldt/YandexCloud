package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Connect() (*sql.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

type Post struct {
	ID             int           `json:"id"`
	Title          string        `json:"title"`
	Content        string        `json:"content"`
	AuthorID       int           `json:"authorId"`
	AuthorUsername string        `json:"authorUsername"`
	Likes          []interface{} `json:"likes"`
}

// FetchPosts возвращает все посты с лайками и информацией об авторе
func FetchPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query(`
        SELECT 
            posts.id, 
            posts.title, 
            posts.content, 
            posts.author_id AS author_id,
            users.username AS author_username,
            COALESCE(
                json_agg(
                    json_build_object(
                        'id', likes.user_id, 
                        'username', liked_users.username
                    )
                ) FILTER (WHERE likes.user_id IS NOT NULL), '[]'
            ) AS likes
        FROM posts
        JOIN users ON posts.author_id = users.id
        LEFT JOIN likes ON posts.id = likes.post_id
        LEFT JOIN users AS liked_users ON likes.user_id = liked_users.id
        GROUP BY posts.id, users.username
		ORDER BY posts.created_at DESC
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var likesJSON string
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.AuthorUsername, &likesJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post row: %w", err)
		}

		// Декодируем JSON-строку likes в массив объектов
		if err := json.Unmarshal([]byte(likesJSON), &post.Likes); err != nil {
			return nil, fmt.Errorf("failed to parse likes JSON: %w", err)
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over rows: %w", err)
	}

	return posts, nil
}

// CreatePost добавляет новый пост в базу данных и возвращает его информацию
func CreatePost(db *sql.DB, title, content string, authorID int) (*Post, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.WithFields(logrus.Fields{
		"title":    title,
		"content":  content,
		"authorID": authorID,
	}).Info("Inserting post into database")

	var post Post
	err := db.QueryRow(`
        WITH inserted_post AS (
            INSERT INTO posts (title, content, author_id)
            VALUES ($1, $2, $3)
            RETURNING id, title, content, author_id
        )
        SELECT 
            inserted_post.id, 
            inserted_post.title, 
            inserted_post.content, 
            inserted_post.author_id, 
            users.username AS author_username
        FROM inserted_post
        JOIN users ON inserted_post.author_id = users.id
    `, title, content, authorID).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.AuthorUsername,
	)
	if err != nil {
		logger.WithError(err).Error("Failed to insert post into database")
		return nil, fmt.Errorf("failed to insert post: %w", err)
	}

	// Новый пост ещё не имеет лайков
	post.Likes = []interface{}{}
	logger.WithFields(logrus.Fields{
		"id":             post.ID,
		"title":          post.Title,
		"content":        post.Content,
		"authorID":       post.AuthorID,
		"authorUsername": post.AuthorUsername,
	}).Info("Post inserted into database successfully")

	return &post, nil
}

// FetchPostByID возвращает пост по ID с информацией о лайках
func FetchPostByID(db *sql.DB, postID int) (*Post, error) {
	var post Post
	var likesJSON string

	err := db.QueryRow(`
        SELECT 
            posts.id, 
            posts.title, 
            posts.content, 
            posts.author_id AS author_id,
            users.username AS author_username,
            COALESCE(
                json_agg(
                    json_build_object(
                        'id', likes.user_id, 
                        'username', liked_users.username
                    )
                ) FILTER (WHERE likes.user_id IS NOT NULL), '[]'
            ) AS likes
        FROM posts
        JOIN users ON posts.author_id = users.id
        LEFT JOIN likes ON posts.id = likes.post_id
        LEFT JOIN users AS liked_users ON likes.user_id = liked_users.id
        WHERE posts.id = $1
        GROUP BY posts.id, users.username
    `, postID).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.AuthorUsername,
		&likesJSON,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Пост не найден
	} else if err != nil {
		return nil, fmt.Errorf("failed to fetch post: %w", err)
	}

	// Декодируем JSON-строку likes в массив объектов
	if err := json.Unmarshal([]byte(likesJSON), &post.Likes); err != nil {
		return nil, fmt.Errorf("failed to parse likes JSON: %w", err)
	}

	return &post, nil
}

// GetPostOwner возвращает ID пользователя, которому принадлежит пост
func GetPostOwner(db *sql.DB, postID int) (int, error) {
	var ownerID int
	err := db.QueryRow("SELECT author_id FROM posts WHERE id = $1", postID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		return 0, nil // Пост не найден
	} else if err != nil {
		return 0, fmt.Errorf("failed to retrieve post owner: %w", err)
	}
	return ownerID, nil
}

// DeletePost удаляет пост по его ID
func DeletePost(db *sql.DB, postID int) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = $1", postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

// FetchUserPosts возвращает список постов конкретного пользователя по его userID.
func FetchUserPosts(db *sql.DB, userID int) ([]Post, error) {
	rows, err := db.Query(`
        SELECT 
            posts.id, 
            posts.title, 
            posts.content, 
            posts.author_id AS author_id,
            users.username AS author_username,
            COALESCE(
                json_agg(
                    json_build_object(
                        'id', likes.user_id, 
                        'username', liker.username
                    )
                ) FILTER (WHERE likes.user_id IS NOT NULL), '[]'
            ) AS likes
        FROM posts
        JOIN users ON posts.author_id = users.id
        LEFT JOIN likes ON posts.id = likes.post_id
        LEFT JOIN users AS liker ON likes.user_id = liker.id
        WHERE posts.author_id = $1
        GROUP BY posts.id, users.username
        ORDER BY posts.created_at DESC
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user posts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var likesJSON string
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.AuthorUsername, &likesJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post row: %w", err)
		}

		if err := json.Unmarshal([]byte(likesJSON), &post.Likes); err != nil {
			return nil, fmt.Errorf("failed to parse likes JSON: %w", err)
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over rows: %w", err)
	}

	return posts, nil
}
