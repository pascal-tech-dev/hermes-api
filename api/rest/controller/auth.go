package controller

import (
	"hermes-api/internal/service"
	"hermes-api/pkg/errors"
	"hermes-api/pkg/logger"
	"hermes-api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Register creates a new user account
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		logger.Error("Failed to parse register request", err)
		appErr := errors.New(errors.ErrorTypeValidation, errors.ErrorCodeInvalidFormat, "Invalid request body")
		options := response.ErrorResponse(appErr, "Invalid request body")
		return response.ApiResponse(ctx, options)
	}

	user, err := c.authService.Register(ctx.Context(), req.Email, req.Username, req.Password, req.FirstName, req.LastName)
	if err != nil {
		logger.Error("Failed to register user", err, zap.String("email", req.Email))
		appErr := errors.New(errors.ErrorTypeValidation, errors.ErrorCodeInvalidValue, err.Error())
		options := response.ErrorResponse(appErr, "Failed to register user")
		return response.ApiResponse(ctx, options)
	}

	options := response.CreatedResponse(user, "User registered successfully")
	return response.ApiResponse(ctx, options)
}

// Login authenticates a user and returns a JWT token
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		logger.Error("Failed to parse login request", err)
		appErr := errors.New(errors.ErrorTypeValidation, errors.ErrorCodeInvalidFormat, "Invalid request body")
		options := response.ErrorResponse(appErr, "Invalid request body")
		return response.ApiResponse(ctx, options)
	}

	token, user, err := c.authService.Login(ctx.Context(), req.Email, req.Password)
	if err != nil {
		logger.Error("Failed to login user", err, zap.String("email", req.Email))
		appErr := errors.New(errors.ErrorTypeValidation, errors.ErrorCodeInvalidValue, "Invalid credentials")
		options := response.ErrorResponse(appErr, "Invalid credentials")
		return response.ApiResponse(ctx, options)
	}

	responseData := fiber.Map{
		"token": token,
		"user":  user,
	}

	options := response.SuccessResponse(responseData, "Login successful")
	return response.ApiResponse(ctx, options)
}

// Me returns the current authenticated user's information
func (c *AuthController) Me(ctx *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	user := ctx.Locals("user")
	if user == nil {
		appErr := errors.New(errors.ErrorTypeUnauthorized, errors.ErrorCodeFiberUnauthorized, "User not authenticated")
		options := response.ErrorResponse(appErr, "User not authenticated")
		return response.ApiResponse(ctx, options)
	}

	options := response.SuccessResponse(user, "User information retrieved successfully")
	return response.ApiResponse(ctx, options)
}
