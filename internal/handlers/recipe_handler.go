package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bexiiiii/smart_food_store/internal/models"
	"github.com/bexiiiii/smart_food_store/internal/repository"
	"github.com/bexiiiii/smart_food_store/internal/services"
)

var recipeService = services.NewRecipeService(&repository.RecipeRepository{})

func RecipeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var recipe models.Recipe
		json.NewDecoder(r.Body).Decode(&recipe)
		json.NewEncoder(w).Encode(recipeService.Create(recipe))
	case http.MethodGet:
		json.NewEncoder(w).Encode(recipeService.GetAll())
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
