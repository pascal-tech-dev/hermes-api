package dto

import (
	"hermes-api/internal/validation"

	"github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r *RegisterRequest) Validate() error {
	validate := validator.New()
	return validation.MapValidationErrors(validate.Struct(r))
}

func (r *LoginRequest) Validate() error {
	validate := validator.New()
	return validation.MapValidationErrors(validate.Struct(r))
}
