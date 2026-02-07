package repository

import (
	"github.com/bexiiiii/smart_food_store/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Find(&products).Error
	return products, err
}

func (r *ProductRepository) GetByCategory(categoryID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}

func (r *ProductRepository) Search(query string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").
		Where("LOWER(name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)", 
			"%"+query+"%", "%"+query+"%").
		Find(&products).Error
	return products, err
}

func (r *ProductRepository) GetByIDs(ids []uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Where("id IN ?", ids).Find(&products).Error
	return products, err
}

func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *ProductRepository) UpdateStock(id uint, quantity float64) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).
		Update("stock", gorm.Expr("stock - ?", quantity)).Error
}

func (r *ProductRepository) GetAllWithStock() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Where("stock > 0").Find(&products).Error
	return products, err
}

// Category methods
type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}
