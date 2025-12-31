package repository

import (
	"github.com/clemilsonazevedo/blog/internal/types"
	"gorm.io/gorm"
)

type Repository[T types.Entity] struct {
	DB *gorm.DB
}

func NewRepository[T types.Entity](db *gorm.DB) *Repository[T] {
	return &Repository[T]{
		DB: db,
	}
}

func (r *Repository[T]) Create(entity T) error {
	return r.DB.Create(&entity).Error
}

func (r *Repository[T]) Update(entity T) error {
	return r.DB.Save(&entity).Error
}

func (r *Repository[T]) GetByID(id any) (*T, error) {
	var entity T
	
	err := r.DB.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) GetAll() ([]*T, error) {
	var entities []*T
	err := r.DB.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}
