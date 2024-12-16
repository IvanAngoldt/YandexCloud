package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"posts_service/internal/database"
	"posts_service/internal/middlewares"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func FetchPosts(db *sql.DB) http.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем посты через функцию FetchPosts из database
		posts, err := database.FetchPosts(db)
		if err != nil {
			logger.WithError(err).Error("Failed to fetch posts from database")
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		// Логируем, если постов нет (информативно, но не ошибка)
		if len(posts) == 0 {
			logger.Info("No posts found, returning empty array")
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(posts); err != nil {
			logger.WithError(err).Error("Failed to encode response")
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// CreatePostRequest представляет запрос на создание поста
type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CreatePost обрабатывает запрос на создание нового поста
func CreatePost(db *sql.DB) http.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Received request to create a post")

		// Получение userID из контекста
		userID, ok := r.Context().Value(middlewares.UserIDKey).(int)
		if !ok {
			logger.Warn("User not authorized")
			http.Error(w, "User not authorized", http.StatusUnauthorized)
			return
		}
		logger.WithField("userID", userID).Info("Authorized user")

		// Декодируем запрос
		var req CreatePostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.WithError(err).Warn("Invalid request body")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		logger.WithFields(logrus.Fields{
			"title":   req.Title,
			"content": req.Content,
		}).Info("Request body decoded")

		// Вставляем пост в базу данных
		post, err := database.CreatePost(db, req.Title, req.Content, userID)
		if err != nil {
			logger.WithError(err).Error("Failed to create post in database")
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}
		logger.WithFields(logrus.Fields{
			"id":             post.ID,
			"title":          post.Title,
			"content":        post.Content,
			"authorID":       post.AuthorID,
			"authorUsername": post.AuthorUsername,
		}).Info("Post created successfully")

		// Возвращаем новый пост
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(post); err != nil {
			logger.WithError(err).Error("Failed to encode response")
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
		logger.Info("Response sent to the client")
	}
}

func FetchPostById(db *sql.DB) http.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		postID, err := atoiParam(vars["id"])
		if err != nil {
			logger.WithField("post_id", vars["id"]).Warn("Invalid post ID")
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Получаем пост из базы данных
		post, err := database.FetchPostByID(db, postID)
		if err != nil {
			logger.WithError(err).Error("Failed to fetch post")
			http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
			return
		}

		if post == nil {
			logger.WithField("post_id", postID).Warn("Post not found")
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		// Возвращаем ответ
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(post); err != nil {
			logger.WithError(err).Error("Failed to encode response")
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func DeletePost(db *sql.DB) http.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return func(w http.ResponseWriter, r *http.Request) {
		// Получение userID из контекста
		userID, ok := r.Context().Value(middlewares.UserIDKey).(int)
		if !ok {
			logger.Warn("User not authorized")
			http.Error(w, "User not authorized", http.StatusUnauthorized)
			return
		}

		// Получение postID из параметров запроса
		vars := mux.Vars(r)
		postID, err := atoiParam(vars["id"])
		if err != nil {
			logger.WithField("post_id", vars["id"]).Warn("Invalid post ID")
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Проверяем, что пост принадлежит данному пользователю
		ownerID, err := database.GetPostOwner(db, postID)
		if err != nil {
			logger.WithError(err).Error("Failed to retrieve post owner")
			http.Error(w, "Failed to retrieve post owner", http.StatusInternalServerError)
			return
		}

		if ownerID == 0 {
			logger.WithField("post_id", postID).Warn("Post not found")
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		if ownerID != userID {
			logger.WithFields(logrus.Fields{
				"post_id":  postID,
				"owner_id": ownerID,
				"user_id":  userID,
			}).Warn("Unauthorized deletion attempt")
			http.Error(w, "You are not authorized to delete this post", http.StatusForbidden)
			return
		}

		// Удаляем пост
		if err := database.DeletePost(db, postID); err != nil {
			logger.WithError(err).Error("Failed to delete post")
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
			return
		}

		// Формируем успешный ответ
		response := map[string]string{"message": "Post deleted successfully"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.WithError(err).Error("Failed to encode response")
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}