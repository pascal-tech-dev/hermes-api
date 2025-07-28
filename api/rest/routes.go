package rest

import (
	"hermes-api/api/rest/controller"
	"hermes-api/internal/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all API routes
func SetupRoutes(api fiber.Router, serviceManager service.ServiceManager) {
	controllerManager := controller.NewControllerManager(serviceManager)
	setupV1Routes(api, controllerManager)
}

// setupV1Routes configures API v1 routes
func setupV1Routes(api fiber.Router, controllerManager *controller.ControllerManager) {
	// Users routes
	setupUserRoutes(api, controllerManager.User())
}

// setupUserRoutes configures user-related routes
func setupUserRoutes(api fiber.Router, userController *controller.UserController) {
	users := api.Group("/users")

	users.Get("/", userController.GetUsers)
	users.Post("/", userController.CreateUser)
	users.Get("/:id", userController.GetUserByID)
	users.Put("/:id", userController.UpdateUser)
	users.Delete("/:id", userController.DeleteUser)
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
