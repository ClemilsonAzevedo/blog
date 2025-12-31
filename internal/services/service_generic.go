package services

import (
	"github.com/clemilsonazevedo/blog/internal/types"
	"github.com/clemilsonazevedo/blog/internal/repository"
)
type Service[T types.Entity] struct {
	repo repository.Repository[T]
}

func NewService[T types.Entity](repo repository.Repository[T]) *Service[T] {
	return &Service[T]{
		repo: repo,
	}
}

func (s *Service[T]) Create(entity T) error {
	return s.repo.Create(entity)
}

func (s *Service[T]) Update(entity T) error {
	return s.repo.Update(entity)
}

func (s *Service[T]) GetByID(id any) (*T, error) {
	return s.repo.GetByID(id)
}

func (s *Service[T]) GetAll() ([]*T, error) {
	return s.repo.GetAll()
}