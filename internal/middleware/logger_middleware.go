package middleware

import (
        "github.com/labstack/echo"
        "github.com/labstack/echo/middleware"
)

// Logger returns a logger middleware
func Logger() echo.MiddlewareFunc {
        // Return simpler Echo middleware
        return middleware.LoggerWithConfig(middleware.LoggerConfig{
                Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
        })
}
