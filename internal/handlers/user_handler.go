package handlers

import (
	"net/http"
	"strconv"

	"github.com/bexiiiii/smart_food_store/internal/middleware"
	"github.com/bexiiiii/smart_food_store/internal/models"
	"github.com/bexiiiii/smart_food_store/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.UserRegisterRequest true "Registration data"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req models.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.userService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.UserLoginRequest true "Login credentials"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.userService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile godoc
// @Summary Get current user profile
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]string
// @Router /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAllUsers godoc (Admin only)
// @Summary Get all users
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.UserResponse
// @Failure 403 {object} map[string]string
// @Router /admin/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc (Admin only)
// @Summary Get user by ID
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 404 {object} map[string]string
// @Router /admin/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserRole godoc (Admin only)
// @Summary Update user role
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param role body map[string]string true "New role"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /admin/users/{id}/role [patch]
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := models.Role(req.Role)
	if role != models.RoleUser && role != models.RoleAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be 'user' or 'admin'"})
		return
	}

	if err := h.userService.UpdateRole(uint(id), role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

// DeleteUser godoc (Admin only)
// @Summary Delete user
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
