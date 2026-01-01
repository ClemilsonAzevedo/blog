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

func (ur *PostRepository) GetAllPosts() ([]*entities.Post, error) {
	var Posts []*entities.Post
	err := ur.DB.Find(&Posts).Error
	if err != nil {
		return nil, err
	}
	return Posts, nil
}