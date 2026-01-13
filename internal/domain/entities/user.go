package entities

import (
	"time"

	"github.com/clemilsonazevedo/blog/internal/domain/enums"
	"github.com/clemilsonazevedo/blog/pkg"
	"golang.org/x/crypto/bcrypt"
)

type Role = enums.Role

type User struct {
	ID pkg.ULID `gorm:"column:id;primaryKey;type:varchar(26);not null" json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`

	UserName  string    `gorm:"column:username;unique;not null" json:"username"`
	Email     string    `gorm:"column:email;unique;not null" json:"email"`
	Password  string    `gorm:"column:password;not null" json:"password"`
	Role      Role      `gorm:"type:user_role;default:'reader'" json:"role" swaggertype:"string" enums:"anonymous,reader,author"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
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
