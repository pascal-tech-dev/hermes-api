package controller

import (
	"hermes-api/internal/model"
	"hermes-api/internal/service"
	"hermes-api/pkg/logger"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// UserController handles HTTP requests for user operations
type UserController struct {
	userService service.UserService
}

// NewUserController creates a new user controller
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetUsers retrieves all users with pagination
func (c *UserController) GetUsers(ctx *fiber.Ctx) error {
	// Parse pagination parameters
	limit := 10 // default limit
	offset := 0 // default offset

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := ctx.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get users from service
	_, err := c.userService.ListUsers(ctx.Context(), limit, offset)
	if err != nil {
		logger.Error("Failed to retrieve users", err)
		// appErr := errors.New(errors.ErrorTypeInternal, errors.ErrorCodeDatabaseError, "Failed to retrieve users")
		// options := response.ErrorResponse(appErr, "Failed to retrieve users")
		// return response.ApiResponse(ctx, options)
		return err
	}

	// options := response.SuccessResponse(fiber.Map{
	// 	"users": users,
	// 	"pagination": fiber.Map{
	// 		"limit":  limit,
	// 		"offset": offset,
	// 		"total":  0,
	// 	},
	// }, "Users retrieved successfully")
	// return response.ApiResponse(ctx, options)
	return nil
}

// CreateUser creates a new user
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var user model.User
	if err := ctx.BodyParser(&user); err != nil {
		logger.Error("Failed to parse user data", err)
		// appErr := errors.New(errors.ErrorTypeValidation, errors.ErrorCodeInvalidFormat, "Invalid request body")
		// options := response.ErrorResponse(appErr, "Invalid request body")
		// return response.ApiResponse(ctx, options)
		return err
	}

	// Create user using service
	if err := c.userService.CreateUser(ctx.Context(), &user); err != nil {
		logger.Error("Failed to create user", err)
		// appErr := errors.New(errors.ErrorTypeInternal, errors.ErrorCodeDatabaseError, "Failed to create user")
		// options := response.ErrorResponse(appErr, "Failed to create user")
		// return response.ApiResponse(ctx, options)
		return err
	}

	// options := response.CreatedResponse(user, "User created successfully")
	// return response.ApiResponse(ctx, options)
	return nil
}

// GetUserByID retrieves a user by ID
func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	// idStr := ctx.Params("id")
	// id, err := strconv.ParseUint(idStr, 10, 32)
	// if err != nil {
	// 	// appErr := errors.New(errors.ErrorTypeValidation, errors.ErrorCodeInvalidValue, "Invalid user ID")
	// 	// options := response.ErrorResponse(appErr, "Invalid user ID")
	// 	// return response.ApiResponse(ctx, options)
	// 	return err
	// }

	// _, err = c.userService.GetUserByID(ctx.Context(), uint(id))
	// if err != nil {
	// 	logger.Error("Failed to retrieve user", err, zap.String("user_id", idStr))
	// 	// appErr := errors.New(errors.ErrorTypeNotFound, errors.ErrorCodeAppNotFound, "User not found")
	// 	// options := response.ErrorResponse(appErr, "User not found")
	// 	// return response.ApiResponse(ctx, options)
	// 	return err
	// }

	// options := response.SuccessResponse(user, "User retrieved successfully")
	// return response.ApiResponse(ctx, options)
	return nil
}

// UpdateUser updates a user
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	// idStr := ctx.Params("id")
	// id, err := strconv.ParseUint(idStr, 10, 32)
	// if err != nil {
	// 	appErr := errors.New(errors.ErrorTypeValidation, errors.ErrorCodeInvalidValue, "Invalid user ID")
	// 	options := response.ErrorResponse(appErr, "Invalid user ID")
	// 	return response.ApiResponse(ctx, options)
	// }

	// var user model.User
	// if err := ctx.BodyParser(&user); err != nil {
	// 	logger.Error("Failed to parse user data", err)
	// 	appErr := errors.New(errors.ErrorTypeValidation, errors.ErrorCodeInvalidFormat, "Invalid request body")
	// 	options := response.ErrorResponse(appErr, "Invalid request body")
	// 	return response.ApiResponse(ctx, options)
	// }

	// user.ID = uint(id)

	// // Update user using service
	// if err := c.userService.UpdateUser(ctx.Context(), &user); err != nil {
	// 	logger.Error("Failed to update user", err, zap.String("user_id", idStr))
	// 	appErr := errors.New(errors.ErrorTypeInternal, errors.ErrorCodeDatabaseError, "Failed to update user")
	// 	options := response.ErrorResponse(appErr, "Failed to update user")
	// 	return response.ApiResponse(ctx, options)
	// }

	// options := response.SuccessResponse(user, "User updated successfully")
	// return response.ApiResponse(ctx, options)
	return nil
}

