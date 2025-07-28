// api/rest/controller/manager.go
package controller

import (
	"hermes-api/internal/service"
)

// ControllerManager manages all controllers and their dependencies
type ControllerManager struct {
	userController *UserController
	// Add other controllers as needed:
	// productController *ProductController
	// orderController   *OrderController
	// authController    *AuthController
}

// NewControllerManager creates a new controller manager
func NewControllerManager(serviceManager service.ServiceManager) *ControllerManager {
	return &ControllerManager{
		userController: NewUserController(serviceManager.User()),
		// Initialize other controllers:
		// productController: NewProductController(serviceManager.Product()),
		// orderController:   NewOrderController(serviceManager.Order()),
		// authController:    NewAuthController(serviceManager.Auth()),
	}
}

// User returns the user controller
func (cm *ControllerManager) User() *UserController {
	return cm.userController
}

// Product returns the product controller (future)
// func (cm *ControllerManager) Product() *ProductController {
// 	return cm.productController
// }

// Order returns the order controller (future)
// func (cm *ControllerManager) Order() *OrderController {
// 	return cm.orderController
// }

// Auth returns the auth controller (future)
// func (cm *ControllerManager) Auth() *AuthController {
// 	return cm.authController
// }
