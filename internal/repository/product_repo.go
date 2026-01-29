package repository

import "github.com/bexiiiii/smart_food_store/internal/models"

type ProductRepository struct{}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	return &models.Product{}, nil
}
