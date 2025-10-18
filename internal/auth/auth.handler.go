package auth

import (
	"net/http"

	"github.com/himtar/go-boilerplate/pkg/response"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Example of successful login response
	userData := map[string]interface{}{
		"user_id": 123,
		"email":   "user@example.com",
		"token":   "jwt-token-here",
	}

	response.SendSuccess(w, "Login successful", userData)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Example of registration with validation error
	if r.URL.Query().Get("error") == "validation" {
		validationErrors := map[string]string{
			"email":    "Email is required",
			"password": "Password must be at least 8 characters",
		}
		response.SendUnprocessableEntity(w, "Validation failed", validationErrors)
		return
	}

	// Example of successful registration
	userData := map[string]interface{}{
		"user_id": 124,
		"email":   "newuser@example.com",
		"message": "Please verify your email",
	}

	response.SendCreated(w, "User registered successfully", userData)
}
