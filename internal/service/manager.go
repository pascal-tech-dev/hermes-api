// internal/service/manager.go
package service

import (
	"hermes-api/internal/repository"
)

// ServiceManager manages all services
type ServiceManager interface {
	User() UserService
	Auth() AuthService
}

// serviceManager implements ServiceManager
type serviceManager struct {
	userService UserService
	authService AuthService
}

// NewServiceManager creates a new service manager using a RepositoryManager
func NewServiceManager(repoManager repository.RepositoryManager, jwtSecret string) ServiceManager {
	return &serviceManager{
		userService: NewUserService(repoManager.User()),
		authService: NewAuthService(repoManager.User(), jwtSecret),
	}
}

// User returns the user service
func (sm *serviceManager) User() UserService {
	return sm.userService
}

// Auth returns the auth service
func (sm *serviceManager) Auth() AuthService {
	return sm.authService
}
