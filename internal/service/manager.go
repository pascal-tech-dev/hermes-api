// internal/service/manager.go
package service

import (
	"hermes-api/internal/repository"
)

// ServiceManager manages all services
type ServiceManager interface {
	User() UserService
	// Add other services here as needed
	// Product() ProductService
	// Order() OrderService
	// Notification() NotificationService
}

// serviceManager implements ServiceManager
type serviceManager struct {
	userService UserService
	// Add other services here
	// productService ProductService
	// orderService OrderService
	// notificationService NotificationService
}

// NewServiceManager creates a new service manager using a RepositoryManager
func NewServiceManager(repoManager repository.RepositoryManager) ServiceManager {
	return &serviceManager{
		userService: NewUserService(repoManager.User()),
		// Add other services
		// productService: NewProductService(repoManager.Product()),
		// orderService: NewOrderService(repoManager.Order(), repoManager.User()),
		// notificationService: NewNotificationService(repoManager.Notification()),
	}
}

// User returns the user service
func (sm *serviceManager) User() UserService {
	return sm.userService
}

// Add other service getters as needed
// func (sm *serviceManager) Product() ProductService {
//     return sm.productService
// }
//
// func (sm *serviceManager) Order() OrderService {
//     return sm.orderService
// }
//
// func (sm *serviceManager) Notification() NotificationService {
//     return sm.notificationService
// }
