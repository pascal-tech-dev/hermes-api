package dto

import (
	"hermes-api/internal/validation"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateApplicationRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required,min=3,max=255"`
}

type CreateApplicationResponse struct {
	ID uuid.UUID `json:"id"`
}

func (r *CreateApplicationRequest) Validate() error {
	validate := validator.New()
	return validation.MapValidationErrors(validate.Struct(r))
}

func (r *CreateApplicationResponse) Validate() error {
	validate := validator.New()
	return validation.MapValidationErrors(validate.Struct(r))
}
