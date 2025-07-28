package main

import (
	"context"
	"fmt"
	"log"

	"hermes-api/config"
	"hermes-api/internal/database"
	"hermes-api/internal/model"
	"hermes-api/internal/repository"
	"hermes-api/internal/service"
	"hermes-api/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	if err := logger.Init(cfg.Logging.Level, cfg.Logging.Format); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Connect to database
	logger.Info("Connecting to database...")
	if err := database.Connect(&cfg.Database); err != nil {
		logger.Fatal("Failed to connect to database", err)
	}

	// Run migrations
	logger.Info("Running database migrations...")
	if err := database.AutoMigrate(); err != nil {
		logger.Fatal("Failed to run migrations", err)
	}

	// Initialize repository and service
	userRepo := repository.NewUserRepository(database.DB)
	userService := service.NewUserService(userRepo)

	ctx := context.Background()

	// Create a test user
	testUser := &model.User{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashed_password_here",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	logger.Info("Creating test user...")
	if err := userService.CreateUser(ctx, testUser); err != nil {
		logger.Fatal("Failed to create user", err)
	}

	logger.Info("User created successfully", zap.Uint("user_id", testUser.ID))

	// Retrieve the user
	logger.Info("Retrieving user...")
	retrievedUser, err := userService.GetUserByID(ctx, testUser.ID)
	if err != nil {
		logger.Fatal("Failed to retrieve user", err)
	}

	fmt.Printf("Retrieved user: %+v\n", retrievedUser)

	// List all users
	logger.Info("Listing all users...")
	users, err := userService.ListUsers(ctx, 10, 0)
	if err != nil {
		logger.Fatal("Failed to list users", err)
	}

	fmt.Printf("Found %d users:\n", len(users))
	for _, user := range users {
		fmt.Printf("- %s (%s)\n", user.Username, user.Email)
	}

	// Get user count
	count, err := userService.GetUserCount(ctx)
	if err != nil {
		logger.Fatal("Failed to get user count", err)
	}

	logger.Info("Database test completed successfully", zap.Int64("total_users", count))

	// Close database connection
	if err := database.Close(); err != nil {
		logger.Error("Failed to close database connection", err)
	} else {
		logger.Info("Database connection closed")
	}
}
