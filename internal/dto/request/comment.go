package request

import (
	"github.com/google/uuid"
)

type CommentCreate struct {
	UserID  uuid.UUID `json:"user_id"`
	PostID  uuid.UUID `json:"post_id"`
	Content string    `json:"content"`
}

type CommentUpdate struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"user_id"`
	PostID  uuid.UUID `json:"post_id"`
	Content string    `json:"content"`
}

type CommentDelete struct {
	ID uuid.UUID `json:"id"`
}
