package repository

import (
	"github.com/bexiiiii/smart_food_store.git/internal/models"
	"sync"
)

type ProductRepository struct {
	Mu       sync.RWMutex
	Products []models.Product
}

var ProductRepo = &ProductRepository{
	Products: []models.Product{
		{ID: 101, Name: "Pasta", Price: 500},
		{ID: 102, Name: "Tomato Sauce", Price: 300},
		{ID: 103, Name: "Cheese", Price: 850},
	},
}
