package controller

import (
	"encoding/json"
	"net/http"

	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/dto/request"
	"github.com/clemilsonazevedo/blog/internal/dto/response"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/go-chi/chi/v5"
	"go.bryk.io/pkg/ulid"
)

type CommentController struct {
	service *service.CommentService
}

func NewCommentController(service *service.CommentService) *CommentController {
	return &CommentController{
		service: service,
	}
}

func (cc *CommentController) CreateComment(w http.ResponseWriter, r *http.Request) {
	var dto request.CommentCreate
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dto.Content == "" || dto.UserID.String() == "" || dto.PostID.String() == "" {
		http.Error(w, "You need to provide all comments data", http.StatusBadRequest)
		return
	}

	commentId, err := ulid.New()
	if err != nil {
		http.Error(w, "Cannot Generate ULID to this Comment", http.StatusInternalServerError)
		return
	}

	Comment := entities.Comment{
		ID:      commentId,
		Content: dto.Content,
		UserID:  dto.UserID,
		PostID:  dto.PostID,
	}

	if err := cc.service.CreateComment(&Comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (cc *CommentController) GetCommentById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := ulid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	Comment, err := cc.service.GetCommentByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := response.CommentResponse{
		ID:        Comment.ID,
		Content:   Comment.Content,
		UserID:    Comment.UserID,
		PostID:    Comment.PostID,
		CreatedAt: Comment.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (cc *CommentController) GetCommentByPostID(w http.ResponseWriter, r *http.Request) {
	postIdStr := chi.URLParam(r, "postID")
	if postIdStr == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postId, err := ulid.Parse(postIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comments, err := cc.service.GetCommentsByPostID(postId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}

func (cc *CommentController) GetCommentByUserID(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "userID")
	if userIdStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userId, err := ulid.Parse(userIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comments, err := cc.service.GetCommentsByUserID(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}

func (cc *CommentController) GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	postIdStr := chi.URLParam(r, "postID")
	if postIdStr == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postId, err := ulid.Parse(postIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comments, err := cc.service.GetCommentsByPostID(postId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}

func (cc *CommentController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentIdStr := chi.URLParam(r, "id")
	if commentIdStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	commentId, err := ulid.Parse(commentIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := cc.service.DeleteComment(commentId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cc *CommentController) UpdateComment(w http.ResponseWriter, r *http.Request) {
	commentIdStr := chi.URLParam(r, "id")
	if commentIdStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var dto request.CommentUpdate
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	commentId, err := ulid.Parse(commentIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment := entities.Comment{
		ID:      commentId,
		UserID:  dto.UserID,
		PostID:  dto.PostID,
		Content: dto.Content,
	}

	if err := cc.service.UpdateComment(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := response.CommentResponse{
		ID:      comment.ID,
		UserID:  comment.UserID,
		PostID:  comment.PostID,
		Content: comment.Content,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
