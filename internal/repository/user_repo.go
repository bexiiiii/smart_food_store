package repository

import (
	"sync"

	"github.com/bexiiiii/smart_food_store/internal/models"
)

type UserRepository struct {
	mu     sync.Mutex
	users  []models.User
	nextID int
}

func (r *UserRepository) Create(user models.User) models.User {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	user.ID = r.nextID
	r.users = append(r.users, user)
	return user
}

func (r *UserRepository) GetAll() []models.User {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.users
}
