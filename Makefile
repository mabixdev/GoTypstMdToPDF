# Makefile for Markdown to PDF Service (Go)

.PHONY: build run clean test deps help build-cli build-api-only

# Default target
help:
	@echo "Available targets:"
	@echo "  build       - Build the full service with web UI"
	@echo "  build-cli   - Build command-line version only"
	@echo "  build-api   - Build API-only version (no web UI)"
	@echo "  run         - Run the full service"
	@echo "  run-cli     - Run CLI version (requires args)"
	@echo "  run-api     - Run API-only service"
	@echo "  dev         - Run in development mode with auto-reload"
	@echo "  deps        - Download dependencies"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  help        - Show this help"

# Build the full Go service (with web UI)
build:
	@echo "🔨 Building full service..."
	go build -o bin/md-pdf-service *.go

# Build CLI version only
build-cli: setup
	@echo "🔨 Building CLI version..."
	go build -o bin/md-pdf-cli cmd/cli/main.go

# Build API-only version (no static files)
build-api: setup
	@echo "🔨 Building API-only service..."
	go build -ldflags="-X main.ApiOnly=true" -o bin/md-pdf-api-only *.go

# Run the full service
run: build
	@echo "🚀 Starting Markdown to PDF Service (Go)..."
	./bin/md-pdf-service

# Run CLI version (example usage)
run-cli: build-cli
	@echo "🚀 Running CLI version (example)..."
	@if [ -f "test.md" ]; then \
		./bin/md-pdf-cli -input test.md -output cli-output.pdf; \
	else \
		echo "Create test.md file first, or use: ./bin/md-pdf-cli -input your-file.md"; \
	fi

# Run API-only version
run-api: build-api
	@echo "🚀 Starting API-only service..."
	./bin/md-pdf-api-only

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
	mkdir -p cmd/cli
	mkdir -p pkg/mdpdf
	mkdir -p examples

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
	# CLI versions
	GOOS=darwin GOARCH=amd64 go build -o bin/md-pdf-cli-darwin-amd64 cmd/cli/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/md-pdf-cli-linux-amd64 cmd/cli/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/md-pdf-cli-windows-amd64.exe cmd/cli/main.go

# Docker targets
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t md-pdf-service-go .

docker-run: docker-build
	@echo "🐳 Running Docker container..."
	docker run -p 3000:3000 -v $(PWD)/public:/app/public md-pdf-service-go

# Docker API-only
docker-build-api:
	@echo "🐳 Building API-only Docker image..."
	docker build -f Dockerfile.api -t md-pdf-service-api .

docker-run-api: docker-build-api
	@echo "🐳 Running API-only Docker container..."
	docker run -p 3000:3000 md-pdf-service-api 