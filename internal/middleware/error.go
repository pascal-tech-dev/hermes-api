package middleware

import (
	"hermes-api/pkg/errors"
	"hermes-api/pkg/logger"
	"hermes-api/pkg/response"
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
	requestID := c.Locals("X-Request-ID")
	requestIDStr := "unknown"
	if requestID != nil {
		if id, ok := requestID.(string); ok {
			requestIDStr = id
		}
	}

	// Handle AppError
	if appErr, ok := err.(*errors.AppError); ok {
		appErr.RequestID = requestIDStr

		logger.Error("Application error occurred", err,
			zap.String("request_id", requestIDStr),
			zap.String("error_type", string(appErr.Type)),
			zap.String("error_code", string(appErr.Code)),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)

		options := response.ErrorResponse(appErr, appErr.Message)
		return response.ApiResponse(c, options)
	}

	// Handle Fiber errors
	if fiberErr, ok := err.(*fiber.Error); ok {
		errorType, errorCode := errors.MapFiberError(fiberErr.Code, fiberErr.Message)
		appErr := errors.New(errorType, errorCode, fiberErr.Message)
		appErr.RequestID = requestIDStr
		appErr.HTTPStatus = fiberErr.Code

		logger.Error("Fiber error occurred", appErr,
			zap.String("request_id", requestIDStr),
			zap.Int("status_code", fiberErr.Code),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("error_message", fiberErr.Message),
		)

		options := response.ErrorResponse(appErr, appErr.Message)
		return response.ApiResponse(c, options)
	}

	// Handle unknown errors
	appErr := errors.New(errors.ErrorTypeInternal, errors.ErrorCodeUnknownError, "Internal server error")
	appErr.RequestID = requestIDStr

	logger.Error("Unknown error occurred", err,
		zap.String("request_id", requestIDStr),
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("stack_trace", string(debug.Stack())),
	)

	options := response.ErrorResponse(appErr, appErr.Message)
	return response.ApiResponse(c, options)
}
