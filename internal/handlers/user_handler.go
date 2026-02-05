package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bexiiiii/smart_food_store/internal/models"
	"github.com/bexiiiii/smart_food_store/internal/repository"
	"github.com/bexiiiii/smart_food_store/internal/services"
)

var userService = services.NewUserService(&repository.UserRepository{})

func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var user models.User
		json.NewDecoder(r.Body).Decode(&user)
		json.NewEncoder(w).Encode(userService.Create(user))
	case http.MethodGet:
		json.NewEncoder(w).Encode(userService.GetAll())
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
