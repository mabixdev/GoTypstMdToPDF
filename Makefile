# Makefile for Markdown to PDF Service (Go)

.PHONY: build run clean test deps help

# Default target
help:
	@echo "Available targets:"
	@echo "  build    - Build the Go binary"
	@echo "  run      - Run the service"
	@echo "  dev      - Run in development mode with auto-reload"
	@echo "  deps     - Download dependencies"
	@echo "  clean    - Clean build artifacts"
	@echo "  test     - Run tests"
	@echo "  help     - Show this help"

# Build the Go binary
build:
	@echo "🔨 Building Go service..."
	go build -o bin/md-pdf-service *.go

# Run the service
run: build
	@echo "🚀 Starting Markdown to PDF Service (Go)..."
	./bin/md-pdf-service

# Development mode (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "🔄 Starting in development mode..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Running without auto-reload..."; \
		go run *.go; \
	fi

# Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -rf temp/
	go clean

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Create necessary directories
setup:
	@echo "📁 Setting up directories..."
	mkdir -p bin
	mkdir -p temp

# Install development tools
install-tools:
	@echo "🔧 Installing development tools..."
	go install github.com/cosmtrek/air@latest

# Build for different platforms
build-all: setup
	@echo "🌍 Building for multiple platforms..."
	GOOS=darwin GOARCH=amd64 go build -o bin/md-pdf-service-darwin-amd64 *.go
	GOOS=darwin GOARCH=arm64 go build -o bin/md-pdf-service-darwin-arm64 *.go
	GOOS=linux GOARCH=amd64 go build -o bin/md-pdf-service-linux-amd64 *.go
	GOOS=windows GOARCH=amd64 go build -o bin/md-pdf-service-windows-amd64.exe *.go

# Docker targets
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t md-pdf-service-go .

docker-run: docker-build
	@echo "🐳 Running Docker container..."
	docker run -p 3000:3000 -v $(PWD)/public:/app/public md-pdf-service-go 