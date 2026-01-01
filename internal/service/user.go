package service

import (
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) CreateUser(user *entities.User) error {
	return s.userRepository.CreateUser(user)
}

func (s *UserService) UpdateUser(user *entities.User) error {
	return s.userRepository.UpdateUser(user)
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.userRepository.DeleteUser(id)
}

func (s *UserService) GetUserByID(id uuid.UUID) (*entities.User, error) {
	return s.userRepository.GetUserByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*entities.User, error) {
	return s.userRepository.GetUserByEmail(email)
}

func (s *UserService) GetUserByName(name string) (*entities.User, error) {
	return s.userRepository.GetUserByName(name)
}

func (s *UserService) GetAllUsers() ([]*entities.User, error) {
	return s.userRepository.GetAllUsers()
}

func (s *UserService) GetUsersByRole(role string) ([]*entities.User, error) {
	return s.userRepository.GetUsersByRole(role)
}
