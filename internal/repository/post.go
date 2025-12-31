package repository

import (
	"time"

	"github.com/clemilsonazevedo/blog/internal/dto/response"
	"github.com/google/uuid"
)

func FindPostBySlug(slug string) (response.PostResponse, error) {
	post := response.PostResponse{
		ID:        uuid.MustParse("f955886a-b45d-401b-bc13-a964db181d53"),
		UserID:    uuid.MustParse("ce70b183-ca52-4189-8fd4-15cef5967843"),
		Title:     "Primeiro Post",
		Content:   "### Titulo",
		Likes:     20,
		Dislikes:  0,
		CreatedAt: time.Now(),
	}
	return post, nil
}
