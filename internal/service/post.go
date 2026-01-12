package service

import (
	"fmt"

	"github.com/clemilsonazevedo/blog/internal/cache"
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/repository"
	"github.com/clemilsonazevedo/blog/pkg"
	"github.com/gosimple/slug"
)

type Post = entities.Post
type PostRepository = repository.PostRepository
type PostService struct {
	PostRepository *PostRepository
	cache          *cache.PostCache
}

func NewPostService(PostRepository *PostRepository, postCache *cache.PostCache) *PostService {
	return &PostService{
		PostRepository: PostRepository,
		cache:          postCache,
	}
}

func (s *PostService) CreatePost(post *Post) error {
	slug, err := s.GenerateUniqueSlug(post.Title)

	if err != nil {
		return err
	}
	post.Slug = slug

	if err := s.PostRepository.CreatePost(post); err != nil {
		return err
	}

	s.cache.InvalidateLists()
	return nil
}

func (s *PostService) UpdatePost(post *Post) error {
	oldPost, _ := s.PostRepository.GetPostByID(post.ID)
	oldSlug := ""
	if oldPost != nil {
		oldSlug = oldPost.Slug
	}

	if err := s.PostRepository.UpdatePost(post); err != nil {
		return err
	}

	s.cache.InvalidatePost(post.ID, oldSlug)
	if post.Slug != oldSlug {
		s.cache.InvalidatePost(post.ID, post.Slug)
	}
	s.cache.InvalidateLists()
	return nil
}

func (s *PostService) DeletePost(id pkg.ULID) error {
	post, _ := s.PostRepository.GetPostByID(id)
	slug := ""
	if post != nil {
		slug = post.Slug
	}

	if err := s.PostRepository.DeletePost(id); err != nil {
		return err
	}

	s.cache.InvalidatePost(id, slug)
	s.cache.InvalidateLists()
	return nil
}

func (s *PostService) GetPostByID(id pkg.ULID) (*Post, error) {
	post, err := s.PostRepository.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetPostBySlug(slug string) (*Post, error) {
	post, err := s.PostRepository.GetPostBySlug(slug)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetAllPosts() ([]*Post, error) {
	if posts, found := s.cache.GetAll(); found {
		return posts, nil
	}

	posts, err := s.PostRepository.GetAllPosts()
	if err != nil {
		return nil, err
	}

	s.cache.SetAll(posts)
	return posts, nil
}

func (s *PostService) GetPaginatedPosts(page, limit int) ([]entities.Post, int64, error) {
	if result, found := s.cache.GetPaginated(page, limit); found {
		return result.Posts, result.Total, nil
	}

	offset := (page - 1) * limit
	posts, total, err := s.PostRepository.FindAllPaginated(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	s.cache.SetPaginated(page, limit, posts, total)
	return posts, total, nil
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
