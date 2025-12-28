# Logger Package

Production-grade structured JSON logger for microservices with Kibana/ELK integration.

## Features

✅ **JSON Formatting** - All logs in structured JSON for Kibana ingestion  
✅ **File Output** - Logs written to files with rotation support  
✅ **Context Propagation** - Thread-safe logger cloning with additional fields  
✅ **Trace ID Support** - Distributed tracing integration  
✅ **Multiple Outputs** - Console and file output simultaneously  
✅ **Thread-Safe** - Safe for concurrent use across goroutines  
✅ **Zero Allocation Cloning** - Efficient context propagation  

## Quick Start

### Basic Usage

```go
import "github.com/himtar/go-boilerplate/pkg/logger"

// Create logger with default config
log, err := logger.NewJSONLogger(nil)
if err != nil {
    panic(err)
}

// Log messages
log.Info("Server started")
log.Error("Database connection failed: %v", err)
log.Warn("High memory usage: %d MB", memUsage)
```

### Configuration

```go
config := &logger.Config{
    LogFilePath:   "logs/myservice.log",
    ServiceName:   "api-gateway",
    EnableConsole: true,  // Log to stdout
    EnableFile:    true,  // Log to file
}

log, err := logger.NewJSONLogger(config)
```

### Context Propagation (Thread-Safe)

```go
// Base logger
baseLogger, _ := logger.NewJSONLogger(nil)

// Add context fields - returns NEW logger instance
requestLogger := baseLogger.WithContext(map[string]interface{}{
    "user_id": "12345",
    "ip": "192.168.1.1",
    "method": "POST",
})

// Add trace ID - returns NEW logger instance
tracedLogger := requestLogger.WithTraceID("trace-abc-123")

// Single field - returns NEW logger instance
orderLogger := tracedLogger.WithField("order_id", "ORDER-001")

// Each logger is independent and thread-safe
go func() {
    orderLogger.Info("Processing order")  
    // Output: {"timestamp":"2024-01-01T10:00:00Z","level":"INFO","message":"Processing order","context":{"user_id":"12345","ip":"192.168.1.1","method":"POST","order_id":"ORDER-001"},"trace_id":"trace-abc-123","service":"go-service"}
}()
```

## HTTP Request Example

```go
func HandleRequest(w http.ResponseWriter, r *http.Request, baseLogger logger.Logger) {
    // Create request-scoped logger
    reqLogger := baseLogger.
        WithTraceID(r.Header.Get("X-Trace-ID")).
        WithContext(map[string]interface{}{
            "method": r.Method,
            "path": r.URL.Path,
            "remote_addr": r.RemoteAddr,
        })

    reqLogger.Info("Request received")

    // Pass to business logic
    if err := ProcessOrder(r, reqLogger); err != nil {
        reqLogger.Error("Order processing failed: %v", err)
        http.Error(w, "Internal error", 500)
        return
    }

    reqLogger.Info("Request completed successfully")
}

func ProcessOrder(r *http.Request, log logger.Logger) error {
    orderID := extractOrderID(r)
    
    // Add order-specific context
    orderLogger := log.WithField("order_id", orderID)
    
    orderLogger.Info("Processing order")
    
    if err := validateOrder(orderID); err != nil {
        orderLogger.Error("Validation failed: %v", err)
        return err
    }
    
    orderLogger.Info("Order validated successfully")
    return nil
}
```

## Log Output Format

### Console Output
```json
{"timestamp":"2024-12-27T10:30:00Z","level":"INFO","message":"Server started","service":"api-gateway"}
{"timestamp":"2024-12-27T10:30:15Z","level":"INFO","message":"Request received","context":{"method":"POST","path":"/api/orders","user_id":"12345"},"trace_id":"trace-abc-123","service":"api-gateway"}
{"timestamp":"2024-12-27T10:30:16Z","level":"ERROR","message":"Database connection failed","context":{"host":"db.example.com","port":5432},"trace_id":"trace-abc-123","service":"api-gateway"}
```

### File Output
Same JSON format written to `logs/app.log` (or configured path)

## Kibana Integration

The JSON format is optimized for ELK stack:

```
filebeat.yml:
  - type: log
    enabled: true
    paths:
      - /path/to/logs/*.log
    json.keys_under_root: true
    json.add_error_key: true
```

**Kibana Query Examples:**
```
level: "ERROR"
trace_id: "trace-abc-123"
context.user_id: "12345"
service: "api-gateway" AND level: "ERROR"
```

## Log Levels

```go
log.Debug("Debug information: %v", data)   // Level: DEBUG
log.Info("Operation completed")             // Level: INFO
log.Warn("Deprecated API called")           // Level: WARN
log.Error("Operation failed: %v", err)      // Level: ERROR
log.Fatal("Cannot start server: %v", err)   // Level: FATAL (exits process)
```

## Thread Safety Guarantee

