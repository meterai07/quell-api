package service

import (
	"quell-api/entity"
	"quell-api/repository"
)

type Service interface {
	CreateUser(user entity.User) error
	GetUserByEmail(email string) (entity.User, error)
	FindUserByEmail(email string) bool
	UpdateUser(user entity.User) error
	DeleteUser(user entity.User) error
}

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *userService {
	return &userService{
		repository: repository,
	}
}

func (s *userService) CreateUser(user entity.User) error {
	result := s.repository.CreateUser(user)
	if result != nil {
		return result
	}
	return nil
}

func (s *userService) GetUserByEmail(email string) (entity.User, error) {
	var user entity.User
	result, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return user, err
	}
	return result, nil
}

func (s *userService) FindUserByEmail(email string) bool {
	result := s.repository.FindUserByEmail(email)
	return result
}

func (s *userService) UpdateUser(user entity.User) error {
	result := s.repository.UpdateUser(user)
	if result != nil {
		return result
	}
	return nil
}

func (s *userService) DeleteUser(user entity.User) error {
	result := s.repository.DeleteUser(user)
	if result != nil {
		return result
	}
	return nil
}
