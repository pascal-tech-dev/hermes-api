package validation

import (
	"hermes-api/pkg/errorx"

	"github.com/go-playground/validator/v10"
)

func MapValidationErrors(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		fieldErrors := make(map[string]any)
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				fieldErrors[e.Field()] = errorx.ErrorCodeRequiredField
			case "email, url, uuid":
				fieldErrors[e.Field()] = errorx.ErrorCodeInvalidFormat
			case "min":
				fieldErrors[e.Field()] = errorx.ErrorCodeFieldTooShort
			case "max":
				fieldErrors[e.Field()] = errorx.ErrorCodeFieldTooLong
			case "len":
				fieldErrors[e.Field()] = errorx.ErrorCodeFieldLengthMismatch
			case "eq":
				fieldErrors[e.Field()] = errorx.ErrorCodeFieldNotEqual
			case "ne":
				fieldErrors[e.Field()] = errorx.ErrorCodeFieldNotEqual
			default:
				fieldErrors[e.Field()] = errorx.ErrorCodeInvalidValue
			}
		}
		appErr := errorx.New(errorx.ErrorTypeValidation, errorx.ErrorCodeInvalidValue, "Invalid request body")
		return appErr.WithDetails(fieldErrors)
	}
	return nil
}
