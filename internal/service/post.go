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

func (ps *PostService) CreatePost(post *Post) error {
	if err := ps.PostRepository.CreatePost(post); err != nil {
		return err
	}

	ps.cache.InvalidateLists()
	return nil
}

func (ps *PostService) UpdatePost(post *Post) error {
	if err := ps.PostRepository.UpdatePost(post); err != nil {
		return err
	}

	ps.cache.InvalidateLists()
	return nil
}

func (ps *PostService) DeletePost(id pkg.ULID) error {
	post, _ := ps.PostRepository.GetPostByID(id)
	slug := ""
	if post != nil {
		slug = post.Slug
	}

	if err := ps.PostRepository.DeletePost(id); err != nil {
		return err
	}

	ps.cache.InvalidatePost(id, slug)
	ps.cache.InvalidateLists()
	return nil
}

func (ps *PostService) GetPostByID(id pkg.ULID) (*Post, error) {
	post, err := ps.PostRepository.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ps *PostService) GetAllPosts() ([]*Post, error) {
	if posts, found := ps.cache.GetAll(); found {
		return posts, nil
	}

	posts, err := ps.PostRepository.GetAllPosts()
	if err != nil {
		return nil, err
	}

	ps.cache.SetAll(posts)
	return posts, nil
}

func (ps *PostService) GetPaginatedPosts(page, limit int) ([]entities.Post, int64, error) {
	if result, found := ps.cache.GetPaginated(page, limit); found {
		return result.Posts, result.Total, nil
	}

	offset := (page - 1) * limit
	posts, total, err := ps.PostRepository.FindAllPaginated(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	ps.cache.SetPaginated(page, limit, posts, total)
	return posts, total, nil
}

func (ps *PostService) GenerateUniqueSlug(title string) (string, error) {
	base := slug.Make(title)
	slug := base
	i := 1

	for {
		exists, err := ps.PostRepository.SlugExists(slug)
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
