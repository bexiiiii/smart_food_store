package database

import (
	"log"

	"github.com/bexiiiii/smart_food_store/internal/config"
	"github.com/bexiiiii/smart_food_store/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) (*gorm.DB, error) {
	var err error

	DB, err = gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return DB, nil
}

func Migrate() error {
	log.Println("Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Recipe{},
		&models.RecipeIngredient{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed")
	return nil
}

func SeedData() error {
	log.Println("Seeding initial data...")

	// Check if categories already exist
	var count int64
	DB.Model(&models.Category{}).Count(&count)
	if count > 0 {
		log.Println("Data already seeded, skipping...")
		return nil
	}

	// Create categories
	categories := []models.Category{
		{Name: "Vegetables"},
		{Name: "Fruits"},
		{Name: "Meat"},
		{Name: "Dairy"},
		{Name: "Grains"},
		{Name: "Spices"},
		{Name: "Beverages"},
		{Name: "Bakery"},
		{Name: "Seafood"},
		{Name: "Frozen"},
	}

	if err := DB.Create(&categories).Error; err != nil {
		return err
	}

	// Create sample products
	products := []models.Product{
		// Vegetables
		{Name: "Tomato", Description: "Fresh red tomatoes", Price: 2.50, Stock: 100, Unit: models.UnitKilogram, CategoryID: 1},
		{Name: "Onion", Description: "Yellow onions", Price: 1.50, Stock: 150, Unit: models.UnitKilogram, CategoryID: 1},
		{Name: "Garlic", Description: "Fresh garlic bulbs", Price: 0.50, Stock: 200, Unit: models.UnitPiece, CategoryID: 1},
		{Name: "Potato", Description: "White potatoes", Price: 1.80, Stock: 200, Unit: models.UnitKilogram, CategoryID: 1},
		{Name: "Carrot", Description: "Fresh carrots", Price: 1.20, Stock: 120, Unit: models.UnitKilogram, CategoryID: 1},
		{Name: "Bell Pepper", Description: "Colorful bell peppers", Price: 3.00, Stock: 80, Unit: models.UnitKilogram, CategoryID: 1},
		{Name: "Cucumber", Description: "Fresh cucumbers", Price: 1.80, Stock: 90, Unit: models.UnitKilogram, CategoryID: 1},

		// Fruits
		{Name: "Apple", Description: "Red apples", Price: 3.00, Stock: 100, Unit: models.UnitKilogram, CategoryID: 2},
		{Name: "Banana", Description: "Fresh bananas", Price: 2.00, Stock: 150, Unit: models.UnitKilogram, CategoryID: 2},
		{Name: "Lemon", Description: "Fresh lemons", Price: 0.30, Stock: 200, Unit: models.UnitPiece, CategoryID: 2},

		// Meat
		{Name: "Chicken Breast", Description: "Boneless chicken breast", Price: 8.00, Stock: 50, Unit: models.UnitKilogram, CategoryID: 3},
		{Name: "Ground Beef", Description: "Fresh ground beef", Price: 10.00, Stock: 40, Unit: models.UnitKilogram, CategoryID: 3},
		{Name: "Lamb", Description: "Fresh lamb meat", Price: 15.00, Stock: 30, Unit: models.UnitKilogram, CategoryID: 3},

		// Dairy
		{Name: "Milk", Description: "Fresh whole milk", Price: 1.50, Stock: 100, Unit: models.UnitLiter, CategoryID: 4},
		{Name: "Eggs", Description: "Farm fresh eggs", Price: 3.50, Stock: 200, Unit: models.UnitPiece, CategoryID: 4},
		{Name: "Butter", Description: "Unsalted butter", Price: 4.00, Stock: 80, Unit: models.UnitGram, CategoryID: 4},
		{Name: "Cheese", Description: "Cheddar cheese", Price: 6.00, Stock: 60, Unit: models.UnitGram, CategoryID: 4},

		// Grains
		{Name: "Rice", Description: "Long grain white rice", Price: 2.50, Stock: 200, Unit: models.UnitKilogram, CategoryID: 5},
		{Name: "Pasta", Description: "Italian spaghetti", Price: 2.00, Stock: 150, Unit: models.UnitGram, CategoryID: 5},
		{Name: "Flour", Description: "All-purpose flour", Price: 1.50, Stock: 100, Unit: models.UnitKilogram, CategoryID: 5},

		// Spices
		{Name: "Salt", Description: "Table salt", Price: 1.00, Stock: 300, Unit: models.UnitGram, CategoryID: 6},
		{Name: "Black Pepper", Description: "Ground black pepper", Price: 3.00, Stock: 150, Unit: models.UnitGram, CategoryID: 6},
		{Name: "Cumin", Description: "Ground cumin", Price: 4.00, Stock: 100, Unit: models.UnitGram, CategoryID: 6},
		{Name: "Paprika", Description: "Sweet paprika", Price: 3.50, Stock: 100, Unit: models.UnitGram, CategoryID: 6},

		// Beverages
		{Name: "Olive Oil", Description: "Extra virgin olive oil", Price: 8.00, Stock: 50, Unit: models.UnitLiter, CategoryID: 7},
		{Name: "Vegetable Oil", Description: "Cooking vegetable oil", Price: 4.00, Stock: 80, Unit: models.UnitLiter, CategoryID: 7},

		// Bakery
		{Name: "Bread", Description: "Fresh white bread", Price: 2.00, Stock: 100, Unit: models.UnitPiece, CategoryID: 8},

		// Seafood
		{Name: "Salmon", Description: "Fresh Atlantic salmon", Price: 18.00, Stock: 25, Unit: models.UnitKilogram, CategoryID: 9},
		{Name: "Shrimp", Description: "Fresh shrimp", Price: 15.00, Stock: 30, Unit: models.UnitKilogram, CategoryID: 9},
	}

	if err := DB.Create(&products).Error; err != nil {
		return err
	}

	log.Println("Initial data seeded successfully")
	return nil
}
