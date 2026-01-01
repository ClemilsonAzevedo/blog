package request

import "github.com/clemilsonazevedo/blog/internal/domain/enums"

type UserRegister struct {
	UserName string `json:"username" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
	Role     int    `json:"role" binding:"required,oneof=0 1 2"`
}

type UserUpdate struct {
	UserName string     `json:"username" binding:"omitempty,min=2,max=100"`
	Email    string     `json:"email" binding:"omitempty,email"`
	Role     enums.Role `json:"role" binding:"omitempty,oneof=0 1 2"`
}

type UserDelete struct {
	ID int `json:"id" binding:"required"`
}
