package service

import (
	"fmt"

	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/repository"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type Post = entities.Post
type PostRepository = repository.PostRepository
type PostService struct {
	PostRepository *PostRepository
}

func NewPostService(PostRepository *PostRepository) *PostService {
	return &PostService{
		PostRepository: PostRepository,
	}
}

func (s *PostService) CreatePost(post *Post) error {
	slug, err := s.GenerateUniqueSlug(post.Title)

	if err != nil {
		return err
	}
	post.Slug = slug
	return s.PostRepository.CreatePost(post)
}

func (s *PostService) UpdatePost(post *Post) error {
	return s.PostRepository.UpdatePost(post)
}

func (s *PostService) DeletePost(id uuid.UUID) error {
	return s.PostRepository.DeletePost(id)
}

func (s *PostService) GetPostByID(id uuid.UUID) (*Post, error) {
	return s.PostRepository.GetPostByID(id)
}

func (s *PostService) GetPostBySlug(slug string) (*Post, error) {
	return s.PostRepository.GetPostBySlug(slug)
}

func (s *PostService) GetAllPosts() ([]*Post, error) {
	return s.PostRepository.GetAllPosts()
}

func (s *PostService) GetPaginatedPosts(page, limit int) ([]entities.Post, int64, error) {
	offset := (page - 1) * limit

	return s.PostRepository.FindAllPaginated(limit, offset)
}

func (s *PostService) GenerateUniqueSlug(title string) (string, error) {
	base := slug.Make(title)
	slug := base
	i := 1

	for {
		exists, err := s.PostRepository.SlugExists(slug)
		if err != nil {
			return "", err
		}

		if !exists {
			break
		}

		slug = fmt.Sprintf("%s-%d", base, i)
		i++
	}

	return slug, nil
}
