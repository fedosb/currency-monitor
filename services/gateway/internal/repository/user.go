package repository

import (
	"errors"
	"sync"

	"github.com/fedosb/currency-monitor/services/gateway/internal/entity"
)

type UserRepository struct {
	Users map[int]entity.User

	mu sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Users: make(map[int]entity.User),
	}
}

func (r *UserRepository) Create(user entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Users[user.ID] = user
	return nil
}

func (r *UserRepository) GetByID(id int) (entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.Users[id]
	if !ok {
		return entity.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) GetByLogin(login string) (entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.Users {
		if user.Login == login {
			return user, nil
		}
	}

	return entity.User{}, errors.New("user not found")
}
