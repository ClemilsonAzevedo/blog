package response

import (
	"github.com/clemilsonazevedo/blog/internal/domain/enums"
	"go.bryk.io/pkg/ulid"
)

type UserLogin struct {
	Token   string `json:"token"`
	Expires int64  `json:"exp"`
}

type UserByID struct {
	ID       ulid.ULID  `json:"id"`
	UserName string     `json:"username"`
	Email    string     `json:"email"`
	Role     enums.Role `json:"role" binding:"required,oneof=anonymous reader author"`
}

type UserLogout struct {
	Message string `json:"message"`
}
