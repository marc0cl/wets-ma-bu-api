package middleware

import (
        "fmt"
        "net/http"
        "strings"

        "github.com/dgrijalva/jwt-go"
        "github.com/labstack/echo"

        "restaurant-api/internal/services"
)

// JWT middleware for handling authentication
func JWT(jwtSecret string) echo.MiddlewareFunc {
        return func(next echo.HandlerFunc) echo.HandlerFunc {
                return func(c echo.Context) error {
                        // Get authorization header
                        authHeader := c.Request().Header.Get("Authorization")
                        if authHeader == "" {
                                return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
                        }

                        // Check if the header has the Bearer prefix
                        if !strings.HasPrefix(authHeader, "Bearer ") {
                                return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization format, expected 'Bearer TOKEN'")
                        }

                        // Extract token
                        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

                        // Parse and validate token
                        token, err := jwt.ParseWithClaims(tokenString, &services.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
                                // Validate signing method
                                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                                        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
                                }
                                return []byte(jwtSecret), nil
                        })

                        if err != nil {
                                return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token: "+err.Error())
                        }

                        // Set parsed token in context
                        c.Set("user", token)
                        return next(c)
                }
        }
}
