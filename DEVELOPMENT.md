# Development Guide

This guide provides comprehensive information for developers working on My Project.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Building](#building)
- [Docker Development](#docker-development)
- [Release Process](#release-process)
- [CI/CD Pipeline](#cicd-pipeline)
- [Contributing](#contributing)

## Getting Started

### Prerequisites

- Go 1.25.0 or higher
- Docker (optional, for containerized development)
- Make (optional, for build automation)
- Git

### Clone the Repository

```bash
git clone https://github.com/xanderbilla/my-project.git
cd my-project
```

## Development Setup

### Install Dependencies

```bash
go mod download
```

### Configuration

1. Copy the example configuration:
   ```bash
   cp config/config.example.yaml config/dev.yaml
   ```

2. Edit `config/dev.yaml` with your local settings:
   ```yaml
   server:
     port: 8080
     read_timeout: 15s
     write_timeout: 15s
     idle_timeout: 60s
   
   log:
     level: debug
     format: json
   ```

### Run Locally

```bash
# Using Go directly
go run cmd/my-project/main.go

# Using Make
make run

# Using Air for hot reload (recommended for development)
air
```

The server will start on `http://localhost:8080`

## Project Structure

```
my-project/
├── cmd/
│   └── my-project/          # Application entry point
│       ├── main.go          # Main function
│       └── main_test.go     # Main tests
├── internal/
│   ├── app/                 # Application initialization
│   ├── config/              # Configuration management
│   ├── constants/           # Application constants
│   ├── handlers/            # HTTP request handlers
│   ├── middleware/          # HTTP middleware
│   ├── server/              # HTTP server setup
│   ├── types/               # Domain types
│   └── utils/               # Utility functions
│       └── response/        # Response helpers
├── config/                  # Configuration files
├── storage/                 # Data storage (if needed)
├── bin/                     # Compiled binaries (generated)
└── .github/                 # CI/CD workflows
```

### Key Components

- **cmd/my-project/main.go**: Application entry point
- **internal/app**: Application initialization and dependency injection
- **internal/server**: HTTP server with multiplexer and middleware chain
- **internal/handlers**: Route handlers for different endpoints
- **internal/middleware**: Request processing middleware (logging, error recovery, request ID)
- **internal/types**: Domain models and data structures
- **internal/utils/response**: Standardized API response helpers

## Development Workflow

### Branch Strategy

- `my-project/dev`: Main development branch
- `main`: Production-ready code (when you merge)
- Feature branches: `feature/feature-name`
- Bug fixes: `fix/bug-description`

### Commit Messages

Follow conventional commits:

```
feat: add user authentication
fix: resolve memory leak in handler
docs: update API documentation
test: add integration tests for users
ci: update workflow triggers
refactor: simplify response utilities
```

### Code Quality

Run these before committing:

```bash
# Format code
go fmt ./...
gofmt -s -w .

# Lint code
golangci-lint run

# Vet code
go vet ./...

# Run all quality checks
make lint
```

## Testing

### Run All Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...

# Using Make
make test
make test-coverage
```

### Run Specific Tests

```bash
# Test a specific package
go test ./internal/handlers

# Test a specific function
go test -run TestCreateUser ./internal/handlers

# Run with race detection
go test -race ./...
```

### Writing Tests

Example test structure:

```go
func TestHandlerName(t *testing.T) {
    tests := []struct {
        name           string
        input          interface{}
        expectedStatus int
        expectedBody   string
    }{
        // Test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Building

### Build Binary

```bash
# Build for current OS
go build -o bin/my-project cmd/my-project/main.go

# Using Make
make build

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o bin/my-project-linux-amd64 cmd/my-project/main.go
```

### Cross-Platform Builds

```bash
# Build for all platforms
make build-all

# Or manually:
GOOS=linux GOARCH=amd64 go build -o bin/my-project-linux-amd64 cmd/my-project/main.go
GOOS=darwin GOARCH=amd64 go build -o bin/my-project-darwin-amd64 cmd/my-project/main.go
GOOS=windows GOARCH=amd64 go build -o bin/my-project-windows-amd64.exe cmd/my-project/main.go
```

## Docker Development

### Development Container

The development Dockerfile includes hot-reload capabilities:

```bash
# Build development image
docker build -f Dockerfile.dev -t my-project:dev .

# Run with hot reload
docker-compose up

# Or using docker run
docker run -p 8080:8080 -v $(pwd):/app my-project:dev
```

### Production Container

```bash
# Build production image
docker build -t my-project:latest .

# Run production image
docker run -p 8080:8080 my-project:latest
```

### Docker Compose

```bash
# Start all services
docker-compose up

# Start in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## Release Process

### Manual Release

1. **Update version** and ensure all changes are committed

2. **Create a tag** following semantic versioning (vMAJOR.MINOR.PATCH):
   ```bash
   git tag -a v1.0.1 -m "Release v1.0.1 - Description of changes"
   git push origin v1.0.1
   ```

3. **Automated release workflow** will:
   - Create GitHub release with changelog
   - Build binaries for all platforms (Linux, macOS, Windows)
   - Build and push Docker images to Docker Hub and GitHub Packages
   - Generate semantic version tags (v1.0.1 → 1.0.1, 1.0, 1, latest)

### Version Guidelines

- **MAJOR** (v2.0.0): Breaking API changes
- **MINOR** (v1.1.0): New features, backward compatible
- **PATCH** (v1.0.1): Bug fixes, backward compatible

### What Triggers Releases

Releases are triggered automatically when you push a tag matching `v*.*.*`:

```bash
git push origin v1.0.1  # ✓ Triggers release
git push origin 1.0.1   # ✗ Does not trigger
git push origin beta    # ✗ Does not trigger
```

## CI/CD Pipeline

### Workflow Triggers

The CI/CD pipeline runs on:

- Push to `main`, `develop`, `my-project/dev` branches
- Pull requests to `main`, `develop`
- Manual workflow dispatch

### What CI/CD Pipeline Skips

To optimize build time, the pipeline **skips** when you only change:
- All `.md` files (README, CONTRIBUTING, etc.)
- LICENSE
- .gitignore
- docs/ folder
- .vscode/ settings

### CI/CD Jobs

1. **Code Quality Check**: Formatting, linting, vetting
2. **Test**: Run all tests with coverage report
3. **Build**: Cross-platform binary builds
4. **Security Scan**: Gosec security analysis
5. **Docker**: Build and push Docker images
6. **Deployment Check**: Verify deployment readiness
7. **Notify**: Build status summary

### Local CI Simulation

Run the same checks locally before pushing:

```bash
# Quality checks
make lint

# Tests
make test

# Build
make build

# All together
make ci
```

## Environment Variables

### Development

Create a `.env` file in the root directory:

```bash
PORT=8080
LOG_LEVEL=debug
LOG_FORMAT=json
```

### Production

Set these environment variables in your deployment:

```bash
PORT=8080
LOG_LEVEL=info
LOG_FORMAT=json
READ_TIMEOUT=15s
WRITE_TIMEOUT=15s
IDLE_TIMEOUT=60s
```

## API Documentation

### Testing Endpoints

```bash
# Health check
curl http://localhost:8080/health

# Get all users
curl http://localhost:8080/api/users

# Get user by ID
curl http://localhost:8080/api/users/1

# Create user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# Update user
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com"}'

# Delete user
curl -X DELETE http://localhost:8080/api/users/1
```

## Performance Testing

### Benchmarking

```bash
# Run benchmarks
go test -bench=. ./...

# With memory profiling
go test -bench=. -benchmem ./...

# Specific benchmark
go test -bench=BenchmarkHandlerName ./internal/handlers
```

### Load Testing

Using [vegeta](https://github.com/tsenart/vegeta):

```bash
# Install vegeta
go install github.com/tsenart/vegeta@latest

# Run load test
echo "GET http://localhost:8080/api/users" | vegeta attack -duration=30s -rate=1000 | vegeta report
```

## Debugging

### Visual Studio Code

Launch configuration in `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/my-project/main.go"
    }
  ]
}
```

### Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug cmd/my-project/main.go
```

## Troubleshooting

### Common Issues

**Port already in use:**
```bash
lsof -ti:8080 | xargs kill -9
```

**Module issues:**
```bash
go mod tidy
go mod verify
```

**Build cache issues:**
```bash
go clean -cache
go clean -modcache
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

### Quick Contribution Checklist

- [ ] Fork the repository
- [ ] Create a feature branch
- [ ] Write tests for new functionality
- [ ] Ensure all tests pass
- [ ] Run code quality checks
- [ ] Update documentation
- [ ] Submit pull request with clear description

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)

## Getting Help

- Open an issue for bugs or feature requests
- Check existing issues and discussions
- Review the [Code of Conduct](CODE_OF_CONDUCT.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Maintainers

- [@xanderbilla](https://github.com/xanderbilla)
- [@emmabites](https://github.com/emmabites) - Core components contributor
