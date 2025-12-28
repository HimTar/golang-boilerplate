package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/himtar/go-boilerplate/pkg/logger"
)

// Common middleware wrappers for easy access

// RequestIDMiddleware adds unique request ID to each request
func RequestIDMiddleware() Middleware {
	return middleware.RequestID
}

// RealIPMiddleware extracts real client IP from headers
func RealIPMiddleware() Middleware {
	return middleware.RealIP
}

func LoggerMiddleware(log logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			// Extract real IP from X-Forwarded-For or X-Real-IP headers, fallback to RemoteAddr
			realIP := r.Header.Get("X-Forwarded-For")
			if realIP == "" {
				realIP = r.Header.Get("X-Real-IP")
			}
			if realIP == "" {
				realIP = r.RemoteAddr
			}

			// Extract request ID if present
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = r.Header.Get("X-Request-ID")
			}

			userAgent := r.UserAgent()

			// Log as JSON
			log.InfoJSON(map[string]interface{}{
				"event":     "http_request",
				"method":    r.Method,
				"path":      r.URL.Path,
				"status":    ww.Status(),
				"size":      ww.BytesWritten(),
				"duration":  duration.String(),
				"ip":        realIP,
				"requestId": requestID,
				"userAgent": userAgent,
			})
		})
	}
}

// RecovererMiddleware recovers from panics and returns 500 Internal Server Error
func RecovererMiddleware() Middleware {
	return middleware.Recoverer
}

// TimeoutMiddleware sets request timeout duration
func TimeoutMiddleware(duration time.Duration) Middleware {
	return middleware.Timeout(duration)
}

// NoCacheMiddleware disables client-side caching
func NoCacheMiddleware() Middleware {
	return middleware.NoCache
}

// CompressMiddleware enables response compression with specified level (0-9)
func CompressMiddleware(level int) Middleware {
	return middleware.Compress(level)
}

// StripSlashesMiddleware removes trailing slashes from request paths
func StripSlashesMiddleware() Middleware {
	return middleware.StripSlashes
}

// RedirectSlashesMiddleware redirects paths with trailing slashes
func RedirectSlashesMiddleware() Middleware {
	return middleware.RedirectSlashes
}

// Custom middleware implementations

// CORSMiddleware adds CORS headers to responses
// For production, consider using a dedicated CORS library like github.com/go-chi/cors
func CORSMiddleware(allowedOrigins, allowedMethods, allowedHeaders string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
			w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ContentTypeMiddleware enforces specific content type for requests
func ContentTypeMiddleware(contentType string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip check for GET requests
			if r.Method == http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			if r.Header.Get("Content-Type") != contentType {
				http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// JSONContentTypeMiddleware enforces application/json content type
func JSONContentTypeMiddleware() Middleware {
	return ContentTypeMiddleware("application/json")
}

// SetContentTypeMiddleware sets content type for all responses
func SetContentTypeMiddleware(contentType string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware provides basic rate limiting
// For production, consider using a dedicated rate limiting library
func RateLimitMiddleware(requestsPerSecond int) Middleware {
	return middleware.Throttle(requestsPerSecond)
}