```go
baseLogger, _ := logger.NewJSONLogger(nil)

// Spawn 1000 concurrent goroutines
for i := 0; i < 1000; i++ {
    go func(id int) {
        // Each goroutine gets its own logger instance
        routineLogger := baseLogger.WithField("goroutine_id", id)
        
        routineLogger.Info("Processing task")
        // No race conditions - each logger is independent
    }(i)
}
```

**How it works:**
- `WithContext()`, `WithTraceID()`, `WithField()` create **new cloned instances**
- Original logger remains unchanged
- No shared mutable state between clones
- Safe to pass logger instances across goroutines

## Advanced Usage

### Middleware Integration

```go
func LoggingMiddleware(baseLogger logger.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Generate trace ID if not present
            traceID := r.Header.Get("X-Trace-ID")
            if traceID == "" {
                traceID = generateTraceID()
            }
            
            // Create request logger
            reqLogger := baseLogger.
                WithTraceID(traceID).
                WithContext(map[string]interface{}{
                    "method": r.Method,
                    "path": r.URL.Path,
                    "user_agent": r.UserAgent(),
                })
            
            // Add to request context
            ctx := context.WithValue(r.Context(), "logger", reqLogger)
            r = r.WithContext(ctx)
            
            reqLogger.Info("Request started")
            
            next.ServeHTTP(w, r)
            
            duration := time.Since(start)
            reqLogger.WithField("duration_ms", duration.Milliseconds()).
                Info("Request completed")
        })
    }
}

// Extract logger from context
func GetLogger(r *http.Request) logger.Logger {
    if log, ok := r.Context().Value("logger").(logger.Logger); ok {
        return log
    }
    return logger.NewDefaultLogger()
}
```

### Structured Error Logging

```go
func HandleDatabaseError(err error, log logger.Logger) {
    errorLogger := log.WithContext(map[string]interface{}{
        "error_type": fmt.Sprintf("%T", err),
        "error_message": err.Error(),
        "stack_trace": debug.Stack(),
    })
    
    errorLogger.Error("Database operation failed")
}
```

### Performance Monitoring

```go
func MeasurePerformance(operation string, log logger.Logger) func() {
    start := time.Now()
    return func() {
        duration := time.Since(start)
        log.WithContext(map[string]interface{}{
            "operation": operation,
            "duration_ms": duration.Milliseconds(),
        }).Info("Performance metric")
    }
}

// Usage
defer MeasurePerformance("database_query", log)()
// ... expensive operation
```

## Configuration Best Practices

### Development
```go
config := &logger.Config{
    LogFilePath:   "logs/dev.log",
    ServiceName:   "myservice-dev",
    EnableConsole: true,   // See logs in terminal
    EnableFile:    false,  // Don't clutter filesystem
}
```

### Production
```go
config := &logger.Config{
    LogFilePath:   "/var/log/myservice/app.log",
    ServiceName:   "myservice-prod",
    EnableConsole: false,  // Reduce stdout noise
    EnableFile:    true,   // Persist for analysis
}
```

### Testing
```go
config := &logger.Config{
    LogFilePath:   "/tmp/test.log",
    ServiceName:   "test",
    EnableConsole: false,
    EnableFile:    true,
}
```

## Log Rotation

The logger creates/appends to files but doesn't handle rotation. Use external tools:

### Using logrotate (Linux)
```
/var/log/myservice/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    postrotate
        systemctl reload myservice
    endscript
}
```

### Using lumberjack (Go)
```go
import "gopkg.in/natefinch/lumberjack.v2"

file := &lumberjack.Logger{
    Filename:   config.LogFilePath,
    MaxSize:    100,  // MB
    MaxBackups: 3,
    MaxAge:     28,   // days
    Compress:   true,
}

// Use as writer in logger creation
```

## Migration Guide

### From Standard log Package
```go
// Before
log.Println("Server started")

// After
logger, _ := logger.NewJSONLogger(nil)
logger.Info("Server started")
```

### From Old Logger Interface
The package maintains backward compatibility with the previous `Logger` interface, so existing code continues to work.

## Design Decisions

1. **JSON Format**: Industry standard for log aggregation systems
2. **Cloning Pattern**: Ensures thread safety without locks on every log call
3. **Interface-based**: Easy to mock for testing
4. **Multiple Writers**: Support console + file + custom outputs
5. **Context Immutability**: Prevent accidental state mutation
6. **UTC Timestamps**: Consistent across time zones
7. **RFC3339 Format**: Sortable and parseable

## Performance Characteristics

- **Clone overhead**: O(n) where n = number of context fields
- **Log write**: O(1) - direct JSON marshaling
- **Memory**: Each clone allocates new context map
- **Concurrency**: Lock-free for writes, read lock only on clone

## Future Enhancements

- [ ] Log level filtering
- [ ] Sampling (log 1/N messages)
- [ ] Async writing with buffering
- [ ] Custom formatters
- [ ] Hooks for external systems (Sentry, DataDog)
- [ ] Metrics integration (log rates, error counts)
