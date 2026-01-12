package request

import "github.com/clemilsonazevedo/blog/pkg"

type CommentCreate struct {
	UserID  pkg.ULID `json:"userId" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	PostID  pkg.ULID `json:"postId" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Content string   `json:"content"`
}

type CommentUpdate struct {
	ID      pkg.ULID `json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	UserID  pkg.ULID `json:"userId" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	PostID  pkg.ULID `json:"postId" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Content string   `json:"content"`
}

type CommentDelete struct {
	ID pkg.ULID `json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
}
