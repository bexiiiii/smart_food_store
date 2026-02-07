package handlers

import (
	"net/http"
	"strconv"

	"github.com/bexiiiii/smart_food_store/internal/models"
	"github.com/bexiiiii/smart_food_store/internal/services"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// GetAllProducts godoc
// @Summary Get all products
// @Tags products
// @Produce json
// @Success 200 {array} models.Product
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.productService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.productService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProductsByCategory godoc
// @Summary Get products by category
// @Tags products
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 200 {array} models.Product
// @Router /products/category/{category_id} [get]
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("category_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	products, err := h.productService.GetByCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// SearchProducts godoc
// @Summary Search products
// @Tags products
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {array} models.Product
// @Router /products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}

	products, err := h.productService.Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetAllCategories godoc
// @Summary Get all categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.Category
// @Router /categories [get]
func (h *ProductHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.productService.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// Admin handlers

// CreateProduct godoc (Admin only)
// @Summary Create a new product
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body models.ProductCreateRequest true "Product data"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]string
// @Router /admin/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc (Admin only)
// @Summary Update a product
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.ProductUpdateRequest true "Product data"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Router /admin/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req models.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.Update(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc (Admin only)
// @Summary Delete a product
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.productService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// CreateCategory godoc (Admin only)
// @Summary Create a new category
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param category body models.CategoryCreateRequest true "Category data"
// @Success 201 {object} models.Category
// @Failure 400 {object} map[string]string
// @Router /admin/categories [post]
func (h *ProductHandler) CreateCategory(c *gin.Context) {
	var req models.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.productService.CreateCategory(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// UpdateCategory godoc (Admin only)
// @Summary Update a category
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body models.CategoryCreateRequest true "Category data"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string
// @Router /admin/categories/{id} [put]
func (h *ProductHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req models.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.productService.UpdateCategory(uint(id), req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory godoc (Admin only)
// @Summary Delete a category
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /admin/categories/{id} [delete]
func (h *ProductHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := h.productService.DeleteCategory(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
