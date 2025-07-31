package middleware

import (
	"hermes-api/internal/service"
	"hermes-api/pkg/errorx"
	"hermes-api/pkg/logger"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// AuthMiddleware creates middleware for JWT authentication
func AuthMiddleware(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			logger.Error("Missing Authorization header", nil)
			appErr := errorx.New(errorx.ErrorTypeUnauthorized, errorx.ErrorCodeFiberUnauthorized, "Missing authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": appErr,
			})
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Error("Invalid authorization header format", nil, zap.String("header", authHeader))
			appErr := errorx.New(errorx.ErrorTypeUnauthorized, errorx.ErrorCodeFiberUnauthorized, "Invalid authorization header format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": appErr,
			})
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token and get user
		user, err := authService.GetUserFromToken(token)
		if err != nil {
			logger.Error("Invalid token", err)
			appErr := errorx.New(errorx.ErrorTypeUnauthorized, errorx.ErrorCodeFiberUnauthorized, "Invalid or expired token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": appErr,
			})
		}

		// Check if user is active
		if !user.IsActive {
			logger.Error("User account is deactivated", nil, zap.String("user_id", user.ID.String()))
			appErr := errorx.New(errorx.ErrorTypeForbidden, errorx.ErrorCodeFiberForbidden, "User account is deactivated")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": appErr,
			})
		}

		// Set user in context
		c.Locals("user", user)

		// Continue to next handler
		return c.Next()
	}
}

// OptionalAuthMiddleware creates middleware for optional JWT authentication
// This allows endpoints to work with or without authentication
func OptionalAuthMiddleware(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// No auth header, continue without user
			return c.Next()
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			// Invalid format, continue without user
			return c.Next()
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Try to validate token and get user
		user, err := authService.GetUserFromToken(token)
		if err != nil {
			// Invalid token, continue without user
			return c.Next()
		}

		// Check if user is active
		if !user.IsActive {
			// Inactive user, continue without user
			return c.Next()
		}

		// Set user in context
		c.Locals("user", user)

		// Continue to next handler
		return c.Next()
	}
}
