package entities

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UserID    uuid.UUID `gorm:"column:user_id;index;not null"`
	PostID    uuid.UUID `gorm:"column:post_id;index;not null"`

	User User `gorm:"foreignKey:UserID;references:ID"`
	Post Post `gorm:"foreignKey:PostID;references:ID"`
}

func (Comment) TableName() string {
	return "comments"
}

func (comment* Comment)GetID() any {
	return comment.ID
}
