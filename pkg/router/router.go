package router

import (
	"net/http"

	"github.com/himtar/go-boilerplate/pkg/errors"
	"github.com/himtar/go-boilerplate/pkg/helpers"
)

type RouterMux struct {
	*http.ServeMux
}

// NewRouterMux creates a new instance of RouterMux.
func NewRouterMux() *RouterMux {
	return &RouterMux{ServeMux: http.NewServeMux()}
}

// AddCustomGetHandler registers a custom GET handler for the specified urlPath.
func (r *RouterMux) AddCustomGetHandler(urlPath string, handler http.HandlerFunc) {
	r.HandleFunc(urlPath, func(w http.ResponseWriter, req *http.Request) {
		// validate method
		err := helpers.ValidateMethod(http.MethodGet, w, req)
		if err != "" {
			return
		}

		// then call the main handler function
		handler(w, req)
	})
}

// AddCustomPostHandler registers a custom POST handler for the specified urlPath.
func (r *RouterMux) AddCustomPostHandler(urlPath string, handler http.HandlerFunc) {
	r.HandleFunc(urlPath, func(w http.ResponseWriter, req *http.Request) {	
		// validate method
		err := helpers.ValidateMethod(http.MethodPost, w, req)
		if err != "" {
			return
		}

		// then call the main handler function
		handler(w, req)
	})
}

// AddCustomPostHandler registers a custom DELETE handler for the specified urlPath.
func (r *RouterMux) AddCustomDeleteHandler(urlPath string, handler http.HandlerFunc) {
	r.HandleFunc(urlPath, func(w http.ResponseWriter, req *http.Request) {	
		// validate method
		err := helpers.ValidateMethod(http.MethodDelete, w, req)
		if err != "" {
			return
		}
		
		// then call the main handler function
		handler(w, req)
	})
}

func (r *RouterMux) AddUnsupportedMethodHandler(urlPath string) {
	r.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {
		errors.MethodNotAllowed(w)
	})
}

// AddCustomHandler registers a custom handler for the specified method and urlPath.
func (r *RouterMux) AddCustomHandler(method, urlPath string, handler http.HandlerFunc) {
	switch method {
	case http.MethodGet:
		r.AddCustomGetHandler(urlPath, handler)
	case http.MethodPost:
		r.AddCustomPostHandler(urlPath, handler)
	case http.MethodDelete:
		r.AddCustomDeleteHandler(urlPath, handler)
	default:
		r.AddUnsupportedMethodHandler(urlPath)
	}
}