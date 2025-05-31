# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/chirpy ./cmd/chirpy

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/bin/chirpy .

# Copy static files
COPY --from=builder /app/web ./web

# Copy SQL files (if needed for migrations)
COPY --from=builder /app/sql ./sql

EXPOSE 8080

CMD ["./chirpy"]
