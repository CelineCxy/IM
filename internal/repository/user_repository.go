package repository

import (
	"context"
	"errors"
	"im/internal/model"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	users map[int64]*model.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: map[int64]*model.User{
			1: {
				ID:   1,
				Name: "小明",
			},
			2: {
				ID:   2,
				Name: "李华",
			},
		},
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	select{
	case <- time.After(time.Second*5):
		user, ok := r.users[id]
		if ok {
			return user,nil
		}
		return nil,ErrUserNotFound
	case <- ctx.Done():
		return nil,ctx.Err()
	}
}
