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

```bash
# 1. Copy environment template
cp .env.example .env

# 2. Start the app
make run-quick
```

That's it! Your server is running on http://localhost:8080

## Common Commands

```bash
make run-quick     # Start app immediately
make dev           # Start with auto-reload (changes apply automatically)
make test          # Run tests
make build         # Build executable
make help          # See all commands
```

💡 **Tip:** All tools install automatically when needed!

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

## Response Format

All responses follow a consistent structure:

**Success Response:**

```json
{
  "success": true,
  "status": 200,
  "timestamp": "2026-03-04T10:30:00Z",
  "requestId": "550e8400-e29b-41d4-a716-446655440000",
  "path": "/api/users",
  "method": "GET",
  "message": "Users retrieved successfully",
  "data": {
    "id": 1,
    "name": "Aman",
    "email": "aman@example.com",
    "age": 22
  },
  "meta": {
    "apiVersion": "v1"
  }
}
```

**Error Response:**

```json
{
  "success": false,
  "status": 422,
  "timestamp": "2026-03-04T10:30:00Z",
  "requestId": "550e8400-e29b-41d4-a716-446655440000",
  "path": "/api/users",
  "method": "POST",
  "error": {
    "type": "VALIDATION_ERROR",
    "code": "INVALID_REQUEST_BODY",
    "userMessage": "Some fields in your request are invalid.",
    "developerMessage": "Request validation failed. Check validationErrors for details.",
    "validationErrors": [
      {
        "field": "Age",
        "message": "Age must be a number, but received a string",
        "rejectedValue": "22"
      }
    ],
    "retryable": false
  }
}
```

## Build

```bash
# Build binary
make build
# OR
go build -o bin/my-project ./cmd/my-project

# Run binary
./bin/my-project
```

## Testing

```bash
# Run all tests
make test
# OR
go test ./... -v

# Run tests with coverage
make test-coverage
# OR
go test ./... -cover
```

## Available Make Commands

```bash
make run            # Run the application
make dev            # Run with auto-reload (requires air)
make test           # Run all tests
make test-coverage  # Generate coverage report
make build          # Build binary to bin/
make clean          # Remove build artifacts
make deps           # Download dependencies
make fmt            # Format code
make lint           # Run linter
make help           # Show all commands
```

## Project Structure

```
my-project/
├── cmd/my-project/          # Application entry point
├── internal/
│   ├── app/                # Application initialization
│   ├── config/             # Configuration management
│   ├── constants/          # Error codes and constants
│   ├── handlers/           # HTTP request handlers
│   ├── middleware/         # HTTP middleware
│   ├── server/             # HTTP server setup
│   ├── types/              # Data models and types
│   └── utils/              # Utility functions
├── config/                 # Configuration files
└── storage/                # Database/storage files
```

## Dependencies

- Go 1.25.0+
- github.com/google/uuid - UUID generation
- github.com/ilyakaznacheev/cleanenv - Configuration management
- github.com/joho/godotenv - Environment file loading

## About

This is a learning project to demonstrate best practices in Go API development. Feel free to use it as a reference for your own projects.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Vikas Singh** ([@xanderbilla](https://github.com/xanderbilla))

[xanderbilla.com](https://xanderbilla.com)
