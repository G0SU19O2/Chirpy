# Chirpy

A Twitter-like web service built with Go that supports user authentication, chirp management, and admin features.

## Project Structure

```
├── cmd/
│   └── chirpy/                 # Main application entry point
│       └── main.go
├── internal/                   # Private application code
│   ├── auth/                  # Authentication utilities
│   │   ├── auth.go
│   │   └── auth_test.go
│   ├── config/                # Application configuration
│   │   └── config.go
│   ├── handlers/              # HTTP request handlers
│   │   ├── chirp.go           # Chirp CRUD operations
│   │   ├── chirp_test.go
│   │   ├── handler_test.go
│   │   ├── metrics.go         # Metrics tracking
│   │   ├── polka.go           # Webhook handlers
│   │   ├── readiness.go       # Health check
│   │   ├── reset.go           # Admin reset functionality
│   │   ├── user.go            # User management & auth
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

- **User Authentication**: JWT-based user registration, login, and token management
- **Chirp Management**: Create, read, update, and delete chirps (tweets)
- **User Management**: User profiles and account updates
- **Admin Dashboard**: Metrics tracking and admin-only operations
- **Health Monitoring**: Health check endpoints for system monitoring
- **Static File Serving**: Efficient static asset delivery with metrics
- **Database Integration**: MySQL database with GORM ORM
- **Webhook Support**: External service integration (Polka webhooks)

## API Endpoints

### Public API

- `GET /api/healthz` - Health check endpoint
- `POST /api/users` - Create a new user account
- `POST /api/login` - User authentication

### Chirp Management

- `POST /api/chirps` - Create a new chirp (requires authentication)
- `GET /api/chirps` - Get all chirps (supports sorting and filtering)
- `GET /api/chirps/{chirpID}` - Get a specific chirp by ID
- `DELETE /api/chirps/{chirpID}` - Delete a chirp (requires authentication & ownership)

### User Management

- `PUT /api/users` - Update user information (requires authentication)
- `POST /api/refresh` - Refresh access token using refresh token
- `POST /api/revoke` - Revoke refresh token (logout)

### Webhooks

- `POST /api/polka/webhooks` - Handle external webhooks (Polka integration)

### Admin API

- `GET /admin/metrics` - View application metrics and statistics
- `POST /admin/reset` - Reset application data (DEV environment only)

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

   ```env
   DB_URL=chirpy_user:chirpy_password@tcp(localhost:3306)/chirpy?charset=utf8mb4&parseTime=True&loc=Local
   PLATFORM=DEV
   JWT_SECRET=your-secret-key-here
   POLKA_API_KEY=your-polka-api-key-here
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
make test
```

Run tests with coverage:

```bash
make test-coverage
```

Run tests for a specific package:

```bash
go test ./internal/handlers
go test ./internal/auth
```

## Authentication & Security

Chirpy uses JWT (JSON Web Tokens) for authentication:

- **Access Tokens**: Short-lived tokens for API authentication
- **Refresh Tokens**: Long-lived tokens for obtaining new access tokens
- **Password Hashing**: Secure password storage using bcrypt
- **Token Validation**: Middleware-based authentication for protected endpoints

### Authentication Flow

1. User registers with `POST /api/users`
2. User logs in with `POST /api/login` to receive access and refresh tokens
3. Include access token in `Authorization: Bearer <token>` header for protected endpoints
4. Use `POST /api/refresh` to get new access tokens when they expire
5. Use `POST /api/revoke` to logout and invalidate refresh tokens

## Development

This project follows Go best practices for project layout:

- `cmd/` contains the main applications for this project
- `internal/` contains private application and library code
- `web/` contains web application specific components

The project uses:

- **GORM** for database operations and ORM
- **MySQL** as the primary database
- **JWT (JSON Web Tokens)** for authentication
- **Go's standard HTTP library** for web server
- **SQLC** for type-safe SQL code generation
- **Docker & Docker Compose** for containerization

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

- **Debugging**: Press `F5` to start debugging the application
- **Tasks**: Use `Ctrl+Shift+P` → "Tasks: Run Task" to access predefined tasks:
  - Build Chirpy
  - Run Chirpy (development mode)
  - Test All
  - Generate sqlc
- **Launch Configurations**: Configured for running and testing the application

Available VS Code Tasks:
- **Build Chirpy** (`make build`)
- **Run Chirpy** (`make dev`) - Background task for development
- **Test All** (`make test`) - Run all tests
- **Generate sqlc** (`make sqlc-generate`) - Regenerate database code
