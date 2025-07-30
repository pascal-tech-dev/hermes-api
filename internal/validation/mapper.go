package validation

import (
	"hermes-api/pkg/errors"

	"github.com/go-playground/validator/v10"
)

func MapValidationErrors(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		fieldErrors := make(map[string]any)
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				fieldErrors[e.Field()] = errors.ErrorCodeRequiredField
			case "email, url, uuid":
				fieldErrors[e.Field()] = errors.ErrorCodeInvalidFormat
			case "min":
				fieldErrors[e.Field()] = errors.ErrorCodeFieldTooShort
			case "max":
				fieldErrors[e.Field()] = errors.ErrorCodeFieldTooLong
			case "len":
				fieldErrors[e.Field()] = errors.ErrorCodeFieldLengthMismatch
			case "eq":
				fieldErrors[e.Field()] = errors.ErrorCodeFieldNotEqual
			case "ne":
				fieldErrors[e.Field()] = errors.ErrorCodeFieldNotEqual
			default:
				fieldErrors[e.Field()] = errors.ErrorCodeInvalidValue
			}
		}
		appErr := errors.New(errors.ErrorTypeValidation, errors.ErrorCodeInvalidValue, "Invalid request body")
		return appErr.WithDetails(fieldErrors)
	}
	return nil
}
