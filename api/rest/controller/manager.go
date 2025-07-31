// api/rest/controller/manager.go
package controller

import (
	"hermes-api/internal/service"
)

// ControllerManager manages all controllers and their dependencies
type ControllerManager struct {
	userController        *UserController
	authController        *AuthController
	applicationController *ApplicationController
	// Add other controllers as needed:
	// productController *ProductController
	// orderController   *OrderController
}

// NewControllerManager creates a new controller manager
func NewControllerManager(serviceManager service.ServiceManager) *ControllerManager {
	return &ControllerManager{
		userController:        NewUserController(serviceManager.User()),
		authController:        NewAuthController(serviceManager.Auth()),
		applicationController: NewApplicationController(serviceManager.Application()),
	}
}

// User returns the user controller
func (cm *ControllerManager) User() *UserController {
	return cm.userController
}

// Auth returns the auth controller
func (cm *ControllerManager) Auth() *AuthController {
	return cm.authController
}

// Application returns the application controller
func (cm *ControllerManager) Application() *ApplicationController {
	return cm.applicationController
}
