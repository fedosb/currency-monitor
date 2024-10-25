package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/fedosb/currency-monitor/services/gateway/internal/entity"
)

type UserRepository struct {
	Users map[string]entity.User

	mu sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Users: make(map[string]entity.User),
	}
}

func (r *UserRepository) Create(_ context.Context, user entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Users[user.Login]; ok {
		return errors.New("user already exists")
	}

	r.Users[user.Login] = user
	return nil
}

func (r *UserRepository) GetByLogin(_ context.Context, login string) (entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if user, ok := r.Users[login]; ok {
		return user, nil
	}

	return entity.User{}, errors.New("user not found")
}
