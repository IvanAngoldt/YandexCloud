// internal/models/notification.go

package models

import "time"

// Notification представляет структуру уведомления
type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`    // ID пользователя, которому адресовано уведомление
	PostID    int       `json:"postId"`    // ID поста, к которому относится уведомление
	Message   string    `json:"message"`   // Текст уведомления
	IsRead    bool      `json:"isRead"`    // Статус прочтения уведомления
	CreatedAt time.Time `json:"createdAt"` // Время создания уведомления
	LikerID   int       `json:"likerId"`   // ID пользователя, который поставил лайк
	Type      string    `json:"type"`      // Тип уведомления (например, "like", "comment")
}
