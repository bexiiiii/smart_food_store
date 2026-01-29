package repository

import "github.com/bexiiiii/smart_food_store/internal/models"

type UserRepository struct{}

func (r *UserRepository) Create(user models.User) error {
	// save user
	return nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	return &models.User{}, nil
}
