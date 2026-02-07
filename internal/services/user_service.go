package services

import (
	"errors"
	"time"

	"github.com/bexiiiii/smart_food_store/internal/config"
	"github.com/bexiiiii/smart_food_store/internal/models"
	"github.com/bexiiiii/smart_food_store/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repository.UserRepository
	cartRepo *repository.CartRepository
	config   *config.Config
}

func NewUserService(userRepo *repository.UserRepository, cartRepo *repository.CartRepository, cfg *config.Config) *UserService {
	return &UserService{
		userRepo: userRepo,
		cartRepo: cartRepo,
		config:   cfg,
	}
}

func (s *UserService) Register(req *models.UserRegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     models.RoleUser,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Create cart for user
	_, err = s.cartRepo.GetOrCreateByUserID(user.ID)
	if err != nil {
		return nil, errors.New("failed to create cart for user")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *UserService) Login(req *models.UserLoginRequest) (*models.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid email or password")
		}
		return nil, errors.New("failed to find user")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *UserService) GetByID(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *UserService) GetAll() ([]models.UserResponse, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var response []models.UserResponse
	for _, user := range users {
		response = append(response, models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		})
	}

	return response, nil
}

func (s *UserService) UpdateRole(id uint, role models.Role) error {
	return s.userRepo.UpdateRole(id, role)
}

func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) generateToken(user *models.User) (string, error) {
	duration, err := time.ParseDuration(s.config.JWTExpiration)
	if err != nil {
		duration = 24 * time.Hour
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(duration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}
