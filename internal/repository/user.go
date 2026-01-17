package repository

import (
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/pkg"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) CreateUser(user *entities.User) error {
	return ur.DB.Create(user).Error
}

func (ur *UserRepository) UpdateUser(user *entities.User) error {
	err := ur.DB.
		Model(&user).
		Select("UserName", "Email").
		Updates(entities.User{UserName: user.UserName, Email: user.Email}).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) DeleteUser(id pkg.ULID) error {
	return ur.DB.Delete(&entities.User{}, id).Error
}

func (ur *UserRepository) GetUserByID(id pkg.ULID) (*entities.User, error) {
	var user entities.User
	err := ur.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := ur.DB.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return &entities.User{}, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
