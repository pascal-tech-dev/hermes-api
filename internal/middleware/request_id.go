package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestID middleware generates and sets a unique request ID
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request ID is already set in headers
		requestID := c.Get("X-Request-ID")

		// If not set, generate a new one
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Set request ID in context for other middleware/handlers to use
		c.Locals("X-Request-ID", requestID)

		// Add request ID to response headers
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}

// generateRequestID creates a unique request ID using UUID v4
func generateRequestID() string {
	return uuid.New().String()
}

// GetRequestID extracts request ID from context
func GetRequestID(c *fiber.Ctx) string {
	return c.Get("X-Request-ID", "unknown")
}
