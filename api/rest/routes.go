package rest

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all API routes
func SetupRoutes(api fiber.Router) {
	setupV1Routes(api)
}

// setupV1Routes configures API v1 routes
func setupV1Routes(api fiber.Router) {
	// Notifications routes
	setupNotificationRoutes(api)

	// Users routes
	setupUserRoutes(api)

	// Auth routes
	setupAuthRoutes(api)
}

// setupNotificationRoutes configures notification-related routes
func setupNotificationRoutes(api fiber.Router) {
	notifications := api.Group("/notifications")

	notifications.Get("/", GetNotifications)
	notifications.Post("/", CreateNotification)
	notifications.Get("/:id", GetNotificationByID)
	notifications.Put("/:id", UpdateNotification)
	notifications.Delete("/:id", DeleteNotification)
	notifications.Post("/:id/read", MarkNotificationAsRead)
}

// setupUserRoutes configures user-related routes
func setupUserRoutes(api fiber.Router) {
	users := api.Group("/users")

	users.Get("/", GetUsers)
	users.Post("/", CreateUser)
	users.Get("/:id", GetUserByID)
	users.Put("/:id", UpdateUser)
	users.Delete("/:id", DeleteUser)
}

// setupAuthRoutes configures authentication-related routes
func setupAuthRoutes(api fiber.Router) {
	auth := api.Group("/auth")

	auth.Post("/login", Login)
	auth.Post("/register", Register)
	auth.Post("/logout", Logout)
	auth.Post("/refresh", RefreshToken)
}

// Health check handler
func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"service":   "hermes-api",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	})
}
