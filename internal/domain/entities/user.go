package entities

import (
	"time"

	"github.com/clemilsonazevedo/blog/internal/domain/enums"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Role = enums.Role

type User struct {
	ID        uuid.UUID `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserName  string    `gorm:"column:username;unique;not null;" json:"username"`
	Email     string    `gorm:"column:email;unique;not null;" json:"email"`
	Password  string    `gorm:"column:password;not null" json:"password"`
	Role      Role      `gorm:"column:role" json:"role"`
	CreatedAt time.Time `gorm:"column:created_at;not null,autoCreateTime" json:"created_at"`

	// Posts []Post `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (User) TableName() string {
	return "users"
}

func (user User) GetID() any {
	return user.ID
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		14,
	)
	return string(bytes), err
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
	return err == nil
}
