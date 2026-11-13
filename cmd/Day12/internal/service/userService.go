package service

import (
	"errors"
	"im/cmd/Day12/internal/model"
	"im/cmd/Day12/internal/repository"
	"strings"
)

var (
	ErrInvalidUser  = errors.New("姓名不能为空，年龄必须在 1 到 150 之间")
	ErrInvalidID    = errors.New("用户 ID 必须是正整数")
	ErrUserNotFound = errors.New("用户不存在")
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) CreateUser(name string, age int) (model.User, error) {
	name = strings.TrimSpace(name)

	if name == "" || age <= 0 || age > 150 {
		return model.User{}, ErrInvalidUser
	}

	user := s.userRepository.CreateUser(name, age)
	return user, nil
}

func (s *UserService) GetUser(id int64) (model.User, error) {
	if id <= 0 {
		return model.User{}, ErrInvalidID
	}
	user, exists := s.userRepository.GetByID(id)
	if !exists {
		return model.User{}, ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) DeleteUser(id int64) error {
	if id <= 0 {
		return ErrInvalidID
	}
	deleted := s.userRepository.DeleteUser(id)
	if deleted {
		return nil
	}
	return ErrUserNotFound
}
