package response

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/clemilsonazevedo/blog/internal/domain/enums"
	"github.com/clemilsonazevedo/blog/pkg"
)

type UserLogin struct {
	Token   string `json:"token"`
	Expires int64  `json:"exp"`
}

type UserResponse struct {
	Message   string    `json:"message"`
	UserID    string    `json:"user_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type SuccessResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
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

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func CreatedUser(w http.ResponseWriter, userID string) {
	WriteJSON(w, http.StatusCreated, UserResponse{
		Message:   "user created successfully",
		UserID:    userID,
		Timestamp: time.Now().UTC(),
	})
}

func OK[T any](w http.ResponseWriter, message string, data T) {
	WriteJSON(w, http.StatusOK, SuccessResponse[T]{
		Message: message,
		Data:    data,
	})
}
