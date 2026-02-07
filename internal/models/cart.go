package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items     []CartItem     `gorm:"foreignKey:CartID" json:"items"`
}

type CartItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CartID    uint           `gorm:"index;not null" json:"cart_id"`
	ProductID uint           `gorm:"index;not null" json:"product_id"`
	Product   *Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  float64        `gorm:"not null" json:"quantity"`
}

type CartItemRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
}

type CartResponse struct {
	ID         uint               `json:"id"`
	Items      []CartItemResponse `json:"items"`
	TotalPrice float64            `json:"total_price"`
	ItemCount  int                `json:"item_count"`
}

type CartItemResponse struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    float64 `json:"quantity"`
	Unit        Unit    `json:"unit"`
	Subtotal    float64 `json:"subtotal"`
}
