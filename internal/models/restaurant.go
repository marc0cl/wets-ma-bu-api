package models

import (
	"time"

	"gorm.io/gorm"
)

// Restaurant represents a restaurant in the system
type Restaurant struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:1000" json:"description"`
	Address     string         `gorm:"size:200" json:"address"`
	Phone       string         `gorm:"size:20" json:"phone"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// RestaurantResponse is a struct for restaurant data that is safe to send in API responses
type RestaurantResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts a Restaurant model to a RestaurantResponse
func (r *Restaurant) ToResponse() RestaurantResponse {
	return RestaurantResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Address:     r.Address,
		Phone:       r.Phone,
		UserID:      r.UserID,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// CreateRestaurantRequest represents the request body for creating a restaurant
type CreateRestaurantRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=1000"`
	Address     string `json:"address" validate:"required,max=200"`
	Phone       string `json:"phone" validate:"max=20"`
}

// UpdateRestaurantRequest represents the request body for updating a restaurant
type UpdateRestaurantRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=100"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	Address     string `json:"address" validate:"omitempty,max=200"`
	Phone       string `json:"phone" validate:"omitempty,max=20"`
}
