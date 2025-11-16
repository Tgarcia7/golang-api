package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

// ============================================================================
// SERVER IMPLEMENTATION
// ============================================================================
// The Server struct encapsulates our HTTP server configuration and router.
// This pattern is called "dependency injection" - we inject the router into the server.

// Server holds the HTTP server configuration.
// By grouping related data and behavior together, we follow the OOP principle
// of encapsulation (even though Go isn't strictly OOP).
type Server struct {
	port   string      // Port to listen on (e.g., ":3000")
	router *Router     // Our custom router instance
	server *http.Server // The underlying HTTP server (for graceful shutdown)
}

// NewServer creates and initializes a new Server instance.
// This is a constructor function - notice it returns a pointer to Server.
//
// Why return a pointer?
// - Avoids copying the entire struct when passing it around
// - Allows methods to modify the server's state
// - More memory efficient for larger structs
// - Convention in Go for "object-like" structs
func NewServer(port string) *Server {
	return &Server{
		port:   port,
		router: newRouter(),
	}
}

// Handle registers a new route with the server.
// This is a method on the Server struct - notice the (s *Server) receiver.
//
// Key concepts:
// - Method receiver: (s *Server) makes this a method on Server
// - Pointer receiver: allows us to modify the Server's state
// - Delegation: we delegate the actual routing logic to our Router
//
// Parameters:
// - method: HTTP method (GET, POST, PUT, DELETE, etc.)
// - path: URL path (e.g., "/api/users")
// - handler: function to handle requests to this route
func (s *Server) Handle(method string, path string, handler http.HandlerFunc) {
	// Normalize the path to remove trailing slashes
	path = normalizePath(path)

	// Ensure the path entry exists in our routing map
	// This is the "lazy initialization" pattern
	if _, exists := s.router.rules[path]; !exists {
		s.router.rules[path] = make(map[string]http.HandlerFunc)
	}

	// Register the handler for this method and path
	s.router.rules[path][method] = handler

	// Log the route registration for debugging
	log.Printf("Registered route: %s %s", method, path)
}

// AddMiddleware chains multiple middleware functions around a handler.
// This is a powerful pattern for adding cross-cutting concerns like:
// - Authentication
// - Logging
// - Rate limiting
// - CORS headers
//
// How middleware chaining works:
// Given: handler H, middlewares M1, M2, M3
// Result: M3(M2(M1(H)))
// Execution order: M3 -> M2 -> M1 -> H -> M1 -> M2 -> M3
// (Each middleware can run code before AND after calling the next)
//
// The ... syntax means "variadic parameters" - accepts any number of arguments.
func (s *Server) AddMiddleware(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	// Apply middlewares in order
	// Each middleware wraps the previous result
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

// Listen starts the HTTP server.
// This is a blocking call - it runs until the server is stopped.
//
// Modern Go practice: Use http.Server struct for more control:
// - Configurable timeouts (prevents slow client attacks)
// - Graceful shutdown support
// - Better error handling
func (s *Server) Listen() error {
	// Create the HTTP server with timeouts
	// These timeouts are CRUCIAL for production:
	// - ReadTimeout: max time to read request (prevents slow-loris attacks)
	// - WriteTimeout: max time to write response
	// - IdleTimeout: max time to keep connection alive
	s.server = &http.Server{
		Addr:         s.port,
		Handler:      s.router, // Our router implements http.Handler
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", s.port)
	log.Printf("Available endpoints:")
	s.printRoutes()

	// ListenAndServe blocks until an error occurs
	// Common errors: port already in use, permission denied
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server.
// This is important for production - it:
// - Stops accepting new connections
// - Waits for existing requests to complete
// - Times out after the given context deadline
//
// Key concept: Context
// Context carries deadlines, cancellation signals, and request-scoped values.
// It's Go's way of handling timeouts and cancellation across goroutines.
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	if s.server == nil {
		return nil
	}

	// Shutdown blocks until all connections are closed or context times out
	return s.server.Shutdown(ctx)
}

// printRoutes logs all registered routes for debugging.
// This is helpful to see what endpoints are available when the server starts.
func (s *Server) printRoutes() {
	for path, methods := range s.router.rules {
		methodList := []string{}
		for method := range methods {
			methodList = append(methodList, method)
		}
		log.Printf("  %s -> [%s]", path, strings.Join(methodList, ", "))
	}
}
