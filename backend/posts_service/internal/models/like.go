package models

type Like struct {
	PostID int `json:"postId"`
	UserID int `json:"userId"`
}
