package repository

import (
	"sync"
import "github.com/bexiiiii/smart_food_store.git/internal/models"

	"github.com/bexiiiii/smart_food_store/internal/models"
)

type UserRepository struct {
	mu     sync.Mutex
	users  []models.User
	nextID int
}
