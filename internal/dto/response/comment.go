package response

import (
	"time"

	"go.bryk.io/pkg/ulid"
)

type CommentResponse struct {
	ID        ulid.ULID `json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	UserID    ulid.ULID `json:"user_id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	PostID    ulid.ULID `json:"post_id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
