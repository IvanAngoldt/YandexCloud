package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func ToggleLike(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var likeRequest struct {
			PostID int `json:"postId"`
			UserID int `json:"userId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&likeRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Валидация входных данных
		if likeRequest.PostID <= 0 || likeRequest.UserID <= 0 {
			http.Error(w, "Invalid PostID or UserID", http.StatusBadRequest)
			return
		}

		// Проверяем, существует ли пост и получаем его автора
		var postAuthorID int
		err := db.QueryRow("SELECT author_id FROM posts WHERE id = $1", likeRequest.PostID).Scan(&postAuthorID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found", http.StatusNotFound)
			} else {
				log.Printf("Failed to check post: %v", err)
				http.Error(w, "Failed to check post", http.StatusInternalServerError)
			}
			return
		}

		notificationServiceURL := os.Getenv("NOTIFICATIONS_SERVICE_URL")
		if notificationServiceURL == "" {
			http.Error(w, "Notifications service URL is not configured", http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodPost:
			// Добавляем лайк
			if err := addLike(db, likeRequest.PostID, likeRequest.UserID); err != nil {
				log.Printf("Failed to add like: %v", err)
				http.Error(w, "Failed to add like", http.StatusInternalServerError)
				return
			}

			// Создаём уведомление через notifications_service
			notification := map[string]interface{}{
				"userId":  postAuthorID,
				"likerId": likeRequest.UserID,
				"postId":  likeRequest.PostID,
				"type":    "like",
				"message": fmt.Sprintf("User %d liked your post %d", likeRequest.UserID, likeRequest.PostID),
			}

			notificationData, err := json.Marshal(notification)
			if err != nil {
				log.Printf("Failed to marshal notification: %v", err)
				// Лайк уже добавлен, продолжаем выполнение.
			} else {
				req, err := http.NewRequest("POST", fmt.Sprintf("%s/notifications", notificationServiceURL), bytes.NewBuffer(notificationData))
				if err != nil {
					log.Printf("Failed to create notification request: %v", err)
					return
				}
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Failed to send notification: %v", err)
				} else {
					defer resp.Body.Close()
					if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
						log.Printf("Notification service responded with status: %s", resp.Status)
					} else {
						log.Println("Notification successfully sent")
					}
				}
			}

		case http.MethodDelete:
			// Удаляем лайк
			if err := removeLike(db, likeRequest.PostID, likeRequest.UserID); err != nil {
				log.Printf("Failed to remove like: %v", err)
				http.Error(w, "Failed to remove like", http.StatusInternalServerError)
				return
			}

			deleteNotificationRequest := map[string]interface{}{
				"userId":  postAuthorID,
				"likerId": likeRequest.UserID,
				"postId":  likeRequest.PostID,
				"type":    "like",
			}

			deleteNotificationData, err := json.Marshal(deleteNotificationRequest)
			if err != nil {
				log.Printf("Failed to marshal notification delete request: %v", err)
				// Удаление лайка успешно, уведомление не удалено.
			} else {
				req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/notifications", notificationServiceURL), bytes.NewBuffer(deleteNotificationData))
				if err != nil {
					log.Printf("Failed to create delete notification request: %v", err)
					return
				}
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Failed to send delete notification request: %v", err)
				} else {
					defer resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						log.Printf("Notification service responded with status: %s", resp.Status)
					} else {
						log.Println("Notification successfully deleted")
					}
				}
			}

		default:
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		// Получаем обновленный список пользователей, лайкнувших пост
		userIDs, err := getLikes(db, likeRequest.PostID)
		if err != nil {
			log.Printf("Failed to fetch likes: %v", err)
			http.Error(w, "Failed to fetch likes", http.StatusInternalServerError)
			return
		}

		users, err := fetchUsersFromUsersService(userIDs)
		if err != nil {
			log.Printf("Failed to fetch user info: %v", err)
			http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func GetLikesForPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := r.URL.Query().Get("postId")
		if postIDStr == "" {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}

		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid Post ID", http.StatusBadRequest)
			return
		}

		userIDs, err := getLikes(db, postID)
		if err != nil {
			http.Error(w, "Failed to fetch likes", http.StatusInternalServerError)
			return
		}

		users, err := fetchUsersFromUsersService(userIDs)
		if err != nil {
			http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func addLike(db *sql.DB, postID, userID int) error {
	_, err := db.Exec(`
		INSERT INTO likes (post_id, user_id) 
		VALUES ($1, $2) 
		ON CONFLICT (post_id, user_id) DO NOTHING
	`, postID, userID)
	return err
}

func removeLike(db *sql.DB, postID, userID int) error {
	_, err := db.Exec(`
		DELETE FROM likes WHERE post_id = $1 AND user_id = $2
	`, postID, userID)
	return err
}

// getLikes возвращает список user_id, лайкнувших пост
func getLikes(db *sql.DB, postID int) ([]int, error) {
	rows, err := db.Query(`
        SELECT user_id
        FROM likes
        WHERE post_id = $1
    `, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var uid int
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, uid)
	}
	return userIDs, nil
}
