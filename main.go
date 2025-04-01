package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"restaurant-api/config"
	"restaurant-api/internal/handlers"
	custommiddleware "restaurant-api/internal/middleware"
	"restaurant-api/internal/models"
	"restaurant-api/internal/repositories"
	"restaurant-api/internal/services"
)

// @title Restaurant Management API
// @version 1.0
// @description API for managing users and restaurants
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Initialize configuration
	config := config.LoadConfig()

	// Initialize database
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate database models
	err = db.AutoMigrate(&models.User{}, &models.Restaurant{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	restaurantRepo := repositories.NewRestaurantRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, config.JWTSecret)
	userService := services.NewUserService(userRepo)
	restaurantService := services.NewRestaurantService(restaurantRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService, authService)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService, authService)

	// Initialize Echo
	e := echo.New()

	// Set up middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(custommiddleware.Logger())
	e.Use(custommiddleware.CORS())

	// API documentation route
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Restaurant API - Welcome to the API Server")
	})

	// Serve Swagger JSON file directly
	e.GET("/swagger.json", func(c echo.Context) error {
		filePath := "docs/swagger.json"
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return c.String(500, "Could not read swagger.json file")
		}
		return c.JSONBlob(200, data)
	})

	// Serve Swagger YAML file directly
	e.GET("/swagger.yaml", func(c echo.Context) error {
		filePath := "docs/swagger.yaml"
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return c.String(500, "Could not read swagger.yaml file")
		}
		return c.Blob(200, "application/yaml", data)
	})

	// Serve Swagger UI HTML
	e.GET("/swagger", func(c echo.Context) error {
		filePath := "docs/swagger-ui.html"
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return c.String(500, "Could not read swagger UI file")
		}
		return c.HTML(200, string(data))
	})

	// Group API routes
	api := e.Group("/api/v1")

	// Auth routes
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	// User routes
	api.GET("/users/:id", userHandler.GetUser, custommiddleware.JWT(config.JWTSecret))
	api.PUT("/users/:id", userHandler.UpdateUser, custommiddleware.JWT(config.JWTSecret))
	api.DELETE("/users/:id", userHandler.DeleteUser, custommiddleware.JWT(config.JWTSecret))

	// Restaurant routes
	api.GET("/users/:userId/restaurants", restaurantHandler.GetUserRestaurants, custommiddleware.JWT(config.JWTSecret))
	api.GET("/users/:userId/restaurants/:id", restaurantHandler.GetUserRestaurant, custommiddleware.JWT(config.JWTSecret))
	api.POST("/restaurants", restaurantHandler.CreateRestaurant, custommiddleware.JWT(config.JWTSecret))
	api.PUT("/restaurants/:id", restaurantHandler.UpdateRestaurant, custommiddleware.JWT(config.JWTSecret))
	api.DELETE("/restaurants/:id", restaurantHandler.DeleteRestaurant, custommiddleware.JWT(config.JWTSecret))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
}
