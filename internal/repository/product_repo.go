package repository

import (
	"sync"

	"github.com/bexiiiii/smart_food_store/internal/models"
)

type ProductRepository struct {
	mu       sync.Mutex
	products []models.Product
	nextID   int
}
