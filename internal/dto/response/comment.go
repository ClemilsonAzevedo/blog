package response

import (
	"time"

	"github.com/clemilsonazevedo/blog/pkg"
)

type CommentResponse struct {
	ID        pkg.ULID  `json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	UserID    pkg.ULID  `json:"user_id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	PostID    pkg.ULID  `json:"post_id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
