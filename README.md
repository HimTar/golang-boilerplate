# Go Boilerplate

A simple and clean Go web server boilerplate built with Chi router and custom middleware.

## Features

- ✅ Clean project structure
- ✅ Environment variable management
- ✅ Custom router with method validation
- ✅ Chi router integration with middleware
- ✅ Hot reload development support with Air
- ✅ Error handling utilities
- ✅ Authentication handlers
- ✅ Groups handlers

## Project Structure

```
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   └── handlers/        # HTTP handlers
├── libraries/
│   └── server/          # Server utilities and configuration
├── pkg/
│   ├── errors/          # Error handling utilities
│   ├── helpers/         # Helper functions
│   └── router/          # Custom router implementation
└── docs/                # Documentation
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
- `GET /auth/login` - Login endpoint
- `GET /auth/register` - Registration endpoint

### Groups
- `GET /group/get` - Get groups
- `POST /group/create` - Create group

### Health Check
- `GET /` - Returns "Hello, World!"

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

This project uses [Air](https://github.com/cosmtrek/air) for hot reload during development. The configuration is in [.air.toml](.air.toml).

To install Air:
```bash
go install github.com/cosmtrek/air@latest
```

### Adding New Routes

1. Create handler functions in `internal/handlers/`
2. Register routes using the custom router in [`pkg/router/router.go`](pkg/router/router.go)
3. Mount the routes in [`cmd/server/main.go`](cmd/server/main.go)

## Error Handling

The project includes custom error handling utilities in [`pkg/errors/errors.go`](pkg/errors/errors.go):

- `BadRequest(w, message)` - 400 Bad Request
- `MethodNotAllowed(w)` - 405 Method Not Allowed  
- `InternalServerError(w, message)` - 500 Internal Server Error

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test your changes
5. Submit a pull request

## License

This project is licensed under the MIT License.