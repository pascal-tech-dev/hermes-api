#!/bin/bash

# Development script for Hermes API

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
}

# Function to start services
start_services() {
    print_status "Starting services with Docker Compose..."
    docker-compose -f .devcontainer/docker-compose.yml up -d
    print_success "Services started successfully!"
    print_status "Database: localhost:5432"
    print_status "RabbitMQ: localhost:5673"
    print_status "RabbitMQ Management: http://localhost:15673"
}

# Function to stop services
stop_services() {
    print_status "Stopping services..."
    docker-compose -f .devcontainer/docker-compose.yml down
    print_success "Services stopped successfully!"
}

# Function to restart services
restart_services() {
    print_status "Restarting services..."
    stop_services
    start_services
}

# Function to show logs
show_logs() {
    print_status "Showing logs..."
    docker-compose -f .devcontainer/docker-compose.yml logs -f
}

# Function to run the API locally
run_local() {
    print_status "Setting up local environment..."
    export APP_ENV=local
    
    print_status "Installing dependencies..."
    go mod tidy
    
    print_status "Starting API server locally..."
    go run cmd/rest-server/main.go
}

# Function to test the API
test_api() {
    local base_url="http://localhost:8000"
    
    if [ "$APP_ENV" = "local" ]; then
        base_url="http://localhost:8080"
    fi
    
    print_status "Testing API at $base_url..."
    
    # Test health endpoint
    print_status "Testing health endpoint..."
    curl -s "$base_url/health" | jq . || echo "Health check failed"
    
    # Test users endpoint
    print_status "Testing users endpoint..."
    curl -s "$base_url/api/v1/users" | jq . || echo "Users endpoint failed"
}

# Function to run database migrations
run_migrations() {
    print_status "Running database migrations..."
    export APP_ENV=local
    go run examples/database_test.go
}

# Function to show help
show_help() {
    echo "Hermes API Development Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  start       Start all services with Docker Compose"
    echo "  stop        Stop all services"
    echo "  restart     Restart all services"
    echo "  logs        Show service logs"
    echo "  local       Run API locally (requires local PostgreSQL)"
    echo "  test        Test API endpoints"
    echo "  migrate     Run database migrations"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 start    # Start Docker Compose services"
    echo "  $0 local    # Run API locally"
    echo "  $0 test     # Test API endpoints"
}

# Main script logic
case "${1:-help}" in
    start)
        check_docker
        start_services
        ;;
    stop)
        check_docker
        stop_services
        ;;
    restart)
        check_docker
        restart_services
        ;;
    logs)
        check_docker
        show_logs
        ;;
    local)
        run_local
        ;;
    test)
        test_api
        ;;
    migrate)
        run_migrations
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac 