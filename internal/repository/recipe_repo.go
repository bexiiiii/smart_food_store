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

func (r *RecipeRepository) Create(recipe models.Recipe) models.Recipe {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	recipe.ID = r.nextID
	r.recipes = append(r.recipes, recipe)
	return recipe
}

func (r *RecipeRepository) GetAll() []models.Recipe {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.recipes
}
