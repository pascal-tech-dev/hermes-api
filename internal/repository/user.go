package repository

import (
	"context"

	"gorm.io/gorm"

	"hermes-api/internal/model"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {

	// Basic CRUD operations
	BaseRepository[model.User]

	// Query operations
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)

	// // Advanced query operations
	// FindByIDs(ctx context.Context, ids []uint) ([]*model.User, error)
	// ExistsByEmail(ctx context.Context, email string) (bool, error)
	// ExistsByUsername(ctx context.Context, username string) (bool, error)

	// // Transaction support
	// WithTransaction(tx interface{}) UserRepository
}

// userRepository implements UserRepository
type userRepository struct {
	BaseRepository[model.User]
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository[model.User](db),
		db:             db,
	}
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
