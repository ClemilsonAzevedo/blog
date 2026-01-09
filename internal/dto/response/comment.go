package response

import (
	"time"

	"go.bryk.io/pkg/ulid"
)

type CommentResponse struct {
	ID        ulid.ULID `json:"id"`
	UserID    ulid.ULID `json:"user_id"`
	PostID    ulid.ULID `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
