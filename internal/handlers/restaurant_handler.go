package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"

	"restaurant-api/internal/models"
	"restaurant-api/internal/services"
	"restaurant-api/internal/utils"
)

// RestaurantHandler handles restaurant-related requests
type RestaurantHandler struct {
	restaurantService *services.RestaurantService
	authService       *services.AuthService
	validator         *validator.Validate
}

// NewRestaurantHandler creates a new RestaurantHandler instance
func NewRestaurantHandler(restaurantService *services.RestaurantService, authService *services.AuthService) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
		authService:       authService,
		validator:         validator.New(),
	}
}

// GetUserRestaurants godoc
// @Summary Get all restaurants for a user
// @Description Get all restaurants owned by a specific user
// @Tags restaurants
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} utils.Response{data=[]models.RestaurantResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /users/{userId}/restaurants [get]
func (h *RestaurantHandler) GetUserRestaurants(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid user ID", err.Error()))
	}

	// Check permissions
	claims, err := h.authService.ExtractTokenClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Invalid token", err.Error()))
	}

	// Users can only view their own restaurants unless they are admins
	if claims.UserID != uint(userID) && claims.Role != "admin" {
		return c.JSON(http.StatusForbidden, utils.NewErrorResponse("Permission denied", "You don't have permission to access this resource"))
	}

	restaurants, err := h.restaurantService.GetRestaurantsByUserID(uint(userID))
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, utils.NewErrorResponse("User not found", "The requested user does not exist"))
		}
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get restaurants", err.Error()))
	}

	// Convert to response objects
	restaurantResponses := make([]models.RestaurantResponse, len(restaurants))
	for i, restaurant := range restaurants {
		restaurantResponses[i] = restaurant.ToResponse()
	}

	return c.JSON(http.StatusOK, utils.NewSuccessResponse("Restaurants retrieved successfully", restaurantResponses))
}

// GetUserRestaurant godoc
// @Summary Get a specific restaurant for a user
// @Description Get a specific restaurant owned by a user
// @Tags restaurants
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param id path int true "Restaurant ID"
// @Success 200 {object} utils.Response{data=models.RestaurantResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /users/{userId}/restaurants/{id} [get]
func (h *RestaurantHandler) GetUserRestaurant(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid user ID", err.Error()))
	}

	restaurantID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid restaurant ID", err.Error()))
	}

	// Check permissions
	claims, err := h.authService.ExtractTokenClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Invalid token", err.Error()))
	}

	// Users can only view their own restaurants unless they are admins
	if claims.UserID != uint(userID) && claims.Role != "admin" {
		return c.JSON(http.StatusForbidden, utils.NewErrorResponse("Permission denied", "You don't have permission to access this resource"))
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(uint(restaurantID), uint(userID))
	if err != nil {
		if errors.Is(err, services.ErrRestaurantNotFound) {
			return c.JSON(http.StatusNotFound, utils.NewErrorResponse("Restaurant not found", "The requested restaurant does not exist"))
		}
		if errors.Is(err, services.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, utils.NewErrorResponse("User not found", "The requested user does not exist"))
		}
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get restaurant", err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewSuccessResponse("Restaurant retrieved successfully", restaurant.ToResponse()))
}

// CreateRestaurant godoc
// @Summary Create a new restaurant
// @Description Create a new restaurant for the authenticated user
// @Tags restaurants
// @Accept json
// @Produce json
// @Param restaurant body models.CreateRestaurantRequest true "Restaurant creation data"
// @Success 201 {object} utils.Response{data=models.RestaurantResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /restaurants [post]
func (h *RestaurantHandler) CreateRestaurant(c echo.Context) error {
	var request models.CreateRestaurantRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.validator.Struct(request); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Validation failed", err.Error()))
	}

	// Get user ID from token
	claims, err := h.authService.ExtractTokenClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Invalid token", err.Error()))
	}

	restaurant, err := h.restaurantService.CreateRestaurant(request, claims.UserID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, utils.NewErrorResponse("User not found", "The user does not exist"))
		}
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to create restaurant", err.Error()))
	}

	return c.JSON(http.StatusCreated, utils.NewSuccessResponse("Restaurant created successfully", restaurant.ToResponse()))
}

// UpdateRestaurant godoc
// @Summary Update a restaurant
// @Description Update a restaurant by ID
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path int true "Restaurant ID"
// @Param restaurant body models.UpdateRestaurantRequest true "Restaurant update data"
// @Success 200 {object} utils.Response{data=models.RestaurantResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /restaurants/{id} [put]
func (h *RestaurantHandler) UpdateRestaurant(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid restaurant ID", err.Error()))
	}

	var request models.UpdateRestaurantRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.validator.Struct(request); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Validation failed", err.Error()))
	}

	// Check permissions
	claims, err := h.authService.ExtractTokenClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Invalid token", err.Error()))
	}

	// Get restaurant to check ownership
	restaurant, err := h.restaurantService.GetRestaurantByIDWithoutUserCheck(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrRestaurantNotFound) {
			return c.JSON(http.StatusNotFound, utils.NewErrorResponse("Restaurant not found", "The requested restaurant does not exist"))
		}
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get restaurant", err.Error()))
	}

	// Users can only update their own restaurants unless they are admins
	if restaurant.UserID != claims.UserID && claims.Role != "admin" {
		return c.JSON(http.StatusForbidden, utils.NewErrorResponse("Permission denied", "You don't have permission to update this restaurant"))
	}

	updatedRestaurant, err := h.restaurantService.UpdateRestaurant(uint(id), request)
	if err != nil {
		if errors.Is(err, services.ErrRestaurantNotFound) {
			return c.JSON(http.StatusNotFound, utils.NewErrorResponse("Restaurant not found", "The requested restaurant does not exist"))
		}
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to update restaurant", err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewSuccessResponse("Restaurant updated successfully", updatedRestaurant.ToResponse()))
}

// DeleteRestaurant godoc
// @Summary Delete a restaurant
// @Description Delete a restaurant by ID
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path int true "Restaurant ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /restaurants/{id} [delete]
func (h *RestaurantHandler) DeleteRestaurant(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid restaurant ID", err.Error()))
	}

	// Check permissions
	claims, err := h.authService.ExtractTokenClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Invalid token", err.Error()))
	}

	// Get restaurant to check ownership
	restaurant, err := h.restaurantService.GetRestaurantByIDWithoutUserCheck(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrRestaurantNotFound) {
			return c.JSON(http.StatusNotFound, utils.NewErrorResponse("Restaurant not found", "The requested restaurant does not exist"))
		}
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get restaurant", err.Error()))
	}

	// Users can only delete their own restaurants unless they are admins
	if restaurant.UserID != claims.UserID && claims.Role != "admin" {
		return c.JSON(http.StatusForbidden, utils.NewErrorResponse("Permission denied", "You don't have permission to delete this restaurant"))
	}

	err = h.restaurantService.DeleteRestaurant(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrRestaurantNotFound) {
			return c.JSON(http.StatusNotFound, utils.NewErrorResponse("Restaurant not found", "The requested restaurant does not exist"))
		}
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to delete restaurant", err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewSuccessResponse("Restaurant deleted successfully", nil))
}
