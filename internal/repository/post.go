package repository

import (
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (ur *PostRepository) CreatePost(Post *entities.Post) error {
	return ur.DB.Create(Post).Error
}

func (ur *PostRepository) UpdatePost(Post *entities.Post) error {
	return ur.DB.Save(Post).Error
}

func (ur *PostRepository) DeletePost(id uuid.UUID) error {
	return ur.DB.Delete(&entities.Post{}, id).Error
}

func (ur *PostRepository) GetPostByID(id uuid.UUID) (*entities.Post, error) {
	var Post entities.Post
	err := ur.DB.First(&Post, id).Error
	if err != nil {
		return nil, err
	}
	return &Post, nil
}

func (r *PostRepository) GetPostBySlug(slug string) (*entities.Post, error) {
	var post entities.Post

	err := r.DB.Where("slug = ?", slug).First(&post).Error

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (ur *PostRepository) GetAllPosts() ([]*entities.Post, error) {
	var Posts []*entities.Post
	err := ur.DB.Find(&Posts).Error
	if err != nil {
		return nil, err
	}
	return Posts, nil
}

func (r *PostRepository) FindAllPaginated(limit, offset int) ([]entities.Post, int64, error) {
	var posts []entities.Post
	var total int64

	if err := r.DB.Model(&entities.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.DB.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error

	return posts, total, err
}

func (r *PostRepository) SlugExists(slug string) (bool, error) {
	var count int64
	err := r.DB.Model(&entities.Post{}).
		Where("slug = ?", slug).
		Count(&count).Error
	return count > 0, err
}
