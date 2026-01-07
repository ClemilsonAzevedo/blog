package request

import (
	"github.com/google/uuid"
)

type CommentCreate struct {
	UserID  uuid.UUID `json:"userId"`
	PostID  uuid.UUID `json:"postId"`
	Content string    `json:"content"`
}

type CommentUpdate struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"userId"`
	PostID  uuid.UUID `json:"postId"`
	Content string    `json:"content"`
}

type CommentDelete struct {
	ID uuid.UUID `json:"id"`
}
