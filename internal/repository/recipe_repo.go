package repository

import (
	"github.com/bexiiiii/smart_food_store/internal/models"
	"gorm.io/gorm"
)

type RecipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) *RecipeRepository {
	return &RecipeRepository{db: db}
}

func (r *RecipeRepository) Create(recipe *models.Recipe) error {
	return r.db.Create(recipe).Error
}

func (r *RecipeRepository) GetByID(id uint) (*models.Recipe, error) {
	var recipe models.Recipe
	err := r.db.Preload("Ingredients.Product").First(&recipe, id).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *RecipeRepository) GetAll() ([]models.Recipe, error) {
	var recipes []models.Recipe
	err := r.db.Preload("Ingredients.Product").Find(&recipes).Error
	return recipes, err
}

func (r *RecipeRepository) Search(query string) ([]models.Recipe, error) {
	var recipes []models.Recipe
	err := r.db.Preload("Ingredients.Product").
		Where("LOWER(name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)",
			"%"+query+"%", "%"+query+"%").
		Find(&recipes).Error
	return recipes, err
}

func (r *RecipeRepository) Update(recipe *models.Recipe) error {
	// Delete existing ingredients and recreate them
	if err := r.db.Where("recipe_id = ?", recipe.ID).Delete(&models.RecipeIngredient{}).Error; err != nil {
		return err
	}
	return r.db.Save(recipe).Error
}

func (r *RecipeRepository) Delete(id uint) error {
	// Delete ingredients first
	if err := r.db.Where("recipe_id = ?", id).Delete(&models.RecipeIngredient{}).Error; err != nil {
		return err
	}
	return r.db.Delete(&models.Recipe{}, id).Error
}

func (r *RecipeRepository) GetByProductIDs(productIDs []uint) ([]models.Recipe, error) {
	var recipes []models.Recipe
	
	// Find recipes that contain any of the specified products
	subQuery := r.db.Model(&models.RecipeIngredient{}).
		Select("DISTINCT recipe_id").
		Where("product_id IN ?", productIDs)
	
	err := r.db.Preload("Ingredients.Product").
		Where("id IN (?)", subQuery).
		Find(&recipes).Error
	
	return recipes, err
}

func (r *RecipeRepository) SaveAIGenerated(recipe *models.Recipe) error {
	recipe.IsAIGenerated = true
	return r.db.Create(recipe).Error
}
