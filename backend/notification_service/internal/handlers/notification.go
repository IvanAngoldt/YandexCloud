// internal/handlers/notification.go

package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"notifications_service/internal/database"
	"notifications_service/internal/models"
	"time"

	"github.com/gorilla/mux"
)

// CreateNotificationRequest представляет структуру входящих данных для создания уведомления
type CreateNotificationRequest struct {
	UserID  int    `json:"userId"`  // ID пользователя, которому адресовано уведомление
	LikerID int    `json:"likerId"` // ID пользователя, который поставил лайк
	PostID  int    `json:"postId"`  // ID поста, к которому относится уведомление
	Type    string `json:"type"`    // Тип уведомления (например, "like", "comment")
	Message string `json:"message"` // Текст уведомления
}

// CreateNotification обрабатывает запросы на создание нового уведомления
func CreateNotification(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateNotificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Валидация входных данных
		if req.UserID <= 0 || req.Type == "" {
			http.Error(w, "Invalid notification data", http.StatusBadRequest)
			return
		}

		notification := models.Notification{
			UserID:    req.UserID,
			Message:   req.Message,
			IsRead:    false,
			Type:      req.Type,
			CreatedAt: time.Now(),
		}

		// Добавляем уведомление в базу данных
		if err := database.AddNotification(db, notification, req.LikerID, req.PostID); err != nil {
			log.Printf("Failed to add notification: %v", err)

			if err.Error() == "notification not added: author cannot send notification to themselves" {
				http.Error(w, "Notification not added: Author cannot send notification to themselves", http.StatusBadRequest)
				return
			}

			http.Error(w, "Failed to add notification", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Println("Notification successfully created")
	}
}

// FetchNotifications обрабатывает запросы на получение уведомлений пользователя
func FetchNotifications(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("userId")
		if userID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		// Извлекаем уведомления из базы данных
		notifications, err := database.GetNotifications(db, userID)
		if err != nil {
			log.Printf("Failed to fetch notifications: %v", err)
			http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
			return
		}

		// Отправляем список уведомлений как JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(notifications)
	}
}

// MarkNotificationAsRead обрабатывает запросы на пометку уведомления как прочитанного
func MarkNotificationAsRead(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Notification ID is required", http.StatusBadRequest)
			return
		}

		// Отметим уведомление как прочитанное
		if err := database.MarkAsRead(db, id); err != nil {
			log.Printf("Failed to mark notification as read: %v", err)
			http.Error(w, "Failed to mark notification as read", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Println("Notification marked as read")
	}
}

// ClearNotifications обрабатывает запросы на очистку всех уведомлений пользователя
func ClearNotifications(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, exists := vars["userId"]
		if !exists || userID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		// Очистим все уведомления для пользователя
		if err := database.ClearNotifications(db, userID); err != nil {
			log.Printf("Failed to clear notifications: %v", err)
			http.Error(w, "Failed to clear notifications", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Println("Notifications cleared")
	}
}

// DeleteNotification обрабатывает запросы на удаление уведомления пользователя
func DeleteNotification(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var deleteRequest struct {
			UserID  int    `json:"userId"`
			LikerID int    `json:"likerId"`
			PostID  int    `json:"postId"`
			Type    string `json:"type"`
		}

		if err := json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Валидация входных данных
		if deleteRequest.UserID <= 0 || deleteRequest.LikerID <= 0 || deleteRequest.PostID <= 0 || deleteRequest.Type == "" {
			http.Error(w, "Invalid or missing notification data", http.StatusBadRequest)
			return
		}

		// Удаляем конкретное уведомление с учетом параметров
		if err := database.DeleteNotification(db, deleteRequest.UserID, deleteRequest.LikerID, deleteRequest.PostID, deleteRequest.Type); err != nil {
			log.Printf("Failed to delete notification: %v", err)
			http.Error(w, "Failed to delete notification", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Println("Notifications deleted")
	}
}
