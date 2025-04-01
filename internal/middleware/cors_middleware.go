package middleware

import (
        "github.com/labstack/echo"
        "github.com/labstack/echo/middleware"
)

// CORS returns a CORS middleware
func CORS() echo.MiddlewareFunc {
        return middleware.CORSWithConfig(middleware.CORSConfig{
                AllowOrigins: []string{"*"},
                AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
                AllowHeaders: []string{
                        echo.HeaderOrigin,
                        echo.HeaderContentType,
                        echo.HeaderAccept,
                        echo.HeaderAuthorization,
                        echo.HeaderXRequestedWith,
                },
                AllowCredentials: true,
        })
}
