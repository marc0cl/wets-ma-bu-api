package repositories

import (
	"errors"

	"gorm.io/gorm"

	"restaurant-api/internal/models"
)

// RestaurantRepository handles database operations for restaurants
type RestaurantRepository struct {
	db *gorm.DB
}

// NewRestaurantRepository creates a new RestaurantRepository instance
func NewRestaurantRepository(db *gorm.DB) *RestaurantRepository {
	return &RestaurantRepository{
		db: db,
	}
}

// Create creates a new restaurant in the database
func (r *RestaurantRepository) Create(restaurant *models.Restaurant) (*models.Restaurant, error) {
	if err := r.db.Create(restaurant).Error; err != nil {
		return nil, err
	}
	return restaurant, nil
}

// GetByID retrieves a restaurant by ID
func (r *RestaurantRepository) GetByID(id uint) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	if err := r.db.First(&restaurant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("restaurant not found")
		}
		return nil, err
	}
	return &restaurant, nil
}

// GetByUserID retrieves all restaurants for a user
func (r *RestaurantRepository) GetByUserID(userID uint) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	if err := r.db.Where("user_id = ?", userID).Find(&restaurants).Error; err != nil {
		return nil, err
	}
	return restaurants, nil
}

// Update updates a restaurant in the database
func (r *RestaurantRepository) Update(restaurant *models.Restaurant) (*models.Restaurant, error) {
	if err := r.db.Save(restaurant).Error; err != nil {
		return nil, err
	}
	return restaurant, nil
}

// Delete deletes a restaurant from the database
func (r *RestaurantRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Restaurant{}, id).Error; err != nil {
		return err
	}
	return nil
}

// List retrieves all restaurants
func (r *RestaurantRepository) List() ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	if err := r.db.Find(&restaurants).Error; err != nil {
		return nil, err
	}
	return restaurants, nil
}
