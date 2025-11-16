package main

import (
	"net/http"
	"strings"
)

// ============================================================================
// CUSTOM ROUTER IMPLEMENTATION
// ============================================================================
// This file implements a simple HTTP router from scratch.
// Understanding this helps you grasp how frameworks like Gin/Echo work internally.

// Router stores our HTTP route rules.
// It implements the http.Handler interface, which is the key to Go's HTTP handling.
//
// The rules map structure: map[path][method] = handler
// Example: rules["/api"]["GET"] = HandlerHome
//
// Key concept: Nested maps
// We use a map of maps to organize routes by path first, then by HTTP method.
// This allows multiple methods for the same path (GET /api, POST /api, etc.)
type Router struct {
	rules map[string]map[string]http.HandlerFunc
}

// newRouter creates and initializes a new Router.
// This is a constructor function - a common Go pattern.
//
// Why use a constructor?
// - Ensures the struct is properly initialized (maps need make())
// - Encapsulates initialization logic
// - Convention: lowercase = unexported (private to this package)
func newRouter() *Router {
	return &Router{
		rules: make(map[string]map[string]http.HandlerFunc),
	}
}

// FindHandler looks up the handler for a given path and method.
// Returns the handler and two booleans:
// - methodExists: true if the method is registered for this path
// - pathExists: true if the path exists at all
//
// This separation allows us to return different HTTP status codes:
// - Path not found -> 404 Not Found
// - Path exists but method not allowed -> 405 Method Not Allowed
func (r *Router) FindHandler(path string, method string) (http.HandlerFunc, bool, bool) {
	// Normalize the path by removing trailing slashes (except for root "/")
	// This makes "/api" and "/api/" match the same route
	path = normalizePath(path)

	// Check if the path exists in our rules
	_, pathExists := r.rules[path]

	// Check if the method exists for this path
	// If path doesn't exist, this will safely return false
	handler, methodExists := r.rules[path][method]

	return handler, methodExists, pathExists
}

// normalizePath cleans up the path for consistent matching.
// This is a helper function to handle common URL variations.
func normalizePath(path string) string {
	// Remove trailing slash, but keep "/" as is
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}
	return path
}

// ServeHTTP implements the http.Handler interface.
// This is the MOST IMPORTANT method - it's called for every HTTP request.
//
// Key concept: The http.Handler interface
// Any type that implements ServeHTTP(ResponseWriter, *Request) can handle HTTP requests.
// This is how Go's standard library achieves flexibility - everything is an interface.
//
// The interface looks like:
//
//	type Handler interface {
//	    ServeHTTP(ResponseWriter, *Request)
//	}
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Find the appropriate handler for this request
	handler, methodExists, pathExists := r.FindHandler(req.URL.Path, req.Method)

	// Case 1: Path doesn't exist -> 404 Not Found
	if !pathExists {
		NotFound(w, "endpoint not found: "+req.URL.Path)
		return
	}

	// Case 2: Path exists but method not allowed -> 405 Method Not Allowed
	if !methodExists {
		w.Header().Set("Allow", r.getAllowedMethods(req.URL.Path))
		JSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   "method " + req.Method + " not allowed for this endpoint",
		})
		return
	}

	// Case 3: Everything is good, call the handler
	handler(w, req)
}

// getAllowedMethods returns a comma-separated list of allowed methods for a path.
// This is used in the 405 response to tell the client what methods ARE allowed.
//
// HTTP spec says you SHOULD include an Allow header in 405 responses.
func (r *Router) getAllowedMethods(path string) string {
	path = normalizePath(path)
	methods := []string{}

	if pathRules, exists := r.rules[path]; exists {
		for method := range pathRules {
			methods = append(methods, method)
		}
	}

	return strings.Join(methods, ", ")
}
