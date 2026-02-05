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
