package services

import "github.com/bexiiiii/smart_food_store.git/internal/models"

type CartService struct {
	carts  []models.Cart
	nextID int
}

func (s *CartService) Create(cart models.Cart) models.Cart {
	s.nextID++
	cart.ID = s.nextID
	s.carts = append(s.carts, cart)
	return cart
}

func (s *CartService) GetAll() []models.Cart {
	return s.carts
func (s *CartService) AddItems(cart *models.Cart, items map[int]float64) {

}
