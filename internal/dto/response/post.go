package response

import (
	"math"
	"net/http"
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

type MetaInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

type ListPostsResponse struct {
	Data      any       `json:"data"`
	Meta      MetaInfo  `json:"meta"`
	Timestamp time.Time `json:"timestamp"`
}

type ShowPostResponse struct {
	Data      any       `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

func ListPosts(w http.ResponseWriter, postsObj any, page, limit, total int) {
	if limit <= 0 {
		limit = 1
	}
	if limit > 25 {
		limit = 25
	}
	if page <= 0 {
		page = 1
	}

	resp := ListPostsResponse{
		Data: postsObj,
		Meta: MetaInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: int(math.Ceil(float64(total) / float64(limit))),
		},
		Timestamp: time.Now().UTC(),
	}

	WriteJSON(w, http.StatusOK, resp)
}

func ShowPost(w http.ResponseWriter, postsObj any) {
	resp := ShowPostResponse{
		Data:      postsObj,
		Timestamp: time.Now().UTC(),
	}

	WriteJSON(w, http.StatusOK, resp)
}
