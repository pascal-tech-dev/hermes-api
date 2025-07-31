package controller

import (
	"hermes-api/internal/dto"
	"hermes-api/internal/service"
	"hermes-api/pkg/context"
	"hermes-api/pkg/errorx"
	"hermes-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// AuthController handles HTTP requests for authentication operations
type AuthController struct {
	authService service.AuthService
}

// NewAuthController creates a new auth controller
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register creates a new user account
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		appErr := errorx.New(errorx.ErrorTypeBadRequest, errorx.ErrorCodeInvalidFormat, err.Error())
		return appErr // return the error to the middleware
	}

	if err := req.Validate(); err != nil {
		return err // return the error to the middleware
	}

	// Create a new context for the service
	serviceCtx, cancel := context.New(ctx).WithDefaultTimeout().Build()
	defer cancel()

	user, err := c.authService.Register(serviceCtx, req.Email, req.Username, req.Password, req.FirstName, req.LastName)
	if err != nil {
		return err // return the error to the middleware
	}

	return response.CreatedResponse(user, "User registered successfully").
		WithRequestID(ctx.Locals("X-Request-ID").(string)).
		Send(ctx)
}

// Login authenticates a user and returns a JWT token
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		appErr := errorx.New(errorx.ErrorTypeBadRequest, errorx.ErrorCodeInvalidFormat, err.Error())
		return appErr
	}

	// Create a new context for the service
	serviceCtx, cancel := context.New(ctx).WithDefaultTimeout().Build()
	defer cancel()

	token, user, err := c.authService.Login(serviceCtx, req.Email, req.Password)
	if err != nil {
		return err
	}

	responseData := fiber.Map{
		"token": token,
		"user":  user,
	}

	return response.SuccessResponse(responseData, "Login successful").
		WithRequestID(ctx.Locals("X-Request-ID").(string)).
		Send(ctx)
}

// Me returns the current authenticated user's information
func (c *AuthController) Me(ctx *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	user := ctx.Locals("user")
	if user == nil {
		appErr := errorx.New(errorx.ErrorTypeUnauthorized, errorx.ErrorCodeFiberUnauthorized, "User not authenticated")
		// options := response.ErrorResponse(appErr, "User not authenticated")
		// return response.ApiResponse(ctx, options)
		return appErr
	}

	// options := response.SuccessResponse(user, "User information retrieved successfully")
	// return response.ApiResponse(ctx, options)
	return nil
}
