package middleware

import (
	"hermes-api/pkg/errorx"
	"hermes-api/pkg/logger"
	"hermes-api/pkg/response"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ErrorHandler returns a Fiber handler that handles errors and returns appropriate responses
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
		if appErr, ok := err.(*errorx.AppError); ok {
			appErr.RequestID = requestIDStr

			logger.Error("Application error occurred", err,
				zap.String("request_id", requestIDStr),
				zap.String("error_type", string(appErr.Type)),
				zap.String("error_code", string(appErr.Code)),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
			)

			return response.ErrorResponse(appErr, appErr.Message).Send(c)
		}

		// Handle Fiber errors
		if fiberErr, ok := err.(*fiber.Error); ok {
			errorType, errorCode := errorx.MapFiberError(fiberErr.Code, fiberErr.Message)
			appErr := errorx.New(errorType, errorCode, fiberErr.Message)
			appErr.RequestID = requestIDStr
			appErr.HTTPStatus = fiberErr.Code

			logger.Error("Fiber error occurred", appErr,
				zap.String("request_id", requestIDStr),
				zap.Int("status_code", fiberErr.Code),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.String("error_message", fiberErr.Message),
			)

			return response.ErrorResponse(appErr, appErr.Message).Send(c)
		}

		// Handle unknown errors
		appErr := errorx.New(errorx.ErrorTypeInternal, errorx.ErrorCodeUnknownError, "Internal server error")
		appErr.RequestID = requestIDStr

		logger.Error("Unknown error occurred", err,
			zap.String("request_id", requestIDStr),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("stack_trace", string(debug.Stack())),
		)

		return response.ErrorResponse(appErr, appErr.Message).Send(c)
	}
}
