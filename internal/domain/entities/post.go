package entities

import (
	"time"

	"github.com/clemilsonazevedo/blog/pkg"
)

type Post struct {
	ID        pkg.ULID  `gorm:"column:id;primaryKey;type:varchar(26);not null" json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Title     string    `gorm:"column:title;index;not null" json:"title"`
	Content   string    `gorm:"column:content;not null" json:"content"`
	Slug      string    `gorm:"uniqueIndex:idx_posts_slug;not null;size:300;" json:"slug"`
	Likes     int       `gorm:"column:likes;not null;default:0" json:"likes"`
	Views     int       `gorm:"column:views;not null;default:0" json:"views"`
	Dislikes  int       `gorm:"column:dislikes;not null;default:0" json:"dislikes"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UserID    pkg.ULID  `gorm:"column:user_id;index;not null" json:"user_id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`

	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (Post) TableName() string {
	return "posts"
}

func (post Post) GetID() any {
	return post.ID
}
