# Router Wrapper Usage Guide

The `pkg/router` package provides a clean wrapper over `go-chi/chi` router for consistent routing across microservices.

## Design Philosophy

This wrapper is designed to:
- **Provide a clean abstraction** over Chi router while maintaining its performance
- **Be library-ready** for extraction into a standalone package
- **Maintain flexibility** with escape hatches for Chi-specific features
- **Support configuration** through functional options pattern
- **Enable easy testing** with mockable interfaces

## Basic Usage

### Creating a Router

```go
import "github.com/himtar/go-boilerplate/pkg/router"

// Simple router without any middleware
r := router.New()

// Router with default middleware (recommended for production)
r := router.New(router.WithDefaults())

// Router with custom configuration
r := router.New(
    router.WithDefaults(),
    router.WithTimeout(30 * time.Second),
    router.WithoutLogger(),
)
```

### Registering Routes

```go
// HTTP methods
r.GET("/users", getUsersHandler)
r.POST("/users", createUserHandler)
r.PUT("/users/:id", updateUserHandler)
r.DELETE("/users/:id", deleteUserHandler)
r.PATCH("/users/:id", patchUserHandler)
r.HEAD("/users/:id", headUserHandler)
r.OPTIONS("/users/:id", optionsUserHandler)

// Generic handler (any method)
r.Handle("/custom", customHandler)
```

### URL Parameters

```go
r.GET("/users/:userID", func(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "userID")
    // Process userID
})

// Multiple parameters
r.GET("/users/:userID/posts/:postID", func(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "userID")
    postID := chi.URLParam(r, "postID")
    // Process both IDs
})
```

### Middleware

```go
// Apply middleware to all routes
r.Use(router.LoggerMiddleware())
r.Use(router.RecovererMiddleware())

// Apply multiple middleware at once
r.Use(
    router.RequestIDMiddleware(),
    router.RealIPMiddleware(),
    router.LoggerMiddleware(),
)

// Custom middleware
r.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Pre-processing
        log.Println("Before request")
        
        next.ServeHTTP(w, r)
        
        // Post-processing
        log.Println("After request")
    })
})
```

### Route Groups

```go
// Group with shared middleware
r.Group(func(r *router.HTTPRouter) {
    r.Use(authMiddleware)
    
    r.GET("/profile", getProfileHandler)
    r.PUT("/profile", updateProfileHandler)
    r.DELETE("/account", deleteAccountHandler)
})

// Named route group (sub-path)
r.Route("/api/v1", func(r *router.HTTPRouter) {
    r.GET("/users", getUsersHandler)
    r.POST("/users", createUserHandler)
    
    // Nested groups
    r.Route("/admin", func(r *router.HTTPRouter) {
        r.Use(adminMiddleware)
        r.GET("/stats", getStatsHandler)
    })
})
```

### Mounting Sub-routers

Perfect for modular route organization:

```go
// Define sub-router in separate file/package
func UserRouter() http.Handler {
    r := router.New()
    
    r.GET("/", listUsersHandler)
    r.POST("/", createUserHandler)
    r.GET("/:id", getUserHandler)
    r.PUT("/:id", updateUserHandler)
    r.DELETE("/:id", deleteUserHandler)
    
    return r.Handler()
}

// Mount in main router
func main() {
    mainRouter := router.New(router.WithDefaults())
    
    mainRouter.Mount("/users", UserRouter())
    mainRouter.Mount("/posts", PostRouter())
    mainRouter.Mount("/comments", CommentRouter())
}
```

## Configuration Options

### WithDefaults()

Applies recommended middleware stack in order:
1. **Request ID** - Adds unique ID to each request
2. **Real IP** - Extracts client's real IP from headers
3. **Logger** - Logs all requests with status and duration
4. **Recoverer** - Recovers from panics, returns 500
5. **Timeout** - 60 second request timeout

```go
r := router.New(router.WithDefaults())
```

### WithTimeout(duration)

Sets custom request timeout:

```go
r := router.New(router.WithTimeout(30 * time.Second))
```

### WithoutLogger()

Disables request logging:

```go
r := router.New(router.WithDefaults(), router.WithoutLogger())
```

### WithoutRecovery()

Disables panic recovery (useful for debugging):

```go
r := router.New(router.WithDefaults(), router.WithoutRecovery())
```

### WithoutRequestID()

Disables request ID generation:

```go
r := router.New(router.WithDefaults(), router.WithoutRequestID())
```

### WithoutRealIP()

Disables real IP extraction:

```go
r := router.New(router.WithDefaults(), router.WithoutRealIP())
```

## Built-in Middleware

### Chi Middleware Wrappers

```go
// Request tracking
router.RequestIDMiddleware()          // Unique ID per request
router.RealIPMiddleware()             // Extract real client IP

// Logging and monitoring
router.LoggerMiddleware()             // Request/response logging

// Error handling
router.RecovererMiddleware()          // Panic recovery

// Performance
router.TimeoutMiddleware(60 * time.Second)  // Request timeout
router.CompressMiddleware(5)                // Gzip compression (level 1-9)
router.RateLimitMiddleware(100)             // Throttle requests/sec

// Caching
router.NoCacheMiddleware()            // Disable client caching

// URL handling
router.StripSlashesMiddleware()       // Remove trailing slashes
router.RedirectSlashesMiddleware()    // Redirect trailing slashes
```

### Custom Middleware

