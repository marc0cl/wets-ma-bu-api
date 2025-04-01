package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents user roles
type Role string

const (
	// AdminRole represents an administrator user
	AdminRole Role = "admin"
	// UserRole represents a regular user
	UserRole Role = "user"
)

// User represents a user in the system
type User struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Email       string         `gorm:"size:100;not null;unique" json:"email"`
	Password    string         `gorm:"size:100;not null" json:"-"`
	Role        string         `gorm:"size:20;not null;default:user" json:"role"`
	Restaurants []Restaurant   `gorm:"foreignKey:UserID" json:"restaurants,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserResponse is a struct for user data that is safe to send in API responses
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts a User model to a UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// RegisterUserRequest represents the request body for user registration
type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginUserRequest represents the request body for user login
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=100"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
	Role     string `json:"role" validate:"omitempty,oneof=admin user"`
}
