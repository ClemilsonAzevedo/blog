package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/domain/exceptions"
	"github.com/clemilsonazevedo/blog/internal/dto/request"
	"github.com/clemilsonazevedo/blog/internal/dto/response"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/clemilsonazevedo/blog/pkg"
	"github.com/clemilsonazevedo/blog/tools"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

type PostController struct {
	service *service.PostService
}

func NewPostController(service *service.PostService) *PostController {
	return &PostController{
		service: service,
	}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Creates a new blog post (Author role required)
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body request.PostCreate true "Post creation data"
// @Success 201 {string} string "Post created"
// @Failure 400 {string} string "You need set Content and AuthorId to create a post"
// @Failure 500 {string} string "Cannot create post"
// @Security CookieAuth
// @Router /post [post]
func (pc *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var dto request.PostCreate
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dto.Content == "" || dto.Title == "" {
		http.Error(w, "You need set Content and AuthorId to create a post", http.StatusBadRequest)
		return
	}

	postId, err := pkg.NewULID()
	if err != nil {
		http.Error(w, "Cannot Generate ULID to this Post", http.StatusInternalServerError)
		return
	}

	Post := entities.Post{
		ID:       postId,
		Title:    dto.Title,
		Content:  dto.Content,
		Likes:    dto.Likes,
		Dislikes: dto.Dislikes,
		AuthorId: dto.AuthorId,
	}
	if err := pc.service.CreatePost(&Post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// CreatePostWithAi godoc
// @Summary Create a post with AI-generated title and hashtags
// @Description Creates a new blog post with AI-generated title and hashtags from content (Author role required)
// @Tags Posts, AI
// @Accept json
// @Produce json
// @Param request body request.AiPostCreate true "AI post creation data"
// @Success 201 {string} string "Post created"
// @Failure 400 {string} string "You need set Content and AuthorId to create a post"
// @Failure 500 {string} string "Cannot create post"
// @Security CookieAuth
// @Router /post-with-ai [post]
func (pc *PostController) CreatePostWithAi(w http.ResponseWriter, r *http.Request) {
	var dto request.AiPostCreate
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dto.Content == "" {
		http.Error(w, "You need set Content and AuthorId to create a post", http.StatusBadRequest)
		return
	}

	aiRes := tools.GeneratePropsOfContent(dto.Content)

	aiPostId, err := pkg.NewULID()
	if err != nil {
		http.Error(w, "Cannot Generate ULID to this Post", http.StatusInternalServerError)
		return
	}

	Post := entities.Post{
		ID:       aiPostId,
		Title:    aiRes.Title,
		Content:  pkg.GeneratePostContent(dto.Content, aiRes.Hashtags),
		AuthorId: dto.AuthorId,
	}

	if err := pc.service.CreatePost(&Post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetPostById godoc
// @Summary Get post by ID
// @Description Retrieves a single post by its ULID
// @Tags Posts
// @Produce json
// @Param id path string true "Post ULID"
// @Success 200 {object} response.PostResponse
// @Failure 400 {string} string "ID is required"
// @Failure 400 {string} string "Cannot parse Id of Post"
// @Failure 500 {string} string "Error retrieving post"
// @Router /post/{id} [get]
func (uc *PostController) GetPostById(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.URL.Query().Get("postId")
	if postIdStr == "" {
		exceptions.BadRequest(w, errors.New("Request Error"), "You need to provide Post Id", postIdStr)
		return
	}

	postId, err := pkg.ParseULID(postIdStr)
	if err != nil {
		exceptions.BadRequest(w, err, "Cannot parse Id of Post", postId)
		return
	}

	post, err := uc.service.GetPostByID(postId)
	if err == gorm.ErrRecordNotFound {
		exceptions.NotFound(w, err, fmt.Sprintf("Post with id %v not found", postId))
		return
	}

	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot get this post", reqId)
		return
	}

	postObj := response.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Slug:      post.Slug,
		Content:   post.Content,
		AuthorId:  post.AuthorId,
		Likes:     post.Likes,
		Dislikes:  post.Dislikes,
		Views:     post.Views,
		CreatedAt: post.CreatedAt,
	}

	response.ShowPost(w, postObj)
}

// GetPostBySlug godoc
// @Summary Get post by slug
// @Description Retrieves a single post by its URL slug
// @Tags Posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} response.PostResponse
// @Failure 400 {string} string "Slug is required"
// @Failure 500 {string} string "Error retrieving post"
// @Router /post/slug/{slug} [get]
func (uc *PostController) GetPostBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "Slug is required", http.StatusBadRequest)
		return
	}

	post, err := uc.service.GetPostBySlug(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := response.PostResponse{
		ID:        post.ID,
		AuthorId:  post.AuthorId,
		Title:     post.Title,
		Content:   post.Content,
		Slug:      post.Slug,
		Likes:     post.Likes,
		Dislikes:  post.Dislikes,
		CreatedAt: post.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetAllPosts godoc
// @Summary Get all posts
// @Description Retrieves all blog posts

// GetPaginatedPosts godoc
// @Summary Get paginated posts
// @Description Retrieves a paginated list of blog posts
// @Tags Posts
// @Produce json
// @Param page query int true "Page number (default: 1)"
// @Param limit query int true "Number of posts per page (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Returns data array and meta object with pagination info"
// @Failure 400 {string} string "Page is required"
// @Failure 400 {string} string "Limit is required"
// @Failure 500 {string} string "Error retrieving posts"
// @Router /posts [get]
func (c *PostController) GetPaginatedPosts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if pageStr == "" || limitStr == "" {
		exceptions.BadRequest(w, errors.New("Request Error"), "You need provide values to page and limit", map[string]string{
			"page":  pageStr,
			"limit": limitStr,
		})
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot convert page into int number", reqId)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot convert limit into int number", reqId)
		return
	}

	if page <= 0 || limit <= 0 || limit > 25 {
		page = 1
		limit = 10
	}

	posts, total, err := c.service.GetPaginatedPosts(page, limit)
	if err != nil {
		reqId := middleware.GetReqID(r.Context())
		exceptions.InternalError(w, err, "Cannot get paginated Posts", reqId)
		return
	}

	postsObj := make([]response.PostResponse, len(posts))
	for i := range len(posts) {
		p := posts[i]
		postsObj[i] = response.PostResponse{
			ID:        p.ID,
			Title:     p.Title,
			Slug:      p.Slug,
			Content:   p.Content,
			Views:     p.Views,
			AuthorId:  p.AuthorId,
			CreatedAt: p.CreatedAt,
		}
	}

	response.ListPosts(w, postsObj, page, limit, int(total))
}

// UpdatePost godoc
// @Summary Update a post
// @Description Updates an existing blog post (Author role required)
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post ULID"
// @Param request body request.PostUpdate true "Post update data"
// @Success 200 {object} response.PostResponse
// @Failure 400 {string} string "ID is required"
// @Failure 500 {string} string "Error updating post"
// @Security CookieAuth
// @Router /post/{id} [put]
func (uc *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postIdStr := chi.URLParam(r, "id")
	if postIdStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	postId, err := pkg.ParseULID(postIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dto request.PostUpdate
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := uc.service.GetPostByID(postId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post.Content = dto.Content
	post.Likes = dto.Likes
	post.Title = dto.Title
	post.Dislikes = dto.Dislikes

	if err := uc.service.UpdatePost(post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := response.PostResponse{
		ID:       post.ID,
		AuthorId: post.AuthorId,
		Title:    post.Title,
		Content:  post.Content,
		Likes:    post.Likes,
		Dislikes: post.Dislikes,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeletePost godoc
// @Summary Delete a post
// @Description Deletes an existing blog post (Author role required)
// @Tags Posts
// @Param id path string true "Post ULID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "ID is required"
// @Failure 500 {string} string "Error deleting post"
// @Security CookieAuth
// @Router /post/{id} [delete]
func (uc *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	postIdStr := chi.URLParam(r, "id")
	if postIdStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	postId, err := pkg.ParseULID(postIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := uc.service.DeletePost(postId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
