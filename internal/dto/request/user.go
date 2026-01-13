package request

import (
	"github.com/clemilsonazevedo/blog/internal/domain/enums"
	"github.com/clemilsonazevedo/blog/pkg"
)

type UserRegister struct {
	UserName string `json:"username" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

type UserUpdate struct {
	ID       pkg.ULID   `json:"id" binding:"required,min=1" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	UserName string     `json:"username" binding:"omitempty,min=2,max=100"`
	Email    string     `json:"email" binding:"omitempty,email"`
	Role     enums.Role `json:"role" binding:"required,oneof=anonymous reader author" swaggertype:"string" enums:"anonymous,reader,author"`
}

type UserDelete struct {
	ID int `json:"id" binding:"required"`
}
