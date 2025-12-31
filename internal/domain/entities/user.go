package entities

import (
	"time"

	"github.com/clemilsonazevedo/blog/internal/domain/enums"
	"github.com/google/uuid"
)

type Role = enums.Role

type User struct {
	ID        uuid.UUID `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	UserName  string    `gorm:"column:username;unique;not null;"`
	Email     string    `gorm:"column:email;unique;not null;"`
	Password  string    `gorm:"column:password;not null"`
	Role      Role      `gorm:"column:role"`
	CreatedAt time.Time `gorm:"column:created_at;not null,autoCreateTime"`

	Posts []Post `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (User) TableName() string {
	return "users"
}

func (user* User)GetID() any {
	return user.ID
}
