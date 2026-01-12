package response

import (
	"github.com/clemilsonazevedo/blog/internal/domain/enums"
	"github.com/clemilsonazevedo/blog/pkg"
)

type UserLogin struct {
	Token   string `json:"token"`
	Expires int64  `json:"exp"`
}

type UserByID struct {
	ID       pkg.ULID   `json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	UserName string     `json:"username"`
	Email    string     `json:"email"`
	Role     enums.Role `json:"role" binding:"required,oneof=anonymous reader author" swaggertype:"string" enums:"anonymous,reader,author"`
}

type UserLogout struct {
	Message string `json:"message"`
}
