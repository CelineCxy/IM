package service

import (
	"context"
	"im/internal/model"
	"im/internal/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	return s.userRepository.GetByID(ctx, id)
}
