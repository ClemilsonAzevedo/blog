package entities

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;index;not null" json:"user_id"`
	PostID    uuid.UUID `gorm:"column:post_id;type:uuid;index;not null" json:"post_id"`

	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Post Post `gorm:"foreignKey:PostID;references:ID" json:"post,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

func (comment Comment) GetID() any {
	return comment.ID
}
