package request

import "github.com/clemilsonazevedo/blog/pkg"

type PostCreate struct {
	Title    string   `json:"title" binding:"required,min=2,max=100"`
	Content  string   `json:"content" binding:"required,min=2,max=1000"`
	AuthorId pkg.ULID `json:"author_id" binding:"required,min=1" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
}

type AiPostCreate struct {
	Content  string   `json:"content" binding:"required,min=2,max=1000"`
	AuthorId pkg.ULID `json:"author_id" binding:"required,min=1" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
}

type PostUpdate struct {
	ID       pkg.ULID `json:"id" binding:"required,min=1" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Title    string   `json:"title" binding:"required,min=2,max=100"`
	Content  string   `json:"content" binding:"required,min=2,max=1000"`
	Likes    int      `json:"likes" binding:"required,min=0"`
	Dislikes int      `json:"dislikes" binding:"required,min=0"`
}

type PostDelete struct {
	ID int `json:"id" binding:"required,min=1"`
}
