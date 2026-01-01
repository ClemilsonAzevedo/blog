package repository

import (
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository struct {
	DB *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

func (ur *CommentRepository) CreateComment(Comment *entities.Comment) error {
	return ur.DB.Create(Comment).Error
}

func (ur *CommentRepository) UpdateComment(Comment *entities.Comment) error {
	return ur.DB.Save(Comment).Error
}

func (ur *CommentRepository) DeleteComment(id uuid.UUID) error {
	return ur.DB.Delete(&entities.Comment{}, id).Error
}

func (ur *CommentRepository) GetCommentByID(id uuid.UUID) (*entities.Comment, error) {
	var Comment entities.Comment
	err := ur.DB.First(&Comment, id).Error
	if err != nil {
		return nil, err
	}
	return &Comment, nil
}

func (ur *CommentRepository) GetAllComments() ([]*entities.Comment, error) {
	var Comments []*entities.Comment
	err := ur.DB.Find(&Comments).Error
	if err != nil {
		return nil, err
	}
	return Comments, nil
}