package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"restaurant-api/internal/models"
	"restaurant-api/internal/repositories"
	"restaurant-api/internal/utils"
)

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// AuthService handles authentication business logic
type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user
func (s *AuthService) Register(request models.RegisterUserRequest) (*models.User, error) {
	// Check if user with email already exists
	existingUser, err := s.userRepo.GetByEmail(request.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
		Role:     "user", // Default role is user
	}

	return s.userRepo.Create(user)
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(request models.LoginUserRequest) (*models.User, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(request.Email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Verify password
	if !utils.CheckPasswordHash(request.Password, user.Password) {
		return nil, "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// ExtractTokenClaims extracts claims from the JWT token in the request
func (s *AuthService) ExtractTokenClaims(c echo.Context) (*JWTClaims, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	return claims, nil
}

// generateToken creates a new JWT token for a user
func (s *AuthService) generateToken(user *models.User) (string, error) {
	// Create claims
	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
