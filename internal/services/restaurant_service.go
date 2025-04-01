package services

import (
	"errors"

	"restaurant-api/internal/models"
	"restaurant-api/internal/repositories"
)

// Common errors
var (
	ErrRestaurantNotFound = errors.New("restaurant not found")
)

// RestaurantService handles restaurant-related business logic
type RestaurantService struct {
	restaurantRepo *repositories.RestaurantRepository
	userRepo       *repositories.UserRepository
}

// NewRestaurantService creates a new RestaurantService instance
func NewRestaurantService(restaurantRepo *repositories.RestaurantRepository, userRepo *repositories.UserRepository) *RestaurantService {
	return &RestaurantService{
		restaurantRepo: restaurantRepo,
		userRepo:       userRepo,
	}
}

// GetRestaurantsByUserID retrieves all restaurants for a user
func (s *RestaurantService) GetRestaurantsByUserID(userID uint) ([]models.Restaurant, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return s.restaurantRepo.GetByUserID(userID)
}

// GetRestaurantByID retrieves a restaurant by ID and checks user ownership
func (s *RestaurantService) GetRestaurantByID(id uint, userID uint) (*models.Restaurant, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	restaurant, err := s.restaurantRepo.GetByID(id)
	if err != nil {
		return nil, ErrRestaurantNotFound
	}

	// Check if restaurant belongs to the user
	if restaurant.UserID != userID {
		return nil, ErrRestaurantNotFound
	}

	return restaurant, nil
}

// GetRestaurantByIDWithoutUserCheck retrieves a restaurant by ID without checking user ownership
func (s *RestaurantService) GetRestaurantByIDWithoutUserCheck(id uint) (*models.Restaurant, error) {
	restaurant, err := s.restaurantRepo.GetByID(id)
	if err != nil {
		return nil, ErrRestaurantNotFound
	}

	return restaurant, nil
}

// CreateRestaurant creates a new restaurant
func (s *RestaurantService) CreateRestaurant(request models.CreateRestaurantRequest, userID uint) (*models.Restaurant, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	restaurant := &models.Restaurant{
		Name:        request.Name,
		Description: request.Description,
		Address:     request.Address,
		Phone:       request.Phone,
		UserID:      userID,
	}

	return s.restaurantRepo.Create(restaurant)
}

// UpdateRestaurant updates a restaurant
func (s *RestaurantService) UpdateRestaurant(id uint, request models.UpdateRestaurantRequest) (*models.Restaurant, error) {
	// Get existing restaurant
	restaurant, err := s.restaurantRepo.GetByID(id)
	if err != nil {
		return nil, ErrRestaurantNotFound
	}

	// Update restaurant fields if provided
	if request.Name != "" {
		restaurant.Name = request.Name
	}

	if request.Description != "" {
		restaurant.Description = request.Description
	}

	if request.Address != "" {
		restaurant.Address = request.Address
	}

	if request.Phone != "" {
		restaurant.Phone = request.Phone
	}

	// Save updated restaurant
	return s.restaurantRepo.Update(restaurant)
}

// DeleteRestaurant deletes a restaurant
func (s *RestaurantService) DeleteRestaurant(id uint) error {
	// Check if restaurant exists
	_, err := s.restaurantRepo.GetByID(id)
	if err != nil {
		return ErrRestaurantNotFound
	}

	// Delete restaurant
	return s.restaurantRepo.Delete(id)
}
