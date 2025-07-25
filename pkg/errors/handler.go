package errors

import (
	"hermes-api/pkg/logger"
	"net/http"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ErrorHandler handles errors and returns appropriate responses
func ErrorHandler(c *fiber.Ctx) error {
	err := c.Next()
	if err == nil {
		return nil
	}

	// Get request ID from context
	requestID := c.Get("X-Request-ID", "unknown")

	// Handle AppError
	if appErr, ok := err.(*AppError); ok {
		appErr.RequestID = requestID

		logger.Error("Application error occurred", err,
			zap.String("request_id", requestID),
			zap.String("error_type", string(appErr.Type)),
			zap.String("error_code", string(appErr.Code)),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)

		return c.Status(appErr.GetHTTPStatus()).JSON(appErr)
	}

	// Handle Fiber errors
	if fiberErr, ok := err.(*fiber.Error); ok {
		errorType, errorCode := MapFiberError(fiberErr.Code, fiberErr.Message)
		appErr := New(errorType, errorCode, fiberErr.Message)
		appErr.RequestID = requestID
		appErr.HTTPStatus = fiberErr.Code

		logger.Error("Fiber error occurred", err,
			zap.String("request_id", requestID),
			zap.Int("status_code", fiberErr.Code),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)

		return c.Status(appErr.GetHTTPStatus()).JSON(appErr)
	}

	// Handle unknown errors
	appErr := New(ErrorTypeInternal, ErrorCodeUnknownError, "Internal server error")
	appErr.RequestID = requestID

	logger.Error("Unknown error occurred", err,
		zap.String("request_id", requestID),
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("stack_trace", string(debug.Stack())),
	)

	return c.Status(http.StatusInternalServerError).JSON(appErr)
}
