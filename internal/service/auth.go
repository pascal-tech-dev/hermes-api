package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"hermes-api/internal/model"
	"hermes-api/internal/repository"
	"hermes-api/pkg/errorx"
	"hermes-api/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AuthService defines the interface for authentication logic
type AuthService interface {
	Register(ctx context.Context, email, username, password, firstName, lastName string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, *model.User, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetUserFromToken(tokenString string) (*model.User, error)
}

// authService implements AuthService
type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account
func (s *authService) Register(ctx context.Context, email, username, password, firstName, lastName string) (*model.User, error) {
	// Check if user with same email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		appErr := errorx.New(
			errorx.ErrorTypeConflict,
			errorx.ErrorCodeUserAlreadyExists,
			fmt.Sprintf("User with email %s already exists", email),
		)
		return nil, appErr
	}

	// Check if user with same username already exists
	existingUser, err = s.userRepo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		appErr := errorx.New(
			errorx.ErrorTypeConflict,
			errorx.ErrorCodeUserAlreadyExists,
			fmt.Sprintf("User with username %s already exists", username),
		)
		return nil, appErr
	}

	// Create new user
	user := &model.User{
		Email:     email,
		Username:  username,
		Password:  password, // Will be hashed by GORM hook
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		appErr := errorx.New(
			errorx.ErrorTypeInternal,
			errorx.ErrorCodeDatabaseError,
			fmt.Sprintf("Failed to create user: %v", err),
		)
		return nil, appErr
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(ctx context.Context, email, password string) (string, *model.User, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appErr := errorx.New(
				errorx.ErrorTypeUnauthorized,
				errorx.ErrorCodeInvalidCredentials,
				fmt.Sprintf("User with email %s not found", email),
			)
			return "", nil, appErr
		}
		appErr := errorx.New(
			errorx.ErrorTypeInternal,
			errorx.ErrorCodeDatabaseError,
			"Failed to fetch user data",
		)
		return "", nil, appErr
	}

	// Verify password
	if !user.CheckPassword(password) {
		appErr := errorx.New(
			errorx.ErrorTypeUnauthorized,
			errorx.ErrorCodeInvalidCredentials,
			"Invalid credentials",
		)
		return "", nil, appErr
	}

	// Check if user is active, if not Forbidden
	if !user.IsActive {
		appErr := errorx.New(
			errorx.ErrorTypeForbidden,
			errorx.ErrorCodeAccountDeactivated,
			"Account is deactivated",
		)
		return "", nil, appErr
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		appErr := errorx.New(
			errorx.ErrorTypeInternal,
			errorx.ErrorCodeTokenGenerationFailed,
			"Failed to generate token",
		)
		return "", nil, appErr
	}

	return token, user, nil
}

// ValidateToken validates a JWT token
func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

// GetUserFromToken extracts user information from a JWT token
func (s *authService) GetUserFromToken(tokenString string) (*model.User, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	logger.Info("User ID from token", zap.String("user_id", userIDStr))
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	// Get user from database
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	user, err := s.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

// generateToken creates a JWT token for a user
func (s *authService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
