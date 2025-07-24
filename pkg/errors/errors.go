package errors

import (
	"fmt"
	"hermes-api/pkg/logger"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeValidation         ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound           ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized       ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden          ErrorType = "FORBIDDEN"
	ErrorTypeConflict           ErrorType = "CONFLICT"
	ErrorTypeInternal           ErrorType = "INTERNAL_ERROR"
	ErrorTypeBadRequest         ErrorType = "BAD_REQUEST"
	ErrorTypeRateLimit          ErrorType = "RATE_LIMIT"
	ErrorTypeServiceUnavailable ErrorType = "SERVICE_UNAVAILABLE"
)

// ErrorCode represents specific error codes
type ErrorCode string

const (
	// App related errors
	ErrorCodeAppNotFound      ErrorCode = "APP_NOT_FOUND"
	ErrorCodeAppAlreadyExists ErrorCode = "APP_ALREADY_EXISTS"
	ErrorCodeInvalidAPIKey    ErrorCode = "INVALID_API_KEY"
	ErrorCodeAppInactive      ErrorCode = "APP_INACTIVE"

	// Notification related errors
	ErrorCodeNotificationNotFound      ErrorCode = "NOTIFICATION_NOT_FOUND"
	ErrorCodeInvalidNotificationType   ErrorCode = "INVALID_NOTIFICATION_TYPE"
	ErrorCodeInvalidRecipient          ErrorCode = "INVALID_RECIPIENT"
	ErrorCodeNotificationQuotaExceeded ErrorCode = "NOTIFICATION_QUOTA_EXCEEDED"

	// Validation errors
	ErrorCodeRequiredField ErrorCode = "REQUIRED_FIELD"
	ErrorCodeInvalidFormat ErrorCode = "INVALID_FORMAT"
	ErrorCodeInvalidValue  ErrorCode = "INVALID_VALUE"
	ErrorCodeFieldTooLong  ErrorCode = "FIELD_TOO_LONG"
	ErrorCodeFieldTooShort ErrorCode = "FIELD_TOO_SHORT"

	// System errors
	ErrorCodeDatabaseError        ErrorCode = "DATABASE_ERROR"
	ErrorCodeRedisError           ErrorCode = "REDIS_ERROR"
	ErrorCodeExternalServiceError ErrorCode = "EXTERNAL_SERVICE_ERROR"
)

// AppError represents a structured application error
type AppError struct {
	Type       ErrorType              `json:"type"`
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	Timestamp  string                 `json:"timestamp"`
	HTTPStatus int                    `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Type, e.Code, e.Message)
}

// Log logs the error using the structured logger
func (e *AppError) Log(msg string, fields ...zap.Field) {
	errorFields := []zap.Field{
		zap.String("error_type", string(e.Type)),
		zap.String("error_code", string(e.Code)),
		zap.String("error_message", e.Message),
		zap.String("timestamp", e.Timestamp),
	}

	if e.RequestID != "" {
		errorFields = append(errorFields, zap.String("request_id", e.RequestID))
	}

	if e.HTTPStatus != 0 {
		errorFields = append(errorFields, zap.Int("http_status", e.HTTPStatus))
	}

	if e.Details != nil {
		errorFields = append(errorFields, zap.Any("details", e.Details))
	}

	// Add any additional fields passed to the method
	errorFields = append(errorFields, fields...)

	logger.Error(msg, e, errorFields...)
}

// New creates a new AppError
func New(errorType ErrorType, code ErrorCode, message string) *AppError {
	return &AppError{
		Type:      errorType,
		Code:      code,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// WithRequestID adds request ID to the error
func (e *AppError) WithRequestID(requestID string) *AppError {
	e.RequestID = requestID
	return e
}

// GetHTTPStatus returns the corresponding HTTP status code
func (e *AppError) GetHTTPStatus() int {
	if e.HTTPStatus != 0 {
		return e.HTTPStatus
	}

	switch e.Type {
	case ErrorTypeValidation, ErrorTypeBadRequest:
		return http.StatusBadRequest
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeRateLimit:
		return http.StatusTooManyRequests
	case ErrorTypeServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
