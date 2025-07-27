package main

import (
	"context"
	"hermes-api/api/rest"
	"hermes-api/config"
	"hermes-api/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"hermes-api/internal/middleware"
)

// loadConfig loads configuration from environment variables
func loadConfig() (*config.Config, error) {
	return config.Load()
}

// setupMiddleware configures Fiber middleware
func setupMiddleware(app *fiber.App, _ *config.Config) {
	// Request ID middleware (FIRST - generates request ID for all subsequent middleware)
	app.Use(middleware.RequestID())

	// Custom recovery middleware (catches panics)
	app.Use(middleware.Recovery())

	// CORS middleware
	app.Use(middleware.CORS())

	// Logger middleware
	app.Use(middleware.Logger())

	// Error handler middleware (before routes)
	app.Use(middleware.ErrorHandler)
}

// setupRoutes configures API routes
func setupRoutes(app *fiber.App) {
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"service":   "hermes-api",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")
	rest.SetupRoutes(api)
}

func main() {
	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Initialize logger
	if err := logger.Init(cfg.Logging.Level, cfg.Logging.Format); err != nil {
		log.Fatalf("❌ Failed to initialize logger: %v", err)
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}()

	logger.Info("🚀 Starting Hermes API server",
		zap.String("environment", cfg.Server.Environment),
		zap.String("port", cfg.Server.Port),
		zap.String("log_level", cfg.Logging.Level),
	)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Hermes API",
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	})

	// Setup middleware
	setupMiddleware(app, cfg)

	// Setup routes
	setupRoutes(app)

	// Start server in a goroutine
	go func() {
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			logger.Fatal("❌ Failed to start server", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", err)
	}

	logger.Info("✅ Server exited gracefully")
}
