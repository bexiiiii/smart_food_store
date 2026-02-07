package services

import (
	"errors"

	"github.com/bexiiiii/smart_food_store/internal/models"
	"github.com/bexiiiii/smart_food_store/internal/repository"
)

type RecipeService struct {
	recipeRepo  *repository.RecipeRepository
	productRepo *repository.ProductRepository
	cartRepo    *repository.CartRepository
}

func NewRecipeService(recipeRepo *repository.RecipeRepository, productRepo *repository.ProductRepository, cartRepo *repository.CartRepository) *RecipeService {
	return &RecipeService{
		recipeRepo:  recipeRepo,
		productRepo: productRepo,
		cartRepo:    cartRepo,
	}
}

func (s *RecipeService) Create(req *models.RecipeCreateRequest) (*models.Recipe, error) {
	recipe := &models.Recipe{
		Name:         req.Name,
		Description:  req.Description,
		Instructions: req.Instructions,
		Servings:     req.Servings,
		PrepTime:     req.PrepTime,
		CookTime:     req.CookTime,
		ImageURL:     req.ImageURL,
	}

	for _, ing := range req.Ingredients {
		// Validate product exists
		_, err := s.productRepo.GetByID(ing.ProductID)
		if err != nil {
			return nil, errors.New("invalid product in ingredients")
		}

		recipe.Ingredients = append(recipe.Ingredients, models.RecipeIngredient{
			ProductID: ing.ProductID,
			Quantity:  ing.Quantity,
			Unit:      ing.Unit,
			Notes:     ing.Notes,
		})
	}

	if err := s.recipeRepo.Create(recipe); err != nil {
		return nil, errors.New("failed to create recipe")
	}

	return s.recipeRepo.GetByID(recipe.ID)
}

func (s *RecipeService) GetByID(id uint) (*models.Recipe, error) {
	return s.recipeRepo.GetByID(id)
}

func (s *RecipeService) GetAll() ([]models.Recipe, error) {
	return s.recipeRepo.GetAll()
}

func (s *RecipeService) Search(query string) ([]models.Recipe, error) {
	return s.recipeRepo.Search(query)
}

func (s *RecipeService) Update(id uint, req *models.RecipeCreateRequest) (*models.Recipe, error) {
	recipe, err := s.recipeRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("recipe not found")
	}

	recipe.Name = req.Name
	recipe.Description = req.Description
	recipe.Instructions = req.Instructions
	recipe.Servings = req.Servings
	recipe.PrepTime = req.PrepTime
	recipe.CookTime = req.CookTime
	recipe.ImageURL = req.ImageURL

	recipe.Ingredients = nil
	for _, ing := range req.Ingredients {
		recipe.Ingredients = append(recipe.Ingredients, models.RecipeIngredient{
			RecipeID:  id,
			ProductID: ing.ProductID,
			Quantity:  ing.Quantity,
			Unit:      ing.Unit,
			Notes:     ing.Notes,
		})
	}

	if err := s.recipeRepo.Update(recipe); err != nil {
		return nil, errors.New("failed to update recipe")
	}

	return s.recipeRepo.GetByID(id)
}

func (s *RecipeService) Delete(id uint) error {
	return s.recipeRepo.Delete(id)
}

func (s *RecipeService) CalculateIngredients(recipeID uint, servings int) ([]models.AIIngredient, float64, error) {
	recipe, err := s.recipeRepo.GetByID(recipeID)
	if err != nil {
		return nil, 0, errors.New("recipe not found")
	}

	var ingredients []models.AIIngredient
	var totalPrice float64

	// Calculate ratio based on servings
	ratio := float64(servings) / float64(recipe.Servings)

	for _, ing := range recipe.Ingredients {
		if ing.Product == nil {
			continue
		}

		adjustedQuantity := ing.Quantity * ratio
		price := ing.Product.Price * adjustedQuantity

		ingredients = append(ingredients, models.AIIngredient{
			ProductID:   ing.ProductID,
			ProductName: ing.Product.Name,
			Quantity:    adjustedQuantity,
			Unit:        ing.Unit,
			Available:   ing.Product.Stock >= adjustedQuantity,
			Price:       price,
		})

		totalPrice += price
	}

	return ingredients, totalPrice, nil
}

func (s *RecipeService) AddRecipeToCart(userID uint, recipeID uint, servings int) (*models.CartResponse, error) {
	ingredients, _, err := s.CalculateIngredients(recipeID, servings)
	if err != nil {
		return nil, err
	}

	cart, err := s.cartRepo.GetOrCreateByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get cart")
	}

	for _, ing := range ingredients {
		if !ing.Available {
			continue // Skip unavailable items
		}

		item := &models.CartItem{
			ProductID: ing.ProductID,
			Quantity:  ing.Quantity,
		}

		if err := s.cartRepo.AddItem(cart.ID, item); err != nil {
			continue
		}
	}

	// Return updated cart
	cart, err = s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return buildCartResponse(cart)
}

func buildCartResponse(cart *models.Cart) (*models.CartResponse, error) {
	var totalPrice float64
	var items []models.CartItemResponse

	for _, item := range cart.Items {
		if item.Product == nil {
			continue
		}

		subtotal := item.Product.Price * item.Quantity
		totalPrice += subtotal

		items = append(items, models.CartItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Price:       item.Product.Price,
			Quantity:    item.Quantity,
			Unit:        item.Product.Unit,
			Subtotal:    subtotal,
		})
	}

	return &models.CartResponse{
		ID:         cart.ID,
		Items:      items,
		TotalPrice: totalPrice,
		ItemCount:  len(items),
	}, nil
}
