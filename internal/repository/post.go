package repository

import (
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (pr *PostRepository) CreatePost(Post *entities.Post) error {
	return pr.DB.Create(Post).Error
}

func (pr *PostRepository) UpdatePost(Post *entities.Post) error {
	return pr.DB.Save(Post).Error
}

func (pr *PostRepository) DeletePost(id uuid.UUID) error {
	return pr.DB.Delete(&entities.Post{}, id).Error
}

func (pr *PostRepository) GetPostByID(postId uuid.UUID) (*entities.Post, error) {
	var Post entities.Post
	err := pr.DB.Model(&entities.Post{}).Where("id = ?", postId).Clauses(clause.Returning{}).UpdateColumn("views", gorm.Expr("views + ?", 1)).Scan(&Post).Error
	if err != nil {
		return nil, err
	}

	return &Post, nil
}

func (pr *PostRepository) GetPostBySlug(slug string) (*entities.Post, error) {
	var post entities.Post
	err := pr.DB.Model(&entities.Post{}).Where("slug = ?", slug).Clauses(clause.Returning{}).UpdateColumn("views", gorm.Expr("views + ?", 1)).Scan(&post).Error
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (pr *PostRepository) GetAllPosts() ([]*entities.Post, error) {
	var Posts []*entities.Post
	err := pr.DB.Find(&Posts).Error
	if err != nil {
		return nil, err
	}
	return Posts, nil
}

func (pr *PostRepository) FindAllPaginated(limit, offset int) ([]entities.Post, int64, error) {
	var posts []entities.Post
	var total int64

	if err := pr.DB.Model(&entities.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := pr.DB.
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
