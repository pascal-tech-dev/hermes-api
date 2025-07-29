package rest

import (
	"hermes-api/api/rest/controller"
	"hermes-api/internal/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all API routes
func SetupRoutes(api fiber.Router, serviceManager service.ServiceManager, authMiddleware fiber.Handler) {
	controllerManager := controller.NewControllerManager(serviceManager)
	setupV1Routes(api, controllerManager, authMiddleware)
}

// setupV1Routes configures API v1 routes
func setupV1Routes(api fiber.Router, controllerManager *controller.ControllerManager, authMiddleware fiber.Handler) {
	// Auth routes (public)
	setupAuthRoutes(api, controllerManager.Auth())

	// Users routes (protected)
	setupUserRoutes(api, controllerManager.User(), authMiddleware)
}

// setupAuthRoutes configures authentication-related routes
func setupAuthRoutes(api fiber.Router, authController *controller.AuthController) {
	auth := api.Group("/auth")

	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Get("/me", authController.Me) // This will need auth middleware
}

// setupUserRoutes configures user-related routes
func setupUserRoutes(api fiber.Router, userController *controller.UserController, authMiddleware fiber.Handler) {
	users := api.Group("/users")

	// Apply auth middleware to all user routes
	users.Use(authMiddleware)

	users.Get("/", userController.GetUsers)
	users.Post("/", userController.CreateUser)
	users.Get("/:id", userController.GetUserByID)
	users.Put("/:id", userController.UpdateUser)
	users.Delete("/:id", userController.DeleteUser)
}

// Health check handler
// To use this function, register it as a route handler in your Fiber app, for example:
//
//	app.Get("/health", healthCheck)
//
// or if using a router group:
//
//	api.Get("/health", healthCheck)
func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"service":   "hermes-api",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	})
}
