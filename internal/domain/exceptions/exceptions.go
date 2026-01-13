package exceptions

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error         string    `json:"error"`
	Message       string    `json:"message"`
	Details       any       `json:"details,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
	CorrelationID string    `json:"correlation_id,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func BadRequest(w http.ResponseWriter, err error, message string, details any) {
	resp := ErrorResponse{
		Error:     err.Error(),
		Message:   message,
		Details:   details,
		Timestamp: time.Now().UTC(),
	}
	writeJSON(w, http.StatusBadRequest, resp)
}

func Unauthorized(w http.ResponseWriter, message string) {
	resp := ErrorResponse{
		Error:     "unauthorized route",
		Message:   message,
		Timestamp: time.Now().UTC(),
	}
	writeJSON(w, http.StatusUnauthorized, resp)
}

func NotFound(w http.ResponseWriter, err error, message string) {
	resp := ErrorResponse{
		Error:     err.Error(),
		Message:   message,
		Timestamp: time.Now().UTC(),
	}
	writeJSON(w, http.StatusNotFound, resp)
}

func Conflict(w http.ResponseWriter, err error, message string) {
	resp := ErrorResponse{
		Error:     err.Error(),
		Message:   message,
		Timestamp: time.Now().UTC(),
	}
	writeJSON(w, http.StatusConflict, resp)
}

func InternalError(w http.ResponseWriter, err error, message string, correlationID string) {
	log.Printf("internal error [%s]: %v", correlationID, err)

	resp := ErrorResponse{
		Error:         err.Error(),
		Message:       message,
		Timestamp:     time.Now().UTC(),
		CorrelationID: correlationID,
	}
	writeJSON(w, http.StatusInternalServerError, resp)
}
