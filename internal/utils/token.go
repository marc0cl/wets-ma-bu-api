package utils

import (
        "errors"
        "fmt"
        "strings"

        "github.com/labstack/echo"
)

// ExtractTokenFromHeader extracts the token from the Authorization header
func ExtractTokenFromHeader(c echo.Context) (string, error) {
        auth := c.Request().Header.Get("Authorization")
        if auth == "" {
                return "", errors.New("authorization header is required")
        }

        parts := strings.SplitN(auth, " ", 2)
        if !(len(parts) == 2 && parts[0] == "Bearer") {
                return "", fmt.Errorf("invalid authorization header format")
        }

        return parts[1], nil
}
