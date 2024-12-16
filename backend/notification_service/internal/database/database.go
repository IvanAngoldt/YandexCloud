package database

import (
	"database/sql"
	"fmt"
	"log"
	"notifications_service/internal/models"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Connect устанавливает соединение с базой данных
func Connect() (*sql.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// GetNotifications извлекает уведомления для указанного пользователя из базы данных
func GetNotifications(db *sql.DB, userId string) ([]models.Notification, error) {
	rows, err := db.Query(`
		SELECT 
			n.id, 
			n.user_id, 
			n.message, 
			n.is_read, 
			n.created_at, 
			n.type,
			COALESCE(nl.liker_id, 0) AS liker_id, 
			COALESCE(nl.post_id, 0) AS post_id,
			COALESCE(u.username, '') AS liker_username
		FROM notifications n
		LEFT JOIN notification_like nl ON n.id = nl.notification_id
		LEFT JOIN users u ON nl.liker_id = u.id
		WHERE n.user_id = $1 
		ORDER BY n.created_at DESC
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		var likerID, postID sql.NullInt64 // Лайки могут быть NULL для других типов уведомлений
		var likerUsername sql.NullString  // Имя пользователя может быть NULL

		if err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Message,
			&notification.IsRead,
			&notification.CreatedAt,
			&notification.Type,
			&likerID,
			&postID,
			&likerUsername,
		); err != nil {
			return nil, err
		}

		// Обработка NULL значений
		if likerID.Valid {
			liker := int(likerID.Int64)
			notification.LikerID = &liker
		}
		if postID.Valid {
			post := int(postID.Int64)
			notification.PostID = &post
		}
		if likerUsername.Valid {
			notification.LikerUsername = likerUsername.String
		}

		notifications = append(notifications, notification)
	}
	return notifications, nil
}

// MarkAsRead помечает уведомление как прочитанное
func MarkAsRead(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE notifications SET is_read = true WHERE id = $1", id)
	return err
}

// DeleteNotification удаляет уведомление пользователя
func DeleteNotification(db *sql.DB, userID int, likerID int, postID int, notificationType string) error {
	_, err := db.Exec(`
		DELETE FROM notifications n
		USING notification_like nl
		WHERE n.id = nl.notification_id
		  AND n.user_id = $1
		  AND nl.liker_id = $2
		  AND nl.post_id = $3
		  AND n.type = $4
	`, userID, likerID, postID, notificationType)
	return err
}

// ClearNotifications очищает все уведомления пользователя
func ClearNotifications(db *sql.DB, userId string) error {
	_, err := db.Exec("DELETE FROM notifications WHERE user_id = $1", userId)
	return err
}

// AddNotification добавляет новое уведомление в базу данных с учетом новых полей
func AddNotification(db *sql.DB, notification models.Notification, likerID int, postID int) error {
	// Проверяем, не является ли лайкер автором поста
	if notification.Type == "like" {
		var postAuthorID int
		err := db.QueryRow("SELECT author_id FROM posts WHERE id = $1", postID).Scan(&postAuthorID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("post not found")
			}
			return fmt.Errorf("failed to retrieve post author: %w", err)
		}

		if likerID == postAuthorID {
			return fmt.Errorf("notification not added: author cannot send notification to themselves")
		}
	}

	// Вставляем запись в таблицу notifications
	var notificationID int
	err := db.QueryRow(`
        INSERT INTO notifications (user_id, message, is_read, created_at, type)
        VALUES ($1, $2, $3, $4, $5) RETURNING id
    `, notification.UserID, notification.Message, notification.IsRead, notification.CreatedAt, notification.Type).Scan(&notificationID)
	if err != nil {
		return fmt.Errorf("failed to insert notification: %w", err)
	}

	// Если уведомление типа "like", добавляем запись в notification_like
	if notification.Type == "like" {
		_, err = db.Exec(`
            INSERT INTO notification_like (notification_id, liker_id, post_id)
            VALUES ($1, $2, $3)
        `, notificationID, likerID, postID)
		if err != nil {
			return fmt.Errorf("failed to insert notification_like: %w", err)
		}

		// Обновляем поля в объекте Notification
		notification.LikerID = &likerID
		notification.PostID = &postID
	}

	log.Printf("Notification successfully added: %+v", notification)
	return nil
}
