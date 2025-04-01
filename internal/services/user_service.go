package services

import (
	"errors"

	"restaurant-api/internal/models"
	"restaurant-api/internal/repositories"
	"restaurant-api/internal/utils"
)

// Common errors
var (
	ErrUserNotFound = errors.New("user not found")
)

// UserService handles user-related business logic
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(id uint, request models.UpdateUserRequest) (*models.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Update user fields if provided
	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Email != "" && request.Email != user.Email {
		// Check if email is already in use
		existingUser, err := s.userRepo.GetByEmail(request.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already in use")
		}
		user.Email = request.Email
	}

	if request.Password != "" {
		// Hash password
		hashedPassword, err := utils.HashPassword(request.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	if request.Role != "" {
		// Validate role
		if request.Role != "admin" && request.Role != "user" {
			return nil, errors.New("invalid role")
		}
		user.Role = request.Role
	}

	// Save updated user
	return s.userRepo.Update(user)
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id uint) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	// Delete user
	return s.userRepo.Delete(id)
}
