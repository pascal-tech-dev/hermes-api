package rest

import (
	"hermes-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// GetUsers retrieves all users
func GetUsers(c *fiber.Ctx) error {
	options := response.SuccessResponse([]string{}, "Users retrieved successfully")
	return response.ApiResponse(c, options)
}

// CreateUser creates a new user
func CreateUser(c *fiber.Ctx) error {
	options := response.CreatedResponse(nil, "User created successfully")
	return response.ApiResponse(c, options)
}

// GetUserByID retrieves a user by ID
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	options := response.SuccessResponse(fiber.Map{"id": id}, "User retrieved successfully")
	return response.ApiResponse(c, options)
}

// UpdateUser updates a user
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	options := response.SuccessResponse(fiber.Map{"id": id}, "User updated successfully")
	return response.ApiResponse(c, options)
}

// DeleteUser deletes a user
func DeleteUser(c *fiber.Ctx) error {
	_ = c.Params("id") // ID parameter for future use
	options := response.NoContentResponse("User deleted successfully")
	return response.ApiResponse(c, options)
}
