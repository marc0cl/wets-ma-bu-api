package handlers

import (
        "errors"
        "net/http"
        "strconv"

        "github.com/labstack/echo"
        "github.com/go-playground/validator/v10"

        "restaurant-api/internal/models"
        "restaurant-api/internal/services"
        "restaurant-api/internal/utils"
)

// UserHandler handles user-related requests
type UserHandler struct {
        userService *services.UserService
        authService *services.AuthService
        validator   *validator.Validate
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService *services.UserService, authService *services.AuthService) *UserHandler {
        return &UserHandler{
                userService: userService,
                authService: authService,
                validator:   validator.New(),
        }
}

// GetUser godoc
// @Summary Get user information
// @Description Get a user's information by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response{data=models.UserResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
        id, err := strconv.ParseUint(c.Param("id"), 10, 32)
        if err != nil {
                return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid user ID", err.Error()))
        }

        // Check permissions
        claims, err := h.authService.ExtractTokenClaims(c)
        if err != nil {
                return c.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Invalid token", err.Error()))
        }

        // Only admins can access other users' data
        if claims.UserID != uint(id) && claims.Role != string(models.AdminRole) {
                return c.JSON(http.StatusForbidden, utils.NewErrorResponse("Permission denied", "You don't have permission to access this resource"))
        }

        user, err := h.userService.GetUserByID(uint(id))
        if err != nil {
                if errors.Is(err, services.ErrUserNotFound) {
                        return c.JSON(http.StatusNotFound, utils.NewErrorResponse("User not found", "The requested user does not exist"))
                }
                return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get user", err.Error()))
        }

        return c.JSON(http.StatusOK, utils.NewSuccessResponse("User retrieved successfully", user.ToResponse()))
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update a user's information by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.UpdateUserRequest true "User update data"
// @Success 200 {object} utils.Response{data=models.UserResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
        id, err := strconv.ParseUint(c.Param("id"), 10, 32)
        if err != nil {
                return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid user ID", err.Error()))
        }

        var request models.UpdateUserRequest
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

        // Check if trying to update role while not being an admin
        if request.Role != "" && claims.Role != string(models.AdminRole) {
                return c.JSON(http.StatusForbidden, utils.NewErrorResponse("Permission denied", "Only admins can change user roles"))
        }

        // Only admins can update other users' data
        if claims.UserID != uint(id) && claims.Role != string(models.AdminRole) {
                return c.JSON(http.StatusForbidden, utils.NewErrorResponse("Permission denied", "You don't have permission to update this user"))
        }

        user, err := h.userService.UpdateUser(uint(id), request)
        if err != nil {
                if errors.Is(err, services.ErrUserNotFound) {
                        return c.JSON(http.StatusNotFound, utils.NewErrorResponse("User not found", "The requested user does not exist"))
                }
                return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to update user", err.Error()))
        }

        return c.JSON(http.StatusOK, utils.NewSuccessResponse("User updated successfully", user.ToResponse()))
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
        id, err := strconv.ParseUint(c.Param("id"), 10, 32)
        if err != nil {
                return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid user ID", err.Error()))
        }

        // Check permissions
        claims, err := h.authService.ExtractTokenClaims(c)
        if err != nil {
                return c.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Invalid token", err.Error()))
        }

        // Only admins can delete other users
        if claims.UserID != uint(id) && claims.Role != string(models.AdminRole) {
                return c.JSON(http.StatusForbidden, utils.NewErrorResponse("Permission denied", "You don't have permission to delete this user"))
        }

        err = h.userService.DeleteUser(uint(id))
        if err != nil {
                if errors.Is(err, services.ErrUserNotFound) {
                        return c.JSON(http.StatusNotFound, utils.NewErrorResponse("User not found", "The requested user does not exist"))
                }
                return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to delete user", err.Error()))
        }

        return c.JSON(http.StatusOK, utils.NewSuccessResponse("User deleted successfully", nil))
}
