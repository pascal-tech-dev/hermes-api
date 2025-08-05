package controller

import (
	"hermes-api/internal/dto"
	"hermes-api/internal/model"
	"hermes-api/internal/service"
	"hermes-api/pkg/context"
	"hermes-api/pkg/errorx"
	"hermes-api/pkg/response"

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

	// Get user from context (set by auth middleware)
	user := ctx.Locals("user").(*model.User)
	if user == nil {
		appErr := errorx.New(errorx.ErrorTypeUnauthorized, errorx.ErrorCodeFiberUnauthorized, "User not found")
		return appErr // return the error to the middleware
	}

	// Create a new context for the service
	serviceCtx, cancel := context.New(ctx).WithDefaultTimeout().Build()
	defer cancel()

	application, err := c.applicationService.CreateApplication(serviceCtx, user.ID, req)
	if err != nil {
		return err
	}

	return response.CreatedResponse(application, "Application created successfully").
		WithRequestID(ctx.Locals("X-Request-ID").(string)).
		Send(ctx)
}
