package request

import (
	"github.com/clemilsonazevedo/blog/internal/domain/enums"
	"github.com/google/uuid"
)

type UserRegister struct {
	UserName string `json:"username" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
	Role     enums.Role    `json:"role" binding:"required,oneof=anonymous reader author"`
}

type UserUpdate struct {
	ID       uuid.UUID `json:"id" binding:"required,min=1"`
	UserName string    `json:"username" binding:"omitempty,min=2,max=100"`
	Email    string     `json:"email" binding:"omitempty,email"`
	Role     enums.Role    `json:"role" binding:"required,oneof=anonymous reader author"`
}

type UserDelete struct {
	ID int `json:"id" binding:"required"`
}
