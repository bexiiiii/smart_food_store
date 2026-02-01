package repository

import "github.com/bexiiiii/smart_food_store.git/internal/models"

type UserRepository struct{}

func (r *UserRepository) Create(user models.User) error {

	return nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	return &models.User{}, nil
}
