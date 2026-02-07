package handlers

import (
	"net/http"

	"github.com/bexiiiii/smart_food_store/internal/middleware"
	"github.com/bexiiiii/smart_food_store/internal/models"
	"github.com/bexiiiii/smart_food_store/internal/services"
	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	aiService   *services.AIService
	cartService *services.CartService
}

func NewAIHandler(aiService *services.AIService, cartService *services.CartService) *AIHandler {
	return &AIHandler{
		aiService:   aiService,
		cartService: cartService,
	}
}

// GetIngredientsForDish godoc
// @Summary Get ingredients for a dish from AI
// @Description Enter a dish name and AI will suggest products from the store with quantities
// @Tags ai
// @Accept json
// @Produce json
// @Param request body models.DishToIngredientsRequest true "Dish name and servings"
// @Success 200 {object} models.DishIngredientsResponse
// @Failure 400 {object} map[string]string
// @Router /ai/dish-to-ingredients [post]
func (h *AIHandler) GetIngredientsForDish(c *gin.Context) {
	var req models.DishToIngredientsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default servings to 2 if not provided
	servings := req.Servings
	if servings < 1 {
		servings = 2
	}

	response, err := h.aiService.GetIngredientsForDish(c.Request.Context(), req.DishName, servings)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetRecipesFromCart godoc
// @Summary Get recipe suggestions from cart items
// @Description AI will suggest what can be cooked from products in your cart
// @Tags ai
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.AIRecipeSuggestion
// @Failure 400 {object} map[string]string
// @Router /ai/cart-to-recipes [get]
func (h *AIHandler) GetRecipesFromCart(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get cart items with names
	cartItems, err := h.cartService.GetCartItemNames(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get cart items"})
		return
	}

	// Get product IDs from cart
	productIDs, err := h.cartService.GetCartProductIDs(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get cart items"})
		return
	}

	if len(productIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	suggestions, err := h.aiService.GetRecipesFromCart(c.Request.Context(), productIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return structured response
	response := models.CartRecipesResponse{
		CartItems: cartItems,
		Recipes:   suggestions,
	}

	c.JSON(http.StatusOK, response)
}

// GetRecipesFromProducts godoc
// @Summary Get recipe suggestions from specific products
// @Description AI will suggest what can be cooked from specified products
// @Tags ai
// @Accept json
// @Produce json
// @Param request body models.IngredientsToRecipesRequest true "Product IDs"
// @Success 200 {array} models.AIRecipeSuggestion
// @Failure 400 {object} map[string]string
// @Router /ai/products-to-recipes [post]
func (h *AIHandler) GetRecipesFromProducts(c *gin.Context) {
	var req models.IngredientsToRecipesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	suggestions, err := h.aiService.GetRecipesFromCart(c.Request.Context(), req.ProductIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, suggestions)
}

// AddAISuggestionToCart godoc
// @Summary Add AI suggestion ingredients to cart
// @Description Add all ingredients from an AI suggestion directly to your cart
// @Tags ai
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.AIRecipeSuggestion true "AI suggestion"
// @Success 200 {object} models.CartResponse
// @Failure 400 {object} map[string]string
// @Router /ai/add-to-cart [post]
func (h *AIHandler) AddAISuggestionToCart(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var suggestion models.AIRecipeSuggestion
	if err := c.ShouldBindJSON(&suggestion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert AI ingredients to cart items
	var items []models.CartItemRequest
	for _, ing := range suggestion.Ingredients {
		if ing.Available {
			items = append(items, models.CartItemRequest{
				ProductID: ing.ProductID,
				Quantity:  ing.Quantity,
			})
		}
	}

	if len(items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No available ingredients to add"})
		return
	}

	cart, err := h.cartService.AddMultipleItems(userID, items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}
