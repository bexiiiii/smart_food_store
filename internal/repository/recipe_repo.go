package repository

import "github.com/bexiiiii/smart_food_store/internal/models"

type RecipeRepository struct{}

func (r *RecipeRepository) GetByID(id int) (*models.Recipe, error) {
	return &models.Recipe{}, nil
}
