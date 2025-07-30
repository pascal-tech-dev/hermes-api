package middleware

import (
	"hermes-api/pkg/errors"
	"hermes-api/pkg/logger"
	"hermes-api/pkg/response"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Recovery middleware that integrates with our error handling system
func Recovery() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Get request ID from context
				requestID := c.Get("X-Request-ID", "unknown")

				// Log the panic with stack trace
				logger.Error("Panic recovered", nil,
					zap.Any("panic", r),
					zap.String("request_id", requestID),
					zap.String("method", c.Method()),
					zap.String("path", c.Path()),
					zap.String("ip", c.IP()),
					zap.String("user_agent", c.Get("User-Agent")),
					zap.String("stack_trace", string(debug.Stack())),
				)

				// Create an internal server error
				appErr := errors.New(errors.ErrorTypeInternal, errors.ErrorCodeUnknownError, "Internal server error")
				appErr.RequestID = requestID
				appErr.HTTPStatus = fiber.StatusInternalServerError

				// Send error response
				_ = response.ErrorResponse(appErr, "Internal server error").Send(c)
			}
		}()

		// Continue to next middleware/handler
		return c.Next()
	}
}
