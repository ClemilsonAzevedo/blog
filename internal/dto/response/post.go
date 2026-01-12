package response

import (
	"time"

	"github.com/clemilsonazevedo/blog/pkg"
)

type PostResponse struct {
	ID        pkg.ULID  `json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	AuthorId  pkg.ULID  `json:"author_id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Slug      string    `json:"slug"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	Views     int       `json:"views"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
