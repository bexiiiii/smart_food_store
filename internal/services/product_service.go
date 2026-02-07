package services

import (
	"errors"

	"github.com/bexiiiii/smart_food_store/internal/models"
	"github.com/bexiiiii/smart_food_store/internal/repository"
)

type ProductService struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
}

func NewProductService(productRepo *repository.ProductRepository, categoryRepo *repository.CategoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// Product methods
func (s *ProductService) Create(req *models.ProductCreateRequest) (*models.Product, error) {
	// Validate category exists
	_, err := s.categoryRepo.GetByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Unit:        req.Unit,
		CategoryID:  req.CategoryID,
		ImageURL:    req.ImageURL,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, errors.New("failed to create product")
	}

	return s.productRepo.GetByID(product.ID)
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {
	return s.productRepo.GetByID(id)
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.productRepo.GetAll()
}

func (s *ProductService) GetByCategory(categoryID uint) ([]models.Product, error) {
	return s.productRepo.GetByCategory(categoryID)
}

func (s *ProductService) Search(query string) ([]models.Product, error) {
	return s.productRepo.Search(query)
}

func (s *ProductService) GetByIDs(ids []uint) ([]models.Product, error) {
	return s.productRepo.GetByIDs(ids)
}

func (s *ProductService) Update(id uint, req *models.ProductUpdateRequest) (*models.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Unit != nil {
		product.Unit = *req.Unit
	}
	if req.CategoryID != nil {
		// Validate category exists
		_, err := s.categoryRepo.GetByID(*req.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
		product.CategoryID = *req.CategoryID
	}
	if req.ImageURL != nil {
		product.ImageURL = *req.ImageURL
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, errors.New("failed to update product")
	}

	return s.productRepo.GetByID(id)
}

func (s *ProductService) Delete(id uint) error {
	return s.productRepo.Delete(id)
}

func (s *ProductService) UpdateStock(id uint, quantity float64) error {
	return s.productRepo.UpdateStock(id, quantity)
}

func (s *ProductService) GetAllWithStock() ([]models.Product, error) {
	return s.productRepo.GetAllWithStock()
}

// Category methods
func (s *ProductService) CreateCategory(req *models.CategoryCreateRequest) (*models.Category, error) {
	category := &models.Category{
		Name: req.Name,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, errors.New("failed to create category")
	}

	return category, nil
}

func (s *ProductService) GetCategoryByID(id uint) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *ProductService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

func (s *ProductService) UpdateCategory(id uint, name string) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	category.Name = name

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, errors.New("failed to update category")
	}

	return category, nil
}

func (s *ProductService) DeleteCategory(id uint) error {
	return s.categoryRepo.Delete(id)
}
