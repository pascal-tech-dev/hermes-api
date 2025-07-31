package controller

import (
	"hermes-api/internal/dto"
	"hermes-api/internal/service"
	"hermes-api/pkg/errorx"

	"github.com/gofiber/fiber/v2"
)

type ApplicationController struct {
	applicationService service.ApplicationService
}

func NewApplicationController(applicationService service.ApplicationService) *ApplicationController {
	return &ApplicationController{
		applicationService: applicationService,
	}
}

// CreateApplication creates a new application
func (c *ApplicationController) CreateApplication(ctx *fiber.Ctx) error {
	var req dto.CreateApplicationRequest
	if err := ctx.BodyParser(&req); err != nil {
		appErr := errorx.New(errorx.ErrorTypeBadRequest, errorx.ErrorCodeInvalidFormat, err.Error())
		return appErr // return the error to the middleware
	}

	if err := req.Validate(); err != nil {
		return err // return the error to the middleware
	}

	return nil
}
