package repository

import (
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"go.bryk.io/pkg/ulid"
	"gorm.io/gorm"
)

type Comment = entities.Comment
type CommentRepository struct {
	DB *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

func (cr *CommentRepository) CreateComment(comment *Comment) error {
	return cr.DB.Create(comment).Error
}

func (cr *CommentRepository) UpdateComment(comment *Comment) error {
	return cr.DB.Save(comment).Error
}

func (cr *CommentRepository) DeleteComment(id ulid.ULID) error {
	return cr.DB.Delete(&Comment{}, id).Error
}

func (cr *CommentRepository) GetCommentByID(id ulid.ULID) (*Comment, error) {
	var comment Comment
	err := cr.DB.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (cr *CommentRepository) GetAllComments() ([]*Comment, error) {
	var comments []*Comment
	err := cr.DB.Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (cr *CommentRepository) GetCommentsByPostID(postID ulid.ULID) ([]*Comment, error) {
	var Comments []*Comment
	err := cr.DB.Where("post_id = ?", postID).Find(&Comments).Error
	if err != nil {
		return nil, err
	}
	return Comments, nil
}

func (cr *CommentRepository) GetCommentsByUserID(userID ulid.ULID) ([]*Comment, error) {
	var Comments []*Comment
	err := cr.DB.Where("user_id = ?", userID).Find(&Comments).Error
	if err != nil {
		return nil, err
	}
	return Comments, nil
}
