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

func (s *CommentService) CreateComment(comment *Comment) error {
	return s.CommentRepository.CreateComment(comment)
}

func (s *CommentService) UpdateComment(comment *Comment) error {
	return s.CommentRepository.UpdateComment(comment)
}

func (s *CommentService) DeleteComment(id pkg.ULID) error {
	return s.CommentRepository.DeleteComment(id)
}

func (s *CommentService) GetCommentByID(id pkg.ULID) (*Comment, error) {
	return s.CommentRepository.GetCommentByID(id)
}

func (s *CommentService) GetAllComments() ([]*Comment, error) {
	return s.CommentRepository.GetAllComments()
}

func (s *CommentService) GetCommentsByPostID(postID pkg.ULID) ([]*Comment, error) {
	return s.CommentRepository.GetCommentsByPostID(postID)
}

func (s *CommentService) GetCommentsByUserID(userID pkg.ULID) ([]*Comment, error) {
	return s.CommentRepository.GetCommentsByUserID(userID)
}
