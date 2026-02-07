package main

import (
	"log"

	"github.com/bexiiiii/smart_food_store/internal/config"
	"github.com/bexiiiii/smart_food_store/internal/database"
	"github.com/bexiiiii/smart_food_store/internal/handlers"
	"github.com/bexiiiii/smart_food_store/internal/middleware"
	"github.com/bexiiiii/smart_food_store/internal/repository"
	"github.com/bexiiiii/smart_food_store/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Seed initial data
	if err := database.SeedData(); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	cartRepo := repository.NewCartRepository(db)
	recipeRepo := repository.NewRecipeRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, cartRepo, cfg)
	productService := services.NewProductService(productRepo, categoryRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	recipeService := services.NewRecipeService(recipeRepo, productRepo, cartRepo)
	aiService, err := services.NewAIService(productRepo, recipeRepo, cfg)
	if err != nil {
		log.Printf("Warning: AI Service not available: %v", err)
	}

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService)
	recipeHandler := handlers.NewRecipeHandler(recipeService)
	aiHandler := handlers.NewAIHandler(aiService, cartService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// Setup Gin router
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "message": "Smart Food Store API is running"})
		})

		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// Products routes (public)
		products := v1.Group("/products")
		{
			products.GET("", productHandler.GetAllProducts)
			products.GET("/:id", productHandler.GetProductByID)
			products.GET("/category/:category_id", productHandler.GetProductsByCategory)
			products.GET("/search", productHandler.SearchProducts)
		}

		// Categories routes (public)
		categories := v1.Group("/categories")
		{
			categories.GET("", productHandler.GetAllCategories)
		}

		// Recipes routes (public read)
		recipes := v1.Group("/recipes")
		{
			recipes.GET("", recipeHandler.GetAllRecipes)
			recipes.GET("/:id", recipeHandler.GetRecipeByID)
			recipes.GET("/search", recipeHandler.SearchRecipes)
			recipes.GET("/:id/calculate", recipeHandler.CalculateIngredients)
		}

		// AI routes (public for dish-to-ingredients)
		ai := v1.Group("/ai")
		{
			ai.POST("/dish-to-ingredients", aiHandler.GetIngredientsForDish)
			ai.POST("/products-to-recipes", aiHandler.GetRecipesFromProducts)
		}

		// Protected routes (requires authentication)
		protected := v1.Group("")
		protected.Use(authMiddleware.AuthRequired())
		{
			// User profile
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetProfile)
			}

			// Cart routes
			cart := protected.Group("/cart")
			{
				cart.GET("", cartHandler.GetCart)
				cart.POST("/items", cartHandler.AddItem)
				cart.POST("/items/bulk", cartHandler.AddMultipleItems)
				cart.PUT("/items/:product_id", cartHandler.UpdateItemQuantity)
				cart.DELETE("/items/:product_id", cartHandler.RemoveItem)
				cart.DELETE("", cartHandler.ClearCart)
			}

			// Recipe - add to cart
			protected.POST("/recipes/:id/add-to-cart", recipeHandler.AddRecipeToCart)

			// AI - cart based suggestions
			protected.GET("/ai/cart-to-recipes", aiHandler.GetRecipesFromCart)
			protected.POST("/ai/add-to-cart", aiHandler.AddAISuggestionToCart)
		}

		// Admin routes (requires admin role)
		admin := v1.Group("/admin")
		admin.Use(authMiddleware.AuthRequired(), authMiddleware.AdminRequired())
		{
			// User management
			admin.GET("/users", userHandler.GetAllUsers)
			admin.GET("/users/:id", userHandler.GetUserByID)
			admin.PATCH("/users/:id/role", userHandler.UpdateUserRole)
			admin.DELETE("/users/:id", userHandler.DeleteUser)

			// Product management
			admin.POST("/products", productHandler.CreateProduct)
			admin.PUT("/products/:id", productHandler.UpdateProduct)
			admin.DELETE("/products/:id", productHandler.DeleteProduct)

			// Category management
			admin.POST("/categories", productHandler.CreateCategory)
			admin.PUT("/categories/:id", productHandler.UpdateCategory)
			admin.DELETE("/categories/:id", productHandler.DeleteCategory)

			// Recipe management
			admin.POST("/recipes", recipeHandler.CreateRecipe)
			admin.PUT("/recipes/:id", recipeHandler.UpdateRecipe)
			admin.DELETE("/recipes/:id", recipeHandler.DeleteRecipe)
		}
	}

	// Start server
	log.Printf("🚀 Smart Food Store API starting on port %s", cfg.ServerPort)
	log.Printf("📚 API endpoints available at http://localhost:%s/api/v1", cfg.ServerPort)
	
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
