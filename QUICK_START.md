# 🚀 Quick Start Guide - PostgreSQL with GORM

## Prerequisites

Make sure you have the following installed:
- Go 1.21+
- PostgreSQL 12+ (or Docker)

## Step 1: Start Services with Docker Compose

### Option A: Using Docker Compose (Recommended)
```bash
# Start all services (PostgreSQL, RabbitMQ, and API)
docker-compose -f .devcontainer/docker-compose.yml up -d

# Or start just the database
docker-compose -f .devcontainer/docker-compose.yml up -d hermes-db
```

### Option B: Local PostgreSQL
```sql
-- Connect to PostgreSQL as superuser
sudo -u postgres psql

-- Create database and user
CREATE DATABASE hermes;
CREATE USER hermes WITH PASSWORD 'hermes';
GRANT ALL PRIVILEGES ON DATABASE hermes TO hermes;
\q
```

## Step 2: Install Dependencies
```bash
go mod tidy
```

## Step 3: Run the Application

### Start the REST Server

#### Inside Docker Compose (recommended)
```bash
# The API will automatically start with the database
docker-compose -f .devcontainer/docker-compose.yml up
```

#### Local Development
```bash
# Set environment for local development
export APP_ENV=local

# Run the server
go run cmd/rest-server/main.go
```

You should see output like:
```
🔌 Connecting to PostgreSQL database
🔄 Running database migrations...
✅ Database setup completed successfully
🚀 Starting Hermes API server
```

### Test the Database Connection
```bash
# Test the health endpoint (Docker Compose)
curl http://localhost:8000/health

# Test the health endpoint (Local development)
curl http://localhost:8080/health

# Expected response:
{
  "status": "ok",
  "service": "hermes-api",
  "timestamp": "2024-01-01T12:00:00Z",
  "version": "1.0.0",
  "database": "ok"
}
```

## Step 4: Test the API

### Create a User
```bash
# Docker Compose
curl -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "username": "john_doe",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'

# Local development
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "username": "john_doe",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Get All Users
```bash
# Docker Compose
curl http://localhost:8000/api/v1/users

# Local development
curl http://localhost:8080/api/v1/users
```

### Get User by ID
```bash
# Docker Compose
curl http://localhost:8000/api/v1/users/1

# Local development
curl http://localhost:8080/api/v1/users/1
```

## Step 5: Run the Database Test Example
```bash
go run examples/database_test.go
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check with database status |
| GET | `/api/v1/users` | List all users (with pagination) |
| POST | `/api/v1/users` | Create a new user |
| GET | `/api/v1/users/:id` | Get user by ID |
| PUT | `/api/v1/users/:id` | Update user |
| DELETE | `/api/v1/users/:id` | Delete user |

## Configuration

### Docker Compose Environment
The database configuration for Docker Compose is in `config/config.yaml`:
```yaml
database:
  host: "hermes-db"
  port: "5432"
  name: "hermes"
  user: "hermes"
  password: "hermes"
  ssl_mode: "disable"
```

### Local Development Environment
For local development, use `config/config.local.yaml`:
```yaml
database:
  host: "localhost"
  port: "5432"
  name: "hermes"
  user: "hermes"
  password: "hermes"
  ssl_mode: "disable"
```

### Environment Variables
You can override these with environment variables:
```bash
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_NAME=hermes
export DATABASE_USER=hermes
export DATABASE_PASSWORD=hermes
export DATABASE_SSL_MODE=disable
```

## Troubleshooting

### Database Connection Issues
1. Check if PostgreSQL is running:
   ```bash
   # Docker
   docker ps | grep postgres
   
   # Local
   sudo systemctl status postgresql
   ```

2. Test connection manually:
   ```bash
   psql -h localhost -p 5432 -U hermes_user -d hermes_dev
   ```

3. Check logs for detailed error messages

### Migration Issues
- Ensure the database exists
- Verify user permissions
- Check GORM logs in the application output

## Next Steps

1. Add more models (Product, Order, etc.)
2. Implement authentication
3. Add database migrations
4. Set up testing
5. Configure for production

## Project Structure

```
hermes-api/
├── cmd/
│   └── rest-server/          # Main application entry point
├── internal/
│   ├── database/             # Database connection and configuration
│   ├── model/                # GORM models
│   ├── repository/           # Data access layer
│   └── service/              # Business logic layer
├── api/rest/                 # HTTP handlers
├── config/                   # Configuration files
├── pkg/                      # Shared packages
└── examples/                 # Example code
```

Your database is now ready to use! 🎉 