package controller

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/dto/request"
	"github.com/clemilsonazevedo/blog/internal/dto/response"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/clemilsonazevedo/blog/pkg"
	"github.com/clemilsonazevedo/blog/tools"
	"github.com/go-chi/chi/v5"
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
		UserID:   dto.UserID,
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
		ID:      aiPostId,
		Title:   aiRes.Title,
		Content: pkg.GeneratePostContent(dto.Content, aiRes.Hashtags),
		UserID:  dto.UserID,
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
	postIdStr := chi.URLParam(r, "id")
	if postIdStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	postId, err := pkg.ParseULID(postIdStr)
	if err != nil {
		http.Error(w, "Cannot parse Id of Post", http.StatusBadRequest)
		return
	}

	post, err := uc.service.GetPostByID(postId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := response.PostResponse{
		AuthorId:  post.UserID,
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Slug:      post.Slug,
		Likes:     post.Likes,
		Dislikes:  post.Dislikes,
		Views:     post.Views,
		CreatedAt: post.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
		AuthorId:  post.UserID,
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
// @Tags Posts
// @Produce json
// @Success 200 {array} entities.Post
// @Failure 500 {string} string "Error retrieving posts"
// @Router /posts/all [get]
func (uc *PostController) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := uc.service.GetAllPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

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
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		http.Error(w, "Page is required", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, "Limit is required", http.StatusBadRequest)
		return
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	posts, total, err := c.service.GetPaginatedPosts(page, limit)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	postsDto := make([]response.PostResponse, len(posts))
	for i, post := range posts {
		postsDto[i] = response.PostResponse{
			ID:        post.ID,
			Content:   post.Content,
			Title:     post.Title,
			Slug:      post.Slug,
			Views:     post.Views,
			AuthorId:  post.UserID,
			CreatedAt: post.CreatedAt,
		}
	}

	response := map[string]any{
		"data": postsDto,
		"meta": map[string]any{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": int(math.Ceil(float64(total) / float64(limit))),
		},
	}

	json.NewEncoder(w).Encode(response)
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
		AuthorId: post.UserID,
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
