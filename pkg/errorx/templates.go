package errorx

import "fmt"

// ErrorTemplates contains predefined error messages
var ErrorTemplates = map[ErrorCode]string{
	// App errors
	ErrorCodeAppNotFound:      "Application with ID '%s' not found",
	ErrorCodeAppAlreadyExists: "Application with name '%s' already exists",
	ErrorCodeInvalidAPIKey:    "Invalid API key provided",
	ErrorCodeAppInactive:      "Application '%s' is inactive",

	// Notification errors
	ErrorCodeNotificationNotFound:      "Notification with ID '%s' not found",
	ErrorCodeInvalidNotificationType:   "Invalid notification type '%s'. Allowed types: email, sms, push, webhook, slack",
	ErrorCodeInvalidRecipient:          "Invalid recipient format: '%s'",
	ErrorCodeNotificationQuotaExceeded: "Notification quota exceeded for app '%s'",

	// Validation errors
	ErrorCodeRequiredField: "Field '%s' is required",
	ErrorCodeInvalidFormat: "Field '%s' has invalid format",
	ErrorCodeInvalidValue:  "Field '%s' has invalid value: '%s'",
	ErrorCodeFieldTooLong:  "Field '%s' is too long. Maximum length: %d",
	ErrorCodeFieldTooShort: "Field '%s' is too short. Minimum length: %d",

	// System errors
	ErrorCodeDatabaseError:        "Database operation failed: %s",
	ErrorCodeRedisError:           "Redis operation failed: %s",
	ErrorCodeExternalServiceError: "External service error: %s",
}

// NewWithTemplate creates a new error using a template
func NewWithTemplate(errorType ErrorType, code ErrorCode, args ...interface{}) *AppError {
	template, exists := ErrorTemplates[code]
	if !exists {
		template = "Unknown error occurred"
	}

	message := fmt.Sprintf(template, args...)
	return New(errorType, code, message)
}

// NewValidationError creates a validation error
func NewValidationError(field, value string) *AppError {
	return NewWithTemplate(ErrorTypeValidation, ErrorCodeInvalidValue, field, value)
}

// NewRequiredFieldError creates a required field error
func NewRequiredFieldError(field string) *AppError {
	return NewWithTemplate(ErrorTypeValidation, ErrorCodeRequiredField, field)
}

// NewAppNotFoundError creates an app not found error
func NewAppNotFoundError(appID string) *AppError {
	return NewWithTemplate(ErrorTypeNotFound, ErrorCodeAppNotFound, appID)
}

// NewInvalidAPIKeyError creates an invalid API key error
func NewInvalidAPIKeyError() *AppError {
	return New(ErrorTypeUnauthorized, ErrorCodeInvalidAPIKey, "Invalid API key provided")
}
