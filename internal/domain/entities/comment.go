package entities

import (
	"time"

	"go.bryk.io/pkg/ulid"
)

type Comment struct {
	ID        ulid.ULID `gorm:"column:id;primaryKey;type:VARCHAR(26);not null" json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UserID    ulid.ULID `gorm:"column:user_id;type:VARCHAR(26);not null" json:"user_id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	PostID    ulid.ULID `gorm:"column:post_id;type:VARCHAR(26);index;not null" json:"post_id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`

	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Post Post `gorm:"foreignKey:PostID;references:ID" json:"post,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

func (comment Comment) GetID() any {
	return comment.ID
}
