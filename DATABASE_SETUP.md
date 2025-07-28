# PostgreSQL Database Setup with GORM

This guide will help you set up PostgreSQL with GORM for the Hermes API project.

## Prerequisites

- Go 1.21 or later
- PostgreSQL 12 or later
- Docker (optional, for running PostgreSQL in a container)

## Quick Start with Docker

### 1. Start PostgreSQL with Docker

```bash
docker run --name hermes-postgres \
  -e POSTGRES_DB=hermes_dev \
  -e POSTGRES_USER=hermes_user \
  -e POSTGRES_PASSWORD=hermes_password \
  -p 5432:5432 \
  -d postgres:15
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Run the Application

```bash
go run cmd/main.go
```

## Manual PostgreSQL Setup

### 1. Install PostgreSQL

#### Ubuntu/Debian:
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

#### macOS:
```bash
brew install postgresql
```

#### Windows:
Download from [PostgreSQL official website](https://www.postgresql.org/download/windows/)

### 2. Create Database and User

```sql
-- Connect to PostgreSQL as superuser
sudo -u postgres psql

-- Create database
CREATE DATABASE hermes_dev;

-- Create user
CREATE USER hermes_user WITH PASSWORD 'hermes_password';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE hermes_dev TO hermes_user;

-- Exit
\q
```

### 3. Update Configuration

Edit `config/config.yaml` with your database credentials:

```yaml
database:
  host: "localhost"
  port: "5432"
  name: "hermes_dev"
  user: "hermes_user"
  password: "hermes_password"
  ssl_mode: "disable"
```

## Project Structure

```
internal/
├── database/
│   └── database.go          # Database connection and configuration
├── model/
│   └── user.go             # GORM models
├── repository/
│   └── user_repository.go  # Data access layer
└── service/
    └── user_service.go     # Business logic layer
```

## Usage Examples

### Creating a User

```go
package main

import (
    "context"
    "log"

    "hermes-api/config"
    "hermes-api/internal/database"
    "hermes-api/internal/model"
    "hermes-api/internal/repository"
    "hermes-api/internal/service"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    // Connect to database
    err = database.Connect(&cfg.Database)
    if err != nil {
        log.Fatal(err)
    }

    // Run migrations
    err = database.AutoMigrate()
    if err != nil {
        log.Fatal(err)
    }

    // Initialize repository and service
    userRepo := repository.NewUserRepository(database.DB)
    userService := service.NewUserService(userRepo)

    // Create a new user
    user := &model.User{
        Email:     "john@example.com",
        Username:  "john_doe",
        Password:  "hashed_password",
        FirstName: "John",
        LastName:  "Doe",
    }

    ctx := context.Background()
    err = userService.CreateUser(ctx, user)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("User created with ID: %d", user.ID)
}
```

### Querying Users

```go
// Get user by ID
user, err := userService.GetUserByID(ctx, 1)
if err != nil {
    log.Fatal(err)
}

// Get user by email
user, err := userService.GetUserByEmail(ctx, "john@example.com")
if err != nil {
    log.Fatal(err)
}

// List users with pagination
users, err := userService.ListUsers(ctx, 10, 0) // limit=10, offset=0
if err != nil {
    log.Fatal(err)
}

// Get total count
count, err := userService.GetUserCount(ctx)
if err != nil {
    log.Fatal(err)
}
```

## GORM Features Used

### Model Tags

- `gorm:"primaryKey"` - Primary key
- `gorm:"uniqueIndex"` - Unique index
- `gorm:"not null"` - NOT NULL constraint
- `gorm:"default:true"` - Default value
- `gorm:"index"` - Index
- `json:"-"` - Exclude from JSON serialization

### Hooks

- `BeforeCreate` - Runs before creating a record
- `BeforeUpdate` - Runs before updating a record
- `BeforeDelete` - Runs before deleting a record

### Soft Deletes

The User model uses `gorm.DeletedAt` for soft deletes. Deleted records are not actually removed from the database but marked as deleted.

## Environment Variables

You can override configuration using environment variables:

```bash
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_NAME=hermes_dev
export DATABASE_USER=hermes_user
export DATABASE_PASSWORD=hermes_password
export DATABASE_SSL_MODE=disable
```

## Troubleshooting

### Connection Issues

1. Check if PostgreSQL is running:
   ```bash
   sudo systemctl status postgresql
   ```

2. Verify connection parameters in `config/config.yaml`

3. Test connection manually:
   ```bash
   psql -h localhost -p 5432 -U hermes_user -d hermes_dev
   ```

### Migration Issues

1. Check if the database exists
2. Verify user permissions
3. Check GORM logs for detailed error messages

### Performance Tips

1. Use connection pooling (already configured)
2. Add indexes for frequently queried fields
3. Use transactions for multiple operations
4. Consider using GORM's query optimization features

## Next Steps

1. Add more models (e.g., Product, Order, etc.)
2. Implement authentication and authorization
3. Add database migrations for schema changes
4. Set up database backups
5. Configure connection pooling for production 