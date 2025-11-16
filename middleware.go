package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// ============================================================================
// MIDDLEWARE FUNCTIONS
// ============================================================================
// Middleware is code that runs BEFORE and/or AFTER your handler.
// Think of it as a pipeline: Request -> Middleware1 -> Middleware2 -> Handler
//
// Common uses:
// - Logging (track all requests)
// - Authentication (verify user identity)
// - Authorization (check permissions)
// - Rate limiting (prevent abuse)
// - CORS (allow cross-origin requests)
// - Panic recovery (gracefully handle crashes)

// Logging creates a middleware that logs request information.
// This is useful for debugging and monitoring.
//
// Key concept: Closure
// A closure is a function that references variables from outside its body.
// Here, we return a function that "closes over" the next handler.
func Logging() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Before handler: record start time
			start := time.Now()

			// Log request details
			log.Printf("Started %s %s from %s",
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
			)

			// Call the next handler in the chain
			next(w, r)

			// After handler: log completion time
			// defer could also be used here for cleanup
			duration := time.Since(start)
			log.Printf("Completed %s %s in %v",
				r.Method,
				r.URL.Path,
				duration,
			)
		}
	}
}

// RequireAuth creates a middleware that checks for authentication.
// This demonstrates a more realistic auth check using bearer tokens.
//
// In production, you would:
// - Verify JWT tokens
// - Check against a database
// - Use proper cryptographic verification
//
// Key concept: Early return pattern
// If auth fails, we return immediately without calling next().
// This "short-circuits" the middleware chain.
func RequireAuth() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")

			// Check if header is present
			if authHeader == "" {
				Unauthorized(w, "Missing Authorization header")
				return // Early return - don't call next()
			}

			// Check if it's a Bearer token
			// Format: "Bearer <token>"
			if !strings.HasPrefix(authHeader, "Bearer ") {
				Unauthorized(w, "Invalid authorization format. Use: Bearer <token>")
				return
			}

			// Extract the token
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate the token (simplified - use proper JWT validation in production)
			// For learning purposes, we accept any non-empty token
			if token == "" {
				Unauthorized(w, "Empty token")
				return
			}

			// For demo: accept "demo-token" or any token starting with "valid-"
			if token != "demo-token" && !strings.HasPrefix(token, "valid-") {
				Unauthorized(w, "Invalid token")
				return
			}

			// Token is valid, log it
			log.Printf("Authenticated request with token: %s...", token[:min(10, len(token))])

			// Call the next handler - authentication passed!
			next(w, r)
		}
	}
}

// CORS creates a middleware that adds Cross-Origin Resource Sharing headers.
// CORS is needed when your API is called from a different domain than where it's hosted.
//
// Example: Your API is at api.example.com, but your frontend is at app.example.com
// Without CORS headers, browsers will block these cross-origin requests.
func CORS() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			// These headers tell the browser which cross-origin requests are allowed

			// Allow requests from any origin (* = all)
			// In production, specify your exact frontend domain
			w.Header().Set("Access-Control-Allow-Origin", "*")

			// Allow these HTTP methods
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			// Allow these headers in requests
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight OPTIONS requests
			// Browsers send OPTIONS first to check if the actual request is allowed
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next(w, r)
		}
	}
}

// RecoverPanic creates a middleware that recovers from panics.
// This prevents your entire server from crashing if one request panics.
//
// Key concept: defer and recover
// - defer: schedules a function to run when the surrounding function returns
// - recover: catches a panic and prevents it from propagating
//
// This is CRUCIAL for production servers - one bad request shouldn't crash everything.
func RecoverPanic() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// defer ensures this runs even if there's a panic
			defer func() {
				// recover() returns the panic value, or nil if no panic
				if err := recover(); err != nil {
					// Log the panic for debugging
					log.Printf("PANIC recovered: %v", err)

					// Send a generic error response
					// Don't expose internal error details to clients
					InternalError(w, "An unexpected error occurred")
				}
			}()

			next(w, r)
		}
	}
}

// RateLimit creates a simple in-memory rate limiter.
// This prevents clients from making too many requests.
//
// Note: This is a simplified example. Production rate limiting should:
// - Use Redis or similar for distributed systems
// - Implement proper token bucket or sliding window algorithms
// - Handle multiple server instances
func RateLimit(requestsPerMinute int) Middleware {
	// Track request counts per IP
	// This map is closed over by the returned function
	requests := make(map[string]int)
	lastReset := time.Now()

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Reset counters every minute
			if time.Since(lastReset) > time.Minute {
				requests = make(map[string]int)
				lastReset = time.Now()
			}

			// Get client IP
			clientIP := r.RemoteAddr

			// Check rate limit
			if requests[clientIP] >= requestsPerMinute {
				w.Header().Set("Retry-After", "60")
				JSON(w, http.StatusTooManyRequests, APIResponse{
					Success: false,
					Error:   "Rate limit exceeded. Please try again later.",
				})
				return
			}

			// Increment counter
			requests[clientIP]++

			next(w, r)
		}
	}
}

// min returns the smaller of two integers.
// Go doesn't have a built-in min for integers (until Go 1.21).
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
