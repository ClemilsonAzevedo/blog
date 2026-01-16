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
	var data request.CommentCreate
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&data); err != nil {
		exceptions.BadRequest(w, err, "Cannot Decode Body", nil)
		return
	}

	if dec.More() {
		exceptions.BadRequest(w, errors.New("Multiple JSON values not allowed"), "multiple JSON values not allowed", nil)
		return
	}

	if data.Content == "" || data.UserID.String() == "" || data.PostID.String() == "" {
		exceptions.BadRequest(w, errors.New("Request Error"), "You need to provide all comments data", &data)
		return
	}

	contextUser, ok := r.Context().Value("user").(*entities.User)
	if !ok {
		exceptions.Unauthorized(w, "unauthorized")
		return
	}

	//Fica aqui até descobrir uma forma melhor de fazer isso
	if contextUser.ID != data.UserID {
		exceptions.Unauthorized(w, "User Does not exists")
		return
	}

	commentId, err := pkg.NewULID()
	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot Generate ULID to this Comment", reqId)
		return
	}

	Comment := entities.Comment{
		ID:      commentId,
		Content: data.Content,
		UserID:  data.UserID,
		PostID:  data.PostID,
	}

	if err := cc.service.CreateComment(&Comment); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			exceptions.BadRequest(w, err, "This post does not exists", data.PostID)
			return
		}

		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot create comment to this post", reqId)
		return
	}

	response.CreatedComment(w, commentId)
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
func (cc *CommentController) GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
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
	// Trocar para pegar o user ID e comment Id do body para poder verificar se o usuario que criou é o que está a tentar deletar | e permitir o autor deletar qualquer comentario
	commentIdStr := r.URL.Query().Get("commentId")
	if commentIdStr == "" {
		exceptions.BadRequest(w, errors.New("Request Error"), "You need Provide Comment Id on route", commentIdStr)
		return
	}

	commentId, err := pkg.ParseULID(commentIdStr)
	if err != nil {
		exceptions.BadRequest(w, err, "Cannot Parse Comment Id", commentId)
		return
	}

	comment, err := cc.service.GetCommentByID(commentId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			exceptions.NotFound(w, err, "Comment Does not exists")
			return
		}

		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot get comment", reqId)
		return
	}

	if err := cc.service.DeleteComment(comment.ID); err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot Delete comment", reqId)
		return
	}

	response.DeletedComment(w, comment.ID)
}
