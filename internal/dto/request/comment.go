package request

import (
	"go.bryk.io/pkg/ulid"
)

type CommentCreate struct {
	UserID  ulid.ULID `json:"userId" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	PostID  ulid.ULID `json:"postId" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Content string    `json:"content"`
}

type CommentUpdate struct {
	ID      ulid.ULID `json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	UserID  ulid.ULID `json:"userId" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	PostID  ulid.ULID `json:"postId" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Content string    `json:"content"`
}

type CommentDelete struct {
	ID ulid.ULID `json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
}
