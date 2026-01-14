package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/domain/exceptions"
	"github.com/clemilsonazevedo/blog/internal/dto/request"
	"github.com/clemilsonazevedo/blog/internal/dto/response"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/clemilsonazevedo/blog/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

type CommentController struct {
	service *service.CommentService
}

func NewCommentController(service *service.CommentService) *CommentController {
	return &CommentController{
		service: service,
	}
}

// CreateComment godoc
// @Summary Create a new comment
// @Description Creates a new comment on a post (authentication required)
// @Tags Comments
// @Accept json
// @Produce json
// @Param request body request.CommentCreate true "Comment creation data"
// @Success 201 {string} string "Comment created"
// @Failure 400 {string} string "You need to provide all comments data"
// @Failure 500 {string} string "Cannot create comment"
// @Security CookieAuth
// @Router /comment [post]
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

	commentId, err := pkg.NewULID()
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

// GetCommentById godoc
// @Summary Get comment by ID
// @Description Retrieves a single comment by its ULID (authentication required)
// @Tags Comments
// @Produce json
// @Param id path string true "Comment ULID"
// @Success 200 {object} response.CommentResponse
// @Failure 400 {string} string "ID is required"
// @Failure 500 {string} string "Error retrieving comment"
// @Security CookieAuth
// @Router /comments [get]
func (cc *CommentController) GetCommentById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := pkg.ParseULID(idStr)
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

// GetCommentByPostID godoc
// @Summary Get comments by post ID
// @Description Retrieves all comments for a specific post
// @Tags Comments
// @Produce json
// @Param postID path string true "Post ULID"
// @Success 200 {array} response.CommentResponse
// @Failure 400 {string} string "Post ID is required"
// @Failure 500 {string} string "Error retrieving comments"
// @Router /comments/{postID} [get]
func (cc *CommentController) GetCommentByPostID(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.URL.Query().Get("postId")
	if postIdStr == "" {
		exceptions.BadRequest(w, errors.New("Request Error"), "You need Provide Post Id on route", postIdStr)
		return
	}

	postId, err := pkg.ParseULID(postIdStr)
	if err != nil {
		exceptions.BadRequest(w, err, "Cannot Parse Post Id", postId)
		return
	}

	comments, err := cc.service.GetCommentsByPostID(postId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			exceptions.NotFound(w, err, fmt.Sprintf("Post with id %v not found", postId))
			return
		}
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot get comments of this post", reqId)
		return
	}

	commentsObj := make([]response.CommentResponse, len(comments))
	for i := range len(comments) {
		p := comments[i]
		commentsObj[i] = response.CommentResponse{
			ID:        p.ID,
			PostID:    p.PostID,
			UserID:    p.UserID,
			Content:   p.Content,
			CreatedAt: p.CreatedAt,
		}
	}

	response.ShowComments(w, commentsObj)
}

// GetCommentByUserID godoc
// @Summary Get comments by user ID
// @Description Retrieves all comments made by a specific user
// @Tags Comments
// @Produce json
// @Param userID path string true "User ULID"
// @Success 200 {array} response.CommentResponse
// @Failure 400 {string} string "User ID is required"
// @Failure 500 {string} string "Error retrieving comments"
// @Security CookieAuth
// @Router /comments/user/{userID} [get]
func (cc *CommentController) GetCommentByUserID(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "userID")
	if userIdStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userId, err := pkg.ParseULID(userIdStr)
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

// GetCommentsByPostID godoc
// @Summary Get comments by post ID (alternative)
// @Description Retrieves all comments for a specific post
// @Tags Comments
// @Produce json
// @Param postID path string true "Post ULID"
// @Success 200 {array} response.CommentResponse
// @Failure 400 {string} string "Post ID is required"
// @Failure 500 {string} string "Error retrieving comments"
// @Router /post/{postID}/comments [get]
func (cc *CommentController) GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	postIdStr := chi.URLParam(r, "postID")
	if postIdStr == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postId, err := pkg.ParseULID(postIdStr)
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

// DeleteComment godoc
// @Summary Delete a comment
// @Description Deletes an existing comment
// @Tags Comments
// @Param id path string true "Comment ULID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "ID is required"
// @Failure 500 {string} string "Error deleting comment"
// @Security CookieAuth
// @Router /comment/{id} [delete]
func (cc *CommentController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentIdStr := chi.URLParam(r, "id")
	if commentIdStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	commentId, err := pkg.ParseULID(commentIdStr)
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

// UpdateComment godoc
// @Summary Update a comment
// @Description Updates an existing comment
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ULID"
// @Param request body request.CommentUpdate true "Comment update data"
// @Success 200 {object} response.CommentResponse
// @Failure 400 {string} string "ID is required"
// @Failure 500 {string} string "Error updating comment"
// @Security CookieAuth
// @Router /comment/{id} [put]
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

	commentId, err := pkg.ParseULID(commentIdStr)
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
