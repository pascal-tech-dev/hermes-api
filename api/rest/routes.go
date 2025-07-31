package rest

import (
	"hermes-api/api/rest/controller"
	"hermes-api/internal/service"

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

	// Applications routes (protected)
	setupApplicationRoutes(api, controllerManager.Application(), authMiddleware)
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

// setupApplicationRoutes configures application-related routes
func setupApplicationRoutes(api fiber.Router, applicationController *controller.ApplicationController, authMiddleware fiber.Handler) {
	applications := api.Group("/applications")

	// Apply auth middleware to all application routes
	applications.Use(authMiddleware)

	applications.Post("/", applicationController.CreateApplication)
}
