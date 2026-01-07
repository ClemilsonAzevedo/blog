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
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PostController struct {
	service *service.PostService
}

func NewPostController(service *service.PostService) *PostController {
	return &PostController{
		service: service,
	}
}

func (uc *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var dto request.PostCreate
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	Post := entities.Post{
		Title:    dto.Title,
		Content:  dto.Content,
		Likes:    dto.Likes,
		Dislikes: dto.Dislikes,
		UserID:   dto.UserID,
	}
	if err := uc.service.CreatePost(&Post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (uc *PostController) GetPostById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	post, err := uc.service.GetPostByID(uuid.MustParse(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := response.PostResponse{
		UserID:    post.UserID,
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
		UserID:    post.UserID,
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

	postsDto := make([]response.PostResponse, len(posts));
	for i, post := range posts {
		postsDto[i] = response.PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
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

func (uc *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	if err := uuid.Validate(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dto request.PostUpdate
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := uc.service.GetPostByID(uuid.MustParse(id))
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
		UserID:   post.UserID,
		Title:    post.Title,
		Content:  post.Content,
		Likes:    post.Likes,
		Dislikes: post.Dislikes,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (uc *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	if err := uuid.Validate(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := uc.service.DeletePost(uuid.MustParse(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
