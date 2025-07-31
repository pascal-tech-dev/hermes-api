package errorx

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
	// Auth related errors
	ErrorCodeInvalidCredentials    ErrorCode = "INVALID_CREDENTIALS"
	ErrorCodeAccountDeactivated    ErrorCode = "ACCOUNT_DEACTIVATED"
	ErrorCodeUserAlreadyExists     ErrorCode = "USER_ALREADY_EXISTS"
	ErrorCodeTokenGenerationFailed ErrorCode = "TOKEN_GENERATION_FAILED"
	ErrorCodeTokenExpired          ErrorCode = "TOKEN_EXPIRED"
	ErrorCodeTokenInvalid          ErrorCode = "TOKEN_INVALID"

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
	ErrorCodeRequiredField       ErrorCode = "REQUIRED_FIELD"
	ErrorCodeInvalidFormat       ErrorCode = "INVALID_FORMAT"
	ErrorCodeInvalidValue        ErrorCode = "INVALID_VALUE"
	ErrorCodeFieldTooLong        ErrorCode = "FIELD_TOO_LONG"
	ErrorCodeFieldTooShort       ErrorCode = "FIELD_TOO_SHORT"
	ErrorCodeFieldLengthMismatch ErrorCode = "FIELD_LENGTH_MISMATCH"
	ErrorCodeFieldNotEqual       ErrorCode = "FIELD_NOT_EQUAL"

	// System errors
	ErrorCodeDatabaseError        ErrorCode = "DATABASE_ERROR"
	ErrorCodeRedisError           ErrorCode = "REDIS_ERROR"
	ErrorCodeExternalServiceError ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrorCodeUnknownError         ErrorCode = "UNKNOWN_ERROR"

	// Fiber/HTTP specific errors (4xx Client Errors)
	ErrorCodeFiberBadRequest                   ErrorCode = "FIBER_BAD_REQUEST"
	ErrorCodeFiberUnauthorized                 ErrorCode = "FIBER_UNAUTHORIZED"
	ErrorCodeFiberPaymentRequired              ErrorCode = "FIBER_PAYMENT_REQUIRED"
	ErrorCodeFiberForbidden                    ErrorCode = "FIBER_FORBIDDEN"
	ErrorCodeFiberNotFound                     ErrorCode = "FIBER_NOT_FOUND"
	ErrorCodeFiberMethodNotAllowed             ErrorCode = "FIBER_METHOD_NOT_ALLOWED"
	ErrorCodeFiberNotAcceptable                ErrorCode = "FIBER_NOT_ACCEPTABLE"
	ErrorCodeFiberProxyAuthRequired            ErrorCode = "FIBER_PROXY_AUTH_REQUIRED"
	ErrorCodeFiberRequestTimeout               ErrorCode = "FIBER_REQUEST_TIMEOUT"
	ErrorCodeFiberConflict                     ErrorCode = "FIBER_CONFLICT"
	ErrorCodeFiberGone                         ErrorCode = "FIBER_GONE"
	ErrorCodeFiberLengthRequired               ErrorCode = "FIBER_LENGTH_REQUIRED"
	ErrorCodeFiberPreconditionFailed           ErrorCode = "FIBER_PRECONDITION_FAILED"
	ErrorCodeFiberRequestEntityTooLarge        ErrorCode = "FIBER_REQUEST_ENTITY_TOO_LARGE"
	ErrorCodeFiberRequestURITooLong            ErrorCode = "FIBER_REQUEST_URI_TOO_LONG"
	ErrorCodeFiberUnsupportedMediaType         ErrorCode = "FIBER_UNSUPPORTED_MEDIA_TYPE"
	ErrorCodeFiberRequestedRangeNotSatisfiable ErrorCode = "FIBER_REQUESTED_RANGE_NOT_SATISFIABLE"
	ErrorCodeFiberExpectationFailed            ErrorCode = "FIBER_EXPECTATION_FAILED"
	ErrorCodeFiberTeapot                       ErrorCode = "FIBER_TEAPOT"
	ErrorCodeFiberMisdirectedRequest           ErrorCode = "FIBER_MISDIRECTED_REQUEST"
	ErrorCodeFiberUnprocessableEntity          ErrorCode = "FIBER_UNPROCESSABLE_ENTITY"
	ErrorCodeFiberLocked                       ErrorCode = "FIBER_LOCKED"
	ErrorCodeFiberFailedDependency             ErrorCode = "FIBER_FAILED_DEPENDENCY"
	ErrorCodeFiberTooEarly                     ErrorCode = "FIBER_TOO_EARLY"
	ErrorCodeFiberUpgradeRequired              ErrorCode = "FIBER_UPGRADE_REQUIRED"
	ErrorCodeFiberPreconditionRequired         ErrorCode = "FIBER_PRECONDITION_REQUIRED"
	ErrorCodeFiberTooManyRequests              ErrorCode = "FIBER_TOO_MANY_REQUESTS"
	ErrorCodeFiberRequestHeaderFieldsTooLarge  ErrorCode = "FIBER_REQUEST_HEADER_FIELDS_TOO_LARGE"
	ErrorCodeFiberUnavailableForLegalReasons   ErrorCode = "FIBER_UNAVAILABLE_FOR_LEGAL_REASONS"

	// Fiber/HTTP specific errors (5xx Server Errors)
	ErrorCodeFiberInternalServerError           ErrorCode = "FIBER_INTERNAL_SERVER_ERROR"
	ErrorCodeFiberNotImplemented                ErrorCode = "FIBER_NOT_IMPLEMENTED"
	ErrorCodeFiberBadGateway                    ErrorCode = "FIBER_BAD_GATEWAY"
	ErrorCodeFiberServiceUnavailable            ErrorCode = "FIBER_SERVICE_UNAVAILABLE"
	ErrorCodeFiberGatewayTimeout                ErrorCode = "FIBER_GATEWAY_TIMEOUT"
	ErrorCodeFiberHTTPVersionNotSupported       ErrorCode = "FIBER_HTTP_VERSION_NOT_SUPPORTED"
	ErrorCodeFiberVariantAlsoNegotiates         ErrorCode = "FIBER_VARIANT_ALSO_NEGOTIATES"
	ErrorCodeFiberInsufficientStorage           ErrorCode = "FIBER_INSUFFICIENT_STORAGE"
	ErrorCodeFiberLoopDetected                  ErrorCode = "FIBER_LOOP_DETECTED"
	ErrorCodeFiberNotExtended                   ErrorCode = "FIBER_NOT_EXTENDED"
	ErrorCodeFiberNetworkAuthenticationRequired ErrorCode = "FIBER_NETWORK_AUTHENTICATION_REQUIRED"
)

