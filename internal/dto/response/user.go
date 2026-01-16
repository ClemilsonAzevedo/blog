package response

import (
	"net/http"
	"time"

	"github.com/clemilsonazevedo/blog/internal/domain/enums"
	"github.com/clemilsonazevedo/blog/pkg"
)

type UserLogin struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"exp"`
}

type UserLogout struct {
	Message string `json:"message"`
}

type UserResponse struct {
	Message   string    `json:"message"`
	UserID    string    `json:"user_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type UserProfile struct {
	ID       pkg.ULID   `json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	UserName string     `json:"username"`
	Email    string     `json:"email"`
	Role     enums.Role `json:"role" binding:"required,oneof=anonymous reader author" swaggertype:"string" enums:"anonymous,reader,author"`
}

type UserDeleted struct {
	ID pkg.ULID `json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
}

func CreatedUser(w http.ResponseWriter, userID string) {
	WriteJSON(w, http.StatusCreated, UserResponse{
		Message:   "user created successfully",
		UserID:    userID,
		Timestamp: time.Now().UTC(),
	})
}
