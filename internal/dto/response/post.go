package response

import (
	"time"

	"go.bryk.io/pkg/ulid"
)

type PostResponse struct {
	ID        ulid.ULID `json:"id"`
	AuthorId  ulid.ULID `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Slug      string    `json:"slug"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	Views     int       `json:"views"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
