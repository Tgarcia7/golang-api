# Go Learning API

A beginner-friendly REST API built with Go (Golang) to demonstrate modern web development practices.

## Prerequisites

Before running this project, you need:

1. **Go 1.21 or higher** installed on your system
   - Download from: https://go.dev/dl/
   - Verify installation: `go version`

2. **curl** (for testing) - usually pre-installed on Mac/Linux
   - Windows users can use PowerShell or install curl

3. **jq** (optional, for pretty JSON output)
   - Mac: `brew install jq`
   - Ubuntu: `sudo apt install jq`
   - Or just remove `| jq .` from examples below

## Quick Start

### 1. Clone and Navigate
```bash
git clone <your-repo-url>
cd golang-api
```

### 2. Run the Server
```bash
# Option A: Run directly (for development)
go run .

# Option B: Build and run (recommended)
go build -o api .
./api
```

### 3. Verify It's Running
Open another terminal and run:
```bash
curl http://localhost:3000/health
```

You should see:
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "status": "ok",
    "version": "1.0.0"
  }
}
```

## Environment Configuration

The API can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `:3000` | Server port (include the colon) |

### Examples:
```bash
# Run on default port 3000
./api

# Run on port 8080
PORT=:8080 ./api

# Run on port 5000 (alternative syntax)
export PORT=:5000
./api
```

## API Endpoints

### Public Endpoints (No Authentication Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check for monitoring |
| GET | `/api` | API information and available endpoints |
| GET | `/api/users` | List all users |
| GET | `/api/user?id=<id>` | Get a specific user by ID |

### Protected Endpoints (Authentication Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/users` | Create a new user |
| PUT | `/api/user?id=<id>` | Update an existing user |
| DELETE | `/api/user?id=<id>` | Delete a user |

## Authentication

Protected endpoints require a Bearer token in the Authorization header.

**For testing/learning, use:**
```
Authorization: Bearer demo-token
```

Or any token starting with `valid-`:
```
Authorization: Bearer valid-mytoken123
```

## Usage Examples

### 1. Get API Information
```bash
curl http://localhost:3000/api | jq .
```

### 2. Create a User (requires auth)
```bash
curl -X POST http://localhost:3000/api/users \
  -H "Authorization: Bearer demo-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "555-1234"
  }' | jq .
```

Response:
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "555-1234"
  }
}
```

### 3. List All Users
```bash
curl http://localhost:3000/api/users | jq .
```

### 4. Get a Specific User
```bash
curl "http://localhost:3000/api/user?id=1" | jq .
```

### 5. Update a User (requires auth)
```bash
curl -X PUT "http://localhost:3000/api/user?id=1" \
  -H "Authorization: Bearer demo-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Updated",
    "email": "john.updated@example.com",
    "phone": "555-9999"
  }' | jq .
```

### 6. Delete a User (requires auth)
```bash
curl -X DELETE "http://localhost:3000/api/user?id=1" \
  -H "Authorization: Bearer demo-token"
```
(Returns 204 No Content on success)

### 7. Test Error Handling
```bash
# Missing authentication
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Test"}' | jq .

# User not found
curl "http://localhost:3000/api/user?id=999" | jq .

# Invalid endpoint
curl http://localhost:3000/invalid | jq .

# Invalid JSON
curl -X POST http://localhost:3000/api/users \
  -H "Authorization: Bearer demo-token" \
  -H "Content-Type: application/json" \
  -d 'not json' | jq .
```

## Project Structure

```
golang-api/
├── main.go        # Entry point, server startup, graceful shutdown
├── server.go      # Server configuration and HTTP setup
├── router.go      # Custom HTTP router implementation
├── handlers.go    # API endpoint handlers (controllers)
├── middleware.go  # Request middleware (auth, logging, etc.)
├── types.go       # Data models and custom types
├── response.go    # JSON response helper functions
├── go.mod         # Go module definition
└── README.md      # This file
```

## Key Go Concepts Demonstrated

Each file contains detailed comments explaining:

- **Structs and Methods** - Object-like patterns in Go
- **Interfaces** - How Go achieves polymorphism
- **Pointers** - Why and when to use them
- **Error Handling** - Go's explicit error handling approach
- **Goroutines** - Lightweight concurrency
- **Channels** - Communication between goroutines
- **Closures** - Functions that capture variables
- **Middleware Pattern** - Decorator/wrapper pattern for HTTP
- **JSON Serialization** - Converting Go structs to/from JSON
- **Context** - Request cancellation and deadlines

## Stopping the Server

Press `Ctrl+C` in the terminal where the server is running. The server will:
1. Stop accepting new connections
2. Wait for existing requests to complete (up to 30 seconds)
3. Shut down gracefully

## Common Issues

### Port Already in Use
```
Error: listen tcp :3000: bind: address already in use
```
**Solution:** Use a different port:
```bash
PORT=:3001 ./api
```

### Permission Denied (Linux/Mac)
```bash
chmod +x ./api
./api
```

### Go Not Found
Make sure Go is installed and in your PATH:
```bash
export PATH=$PATH:/usr/local/go/bin
go version
```

## Next Steps for Learning

1. Read through each `.go` file - comments explain Go concepts
2. Try modifying handlers to add new functionality
3. Add new middleware (e.g., request ID tracking)
4. Implement data persistence (file or database)
5. Add input validation with regex
6. Explore popular Go web frameworks (Gin, Echo, Chi)

## License

MIT License - Feel free to use for learning!
