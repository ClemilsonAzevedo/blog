package request

import (
	"go.bryk.io/pkg/ulid"
)

type CommentCreate struct {
	UserID  ulid.ULID `json:"userId"`
	PostID  ulid.ULID `json:"postId"`
	Content string    `json:"content"`
}

type CommentUpdate struct {
	ID      ulid.ULID `json:"id"`
	UserID  ulid.ULID `json:"userId"`
	PostID  ulid.ULID `json:"postId"`
	Content string    `json:"content"`
}

type CommentDelete struct {
	ID ulid.ULID `json:"id"`
}
