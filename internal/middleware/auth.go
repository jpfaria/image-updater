package middleware

import (
	"strings"

	"github.com/jpfaria/image-updater/internal/auth"
	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/log"
)

// JWTMiddleware is a middleware that validates JWT tokens
func JWTMiddleware(authService *auth.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip authentication for login and health check endpoints
			if c.Path() == "/api/auth/login" || c.Path() == "/health" || c.Path() == "/api/webhooks/docker" {
				return next(c)
			}

			// Get token from header
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return echo.NewHTTPError(401, "No token provided")
			}

			// Remove "Bearer " prefix if present
			if len(tokenString) > 7 && strings.HasPrefix(tokenString, "Bearer ") {
				tokenString = tokenString[7:]
			}

			// Validate token
			user, err := authService.ValidateToken(c.Request().Context(), tokenString)
			if err != nil {
				log.Errorf("Invalid token: %v", err)
				return echo.NewHTTPError(401, "Invalid token")
			}

			// Set user in context
			c.Set("user", user)

			return next(c)
		}
	}
}
