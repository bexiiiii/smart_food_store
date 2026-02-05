package repository

import (
	"sync"

	"github.com/bexiiiii/smart_food_store/internal/models"
)

type ProductRepository struct {
	mu       sync.Mutex
	products []models.Product
	nextID   int
}

func (r *ProductRepository) Create(p models.Product) models.Product {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	p.ID = r.nextID
	r.products = append(r.products, p)
	return p
}

func (r *ProductRepository) GetAll() []models.Product {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.products
}
