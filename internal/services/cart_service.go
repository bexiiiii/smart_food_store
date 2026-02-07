package services

import (
	"errors"

	"github.com/bexiiiii/smart_food_store/internal/models"
	"github.com/bexiiiii/smart_food_store/internal/repository"
)

type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewCartService(cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *CartService) GetCart(userID uint) (*models.CartResponse, error) {
	cart, err := s.cartRepo.GetOrCreateByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get cart")
	}

	return s.buildCartResponse(cart)
}

func (s *CartService) AddItem(userID uint, req *models.CartItemRequest) (*models.CartResponse, error) {
	// Validate product exists
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check stock
	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	cart, err := s.cartRepo.GetOrCreateByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get cart")
	}

	item := &models.CartItem{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := s.cartRepo.AddItem(cart.ID, item); err != nil {
		return nil, errors.New("failed to add item to cart")
	}

	// Refresh cart data
	cart, err = s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get updated cart")
	}

	return s.buildCartResponse(cart)
}

func (s *CartService) UpdateItemQuantity(userID uint, productID uint, quantity float64) (*models.CartResponse, error) {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	// Validate product exists and has enough stock
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.Stock < quantity {
		return nil, errors.New("insufficient stock")
	}

	if quantity <= 0 {
		// Remove item if quantity is 0 or negative
		if err := s.cartRepo.RemoveItem(cart.ID, productID); err != nil {
			return nil, errors.New("failed to remove item")
		}
	} else {
		if err := s.cartRepo.UpdateItemQuantity(cart.ID, productID, quantity); err != nil {
			return nil, errors.New("failed to update item quantity")
		}
	}

	// Refresh cart data
	cart, err = s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get updated cart")
	}

	return s.buildCartResponse(cart)
}

func (s *CartService) RemoveItem(userID uint, productID uint) (*models.CartResponse, error) {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	if err := s.cartRepo.RemoveItem(cart.ID, productID); err != nil {
		return nil, errors.New("failed to remove item")
	}

	// Refresh cart data
	cart, err = s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get updated cart")
	}

	return s.buildCartResponse(cart)
}

func (s *CartService) ClearCart(userID uint) error {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return errors.New("cart not found")
	}

	return s.cartRepo.ClearCart(cart.ID)
}

func (s *CartService) AddMultipleItems(userID uint, items []models.CartItemRequest) (*models.CartResponse, error) {
	cart, err := s.cartRepo.GetOrCreateByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get cart")
	}

	for _, req := range items {
		// Validate product exists
		product, err := s.productRepo.GetByID(req.ProductID)
		if err != nil {
			continue // Skip invalid products
		}

		// Check stock (add what's available)
		quantity := req.Quantity
		if product.Stock < quantity {
			quantity = product.Stock
		}

		if quantity > 0 {
			item := &models.CartItem{
				ProductID: req.ProductID,
				Quantity:  quantity,
			}
			s.cartRepo.AddItem(cart.ID, item)
		}
	}

	// Refresh cart data
	cart, err = s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get updated cart")
	}

	return s.buildCartResponse(cart)
}

func (s *CartService) GetCartProductIDs(userID uint) ([]uint, error) {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var productIDs []uint
	for _, item := range cart.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	return productIDs, nil
}

func (s *CartService) GetCartItemNames(userID uint) ([]string, error) {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, item := range cart.Items {
		if item.Product.Name != "" {
			names = append(names, item.Product.Name)
		}
	}

	return names, nil
}

func (s *CartService) buildCartResponse(cart *models.Cart) (*models.CartResponse, error) {
	var totalPrice float64
	var items []models.CartItemResponse

	for _, item := range cart.Items {
		if item.Product == nil {
			continue
		}

		subtotal := item.Product.Price * item.Quantity
		totalPrice += subtotal

		items = append(items, models.CartItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Price:       item.Product.Price,
			Quantity:    item.Quantity,
			Unit:        item.Product.Unit,
			Subtotal:    subtotal,
		})
	}

	return &models.CartResponse{
		ID:         cart.ID,
		Items:      items,
		TotalPrice: totalPrice,
		ItemCount:  len(items),
	}, nil
}
