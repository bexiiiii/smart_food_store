package handlers

import (
	"net/http"
	"strconv"

	"github.com/bexiiiii/smart_food_store/internal/middleware"
	"github.com/bexiiiii/smart_food_store/internal/models"
	"github.com/bexiiiii/smart_food_store/internal/services"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService *services.CartService
}

func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

// GetCart godoc
// @Summary Get user's cart
// @Tags cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.CartResponse
// @Failure 401 {object} map[string]string
// @Router /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddItem godoc
// @Summary Add item to cart
// @Tags cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param item body models.CartItemRequest true "Cart item"
// @Success 200 {object} models.CartResponse
// @Failure 400 {object} map[string]string
// @Router /cart/items [post]
func (h *CartHandler) AddItem(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartService.AddItem(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// UpdateItemQuantity godoc
// @Summary Update item quantity in cart
// @Tags cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param quantity body map[string]float64 true "Quantity"
// @Success 200 {object} models.CartResponse
// @Failure 400 {object} map[string]string
// @Router /cart/items/{product_id} [put]
func (h *CartHandler) UpdateItemQuantity(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req struct {
		Quantity float64 `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartService.UpdateItemQuantity(userID, uint(productID), req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// RemoveItem godoc
// @Summary Remove item from cart
// @Tags cart
// @Security BearerAuth
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} models.CartResponse
// @Failure 400 {object} map[string]string
// @Router /cart/items/{product_id} [delete]
func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	cart, err := h.cartService.RemoveItem(userID, uint(productID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// ClearCart godoc
// @Summary Clear all items from cart
// @Tags cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /cart [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.cartService.ClearCart(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}

// AddMultipleItems godoc
// @Summary Add multiple items to cart
// @Tags cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param items body []models.CartItemRequest true "Cart items"
// @Success 200 {object} models.CartResponse
// @Failure 400 {object} map[string]string
// @Router /cart/items/bulk [post]
func (h *CartHandler) AddMultipleItems(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var items []models.CartItemRequest
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartService.AddMultipleItems(userID, items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}
