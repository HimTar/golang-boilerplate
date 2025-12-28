package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/himtar/go-boilerplate/internal/libraries/env"
	"github.com/himtar/go-boilerplate/pkg/logger"
)

type AuthModule struct {
	Config *env.Variables
	Logger logger.Logger
	// Add domain services, repositories, etc. here as fields for DDD
}

// NewAuthModule constructs the AuthModule and injects dependencies.
func NewAuthModule(config *env.Variables, log logger.Logger) *AuthModule {
	return &AuthModule{
		Config: config,
		Logger: log,
	}
}

// Router returns the chi.Router with all authentication routes registered.
func (a *AuthModule) Router() chi.Router {
	router := chi.NewRouter()

	router.Post("/login", a.loginHandler)
	router.Post("/register", a.registerHandler)

	return router
}

// loginHandler is the HTTP handler for the /login endpoint.
func (a *AuthModule) loginHandler(w http.ResponseWriter, r *http.Request) {
	a.Logger.Info("AuthModule: /login called")
	// TODO: Implement login logic using a.Config, a.Logger, etc.
}

// registerHandler is the HTTP handler for the /register endpoint.
func (a *AuthModule) registerHandler(w http.ResponseWriter, r *http.Request) {
	a.Logger.Info("AuthModule: /register called")
	// TODO: Implement registration logic using a.Config, a.Logger, etc.
}
