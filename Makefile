# AI Context Gap Tracker - Makefile

.PHONY: help build run test clean docker-build docker-run docker-stop dev-setup

# Default target
help:
	@echo "AI Context Gap Tracker - Available Commands"
	@echo "==========================================="
	@echo "  build        - Build the Go application"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker images"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  docker-stop  - Stop Docker containers"
	@echo "  dev-setup    - Set up development environment"
	@echo "  system-test  - Run system integration tests"
	@echo "  example      - Run example usage"

# Build Go application
build:
	@echo "🔨 Building Go application..."
	go mod tidy
	go build -o bin/ai-context-tracker ./cmd/main.go
	@echo "✅ Build completed successfully!"

# Run application locally
run: build
	@echo "🚀 Starting AI Context Gap Tracker..."
	./bin/ai-context-tracker

# Run tests
test:
	@echo "🧪 Running Go tests..."
	go test -v ./...
	@echo "✅ Tests completed!"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -f /tmp/response.json
	@echo "✅ Clean completed!"

# Build Docker images
docker-build:
	@echo "🐳 Building Docker images..."
	docker-compose build
	@echo "✅ Docker images built successfully!"

# Run with Docker Compose
docker-run:
	@echo "🚀 Starting services with Docker Compose..."
	docker-compose up -d
	@echo "✅ Services started!"
	@echo "📋 Service URLs:"
	@echo "   - Main API: http://localhost:8080"
	@echo "   - NLP Service: http://localhost:5000"
	@echo "   - PostgreSQL: localhost:5432"
	@echo "   - Redis: localhost:6379"

# Stop Docker containers
docker-stop:
	@echo "🛑 Stopping Docker containers..."
	docker-compose down
	@echo "✅ Containers stopped!"

# Set up development environment
dev-setup:
	@echo "🔧 Setting up development environment..."
	go mod tidy
	@echo "📦 Installing Python dependencies..."
	cd python-nlp && pip install -r requirements.txt
	@echo "🐳 Building Docker images..."
	docker-compose build
	@echo "✅ Development environment ready!"

# Run system integration tests
system-test:
	@echo "🧪 Running system integration tests..."
	./scripts/test_system.sh
	@echo "✅ System tests completed!"

# Run example usage
example: build
	@echo "📚 Running example usage..."
	go run examples/example_usage.go
	@echo "✅ Example completed!"

# Database migrations (if needed)
migrate:
	@echo "🗄️ Running database migrations..."
	# Migrations are handled automatically by the application
	@echo "✅ Migrations completed!"

# Development server with hot reload
dev:
	@echo "🔥 Starting development server..."
	@echo "Note: Install 'air' for hot reload: go install github.com/cosmtrek/air@latest"
	air -c .air.toml || go run cmd/main.go

# Lint code
lint:
	@echo "🔍 Running linters..."
	golangci-lint run || echo "Install golangci-lint for better linting"
	@echo "✅ Linting completed!"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted!"

# Update dependencies
update-deps:
	@echo "📦 Updating dependencies..."
	go mod tidy
	go get -u all
	@echo "✅ Dependencies updated!"

# Production build
prod-build:
	@echo "🏗️ Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/ai-context-tracker-prod ./cmd/main.go
	@echo "✅ Production build completed!"

# Full setup from scratch
full-setup: clean dev-setup docker-build
	@echo "🎉 Full setup completed!"
	@echo "Run 'make docker-run' to start the services"

# Quick test cycle
quick-test: build test
	@echo "⚡ Quick test cycle completed!"

# Health check
health-check:
	@echo "🏥 Checking service health..."
	curl -s http://localhost:8080/api/v1/health | jq . || echo "Service not running"
	curl -s http://localhost:5000/health | jq . || echo "NLP service not running"