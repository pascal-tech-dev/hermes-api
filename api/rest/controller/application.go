package controller

import (
	"hermes-api/internal/dto"
	"hermes-api/internal/model"
	"hermes-api/internal/service"
	"hermes-api/pkg/context"
	"hermes-api/pkg/errorx"
	"hermes-api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

// GetApplications gets all applications
func (c *ApplicationController) GetApplications(ctx *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	user := ctx.Locals("user").(*model.User)
	if user == nil {
		appErr := errorx.New(errorx.ErrorTypeUnauthorized, errorx.ErrorCodeFiberUnauthorized, "User not found")
		return appErr // return the error to the middleware
	}

	// Create a new context for the service
	serviceCtx, cancel := context.New(ctx).WithDefaultTimeout().Build()
	defer cancel()

	applications, err := c.applicationService.GetApplications(serviceCtx, user.ID)
	if err != nil {
		return err
	}

	return response.SuccessResponse(applications, "Applications retrieved successfully").
		WithRequestID(ctx.Locals("X-Request-ID").(string)).
		Send(ctx)
}

// GetApplicationsWithPagination gets applications with pagination support
func (c *ApplicationController) GetApplicationsWithPagination(ctx *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	user := ctx.Locals("user").(*model.User)
	if user == nil {
		appErr := errorx.New(errorx.ErrorTypeUnauthorized, errorx.ErrorCodeFiberUnauthorized, "User not found")
		return appErr // return the error to the middleware
	}

	// Get pagination parameters from query string
	limit := ctx.QueryInt("limit", 10)  // Default limit: 10
	offset := ctx.QueryInt("offset", 0) // Default offset: 0

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 10 // Reset to default if invalid
	}
	if offset < 0 {
		offset = 0 // Reset to default if invalid
	}

	// Create a new context for the service
	serviceCtx, cancel := context.New(ctx).WithDefaultTimeout().Build()
	defer cancel()

	applications, err := c.applicationService.GetApplicationsWithPagination(serviceCtx, user.ID, limit, offset)
	if err != nil {
		return err
	}

	// Create pagination response
	paginationResponse := map[string]any{
		"applications": applications,
		"pagination": map[string]any{
			"limit":  limit,
			"offset": offset,
			"count":  len(applications),
		},
	}

	return response.SuccessResponse(paginationResponse, "Applications retrieved successfully").
		WithRequestID(ctx.Locals("X-Request-ID").(string)).
		Send(ctx)
}

// GetApplicationByID gets an application by its ID
func (c *ApplicationController) GetApplicationByID(ctx *fiber.Ctx) error {
	appIDStr := ctx.Params("id")

	// Parse UUID
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		appErr := errorx.New(errorx.ErrorTypeBadRequest, errorx.ErrorCodeInvalidFormat, "Invalid application ID format")
		return appErr
	}

	// Create a new context for the service
	serviceCtx, cancel := context.New(ctx).WithDefaultTimeout().Build()
	defer cancel()

	application, err := c.applicationService.GetApplicationByID(serviceCtx, appID)
	if err != nil {
		return err
	}

	return response.SuccessResponse(application, "Application retrieved successfully").
		WithRequestID(ctx.Locals("X-Request-ID").(string)).
		Send(ctx)
}

// UpdateApplication updates an application
func (c *ApplicationController) UpdateApplication(ctx *fiber.Ctx) error {
	appIDStr := ctx.Params("id")

	// Parse UUID
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		appErr := errorx.New(errorx.ErrorTypeBadRequest, errorx.ErrorCodeInvalidFormat, "Invalid application ID format")
		return appErr
	}

	// Parse request body
	var req dto.UpdateApplicationRequest
	if err := ctx.BodyParser(&req); err != nil {
		appErr := errorx.New(errorx.ErrorTypeBadRequest, errorx.ErrorCodeInvalidFormat, err.Error())
		return appErr
	}

	if err := req.Validate(); err != nil {
		return err
	}

	// Create a new context for the service
	serviceCtx, cancel := context.New(ctx).WithDefaultTimeout().Build()
	defer cancel()

	// Get current application
	application, err := c.applicationService.GetApplicationByID(serviceCtx, appID)
	if err != nil {
		return err
	}

	// Update fields
	if req.Name != "" {
		application.Name = req.Name
	}
	if req.Description != "" {
		application.Description = req.Description
	}
	if req.Status != "" {
		application.Status = model.ApplicationStatus(req.Status)
	}

	err = c.applicationService.UpdateApplication(serviceCtx, application)
	if err != nil {
		return err
	}

	return response.SuccessResponse(application, "Application updated successfully").
		WithRequestID(ctx.Locals("X-Request-ID").(string)).
		Send(ctx)
}

// DeleteApplication deletes an application
func (c *ApplicationController) DeleteApplication(ctx *fiber.Ctx) error {
	appIDStr := ctx.Params("id")

	// Parse UUID
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		appErr := errorx.New(errorx.ErrorTypeBadRequest, errorx.ErrorCodeInvalidFormat, "Invalid application ID format")
		return appErr
	}

	// Create a new context for the service
	serviceCtx, cancel := context.New(ctx).WithDefaultTimeout().Build()
	defer cancel()

	err = c.applicationService.DeleteApplication(serviceCtx, appID)
	if err != nil {
		return err
	}

	return response.SuccessResponse(nil, "Application deleted successfully").
		WithRequestID(ctx.Locals("X-Request-ID").(string)).
		Send(ctx)
}
