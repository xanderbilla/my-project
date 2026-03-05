# Contributing

This is a learning project to demonstrate Go API development patterns. Contributions and suggestions are welcome!

## How to Contribute

### Reporting Issues

Found a bug or have a suggestion? Open an issue with:

- Description of the issue or idea
- Steps to reproduce (for bugs)
- Your environment (Go version, OS)

### Pull Requests

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test your changes (`go test ./...`)
5. Commit with clear messages
6. Open a pull request

## Code Guidelines

### Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` to format your code
- Run `go vet` to check for common issues
- Add comments for exported functions and types
- Keep functions focused and single-purpose

### Comments

This project emphasizes educational comments. When adding new code:

- Follow standard Go conventions
- Use `gofmt` to format your code
- Add comments for clarity (this is a learning project)
- Keep code simple and readable

### Commit Messages

Use clear commit messages:

- Start with a verb (Add, Fix, Update, Remove)
- Keep it concise
- Example: "Add email validation to user handler"bash
  cp .env.example .env

````

3. Install dependencies:
```bash
go mod download
````

4. Run the application:

```bash
go run cmd/my-project/main.go
```

5. Run tests:

```bash
go test ./...
```

## Project Structure

Understanding the project structure helps when contributing:

````
my-project/
├── cmd/my-project/          # Entry point
├── internal/
│   ├── app/                # Application setup
│   ├── config/             # Configuration loading
│   ├── constants/          # Error codes and constants
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # HTTP middleware
│   ├── server/             # Server setup
│   ├── types/              # Data models
│  Getting Started

```bash
# Clone and setup
git clone https://github.com/xanderbilla/my-project.git
cd my-project
cp .env.example .env

# Run
go run cmd/my-project/main.go

# Test
go test ./...
````

## Learning Focus

This project demonstrates:

- Clean project structure
- Error handling patterns
- Middleware implementation
- Environment-based configuration
- Structured logging

## Questions?

Open an issue for discussion or questions about the code.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