// MapFiberError maps Fiber HTTP status codes to semantic error codes and types
func MapFiberError(statusCode int, message string) (ErrorType, ErrorCode) {
	switch statusCode {
	// 4xx Client Errors
	case http.StatusBadRequest:
		return ErrorTypeBadRequest, ErrorCodeFiberBadRequest
	case http.StatusUnauthorized:
		return ErrorTypeUnauthorized, ErrorCodeFiberUnauthorized
	case http.StatusPaymentRequired:
		return ErrorTypeBadRequest, ErrorCodeFiberPaymentRequired
	case http.StatusForbidden:
		return ErrorTypeForbidden, ErrorCodeFiberForbidden
	case http.StatusNotFound:
		return ErrorTypeNotFound, ErrorCodeFiberNotFound
	case http.StatusMethodNotAllowed:
		return ErrorTypeBadRequest, ErrorCodeFiberMethodNotAllowed
	case http.StatusNotAcceptable:
		return ErrorTypeBadRequest, ErrorCodeFiberNotAcceptable
	case http.StatusProxyAuthRequired:
		return ErrorTypeUnauthorized, ErrorCodeFiberProxyAuthRequired
	case http.StatusRequestTimeout:
		return ErrorTypeBadRequest, ErrorCodeFiberRequestTimeout
	case http.StatusConflict:
		return ErrorTypeConflict, ErrorCodeFiberConflict
	case http.StatusGone:
		return ErrorTypeNotFound, ErrorCodeFiberGone
	case http.StatusLengthRequired:
		return ErrorTypeBadRequest, ErrorCodeFiberLengthRequired
	case http.StatusPreconditionFailed:
		return ErrorTypeBadRequest, ErrorCodeFiberPreconditionFailed
	case http.StatusRequestEntityTooLarge:
		return ErrorTypeBadRequest, ErrorCodeFiberRequestEntityTooLarge
	case http.StatusRequestURITooLong:
		return ErrorTypeBadRequest, ErrorCodeFiberRequestURITooLong
	case http.StatusUnsupportedMediaType:
		return ErrorTypeBadRequest, ErrorCodeFiberUnsupportedMediaType
	case http.StatusRequestedRangeNotSatisfiable:
		return ErrorTypeBadRequest, ErrorCodeFiberRequestedRangeNotSatisfiable
	case http.StatusExpectationFailed:
		return ErrorTypeBadRequest, ErrorCodeFiberExpectationFailed
	case http.StatusTeapot:
		return ErrorTypeBadRequest, ErrorCodeFiberTeapot
	case http.StatusMisdirectedRequest:
		return ErrorTypeBadRequest, ErrorCodeFiberMisdirectedRequest
	case http.StatusUnprocessableEntity:
		return ErrorTypeValidation, ErrorCodeFiberUnprocessableEntity
	case http.StatusLocked:
		return ErrorTypeForbidden, ErrorCodeFiberLocked
	case http.StatusFailedDependency:
		return ErrorTypeBadRequest, ErrorCodeFiberFailedDependency
	case http.StatusTooEarly:
		return ErrorTypeBadRequest, ErrorCodeFiberTooEarly
	case http.StatusUpgradeRequired:
		return ErrorTypeBadRequest, ErrorCodeFiberUpgradeRequired
	case http.StatusPreconditionRequired:
		return ErrorTypeBadRequest, ErrorCodeFiberPreconditionRequired
	case http.StatusTooManyRequests:
		return ErrorTypeRateLimit, ErrorCodeFiberTooManyRequests
	case http.StatusRequestHeaderFieldsTooLarge:
		return ErrorTypeBadRequest, ErrorCodeFiberRequestHeaderFieldsTooLarge
	case http.StatusUnavailableForLegalReasons:
		return ErrorTypeForbidden, ErrorCodeFiberUnavailableForLegalReasons

	// 5xx Server Errors
	case http.StatusInternalServerError:
		return ErrorTypeInternal, ErrorCodeFiberInternalServerError
	case http.StatusNotImplemented:
		return ErrorTypeInternal, ErrorCodeFiberNotImplemented
	case http.StatusBadGateway:
		return ErrorTypeServiceUnavailable, ErrorCodeFiberBadGateway
	case http.StatusServiceUnavailable:
		return ErrorTypeServiceUnavailable, ErrorCodeFiberServiceUnavailable
	case http.StatusGatewayTimeout:
		return ErrorTypeServiceUnavailable, ErrorCodeFiberGatewayTimeout
	case http.StatusHTTPVersionNotSupported:
		return ErrorTypeBadRequest, ErrorCodeFiberHTTPVersionNotSupported
	case http.StatusVariantAlsoNegotiates:
		return ErrorTypeInternal, ErrorCodeFiberVariantAlsoNegotiates
	case http.StatusInsufficientStorage:
		return ErrorTypeServiceUnavailable, ErrorCodeFiberInsufficientStorage
	case http.StatusLoopDetected:
		return ErrorTypeInternal, ErrorCodeFiberLoopDetected
	case http.StatusNotExtended:
		return ErrorTypeBadRequest, ErrorCodeFiberNotExtended
	case http.StatusNetworkAuthenticationRequired:
		return ErrorTypeUnauthorized, ErrorCodeFiberNetworkAuthenticationRequired

	default:
		// For unknown status codes, default to internal server error
		return ErrorTypeInternal, ErrorCodeFiberInternalServerError
	}
}

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
	case ErrorTypeBadRequest:
		return http.StatusBadRequest
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeValidation:
		return http.StatusUnprocessableEntity
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
