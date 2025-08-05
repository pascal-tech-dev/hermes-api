package service

import (
	"context"
	"fmt"

	"hermes-api/internal/model"
	"hermes-api/internal/repository"

	"github.com/google/uuid"
)

// UserService defines the interface for user business logic
type UserService interface {
	// Basic CRUD operations
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*model.User, error)
	GetUserCount(ctx context.Context) (int64, error)

	// Business operations
	// CreateUserWithProfile(ctx context.Context, user *model.User, profile *model.UserProfile) error
	// GetUserWithOrders(ctx context.Context, id uint) (*UserWithOrders, error)
	// RegisterUser(ctx context.Context, email, username, password, firstName, lastName string) (*model.User, error)
	// ActivateUser(ctx context.Context, id uint) error
	// DeactivateUser(ctx context.Context, id uint) error
}

// userService implements UserService
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	// Check if user with same email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Check if user with same username already exists
	existingUser, err = s.userRepo.GetByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return fmt.Errorf("user with username %s already exists", user.Username)
	}

	return s.userRepo.Create(ctx, user)
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetUserByEmail retrieves a user by email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(ctx context.Context, user *model.User) error {
	// Check if user exists
	existingUser, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Update the user
	user.CreatedAt = existingUser.CreatedAt // Preserve creation time
	return s.userRepo.Update(ctx, user)
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

// ListUsers retrieves a list of users with pagination
func (s *userService) ListUsers(ctx context.Context, limit, offset int) ([]*model.User, error) {
	return s.userRepo.List(ctx, limit, offset)
}

// GetUserCount returns the total number of users
func (s *userService) GetUserCount(ctx context.Context) (int64, error) {
	return s.userRepo.Count(ctx)
}
