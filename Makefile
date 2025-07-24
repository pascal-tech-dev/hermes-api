# Hermes API Makefile

# Variables
BINARY_NAME=hermes-api
BUILD_DIR=bin
MAIN_PATH=cmd/rest-server/main.go

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "🔨 Building Hermes API..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build completed!"

# Run the application
.PHONY: run
run:
	@echo "🚀 Starting Hermes API..."
	@go run $(MAIN_PATH)

# Run the application in development mode with hot reload
.PHONY: dev
dev:
	@echo "🔥 Starting Hermes API in development mode..."
	@go install github.com/cosmtrek/air@latest
	@air

# Clean build artifacts
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "✅ Clean completed!"

# Install dependencies
.PHONY: deps
deps:
	@echo "📦 Installing dependencies..."
	@go mod tidy
	@go mod download
	@echo "✅ Dependencies installed!"

# Run tests
.PHONY: test
test:
	@echo "🧪 Running tests..."
	@go test ./...
	@echo "✅ Tests completed!"

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated!"

# Lint the code
.PHONY: lint
lint:
	@echo "🔍 Linting code..."
	@go install golang.org/x/lint/golint@latest
	@golint ./...
	@echo "✅ Linting completed!"

# Format the code
.PHONY: fmt
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Formatting completed!"

# Generate API documentation (placeholder for future)
.PHONY: docs
docs:
	@echo "📚 Generating API documentation..."
	@echo "⚠️  Documentation generation not implemented yet"
	@echo "✅ Documentation placeholder created!"

# Docker build
.PHONY: docker-build
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t hermes-api .
	@echo "✅ Docker image built!"

# Docker run
.PHONY: docker-run
docker-run:
	@echo "🐳 Running Docker container..."
	@docker run -p 8080:8080 hermes-api
	@echo "✅ Docker container started!"

# Help
.PHONY: help
help:
	@echo "Hermes API - Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Run with hot reload (requires air)"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  lint          - Lint the code"
	@echo "  fmt           - Format the code"
	@echo "  docs          - Generate API documentation"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  help          - Show this help message"
