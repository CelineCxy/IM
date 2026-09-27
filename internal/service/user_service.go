package service

import (
	"im/internal/model"
	"im/internal/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService{
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUser(id int64) (*model.User,bool){
	return s.userRepository.GetByID(id)
}