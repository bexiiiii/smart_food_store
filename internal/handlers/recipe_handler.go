package handlers

import (
	"net/http"
	"strconv"

	"github.com/bexiiiii/smart_food_store/internal/middleware"
	"github.com/bexiiiii/smart_food_store/internal/models"
	"github.com/bexiiiii/smart_food_store/internal/services"
	"github.com/gin-gonic/gin"
)

type RecipeHandler struct {
	recipeService *services.RecipeService
}

func NewRecipeHandler(recipeService *services.RecipeService) *RecipeHandler {
	return &RecipeHandler{recipeService: recipeService}
}

// GetAllRecipes godoc
// @Summary Get all recipes
// @Tags recipes
// @Produce json
// @Success 200 {array} models.Recipe
// @Router /recipes [get]
func (h *RecipeHandler) GetAllRecipes(c *gin.Context) {
	recipes, err := h.recipeService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recipes"})
		return
	}

	c.JSON(http.StatusOK, recipes)
}

// GetRecipeByID godoc
// @Summary Get recipe by ID
// @Tags recipes
// @Produce json
// @Param id path int true "Recipe ID"
// @Success 200 {object} models.Recipe
// @Failure 404 {object} map[string]string
// @Router /recipes/{id} [get]
func (h *RecipeHandler) GetRecipeByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	recipe, err := h.recipeService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// SearchRecipes godoc
// @Summary Search recipes
// @Tags recipes
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {array} models.Recipe
// @Router /recipes/search [get]
func (h *RecipeHandler) SearchRecipes(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}

	recipes, err := h.recipeService.Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	c.JSON(http.StatusOK, recipes)
}

// CalculateIngredients godoc
// @Summary Calculate ingredients for recipe with custom servings
// @Tags recipes
// @Produce json
// @Param id path int true "Recipe ID"
// @Param servings query int true "Number of servings"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /recipes/{id}/calculate [get]
func (h *RecipeHandler) CalculateIngredients(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	servingsStr := c.Query("servings")
	servings, err := strconv.Atoi(servingsStr)
	if err != nil || servings < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid servings required (min 1)"})
		return
	}

	ingredients, totalPrice, err := h.recipeService.CalculateIngredients(uint(id), servings)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recipe_id":    id,
		"servings":     servings,
		"ingredients":  ingredients,
		"total_price":  totalPrice,
	})
}

// AddRecipeToCart godoc
// @Summary Add all recipe ingredients to cart
// @Tags recipes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Recipe ID"
// @Param request body models.AddRecipeToCartRequest true "Request with servings"
// @Success 200 {object} models.CartResponse
// @Failure 400 {object} map[string]string
// @Router /recipes/{id}/add-to-cart [post]
func (h *RecipeHandler) AddRecipeToCart(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	var req models.AddRecipeToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.recipeService.AddRecipeToCart(userID, uint(id), req.Servings)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// CreateRecipe godoc (Admin only)
// @Summary Create a new recipe
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param recipe body models.RecipeCreateRequest true "Recipe data"
// @Success 201 {object} models.Recipe
// @Failure 400 {object} map[string]string
// @Router /admin/recipes [post]
func (h *RecipeHandler) CreateRecipe(c *gin.Context) {
	var req models.RecipeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe, err := h.recipeService.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, recipe)
}

// UpdateRecipe godoc (Admin only)
// @Summary Update a recipe
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Recipe ID"
// @Param recipe body models.RecipeCreateRequest true "Recipe data"
// @Success 200 {object} models.Recipe
// @Failure 400 {object} map[string]string
// @Router /admin/recipes/{id} [put]
func (h *RecipeHandler) UpdateRecipe(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	var req models.RecipeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe, err := h.recipeService.Update(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// DeleteRecipe godoc (Admin only)
// @Summary Delete a recipe
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "Recipe ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/recipes/{id} [delete]
func (h *RecipeHandler) DeleteRecipe(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	if err := h.recipeService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}
