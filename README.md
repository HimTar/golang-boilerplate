# Go Boilerplate

A simple and clean Go web server boilerplate built with Chi router, custom middleware, and a structured logger.

## Features

- ✅ Clean project structure
- ✅ Environment variable management
- ✅ Custom router with method validation
- ✅ Chi router integration with middleware
- ✅ Hot reload development support with Air
- ✅ Comprehensive response helper library
- ✅ Consistent JSON API responses
- ✅ Error handling utilities
- ✅ Authentication handlers with JWT example
- ✅ CRUD operations example
- ✅ Windows PowerShell compatibility
- ✅ Structured logging with pluggable logger

## Project Structure

```
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── auth/            # Authentication handlers and routes
│   └── order/           # Order-related handlers
├── pkg/
│   ├── env/             # Environment configuration
│   ├── logger/          # Structured logger implementation
│   ├── middleware/      # HTTP middleware
│   ├── server/          # Server utilities, custom router and configuration
├── docs/                # Documentation
│   ├── response-helpers.md  # Response helper usage guide
│   ├── logger-usage.md      # Logger usage and configuration
│   └── router-wrapper.md    # Router wrapper and middleware
└── logs/                # Application logs
```

## Getting Started

### Prerequisites

- Go 1.21.5 or higher
- Air (for hot reload development)

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd golang-boilerplate
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

### Running the Application

#### Development (with hot reload):

**Linux/Mac:**
```bash
make dev
```

**Windows (PowerShell):**
```powershell
air
```

**Note for Windows users**: If `make` is not available, install it via:
- Chocolatey: `choco install make`
- Scoop: `scoop install make`
- winget: `winget install GnuWin32.Make`

#### Production:

**Linux/Mac:**
```bash
make run
```

**Windows (PowerShell):**
```powershell
go run ./cmd/server/main.go
```

The server will start on the port specified in your `.env` file (default: `:8000`).

## API Endpoints

### Authentication
- `GET /auth/login` - Login endpoint (returns JWT token example)
- `GET /auth/register` - Registration endpoint
- `GET /auth/register?error=validation` - Test validation errors

### Groups
- `GET /group/get` - Get groups with pagination
- `POST /group/create` - Create group
- `POST /group/create?duplicate=true` - Test conflict error
- `POST /group/create?invalid=true` - Test validation error

### Health Check
- `GET /` - Returns "Hello, World!"

## Testing API Endpoints

### Using PowerShell (Windows):
```powershell
# Test successful login
Invoke-WebRequest -Uri "http://localhost:8000/auth/login"

# Test registration
Invoke-WebRequest -Uri "http://localhost:8000/auth/register"

# Test validation errors
Invoke-WebRequest -Uri "http://localhost:8000/auth/register?error=validation"

# Test groups
Invoke-WebRequest -Uri "http://localhost:8000/group/get"
```

### Using curl (Linux/Mac or after installing real curl on Windows):
```bash
# Test successful login
curl "http://localhost:8000/auth/login"

# Test registration with validation error
curl "http://localhost:8000/auth/register?error=validation"

# Test groups
curl "http://localhost:8000/group/get"
```

## Configuration

The application uses environment variables for configuration. Create a `.env` file in the root directory:

```env
PORT=":8000"
DB="postgres"
DB_URI="your-database-uri"
ENV="development"
```

## Development

### Hot Reload

This project uses [Air](https://github.com/air-verse/air) for hot reload during development. The configuration is in [.air.toml](.air.toml).

To install Air:
```bash
go install github.com/air-verse/air@latest
```

**Note**: The Air project has moved from `cosmtrek/air` to `air-verse/air`.

### Adding New Routes

1. Create handler functions in `internal/[module]/` (e.g., `internal/auth/`)
2. Use the response helper library from `pkg/response/` for consistent API responses
3. Register routes using the custom router in `pkg/router/`
4. Mount the routes in [`cmd/server/main.go`](cmd/server/main.go)

Example handler with response helpers:
```go
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    // Validate input
    if err := validateInput(r); err != nil {
        response.SendBadRequest(w, "Invalid input", err.Error())
        return
    }
    
    // Create user
    user, err := createUser(userData)
    if err != nil {
        response.SendInternalServerError(w, "Failed to create user")
        return
    }
    
    // Success response
    response.SendCreated(w, "User created successfully", user)
}
```

## Response Handling

The project includes a comprehensive response helper library in [`pkg/response/response.go`](pkg/response/response.go):

### Success Responses (2xx)
- `SendSuccess(w, message, data)` - 200 OK
- `SendCreated(w, message, data)` - 201 Created
- `SendAccepted(w, message)` - 202 Accepted
- `SendNoContent(w)` - 204 No Content

### Client Error Responses (4xx)
- `SendBadRequest(w, message, error)` - 400 Bad Request
- `SendUnauthorized(w, message)` - 401 Unauthorized
- `SendForbidden(w, message)` - 403 Forbidden
- `SendNotFound(w, message)` - 404 Not Found
- `SendMethodNotAllowed(w, methods...)` - 405 Method Not Allowed
- `SendConflict(w, message)` - 409 Conflict
- `SendUnprocessableEntity(w, message, validationErrors)` - 422 Unprocessable Entity
- `SendTooManyRequests(w, message)` - 429 Too Many Requests

### Server Error Responses (5xx)
- `SendInternalServerError(w, message)` - 500 Internal Server Error
- `SendNotImplemented(w, message)` - 501 Not Implemented
- `SendBadGateway(w, message)` - 502 Bad Gateway
- `SendServiceUnavailable(w, message)` - 503 Service Unavailable
- `SendGatewayTimeout(w, message)` - 504 Gateway Timeout

### Response Format
All responses follow a consistent JSON structure:
```json
{
  "success": true/false,
  "message": "Human readable message",
  "data": {}, // Optional, for successful responses
  "error": "Error description" // Optional, for error responses
}
```

For detailed usage examples, see [`docs/response-helpers.md`](docs/response-helpers.md).

## Logging

The project uses a structured logger in [`pkg/logger/`](pkg/logger/) for consistent, pluggable logging across the application.  
See [`docs/logger-usage.md`](docs/logger-usage.md) for configuration and usage examples.

## Legacy Error Handling

The project also includes legacy error handling utilities in [`pkg/errors/errors.go`](pkg/errors/errors.go) for backward compatibility.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test your changes
5. Submit a pull request

## License

This project is licensed under the MIT License.