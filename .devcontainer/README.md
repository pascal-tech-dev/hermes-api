# Hermes API Dev Container

This directory contains the development container configuration for the Hermes API project.

## Features

### Enhanced Go Development Environment
- **Go 1.24.5** with all essential development tools
- **gopls** language server for enhanced IntelliSense
- **golangci-lint** for code linting
- **delve** debugger for debugging
- **goimports** for automatic import management
- **staticcheck** for static analysis
- **gotests** for test generation
- **gomodifytags** for struct tag management

### Pre-installed Extensions
- **Go** - Official Go extension
- **Docker** - Docker support
- **Prettier** - Code formatting
- **GitLens** - Enhanced Git capabilities
- **YAML** - YAML file support
- **JSON** - JSON file support
- **Makefile Tools** - Makefile support
- **Remote Development** - Container, SSH, and WSL support

### Database & Message Queue
- **PostgreSQL 15** with Alpine Linux
- **RabbitMQ 3** with management interface

### Performance Optimizations
- **Go module cache** mounted as volume for faster builds
- **Cached workspace mounting** for better performance
- **Automatic go mod download** on container creation
- **Automatic go mod tidy** on container start

## Usage

1. Open the project in VS Code
2. When prompted, click "Reopen in Container"
3. Wait for the container to build and start
4. The development environment will be ready with:
   - Go tools pre-installed
   - Database running on port 5432
   - RabbitMQ running on ports 5672 (AMQP) and 15672 (Management UI)
   - API server ports 8080 (REST) and 50051 (gRPC) forwarded

## Environment Variables

- `GOPATH=/go`
- `GOROOT=/usr/local/go`
- `GO111MODULE=on`
- `CGO_ENABLED=1`

## Ports

- **8080** - REST API
- **50051** - gRPC API
- **5432** - PostgreSQL
- **5672** - RabbitMQ AMQP
- **15672** - RabbitMQ Management UI

## Database Credentials

- **Username**: hermes
- **Password**: hermes
- **Database**: hermes

## RabbitMQ Credentials

- **Username**: hermes
- **Password**: hermes
- **Management UI**: http://localhost:15672

## Recent Improvements

1. **Enhanced Go Settings**: Added comprehensive Go language server configuration
2. **Better Tooling**: Pre-installed essential Go development tools
3. **Performance**: Added module cache mounting and optimized workspace mounting
4. **Security**: Improved user setup and package management
5. **Developer Experience**: Added useful VS Code extensions and settings
6. **Automation**: Added post-create and post-start commands for setup automation

## Troubleshooting

If you encounter issues:

1. **Rebuild the container**: Command Palette → "Dev Containers: Rebuild Container"
2. **Check logs**: View the Dev Container logs in the Output panel
3. **Clear cache**: Delete the `.devcontainer/go-cache/` directory and rebuild
4. **Update tools**: Run `go install -u all` inside the container to update Go tools 