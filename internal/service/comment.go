package service

import (
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/repository"
	"github.com/clemilsonazevedo/blog/pkg"
)

type Comment = entities.Comment
type CommentRepository = repository.CommentRepository
type CommentService struct {
	CommentRepository *CommentRepository
}

func NewCommentService(commentRepository *repository.CommentRepository) *CommentService {
	return &CommentService{
		CommentRepository: commentRepository,
	}
}

func (cs *CommentService) CreateComment(comment *Comment) error {
	return cs.CommentRepository.CreateComment(comment)
}

func (cs *CommentService) UpdateComment(comment *Comment) error {
	return cs.CommentRepository.UpdateComment(comment)
}

func (cs *CommentService) DeleteComment(id pkg.ULID) error {
	return cs.CommentRepository.DeleteComment(id)
}

func (cs *CommentService) GetCommentByID(id pkg.ULID) (*Comment, error) {
	return cs.CommentRepository.GetCommentByID(id)
}

func (cs *CommentService) GetAllComments() ([]*Comment, error) {
	return cs.CommentRepository.GetAllComments()
}

func (cs *CommentService) GetCommentsByPostID(postID pkg.ULID) ([]*Comment, error) {
	return cs.CommentRepository.GetCommentsByPostID(postID)
}

func (cs *CommentService) GetCommentsByUserID(userID pkg.ULID) ([]*Comment, error) {
	return cs.CommentRepository.GetCommentsByUserID(userID)
}
