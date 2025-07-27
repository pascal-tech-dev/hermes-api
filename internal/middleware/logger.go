package middleware

import (
	"hermes-api/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Logger middleware that integrates with our logging system
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request details
		method := c.Method()
		path := c.Path()
		status := c.Response().StatusCode()
		ip := c.IP()
		userAgent := c.Get("User-Agent")
		requestID := c.Locals("X-Request-ID")
		requestIDStr := "unknown"
		if requestID != nil {
			if id, ok := requestID.(string); ok {
				requestIDStr = id
			}
		}

		// Only log successful requests (2xx, 3xx) and client errors (4xx) that aren't 404s
		// Let the ErrorHandler handle 404s and server errors (5xx)
		if status < 400 || (status >= 400 && status < 500 && status != 404) {
			logger.Info("HTTP Request",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", status),
				zap.Duration("latency", latency),
				zap.String("ip", ip),
				zap.String("user_agent", userAgent),
				zap.String("request_id", requestIDStr),
			)
		}

		return err
	}
}
