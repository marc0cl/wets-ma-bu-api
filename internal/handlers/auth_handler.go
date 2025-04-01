package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"

	"restaurant-api/internal/models"
	"restaurant-api/internal/services"
	"restaurant-api/internal/utils"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *services.AuthService
	validator   *validator.Validate
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterUserRequest true "User registration data"
// @Success 201 {object} utils.Response{data=models.UserResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var request models.RegisterUserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.validator.Struct(request); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Validation failed", err.Error()))
	}

	user, err := h.authService.Register(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to register user", err.Error()))
	}

	return c.JSON(http.StatusCreated, utils.NewSuccessResponse("User registered successfully", user.ToResponse()))
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.LoginUserRequest true "User login data"
// @Success 200 {object} utils.Response{data=map[string]interface{}}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var request models.LoginUserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.validator.Struct(request); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Validation failed", err.Error()))
	}

	user, token, err := h.authService.Login(request)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Authentication failed", err.Error()))
	}

	response := map[string]interface{}{
		"user":  user.ToResponse(),
		"token": token,
	}

	return c.JSON(http.StatusOK, utils.NewSuccessResponse("Login successful", response))
}
