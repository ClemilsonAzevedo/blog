package controller

import (
	"encoding/json"
	"net/http"
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/dto/request"
	"github.com/clemilsonazevedo/blog/internal/dto/response"
	"github.com/clemilsonazevedo/blog/internal/service"
	"github.com/go-chi/chi"
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
		Title: dto.Title,
		Content: dto.Content,
		Likes: dto.Likes,
		Dislikes: dto.Dislikes,
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
	
	Post, err := uc.service.GetPostByID(uuid.MustParse(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := response.PostResponse{
		ID:       Post.ID,
		UserID:   Post.UserID,
		Title:    Post.Title,
		Content:  Post.Content,
		Likes:    Post.Likes,
		Dislikes: Post.Dislikes,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (uc *PostController) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	Posts, err := uc.service.GetAllPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Posts)
}

func (uc *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
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

	if err := uc.service.DeletePost(uuid.MustParse(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
