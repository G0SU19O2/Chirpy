# Chirpy

A simple Twitter-like web service built with Go.

## Project Structure

```
├── cmd/
│   └── chirpy/                 # Main application entry point
│       └── main.go
├── internal/                   # Private application code
│   ├── config/                # Application configuration
│   │   └── config.go
│   ├── handlers/              # HTTP request handlers
│   │   ├── chirp.go
│   │   ├── chirp_test.go
│   │   ├── metrics.go
│   │   ├── readiness.go
│   │   ├── reset.go
│   │   ├── user.go
│   │   └── utils.go
│   ├── middleware/            # HTTP middleware
│   │   └── metrics.go
│   ├── models/                # Application models
│   │   └── models.go
│   └── router/                # Route configuration
│       └── router.go
└── web/                       # Static web assets
    └── static/
        ├── index.html
        └── assets/
            └── logo.png
```

## Features

- User management (create users)
- Chirp validation with profanity filtering
- Admin metrics dashboard
- Health check endpoint
- Static file serving

## API Endpoints

### Public API

- `GET /api/healthz` - Health check
- `POST /api/validate_chirp` - Validate and clean chirp content
- `POST /api/users` - Create a new user

### Admin API

- `GET /admin/metrics` - View application metrics
- `POST /admin/reset` - Reset application data (DEV only)

### Static Files

- `/app/*` - Static file serving with metrics tracking

## Getting Started

### Quick Setup

Run the setup script for automatic environment configuration:

```bash
./setup.sh
```

### Manual Setup

1. Set up your environment variables in `.env`:

   ```
   DB_URL=chirpy_user:chirpy_password@tcp(localhost:3306)/chirpy?charset=utf8mb4&parseTime=True&loc=Local
   PLATFORM=DEV
   ```

2. Install dependencies:

   ```bash
   make deps
   ```

3. Generate database code (requires sqlc):

   ```bash
   make sqlc-generate
   ```

4. Build and run the application:

   ```bash
   make build
   make run
   ```

   Or run directly in development mode:

   ```bash
   make dev
   ```

5. The server will start on port 8080.

### Docker Setup

Run with Docker Compose (includes MySQL database):

```bash
make docker-run
```

Stop the Docker setup:

```bash
make docker-down
```

## Testing

Run all tests:

```bash
go test ./...
```

Run tests for a specific package:

```bash
go test ./internal/handlers
```

## Development

This project follows Go best practices for project layout:

- `cmd/` contains the main applications for this project
- `internal/` contains private application and library code
- `web/` contains web application specific components

The project uses:

- GORM for database operations
- Standard library HTTP server

### Available Make Commands

```bash
make build          # Build the application
make run            # Run the built application
make dev            # Run in development mode
make test           # Run all tests
make test-coverage  # Run tests with coverage
make clean          # Clean build artifacts
make deps           # Install dependencies
make sqlc-generate  # Generate database code
make fmt            # Format code
make lint           # Lint code (requires golangci-lint)
make docker-build   # Build Docker image
make docker-run     # Run with Docker Compose
make docker-down    # Stop Docker containers
make docker-clean   # Clean Docker resources
```

### VS Code Integration

The project includes VS Code configuration for:

- Debugging (`F5` to start debugging)
- Tasks (Ctrl+Shift+P → "Tasks: Run Task")
- Launch configurations for running and testing
