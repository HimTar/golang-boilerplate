package errors

import "net/http"

func BadRequest(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Bad Request !"
	}

	http.Error(w, message, http.StatusBadRequest)
}

func MethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "Method Not Allowed !", http.StatusMethodNotAllowed)
}

func InternalServerError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Internal Server Error !"
	}

	http.Error(w, message, http.StatusInternalServerError)
}
