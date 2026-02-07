package repository

import (
	"github.com/bexiiiii/smart_food_store/internal/models"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetOrCreateByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error
	
	if err == gorm.ErrRecordNotFound {
		cart = models.Cart{UserID: userID}
		if err := r.db.Create(&cart).Error; err != nil {
			return nil, err
		}
		return &cart, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	return &cart, nil
}

func (r *CartRepository) GetByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) AddItem(cartID uint, item *models.CartItem) error {
	// Check if item already exists in cart
	var existingItem models.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, item.ProductID).First(&existingItem).Error
	
	if err == nil {
		// Item exists, update quantity
		existingItem.Quantity += item.Quantity
		return r.db.Save(&existingItem).Error
	}
	
	if err == gorm.ErrRecordNotFound {
		// Item doesn't exist, create new
		item.CartID = cartID
		return r.db.Create(item).Error
	}
	
	return err
}

func (r *CartRepository) UpdateItemQuantity(cartID uint, productID uint, quantity float64) error {
	return r.db.Model(&models.CartItem{}).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		Update("quantity", quantity).Error
}

func (r *CartRepository) RemoveItem(cartID uint, productID uint) error {
	return r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).
		Delete(&models.CartItem{}).Error
}

func (r *CartRepository) ClearCart(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}

func (r *CartRepository) GetCartItems(cartID uint) ([]models.CartItem, error) {
	var items []models.CartItem
	err := r.db.Preload("Product").Where("cart_id = ?", cartID).Find(&items).Error
	return items, err
}

func (r *CartRepository) AddMultipleItems(cartID uint, items []models.CartItem) error {
	for _, item := range items {
		if err := r.AddItem(cartID, &item); err != nil {
			return err
		}
	}
	return nil
}
