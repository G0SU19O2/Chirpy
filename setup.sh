#!/bin/bash

# Chirpy Development Setup Script

set -e

echo "🐦 Setting up Chirpy development environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Go version: $GO_VERSION"

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy
go mod download

# Check if .env exists, if not create a template
if [ ! -f .env ]; then
    echo "📝 Creating .env template..."
    cat > .env << EOF
# Database connection string
DB_URL=chirpy_user:chirpy_password@tcp(localhost:3306)/chirpy?charset=utf8mb4&parseTime=True&loc=Local

# Platform (DEV or PROD)
PLATFORM=DEV
EOF
    echo "⚠️  Please update .env with your actual database configuration"
fi

# Check if sqlc is installed
if command -v sqlc &> /dev/null; then
    echo "🔄 Generating database code..."
    make sqlc-generate
else
    echo "⚠️  sqlc not found. Install it with: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
fi

# Run tests
echo "🧪 Running tests..."
make test

# Build the application
echo "🔨 Building application..."
make build

echo ""
echo "🎉 Setup complete!"
echo ""
echo "Next steps:"
echo "1. Update .env with your database configuration"
echo "2. Start your database (or use: make docker-run)"
echo "3. Run the application: make dev"
echo ""
echo "Available commands:"
echo "  make dev          - Run in development mode"
echo "  make build        - Build the application"
echo "  make test         - Run tests"
echo "  make docker-run   - Run with Docker Compose"
echo "  make clean        - Clean build artifacts"
