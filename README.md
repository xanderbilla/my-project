# My Project (Go)

A REST API learning project built with Go demonstrating proper error handling, structured logging, and middleware patterns.

## Features

- Structured error responses with user-friendly messages
- JSON logging with configurable levels
- Request ID tracking for debugging
- Input validation
- Graceful shutdown
- Environment-based configuration
- Panic recovery

## Quick Start

### Using Docker (Recommended)

Pull and run the latest image from Docker Hub or GitHub Container Registry:

```bash
# Using Docker Hub
docker pull docker.io/xanderbilla/my-project:latest
docker run -p 8080:8080 docker.io/xanderbilla/my-project:latest

# OR using GitHub Container Registry
docker pull ghcr.io/xanderbilla/my-project:latest
docker run -p 8080:8080 ghcr.io/xanderbilla/my-project:latest

# Run with custom environment variables
docker run -p 8080:8080 \
  -e LOG_LEVEL=debug \
  -e PORT=8080 \
  ghcr.io/xanderbilla/my-project:latest

# Use a specific version (recommended for production)
docker pull ghcr.io/xanderbilla/my-project:v1.0.1
docker run -p 8080:8080 ghcr.io/xanderbilla/my-project:v1.0.1
```

**Available tags:** `latest`, `v1.0.1`, `v1.0.0`, `1.0.1`, `1.0`, `1`

See all [releases](https://github.com/xanderbilla/my-project/releases) for version history.

### Using Local Development

```bash
# 1. Copy environment template
cp .env.example .env

# 2. Start the app
make run-quick
```

That's it! Your server is running on http://localhost:8080

### Using Docker Compose

```bash
# Start all services
docker-compose up

# Start in detached mode
docker-compose up -d
```

## Configuration

All configuration is managed via environment variables. See [.env.example](.env.example) for available options.

Required variables:

- `HTTP_ADDRESS` - Server bind address (e.g., `localhost:8080`)
- `STORAGE_PATH` - Path to storage file (e.g., `storage/dev-storage.db`)

Optional variables:

- `ENV` - Environment: dev, staging, prod (default: prod)
- `LOG_LEVEL` - Log level: DEBUG, INFO, WARN, ERROR (default: INFO)

For detailed configuration guide, see [CONFIGURATION.md](CONFIGURATION.md).

## API Endpoints

| Method | Endpoint         | Description     |
| ------ | ---------------- | --------------- |
| GET    | `/api/users`     | Get all users   |
| GET    | `/api/users/:id` | Get user by ID  |
| POST   | `/api/users`     | Create new user |
| PUT    | `/api/users/:id` | Update user     |
| DELETE | `/api/users/:id` | Delete user     |

## Example Usage

```bash
# Get all users
curl http://localhost:8080/api/users

# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Aman","email":"aman@example.com","age":22}'

# Get specific user
curl http://localhost:8080/api/users/1

# Update user
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Aman Updated","email":"aman@example.com","age":23}'

# Delete user
curl -X DELETE http://localhost:8080/api/users/1
```

For detailed API response format, error handling, and examples, see [DEVELOPMENT.md](DEVELOPMENT.md#api-documentation).

## For Developers

All technical documentation, including setup, testing, building, and release processes, is in:

**[📖 Development Guide](DEVELOPMENT.md)**

Quick links:

- [Development Setup](DEVELOPMENT.md#development-setup) - Local setup and configuration
- [Configuration](DEVELOPMENT.md#configuration) - Environment setup
- [Project Structure](DEVELOPMENT.md#project-structure) - Code organization
- [Testing](DEVELOPMENT.md#testing) - Running tests and coverage
- [Building](DEVELOPMENT.md#building) - Creating binaries
- [Docker Development](DEVELOPMENT.md#docker-development) - Container workflows
- [API Documentation](DEVELOPMENT.md#api-documentation) - Response formats and endpoints
- [Release Process](DEVELOPMENT.md#release-process) - How to release
- [CI/CD Pipeline](DEVELOPMENT.md#cicd-pipeline) - Automated workflows
- [Troubleshooting](DEVELOPMENT.md#troubleshooting) - Common issues

## About

This is a learning project to demonstrate best practices in Go API development. Feel free to use it as a reference for your own projects.

## Contributing

We welcome contributions! Please see:

- [DEVELOPMENT.md](DEVELOPMENT.md) - Development guide and setup
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) - Community guidelines

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Vikas Singh** ([@xanderbilla](https://github.com/xanderbilla))

[xanderbilla.com](https://xanderbilla.com)
