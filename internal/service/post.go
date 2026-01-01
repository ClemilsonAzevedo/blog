package service

import (
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/repository"
	"github.com/google/uuid"
)

type PostService struct {
	PostRepository *repository.PostRepository
}

func NewPostService(PostRepository *repository.PostRepository) *PostService {
	return &PostService{
		PostRepository: PostRepository,
	}
}

func (s *PostService) CreatePost(post *entities.Post) error {
	return s.PostRepository.CreatePost(post)
}

func (s *PostService) UpdatePost(post *entities.Post) error {
	return s.PostRepository.UpdatePost(post)
}

func (s *PostService) DeletePost(id uuid.UUID) error {
	return s.PostRepository.DeletePost(id)
}

func (s *PostService) GetPostByID(id uuid.UUID) (*entities.Post, error) {
	return s.PostRepository.GetPostByID(id)
}

func (s *PostService) GetAllPosts() ([]*entities.Post, error) {
	return s.PostRepository.GetAllPosts()
}

