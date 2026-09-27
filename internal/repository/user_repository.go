package repository

import "im/internal/model"

type UserRepository struct {
	users map[int64]*model.User
}

func NewUserRepository() *UserRepository{
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

func (r *UserRepository) GetByID(id int64) (*model.User,bool){
	user,ok := r.users[id]
	return user,ok
}