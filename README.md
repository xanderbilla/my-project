# Simple User API (Go)

A basic REST API built using Go’s `net/http`.

## Endpoints

* `GET /api/users` – Returns a simple message
* `POST /api/users` – Creates a new user (JSON required)

## User JSON Example

```json
{
  "name": "Aman",
  "email": "aman@example.com",
  "age": 22
}
```

## Run the Project

```bash
go run cmd/main.go
```

## Test with curl

```bash
curl http://localhost:8080/api/users
```

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
```

## Author

[Xander Billa](https://github.com/xanderbilla)
