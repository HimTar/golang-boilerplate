package response

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SendJSON sends a JSON response with the specified status code
func SendJSON(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Success responses (2xx)

// SendSuccess sends a 200 OK response with data
func SendSuccess(w http.ResponseWriter, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	SendJSON(w, http.StatusOK, response)
}

// SendCreated sends a 201 Created response
func SendCreated(w http.ResponseWriter, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	SendJSON(w, http.StatusCreated, response)
}

// SendAccepted sends a 202 Accepted response
func SendAccepted(w http.ResponseWriter, message string) {
	response := Response{
		Success: true,
		Message: message,
	}
	SendJSON(w, http.StatusAccepted, response)
}

// SendNoContent sends a 204 No Content response
func SendNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Client error responses (4xx)

// SendBadRequest sends a 400 Bad Request response
func SendBadRequest(w http.ResponseWriter, message string, errors ...string) {
	if message == "" {
		message = "Bad Request"
	}

	var errorStr string
	if len(errors) > 0 {
		errorStr = errors[0]
	}

	response := Response{
		Success: false,
		Message: message,
		Error:   errorStr,
	}
	SendJSON(w, http.StatusBadRequest, response)
}

// SendUnauthorized sends a 401 Unauthorized response
func SendUnauthorized(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Unauthorized"
	}

	response := Response{
		Success: false,
		Message: message,
		Error:   "Authentication required",
	}
	SendJSON(w, http.StatusUnauthorized, response)
}

// SendForbidden sends a 403 Forbidden response
func SendForbidden(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Forbidden"
	}

	response := Response{
		Success: false,
		Message: message,
		Error:   "Access denied",
	}
	SendJSON(w, http.StatusForbidden, response)
}

// SendNotFound sends a 404 Not Found response
func SendNotFound(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Not Found"
	}

	response := Response{
		Success: false,
		Message: message,
		Error:   "Resource not found",
	}
	SendJSON(w, http.StatusNotFound, response)
}

// SendMethodNotAllowed sends a 405 Method Not Allowed response
func SendMethodNotAllowed(w http.ResponseWriter, allowedMethods ...string) {
	message := "Method Not Allowed"

	// Set Allow header if methods are provided
	if len(allowedMethods) > 0 {
		w.Header().Set("Allow", joinMethods(allowedMethods))
	}

	response := Response{
		Success: false,
		Message: message,
		Error:   "HTTP method not allowed for this endpoint",
	}
	SendJSON(w, http.StatusMethodNotAllowed, response)
}

// SendConflict sends a 409 Conflict response
func SendConflict(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Conflict"
	}

	response := Response{
		Success: false,
		Message: message,
		Error:   "Resource conflict",
	}
	SendJSON(w, http.StatusConflict, response)
}

// SendUnprocessableEntity sends a 422 Unprocessable Entity response
func SendUnprocessableEntity(w http.ResponseWriter, message string, validationErrors interface{}) {
	if message == "" {
		message = "Unprocessable Entity"
	}

	response := Response{
		Success: false,
		Message: message,
		Data:    validationErrors,
		Error:   "Validation failed",
	}
	SendJSON(w, http.StatusUnprocessableEntity, response)
}

// SendTooManyRequests sends a 429 Too Many Requests response
func SendTooManyRequests(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Too Many Requests"
	}

	response := Response{
		Success: false,
		Message: message,
		Error:   "Rate limit exceeded",
	}
	SendJSON(w, http.StatusTooManyRequests, response)
}

// Server error responses (5xx)

// SendInternalServerError sends a 500 Internal Server Error response
func SendInternalServerError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Internal Server Error"
	}

	response := Response{
		Success: false,
		Message: message,
		Error:   "An internal server error occurred",
	}
	SendJSON(w, http.StatusInternalServerError, response)
}

// SendNotImplemented sends a 501 Not Implemented response
func SendNotImplemented(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Not Implemented"
	}

	response := Response{
		Success: false,
		Message: message,
		Error:   "Feature not implemented",
	}
	SendJSON(w, http.StatusNotImplemented, response)
}

// SendBadGateway sends a 502 Bad Gateway response
func SendBadGateway(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Bad Gateway"
	}

	response := Response{
		Success: false,
		Message: message,
		Error:   "Bad gateway",
	}
	SendJSON(w, http.StatusBadGateway, response)
}

// SendServiceUnavailable sends a 503 Service Unavailable response
func SendServiceUnavailable(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Service Unavailable"
	}

	response := Response{
		Success: false,
		Message: message,
		Error:   "Service temporarily unavailable",
	}
	SendJSON(w, http.StatusServiceUnavailable, response)
}

// SendGatewayTimeout sends a 504 Gateway Timeout response
func SendGatewayTimeout(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Gateway Timeout"
	}

	response := Response{
		Success: false,
		Message: message,
		Error:   "Gateway timeout",
	}
	SendJSON(w, http.StatusGatewayTimeout, response)
}

// Utility functions

// SendCustom sends a custom response with any status code
func SendCustom(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}, errorMsg string) {
	response := Response{
		Success: success,
		Message: message,
		Data:    data,
		Error:   errorMsg,
	}
	SendJSON(w, statusCode, response)
}

// Helper function to join allowed methods
func joinMethods(methods []string) string {
	result := ""
	for i, method := range methods {
		if i > 0 {
			result += ", "
		}
		result += method
	}
	return result
}
