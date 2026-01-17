package service

import (
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/repository"
	"github.com/clemilsonazevedo/blog/pkg"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (us *UserService) CreateUser(user *entities.User) error {
	return us.userRepository.CreateUser(user)
}

func (us *UserService) UpdateUser(user *entities.User) error {
	return us.userRepository.UpdateUser(user)
}

func (us *UserService) DeleteUser(id pkg.ULID) error {
	return us.userRepository.DeleteUser(id)
}

func (us *UserService) GetUserByID(id pkg.ULID) (*entities.User, error) {
	return us.userRepository.GetUserByID(id)
}

func (us *UserService) GetUserByEmail(email string) (*entities.User, error) {
	return us.userRepository.GetUserByEmail(email)
}

func (us *UserService) GetUserByName(name string) (*entities.User, error) {
	return us.userRepository.GetUserByName(name)
}

func (us *UserService) GetAllUsers() ([]*entities.User, error) {
	return us.userRepository.GetAllUsers()
}
