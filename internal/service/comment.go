package service

import (
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/repository"
	"github.com/google/uuid"
)

type CommentService struct {
	CommentRepository *repository.CommentRepository
}

func NewCommentService(commentRepository *repository.CommentRepository) *CommentService {
	return &CommentService{
		CommentRepository: commentRepository,
	}
}

func (s *CommentService) CreateComment(comment *entities.Comment) error {
	return s.CommentRepository.CreateComment(comment)
}

func (s *CommentService) UpdateComment(comment *entities.Comment) error {
	return s.CommentRepository.UpdateComment(comment)
}

func (s *CommentService) DeleteComment(id uuid.UUID) error {
	return s.CommentRepository.DeleteComment(id)
}

func (s *CommentService) GetCommentByID(id uuid.UUID) (*entities.Comment, error) {
	return s.CommentRepository.GetCommentByID(id)
}

func (s *CommentService) GetAllComments() ([]*entities.Comment, error) {
	return s.CommentRepository.GetAllComments()
}

