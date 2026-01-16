package response

import (
	"net/http"
	"time"

	"github.com/clemilsonazevedo/blog/pkg"
)

type CreatedCommentResponse struct {
	Message   string    `json:"message"`
	CommentId pkg.ULID  `json:"comment_id"`
	Timestamp time.Time `json:"timestamp"`
}

type DeletedCommentResponse struct {
	Message   string    `json:"message"`
	CommentId pkg.ULID  `json:"comment_id"`
	Timestamp time.Time `json:"timestamp"`
}

type CommentResponse struct {
	ID        pkg.ULID  `json:"id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	UserID    pkg.ULID  `json:"user_id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	PostID    pkg.ULID  `json:"post_id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ShowCommentResponse struct {
	Data      any       `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

func ShowComments(w http.ResponseWriter, commentsObj any) {
	OK(w, "success", commentsObj)
}

func CreatedComment(w http.ResponseWriter, commentId pkg.ULID) {
	WriteJSON(w, http.StatusCreated, CreatedCommentResponse{
		Message:   "Comment created successfully",
		CommentId: commentId,
		Timestamp: time.Now().UTC(),
	})
}

func DeletedComment(w http.ResponseWriter, commentId pkg.ULID) {
	WriteJSON(w, http.StatusOK, DeletedCommentResponse{
		Message:   "Comment deleted successfully",
		CommentId: commentId,
		Timestamp: time.Now().UTC(),
	})
}