```go
// CORS support
r.Use(router.CORSMiddleware(
    "*",                          // Allowed origins
    "GET,POST,PUT,DELETE",       // Allowed methods
    "Content-Type,Authorization", // Allowed headers
))

// Enforce JSON content type for POST/PUT/PATCH
r.Use(router.JSONContentTypeMiddleware())

// Set response content type
r.Use(router.SetContentTypeMiddleware("application/json"))
```

## Advanced Usage

### Direct Chi Access (Escape Hatch)

When you need Chi-specific features not exposed by the wrapper:

```go
r := router.New()

// Get underlying Chi router
chiRouter := r.Chi()

// Use Chi-specific features
chiRouter.Use(middleware.Heartbeat("/healthz"))
chiRouter.NotFound(custom404Handler)
chiRouter.MethodNotAllowed(custom405Handler)
```

### Complete Example

```go
package main

import (
    "net/http"
    "time"
    
    "github.com/himtar/go-boilerplate/pkg/router"
)

func main() {
    // Create router with defaults
    r := router.New(
        router.WithDefaults(),
        router.WithTimeout(30 * time.Second),
    )
    
    // Global middleware
    r.Use(router.CORSMiddleware("*", "GET,POST,PUT,DELETE", "Content-Type"))
    
    // Public routes
    r.GET("/", homeHandler)
    r.GET("/health", healthCheckHandler)
    
    // API v1
    r.Route("/api/v1", func(r *router.HTTPRouter) {
        // Auth routes (public)
        r.POST("/login", loginHandler)
        r.POST("/register", registerHandler)
        
        // Protected routes
        r.Group(func(r *router.HTTPRouter) {
            r.Use(authMiddleware)
            
            r.GET("/profile", getProfileHandler)
            r.PUT("/profile", updateProfileHandler)
            
            // Mount sub-routers
            r.Mount("/users", UserRouter())
            r.Mount("/posts", PostRouter())
        })
    })
    
    // Start server
    http.ListenAndServe(":8000", r.Handler())
}
```

## Testing

### Unit Testing Handlers

```go
func TestGetUserHandler(t *testing.T) {
    // Create test router
    r := router.New()
    r.GET("/users/:id", getUserHandler)
    
    // Create test request
    req := httptest.NewRequest("GET", "/users/123", nil)
    rr := httptest.NewRecorder()
    
    // Execute request
    r.Handler().ServeHTTP(rr, req)
    
    // Assert response
    if rr.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", rr.Code)
    }
}
```

### Testing with Middleware

```go
func TestProtectedRoute(t *testing.T) {
    r := router.New()
    r.Use(authMiddleware)
    r.GET("/protected", protectedHandler)
    
    req := httptest.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer token123")
    rr := httptest.NewRecorder()
    
    r.Handler().ServeHTTP(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
}
```

## Migration from Direct Chi

### Before (Direct Chi Usage)

```go
import "github.com/go-chi/chi/v5"

mux := chi.NewRouter()
mux.Use(middleware.Logger)
mux.Get("/users", getUsersHandler)
```

### After (Using Wrapper)

```go
import "github.com/himtar/go-boilerplate/pkg/router"

r := router.New(router.WithDefaults())
r.GET("/users", getUsersHandler)
```

## Best Practices

1. **Always use `WithDefaults()`** for production routers
2. **Group related routes** using `Route()` for better organization
3. **Apply authentication middleware** at group level, not per-route
4. **Mount modular sub-routers** for clean separation of concerns
5. **Use escape hatch sparingly** - only for Chi-specific features
6. **Set appropriate timeouts** based on your use case
7. **Enable compression** for API responses to reduce bandwidth
8. **Use middleware composition** instead of monolithic middleware

## Future: Library Extraction

This package is designed to be extracted into a standalone library for use across microservices:

```bash
# Future usage
go get github.com/yourorg/gorouter

# In your microservices
import "github.com/yourorg/gorouter"

r := gorouter.New(gorouter.WithDefaults())
```

All microservices will share the same routing interface, middleware, and conventions.

## Comparison with Direct Chi

| Feature | Direct Chi | Router Wrapper | Notes |
|---------|-----------|----------------|-------|
| Performance | ✅ Fast | ✅ Fast | Zero overhead |
| Type Safety | ✅ Yes | ✅ Yes | Same types |
| Middleware | ✅ Full | ✅ Full | All Chi middleware available |
| Customization | ✅ Full | ✅ Full | Escape hatch available |
| Configuration | ❌ Manual | ✅ Options | Functional options pattern |
| Library Ready | ❌ No | ✅ Yes | Designed for extraction |
| Learning Curve | 📚 Chi docs | 📚 This doc | Wrapper adds convenience |

## FAQ

**Q: Does this wrapper affect performance?**  
A: No, it's a zero-cost abstraction. The wrapper directly delegates to Chi's methods.

**Q: Can I still use Chi middleware?**  
A: Yes! All Chi middleware work seamlessly. Use the escape hatch if needed.

**Q: Why not use Chi directly?**  
A: The wrapper provides:
- Consistent API across microservices
- Easy configuration with functional options
- Library extraction readiness
- Standardized middleware stack

**Q: What if I need a Chi feature not exposed?**  
A: Use `r.Chi()` to get the underlying router and access all Chi features.

**Q: Can I mix wrapper and direct Chi usage?**  
A: Yes, but it's recommended to use the wrapper consistently for maintainability.

## References

- [Chi Router Documentation](https://github.com/go-chi/chi)
- [Chi Middleware](https://github.com/go-chi/chi/tree/master/middleware)
- [HTTP Handler Testing](https://golang.org/pkg/net/http/httptest/)
