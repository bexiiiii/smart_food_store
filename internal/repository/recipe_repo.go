package repository

import (
	"github.com/bexiiiii/smart_food_store.git/internal/models"
	"sync"
)

type RecipeRepository struct {
	Mu      sync.RWMutex
	Recipes []models.Recipe
}

var RecipeRepo = &RecipeRepository{
	Recipes: []models.Recipe{
		{
			ID:   1,
			Name: "Pasta Carbonara",
			Ingredients: []models.Product{
				{ID: 101, Name: "Pasta", Price: 500},
				{ID: 104, Name: "Bacon", Price: 1200},
			},
			Portions: 2,
		},
	},
}
