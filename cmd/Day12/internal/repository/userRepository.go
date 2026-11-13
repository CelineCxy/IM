package repository

import (
	"im/cmd/Day12/internal/model"
	"sync"
)

type UserRepository struct {
	mu     sync.Mutex
	users  map[int64]model.User
	nextID int64
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:  make(map[int64]model.User),
		nextID: 1,
	}
}

func (r *UserRepository) CreateUser(name string, age int) model.User {
	r.mu.Lock()
	user := model.User{
		ID:   r.nextID,
		Name: name,
		Age:  age,
	}
	r.users[user.ID] = user
	r.nextID++
	r.mu.Unlock()
	return user
}
func (r *UserRepository) GetByID(id int64) (model.User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, exists := r.users[id]
	return user, exists

}
func (r *UserRepository) DeleteUser(id int64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, exists := r.users[id]
	if exists {
		delete(r.users, id)
		return true
	}
	return false
}
