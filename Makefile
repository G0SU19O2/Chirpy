.PHONY: build run test clean dev

# Build the application
build:
	go build -o bin/chirpy ./cmd/chirpy

# Run the application
run: build
	./bin/chirpy

# Run in development mode (without building binary)
dev:
	go run ./cmd/chirpy

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Generate database code with sqlc
sqlc-generate:
	sqlc generate

# Docker commands
docker-build:
	docker build -t chirpy .

docker-run:
	docker-compose up --build

docker-down:
	docker-compose down

docker-clean:
	docker-compose down -v
	docker rmi chirpy 2>/dev/null || true

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Create bin directory
bin:
	mkdir -p bin