// DeleteUser deletes a user
func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	// idStr := ctx.Params("id")
	// id, err := strconv.ParseUint(idStr, 10, 32)
	// if err != nil {
	// 	appErr := errors.New(errors.ErrorTypeValidation, errors.ErrorCodeInvalidValue, "Invalid user ID")
	// 	options := response.ErrorResponse(appErr, "Invalid user ID")
	// 	return response.ApiResponse(ctx, options)
	// }

	// // Delete user using service
	// if err := c.userService.DeleteUser(ctx.Context(), uint(id)); err != nil {
	// 	logger.Error("Failed to delete user", err, zap.String("user_id", idStr))
	// 	appErr := errors.New(errors.ErrorTypeInternal, errors.ErrorCodeDatabaseError, "Failed to delete user")
	// 	options := response.ErrorResponse(appErr, "Failed to delete user")
	// 	return response.ApiResponse(ctx, options)
	// }

	// options := response.NoContentResponse("User deleted successfully")
	// return response.ApiResponse(ctx, options)
	return nil
}

// GetUserByEmail retrieves a user by email
func (c *UserController) GetUserByEmail(ctx *fiber.Ctx) error {
	// email := ctx.Query("email")
	// if email == "" {
	// 	appErr := errors.New(errors.ErrorTypeValidation, errors.ErrorCodeInvalidValue, "Email parameter is required")
	// 	options := response.ErrorResponse(appErr, "Email parameter is required")
	// 	return response.ApiResponse(ctx, options)
	// }

	// user, err := c.userService.GetUserByEmail(ctx.Context(), email)
	// if err != nil {
	// 	logger.Error("Failed to retrieve user by email", err, zap.String("email", email))
	// 	appErr := errors.New(errors.ErrorTypeNotFound, errors.ErrorCodeAppNotFound, "User not found")
	// 	options := response.ErrorResponse(appErr, "User not found")
	// 	return response.ApiResponse(ctx, options)
	// }

	// options := response.SuccessResponse(user, "User retrieved successfully")
	// return response.ApiResponse(ctx, options)
	return nil
}

// ActivateUser activates a user account
func (c *UserController) ActivateUser(ctx *fiber.Ctx) error {
	// idStr := ctx.Params("id")
	// id, err := strconv.ParseUint(idStr, 10, 32)
	// if err != nil {
	// 	appErr := errors.New(errors.ErrorTypeValidation, errors.ErrorCodeInvalidValue, "Invalid user ID")
	// 	options := response.ErrorResponse(appErr, "Invalid user ID")
	// 	return response.ApiResponse(ctx, options)
	// }

	// if err := c.userService.ActivateUser(ctx.Context(), uint(id)); err != nil {
	// 	logger.Error("Failed to activate user", err, zap.String("user_id", idStr))
	// 	appErr := errors.New(errors.ErrorTypeInternal, errors.ErrorCodeDatabaseError, "Failed to activate user")
	// 	options := response.ErrorResponse(appErr, "Failed to activate user")
	// 	return response.ApiResponse(ctx, options)
	// }

	// options := response.SuccessResponse(nil, "User activated successfully")
	// return response.ApiResponse(ctx, options)
	return nil
}

// DeactivateUser deactivates a user account
func (c *UserController) DeactivateUser(ctx *fiber.Ctx) error {
	// idStr := ctx.Params("id")
	// id, err := strconv.ParseUint(idStr, 10, 32)
	// if err != nil {
	// 	appErr := errors.New(errors.ErrorTypeValidation, errors.ErrorCodeInvalidValue, "Invalid user ID")
	// 	options := response.ErrorResponse(appErr, "Invalid user ID")
	// 	return response.ApiResponse(ctx, options)
	// }

	// if err := c.userService.DeactivateUser(ctx.Context(), uint(id)); err != nil {
	// 	logger.Error("Failed to deactivate user", err, zap.String("user_id", idStr))
	// 	appErr := errors.New(errors.ErrorTypeInternal, errors.ErrorCodeDatabaseError, "Failed to deactivate user")
	// 	options := response.ErrorResponse(appErr, "Failed to deactivate user")
	// 	return response.ApiResponse(ctx, options)
	// }

	// options := response.SuccessResponse(nil, "User deactivated successfully")
	// return response.ApiResponse(ctx, options)
	return nil
}
