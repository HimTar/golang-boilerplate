package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/himtar/go-boilerplate/pkg/logger"
)

// HTTPRouter is the main router interface
type HTTPRouter struct {
	chi    chi.Router
	logger *logger.Logger
}

// New creates a new router with provided options
func New() *HTTPRouter {
	mux := chi.NewRouter()

	return &HTTPRouter{chi: mux}
}

// HTTP Methods

// GET registers a GET route
func (r *HTTPRouter) GET(pattern string, handler http.HandlerFunc) {
	r.chi.Get(pattern, handler)
}

// POST registers a POST route
func (r *HTTPRouter) POST(pattern string, handler http.HandlerFunc) {
	r.chi.Post(pattern, handler)
}

// PUT registers a PUT route
func (r *HTTPRouter) PUT(pattern string, handler http.HandlerFunc) {
	r.chi.Put(pattern, handler)
}

// DELETE registers a DELETE route
func (r *HTTPRouter) DELETE(pattern string, handler http.HandlerFunc) {
	r.chi.Delete(pattern, handler)
}

// PATCH registers a PATCH route
func (r *HTTPRouter) PATCH(pattern string, handler http.HandlerFunc) {
	r.chi.Patch(pattern, handler)
}

// HEAD registers a HEAD route
func (r *HTTPRouter) HEAD(pattern string, handler http.HandlerFunc) {
	r.chi.Head(pattern, handler)
}

// OPTIONS registers an OPTIONS route
func (r *HTTPRouter) OPTIONS(pattern string, handler http.HandlerFunc) {
	r.chi.Options(pattern, handler)
}

// Middleware management

// Use appends one or more middlewares to the router
func (r *HTTPRouter) Use(middlewares ...Middleware) {
	for _, mw := range middlewares {
		r.chi.Use(mw)
	}
}

// Grouping and mounting

// Group creates a new inline-Middleware Router that inherits all middleware from parent router
func (r *HTTPRouter) Group(fn func(r *HTTPRouter)) {
	r.chi.Group(func(mux chi.Router) {
		fn(&HTTPRouter{chi: mux})
	})
}

// Route mounts a sub-Router along a `pattern` string
func (r *HTTPRouter) Route(pattern string, fn func(r *HTTPRouter)) {
	r.chi.Route(pattern, func(mux chi.Router) {
		fn(&HTTPRouter{chi: mux})
	})
}

// Mount attaches another http.Handler along `pattern`
func (r *HTTPRouter) Mount(pattern string, handler http.Handler) {
	r.chi.Mount(pattern, handler)
}

// Handle registers a generic http.Handler
func (r *HTTPRouter) Handle(pattern string, handler http.Handler) {
	r.chi.Handle(pattern, handler)
}

// Handler returns the underlying http.Handler for server startup
func (r *HTTPRouter) Handler() http.Handler {
	return r.chi
}

// Chi returns direct access to underlying Chi router (escape hatch)
// Use this when you need Chi-specific features not exposed by the wrapper
func (r *HTTPRouter) Chi() chi.Router {
	return r.chi
}

// Type aliases for convenience
type (
	// Middleware is a function that wraps an http.Handler
	Middleware func(http.Handler) http.Handler

	// HandlerFunc is an alias for http.HandlerFunc
	HandlerFunc = http.HandlerFunc
)
