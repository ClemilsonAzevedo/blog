package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	Title     string    `gorm:"column:title;index;not null"`
	Content   string    `gorm:"column:content;not null"`
	Likes     int       `gorm:"column:likes;not null;default:0"`
	Dislikes  int       `gorm:"column:dislikes;not null;default:0"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UserID    uuid.UUID `gorm:"column:user_id;index;not null"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}

func (Post) TableName() string {
	return "posts"
}

func (post* Post)GetID() any {
	return post.ID
}
