# Response Helper Library Usage

The response helper library provides a consistent way to send HTTP responses with proper status codes and JSON structure.

## Response Structure

All responses follow this structure:

```json
{
  "success": true/false,
  "message": "Human readable message",
  "data": {}, // Optional, only for successful responses
  "error": "Error description" // Optional, only for error responses
}
```

## Usage Examples

### Success Responses (2xx)

```go
import "github.com/himtar/go-boilerplate/pkg/helpers"

// 200 OK - Send data
userData := map[string]interface{}{
    "id": 123,
    "name": "John Doe",
}
helpers.SendSuccess(w, "User retrieved successfully", userData)

// 201 Created - Resource created
helpers.SendCreated(w, "User created successfully", userData)

// 202 Accepted - Request accepted for processing
helpers.SendAccepted(w, "Request accepted for processing")

// 204 No Content - Success with no response body
helpers.SendNoContent(w)
```

### Client Error Responses (4xx)

```go
// 400 Bad Request
helpers.SendBadRequest(w, "Invalid input data", "Email is required")

// 401 Unauthorized
helpers.SendUnauthorized(w, "Please login to access this resource")

// 403 Forbidden
helpers.SendForbidden(w, "You don't have permission to access this resource")

// 404 Not Found
helpers.SendNotFound(w, "User not found")

// 405 Method Not Allowed
helpers.SendMethodNotAllowed(w, "GET", "POST") // Allowed methods

// 409 Conflict
helpers.SendConflict(w, "Email already exists")

// 422 Unprocessable Entity (Validation errors)
validationErrors := map[string]string{
    "email": "Email format is invalid",
    "password": "Password must be at least 8 characters",
}
helpers.SendUnprocessableEntity(w, "Validation failed", validationErrors)

// 429 Too Many Requests
helpers.SendTooManyRequests(w, "Rate limit exceeded. Try again later")
```

### Server Error Responses (5xx)

```go
// 500 Internal Server Error
helpers.SendInternalServerError(w, "Database connection failed")

// 501 Not Implemented
helpers.SendNotImplemented(w, "This feature is not yet implemented")

// 502 Bad Gateway
helpers.SendBadGateway(w, "Upstream service is unavailable")

// 503 Service Unavailable
helpers.SendServiceUnavailable(w, "Service is temporarily down for maintenance")

// 504 Gateway Timeout
helpers.SendGatewayTimeout(w, "Upstream service timed out")
```

### Custom Responses

```go
// Send custom response with any status code
helpers.SendCustom(w, 418, false, "I'm a teapot", nil, "Cannot brew coffee")

// Send JSON response with custom structure
response := helpers.Response{
    Success: true,
    Message: "Custom response",
    Data: map[string]string{"custom": "data"},
}
helpers.SendJSON(w, http.StatusOK, response)
```

## Error Handling in Handlers

Here's how to structure your handlers with proper error handling:

```go
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    // Validate input
    if r.Header.Get("Content-Type") != "application/json" {
        helpers.SendBadRequest(w, "Content-Type must be application/json", "")
        return
    }
    
    // Parse request body
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        helpers.SendBadRequest(w, "Invalid JSON format", err.Error())
        return
    }
    
    // Validate fields
    if user.Email == "" {
        helpers.SendUnprocessableEntity(w, "Validation failed", map[string]string{
            "email": "Email is required",
        })
        return
    }
    
    // Check if user exists
    if userExists(user.Email) {
        helpers.SendConflict(w, "User with this email already exists")
        return
    }
    
    // Create user
    createdUser, err := createUser(user)
    if err != nil {
        helpers.SendInternalServerError(w, "Failed to create user")
        return
    }
    
    // Success response
    helpers.SendCreated(w, "User created successfully", createdUser)
}
```

## Migration from Old Error Package

If you're migrating from the old `pkg/errors` package:

```go
// Old way
errors.BadRequest(w, "Invalid input")
errors.MethodNotAllowed(w)
errors.InternalServerError(w, "Something went wrong")

// New way
helpers.SendBadRequest(w, "Invalid input", "")
helpers.SendMethodNotAllowed(w)
helpers.SendInternalServerError(w, "Something went wrong")
```

## Best Practices

1. **Consistent Messages**: Use clear, user-friendly messages
2. **Include Data**: Always include relevant data in success responses
3. **Validation Errors**: Use `SendUnprocessableEntity` for validation errors with details
4. **Security**: Don't expose sensitive error details to clients
5. **HTTP Status Codes**: Use appropriate status codes for different scenarios
6. **Error Context**: Include error context in server error responses for debugging