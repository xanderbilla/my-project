.PHONY: run run-quick dev test test-coverage build clean deps fmt lint lint-warn check check-dev install-tools install-air install-golangci-lint docker-check docker-up docker-up-dev docker-down docker-build docker-rebuild docker-logs docker-clean docker-ps help

# Full production workflow: clean, prepare, validate, build, and run
run: clean deps fmt lint-warn test build
	@echo ""
	@echo "✓ All checks passed"
	@echo ""
	@echo "Starting application..."
	@echo "💡 Tip: Press Ctrl+C to stop the server"
	@echo ""
	@lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	@./bin/my-project

# Quick run without checks (for development)
run-quick:
	@echo "Checking for existing server..."
	@lsof -ti:8080 | xargs kill -9 2>/dev/null && echo "✓ Stopped previous server" || echo "✓ No server running"
	@echo "💡 Tip: Press Ctrl+C to stop the server"
	@echo ""
	@go run ./cmd/my-project

# Development mode: prepare and run with hot reload
dev: deps fmt
	@if ! command -v air > /dev/null; then \
		echo "Air not found. Installing..."; \
		$(MAKE) install-air; \
	fi
	@echo "Starting development server with hot reload..."
	@air

# Run all checks (lint + test) - strict mode
check: fmt lint test
	@echo "✓ All checks passed"

# Quick check without failing on lint issues (for development)
check-dev: fmt lint-warn test
	@echo "✓ Development checks complete"

# Run all tests (always fresh, no cache)
test:
	@echo "Clearing test cache..."
	@go clean -testcache
	@go test ./... -v

# Run tests with coverage (always fresh, no cache)
test-coverage:
	@echo "Clearing test cache..."
	@go clean -testcache
	@go test ./... -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build the application
build:
	@go build -o bin/my-project ./cmd/my-project
	@echo "Binary created: bin/my-project"

# Clean build artifacts and caches
clean:
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean -cache -testcache -modcache 2>/dev/null || true
	@echo "Cleaned build artifacts and caches"

# Install dependencies
deps:
	@go mod download
	@go mod tidy

# Format code
fmt:
	@go fmt ./...

# Run linter (auto-installs if missing)
lint:
	@if ! command -v golangci-lint > /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		$(MAKE) install-golangci-lint; \
	fi
	@golangci-lint run

# Run linter with warnings only (doesn't fail)
lint-warn:
	@if ! command -v golangci-lint > /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		$(MAKE) install-golangci-lint; \
	fi
	@golangci-lint run || true

# Install air for hot reload
install-air:
	@echo "Installing air..."
	@go install github.com/air-verse/air@latest
	@echo "✓ Air installed successfully"

# Install golangci-lint
install-golangci-lint:
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✓ golangci-lint installed successfully"

# Install all development tools
install-tools: install-air install-golangci-lint
	@echo "✓ All tools installed"

# Check if Docker is running, start if not (macOS)
docker-check:
	@echo "Checking Docker status..."
	@if ! docker info > /dev/null 2>&1; then \
		echo "Docker is not running. Starting Docker..."; \
		if command -v open > /dev/null; then \
			open -a Docker; \
			echo "Waiting for Docker to start..."; \
			WAIT_COUNT=0; \
			while ! docker info > /dev/null 2>&1 && [ $$WAIT_COUNT -lt 30 ]; do \
				printf "."; \
				sleep 2; \
				WAIT_COUNT=$$((WAIT_COUNT + 1)); \
			done; \
			echo ""; \
			if docker info > /dev/null 2>&1; then \
				echo "✓ Docker is ready"; \
			else \
				echo "✗ Failed to start Docker. Please start Docker Desktop manually."; \
				exit 1; \
			fi; \
		else \
			echo "✗ Cannot auto-start Docker. Please start Docker manually."; \
			exit 1; \
		fi; \
	else \
		echo "✓ Docker is already running"; \
	fi

# Start production services
docker-up: docker-check
	@echo "Checking for existing containers and port conflicts..."
	@docker-compose down 2>/dev/null || true
	@lsof -ti:8080 | xargs kill -9 2>/dev/null && echo "✓ Killed process on port 8080" || true
	@echo "Starting production services..."
	@docker-compose up -d app
	@echo "✓ Services started"
	@echo ""
	@echo "🚀 Application: http://localhost:8080"
	@echo "📊 Health: http://localhost:8080/health"
	@echo ""
	@echo "View logs: make docker-logs"
	@echo "Stop: make docker-down"

# Start development services with hot reload
docker-up-dev: docker-check
	@echo "Checking for existing containers and port conflicts..."
	@docker-compose down 2>/dev/null || true
	@lsof -ti:8080 | xargs kill -9 2>/dev/null && echo "✓ Killed process on port 8080" || true
	@echo "Starting development services with hot reload..."
	@docker-compose --profile dev up -d app-dev
	@echo "✓ Development services started"
	@echo ""
	@echo "🚀 Application: http://localhost:8080"
	@echo "📊 Health: http://localhost:8080/health"
	@echo "🔥 Hot reload enabled"
	@echo ""
	@echo "View logs: make docker-logs"
	@echo "Stop: make docker-down"

# Stop all services
docker-down:
	@echo "Stopping services..."
	@docker-compose down
	@echo "✓ Services stopped"

# Build Docker image
docker-build: docker-check
	@echo "Building Docker image..."
	@docker-compose build
	@echo "✓ Image built successfully"

# Rebuild Docker image (no cache)
docker-rebuild: docker-check
	@echo "Rebuilding Docker image (no cache)..."
	@docker-compose build --no-cache
	@echo "✓ Image rebuilt successfully"

# View logs
docker-logs:
	@docker-compose logs -f

# Show running containers
docker-ps:
	@docker-compose ps

# Clean Docker resources
docker-clean: docker-down
	@echo "Cleaning Docker resources..."
	@docker-compose down -v --remove-orphans
	@docker system prune -f
	@echo "✓ Docker resources cleaned"

# Show help
help:
	@echo "My Project - Available Commands"
	@echo ""
	@echo "Getting Started:"
	@echo "  make run-quick        - Start app immediately (fastest way)"
	@echo "  make dev              - Start with auto-reload (saves time while coding)"
	@echo ""
	@echo "Development:"
	@echo "  make test             - Run your tests"
	@echo "  make test-coverage    - See which code is tested"
	@echo "  make fmt              - Auto-format your code"
	@echo "  make lint             - Check code quality"
	@echo ""
	@echo "Building:"
	@echo "  make build            - Create executable file"
	@echo "  make clean            - Clean up generated files"
	@echo ""
	@echo "Production:"
	@echo "  make run              - Full check & build (use before deploying)"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up        - Start app in Docker (production)"
	@echo "  make docker-up-dev    - Start app in Docker (development + hot reload)"
	@echo "  make docker-down      - Stop Docker services"
	@echo "  make docker-logs      - View Docker logs"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-rebuild   - Rebuild Docker image (no cache)"
	@echo "  make docker-clean     - Clean Docker resources"
	@echo ""
	@echo "💡 Tip: Tools install automatically when needed!"
	@echo "💡 Docker: Auto-starts if not running!"
	@echo "💡 Port conflicts: Automatically cleared before starting!"
