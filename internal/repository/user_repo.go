package repository

import (
	"github.com/bexiiiii/smart_food_store.git/internal/models"
	"sync"
)

type UserRepository struct {
	Mu    sync.RWMutex
	Users map[int]models.User
}

var UserRepo = &UserRepository{
	Users: make(map[int]models.User),
}

func (r *UserRepository) CreateUser(u models.User) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.Users[u.ID] = u
}
