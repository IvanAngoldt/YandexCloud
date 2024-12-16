package models

import "time"

// Notification представляет структуру уведомления
type Notification struct {
	ID            int       `json:"id"`
	UserID        int       `json:"userId"`
	PostID        *int      `json:"postId,omitempty"`
	Message       string    `json:"message"`
	IsRead        bool      `json:"isRead"`
	CreatedAt     time.Time `json:"createdAt"`
	LikerID       *int      `json:"likerId,omitempty"`
	LikerUsername string    `json:"likerUsername,omitempty"`
	Type          string    `json:"type"`
}
