package main

import (
	"context"
	"hermes-api/config"
	"hermes-api/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

// loadConfig loads configuration from environment variables
func loadConfig() (*config.Config, error) {
	return config.Load()
}

// setupMiddleware configures Fiber middleware
func setupMiddleware(app *fiber.App, cfg *config.Config) {
	// Recovery middleware
	app.Use(recover.New())

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Logger middleware
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))
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

	// Notifications routes (placeholder for now)
	notifications := api.Group("/notifications")
	notifications.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Notifications endpoint - coming soon",
			"status":  "not implemented",
		})
	})

	notifications.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Create notification endpoint - coming soon",
			"status":  "not implemented",
		})
	})

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Not Found",
			"message": "The requested resource was not found",
			"path":    c.Path(),
		})
	})
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

	log := logger.GetLogger()

	log.Info("🚀 Starting Hermes API server",
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
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Error("Request error", err,
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.String("ip", c.IP()),
			)

			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		},
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
