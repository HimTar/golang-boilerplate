package helpers

import (
	"net/http"

	"github.com/himtar/go-boilerplate/pkg/errors"
)

func ValidateMethod(method string, w http.ResponseWriter, req *http.Request) string  {
	if req.Method != method {
		errors.MethodNotAllowed(w)
		return "Method Not Allowed"
	}

	return ""
}