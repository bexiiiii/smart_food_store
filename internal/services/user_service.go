package services

import (
	"github.com/bexiiiii/smart_food_store/internal/models"
	"github.com/bexiiiii/smart_food_store/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(user models.User) models.User {
	return s.repo.Create(user)
}

func (s *UserService) GetAll() []models.User {
	return s.repo.GetAll()
}
