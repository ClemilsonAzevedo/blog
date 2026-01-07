package request

import "github.com/google/uuid"

type PostCreate struct {
	Title    string    `json:"title" binding:"required,min=2,max=100"`
	Content  string    `json:"content" binding:"required,min=2,max=1000"`
	Likes    int       `json:"likes" binding:"required,min=0"`
	Dislikes int       `json:"dislikes" binding:"required,min=0"`
	UserID   uuid.UUID `json:"userId" binding:"required,min=1"`
}

type PostUpdate struct {
	ID       uuid.UUID `json:"id" binding:"required,min=1"`
	Title    string    `json:"title" binding:"required,min=2,max=100"`
	Content  string    `json:"content" binding:"required,min=2,max=1000"`
	Likes    int       `json:"likes" binding:"required,min=0"`
	Dislikes int       `json:"dislikes" binding:"required,min=0"`
}

type PostDelete struct {
	ID int `json:"id" binding:"required,min=1"`
}
