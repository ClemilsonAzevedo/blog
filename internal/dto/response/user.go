package response

import "github.com/google/uuid"

type UserLogin struct {
	Token   string `json:"token"`
	Expires int64  `json:"exp"`
}

type UserByID struct {
	ID       uuid.UUID `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type UserLogout struct {
	Message string `json:"message"`
}
