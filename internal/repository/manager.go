package repository

import (
	"context"

	"gorm.io/gorm"
)

// RepositoryManager manages all repositories
type RepositoryManager interface {
	User() UserRepository
	Application() ApplicationRepository

	// Transaction support
	WithTransaction(ctx context.Context, fn func(RepositoryManager) error) error
}

// repositoryManager implements RepositoryManager
type repositoryManager struct {
	db          *gorm.DB
	user        UserRepository
	application ApplicationRepository
}

// NewRepositoryManager creates a new repository manager
func NewRepositoryManager(db *gorm.DB) RepositoryManager {
	return &repositoryManager{
		db:          db,
		user:        NewUserRepository(db),
		application: NewApplicationRepository(db),
	}
}

// User returns the user repository
func (rm *repositoryManager) User() UserRepository {
	return rm.user
}

// Application returns the application repository
func (rm *repositoryManager) Application() ApplicationRepository {
	return rm.application
}

// WithTransaction executes a function within a database transaction
func (rm *repositoryManager) WithTransaction(ctx context.Context, fn func(RepositoryManager) error) error {
	return rm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txManager := &repositoryManager{
			db:   tx,
			user: NewUserRepository(tx),
		}
		return fn(txManager)
	})
}
