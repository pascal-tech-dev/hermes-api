package rest

import (
	"hermes-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// Login handles user login
func Login(c *fiber.Ctx) error {
	options := response.SuccessResponse(fiber.Map{"token": "jwt-token-here"}, "Login successful")
	return response.ApiResponse(c, options)
}

// Register handles user registration
func Register(c *fiber.Ctx) error {
	options := response.CreatedResponse(fiber.Map{"user_id": "new-user-id"}, "User registered successfully")
	return response.ApiResponse(c, options)
}

// Logout handles user logout
func Logout(c *fiber.Ctx) error {
	options := response.SuccessResponse(nil, "Logout successful")
	return response.ApiResponse(c, options)
}

// RefreshToken handles token refresh
func RefreshToken(c *fiber.Ctx) error {
	options := response.SuccessResponse(fiber.Map{"token": "new-jwt-token"}, "Token refreshed successfully")
	return response.ApiResponse(c, options)
}
