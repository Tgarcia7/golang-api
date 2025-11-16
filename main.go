package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ============================================================================
// APPLICATION ENTRY POINT
// ============================================================================
// main() is where Go programs start execution.
// This file demonstrates modern Go practices for building a web server.

func main() {
	// ========================================================================
	// CONFIGURATION
	// ========================================================================
	// Best practice: Use environment variables for configuration.
	// This follows the 12-factor app methodology.
	// If env var isn't set, use a default value.

	port := getEnv("PORT", ":3000")

	// Create our server instance
	server := NewServer(port)

	// ========================================================================
	// ROUTE REGISTRATION
	// ========================================================================
	// Register all API endpoints.
	// Routes are organized by functionality.

	// Health check endpoint (no middleware - always accessible)
	// This is typically used by load balancers and monitoring systems
	server.Handle("GET", "/health", HealthCheck)

	// API info endpoint with basic middleware stack
	server.Handle("GET", "/api",
		server.AddMiddleware(GetAPIInfo,
			RecoverPanic(), // Catch any panics
			Logging(),      // Log all requests
			CORS(),         // Enable CORS for browser clients
		),
	)

	// ========================================================================
	// USER CRUD ENDPOINTS
	// ========================================================================
	// These endpoints demonstrate RESTful API patterns.
	// Each operation (Create, Read, Update, Delete) maps to an HTTP method.

	// List all users (public endpoint with logging)
	server.Handle("GET", "/api/users",
		server.AddMiddleware(ListUsers,
			RecoverPanic(),
			Logging(),
			CORS(),
		),
	)

	// Create a new user (requires authentication)
	server.Handle("POST", "/api/users",
		server.AddMiddleware(CreateUser,
			RecoverPanic(),
			Logging(),
			CORS(),
			RequireAuth(), // Only authenticated users can create
		),
	)

	// Get a specific user by ID (public endpoint)
	// Usage: GET /api/user?id=1
	server.Handle("GET", "/api/user",
		server.AddMiddleware(GetUser,
			RecoverPanic(),
			Logging(),
			CORS(),
		),
	)

	// Update an existing user (requires authentication)
	// Usage: PUT /api/user?id=1 with JSON body
	server.Handle("PUT", "/api/user",
		server.AddMiddleware(UpdateUser,
			RecoverPanic(),
			Logging(),
			CORS(),
			RequireAuth(),
		),
	)

	// Delete a user (requires authentication)
	// Usage: DELETE /api/user?id=1
	server.Handle("DELETE", "/api/user",
		server.AddMiddleware(DeleteUser,
			RecoverPanic(),
			Logging(),
			CORS(),
			RequireAuth(),
		),
	)

	// ========================================================================
	// GRACEFUL SHUTDOWN
	// ========================================================================
	// This is a modern Go pattern for handling server shutdown properly.
	// It ensures that:
	// 1. We catch interrupt signals (Ctrl+C, kill commands)
	// 2. We give in-flight requests time to complete
	// 3. We clean up resources properly

	// Create a channel to listen for OS signals
	// Channels are Go's way of communicating between goroutines
	quit := make(chan os.Signal, 1)

	// Notify this channel when we receive SIGINT (Ctrl+C) or SIGTERM (kill)
	// signal.Notify routes incoming signals to this channel
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine (lightweight thread)
	// This allows main() to continue and wait for the shutdown signal
	go func() {
		// server.Listen() blocks, so it runs in this goroutine
		if err := server.Listen(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	log.Println("Server is ready to handle requests")
	log.Println("Press Ctrl+C to stop")

	// Block until we receive a signal
	// This is where main() waits for shutdown signal
	<-quit
	log.Println("Received shutdown signal")

	// Create a deadline for graceful shutdown
	// Give 30 seconds for in-flight requests to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Always call cancel to release resources

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shutdown complete")
}

// getEnv retrieves an environment variable or returns a default value.
// This is a helper function that makes configuration cleaner.
//
// Key concept: Helper functions
// Small, focused functions that do one thing well.
// They make your code more readable and maintainable.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
