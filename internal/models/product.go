package models

import (
	"time"

	"gorm.io/gorm"
)

type Unit string

const (
	UnitGram      Unit = "g"
	UnitKilogram  Unit = "kg"
	UnitLiter     Unit = "l"
	UnitMilliliter Unit = "ml"
	UnitPiece     Unit = "pcs"
)

type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Products  []Product      `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:150;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	Stock       float64        `gorm:"default:0" json:"stock"`
	Unit        Unit           `gorm:"size:10;default:g" json:"unit"`
	CategoryID  uint           `gorm:"index" json:"category_id"`
	Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	ImageURL    string         `gorm:"size:255" json:"image_url"`
}

type ProductCreateRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=150"`
	Description string  `json:"description" binding:"max=500"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       float64 `json:"stock" binding:"gte=0"`
	Unit        Unit    `json:"unit" binding:"required"`
	CategoryID  uint    `json:"category_id" binding:"required"`
	ImageURL    string  `json:"image_url"`
}

type ProductUpdateRequest struct {
	Name        *string  `json:"name" binding:"omitempty,min=2,max=150"`
	Description *string  `json:"description" binding:"omitempty,max=500"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       *float64 `json:"stock" binding:"omitempty,gte=0"`
	Unit        *Unit    `json:"unit"`
	CategoryID  *uint    `json:"category_id"`
	ImageURL    *string  `json:"image_url"`
}

type CategoryCreateRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}
