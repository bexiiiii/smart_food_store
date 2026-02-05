package services

import (
	"github.com/bexiiiii/smart_food_store/internal/models"
	"github.com/bexiiiii/smart_food_store/internal/repository"
)

type RecipeService struct {
	repo *repository.RecipeRepository
}

func NewRecipeService(repo *repository.RecipeRepository) *RecipeService {
	return &RecipeService{repo: repo}
}

func (s *RecipeService) Create(recipe models.Recipe) models.Recipe {
	return s.repo.Create(recipe)
}

func (s *RecipeService) GetAll() []models.Recipe {
	return s.repo.GetAll()
}
