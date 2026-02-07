package models

import (
	"time"

	"gorm.io/gorm"
)

type Recipe struct {
	ID           uint               `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	DeletedAt    gorm.DeletedAt     `gorm:"index" json:"-"`
	Name         string             `gorm:"size:200;not null" json:"name"`
	Description  string             `gorm:"type:text" json:"description"`
	Instructions string             `gorm:"type:text" json:"instructions"`
	Servings     int                `gorm:"default:1" json:"servings"`
	PrepTime     int                `json:"prep_time"`  // in minutes
	CookTime     int                `json:"cook_time"`  // in minutes
	ImageURL     string             `gorm:"size:255" json:"image_url"`
	Ingredients  []RecipeIngredient `gorm:"foreignKey:RecipeID" json:"ingredients"`
	IsAIGenerated bool              `gorm:"default:false" json:"is_ai_generated"`
}

type RecipeIngredient struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	RecipeID  uint           `gorm:"index;not null" json:"recipe_id"`
	ProductID uint           `gorm:"index;not null" json:"product_id"`
	Product   *Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  float64        `gorm:"not null" json:"quantity"`
	Unit      Unit           `gorm:"size:10" json:"unit"`
	Notes     string         `gorm:"size:100" json:"notes"` // e.g., "chopped", "melted"
}

// AI Request/Response models
type DishToIngredientsRequest struct {
	DishName string `json:"dish_name" binding:"required"`
	Servings int    `json:"servings"` // Optional, defaults to 2 if not provided
}

type IngredientsToRecipesRequest struct {
	ProductIDs []uint `json:"product_ids" binding:"required,min=1"`
}

type AIIngredient struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Unit        Unit    `json:"unit"`
	Available   bool    `json:"available"`
	Price       float64 `json:"price"`
}

// RequiredIngredient - ingredient needed for a dish
type RequiredIngredient struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

// MatchedProduct - product from store that matches ingredient
type MatchedProduct struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Unit  string  `json:"unit"`
}

// DishIngredientsResponse - response for dish-to-ingredients endpoint
type DishIngredientsResponse struct {
	DishName            string               `json:"dish_name"`
	Description         string               `json:"description"`
	Servings            int                  `json:"servings"`
	RequiredIngredients []RequiredIngredient `json:"required_ingredients"`
	MatchedProducts     []MatchedProduct     `json:"matched_products"`
	CookingTips         string               `json:"cooking_tips"`
	TotalPrice          float64              `json:"total_price"`
}

type AIRecipeSuggestion struct {
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Instructions string         `json:"instructions"`
	PrepTime     int            `json:"prep_time"`
	CookTime     int            `json:"cook_time"`
	Servings     int            `json:"servings"`
	Ingredients  []AIIngredient `json:"ingredients"`
	TotalPrice   float64        `json:"total_price"`
	Confidence   float64        `json:"confidence"` // AI confidence score
}

// CartRecipesResponse - response for cart-to-recipes endpoint
type CartRecipesResponse struct {
	CartItems []string             `json:"cart_items"`
	Recipes   []AIRecipeSuggestion `json:"recipes"`
}

type AddRecipeToCartRequest struct {
	RecipeID uint `json:"recipe_id"`
	Servings int  `json:"servings" binding:"required,min=1"`
}

// Recipe creation request for admins
type RecipeCreateRequest struct {
	Name         string                        `json:"name" binding:"required,min=2,max=200"`
	Description  string                        `json:"description"`
	Instructions string                        `json:"instructions" binding:"required"`
	Servings     int                           `json:"servings" binding:"required,min=1"`
	PrepTime     int                           `json:"prep_time" binding:"gte=0"`
	CookTime     int                           `json:"cook_time" binding:"gte=0"`
	ImageURL     string                        `json:"image_url"`
	Ingredients  []RecipeIngredientCreateRequest `json:"ingredients" binding:"required,min=1"`
}

type RecipeIngredientCreateRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	Unit      Unit    `json:"unit" binding:"required"`
	Notes     string  `json:"notes"`
}
