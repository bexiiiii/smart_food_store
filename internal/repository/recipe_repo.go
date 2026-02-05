package repository

import (
	"sync"

	"github.com/bexiiiii/smart_food_store/internal/models"
)

type RecipeRepository struct {
	mu      sync.Mutex
	recipes []models.Recipe
	nextID  int
}
