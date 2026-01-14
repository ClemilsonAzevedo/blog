package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func OK[T any](w http.ResponseWriter, message string, data T) {
	WriteJSON(w, http.StatusOK, SuccessResponse[T]{
		Message: message,
		Data:    data,
	})
}
