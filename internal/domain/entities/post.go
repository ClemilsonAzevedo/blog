package entities

import (
	"time"

	"github.com/clemilsonazevedo/blog/pkg"
)

type Post struct {
	ID pkg.ULID `gorm:"column:id;primaryKey;type:varchar(26);not null" json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`

	Title   string `gorm:"column:title;index;not null" json:"title"`
	Slug    string `gorm:"column:slug;uniqueIndex:idx_posts_slug;not null;size:300" json:"slug"`
	Content string `gorm:"column:content;not null" json:"content"`

	Likes    int `gorm:"column:likes;not null;default:0" json:"likes"`
	Views    int `gorm:"column:views;not null;default:0" json:"views"`
	Dislikes int `gorm:"column:dislikes;not null;default:0" json:"dislikes"`

	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`

	AuthorId pkg.ULID `gorm:"column:author_id;type:varchar(26);index;not null" json:"author_id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Author   User     `gorm:"foreignKey:AuthorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
}

func (Post) TableName() string {
	return "posts"
}

func (post Post) GetID() any {
	return post.ID
}
