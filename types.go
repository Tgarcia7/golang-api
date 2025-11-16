package main

import (
	"encoding/json"
	"net/http"
)

// ============================================================================
// MIDDLEWARE TYPE
// ============================================================================

// Middleware is a function type that wraps an HTTP handler.
// This pattern is called the "decorator pattern" - it lets you add behavior
// to handlers without modifying them directly.
//
// How it works:
// 1. Takes an http.HandlerFunc as input
// 2. Returns a NEW http.HandlerFunc that wraps the original
// 3. The wrapper can execute code before/after calling the original handler
//
// Example flow: Request -> Middleware1 -> Middleware2 -> Handler -> Response
type Middleware func(http.HandlerFunc) http.HandlerFunc

// ============================================================================
// DATA MODELS
// ============================================================================

// User represents a user in our system.
// The `json:"..."` tags tell Go how to serialize/deserialize JSON.
//
// Important concepts:
// - struct tags: metadata that other packages (like encoding/json) can read
// - json tags define the JSON field names (lowercase for conventions)
// - omitempty: if the field is empty, don't include it in JSON output
type User struct {
	ID    int    `json:"id,omitempty"` // omitempty = skip if zero value
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// ToJSON converts the User struct to JSON bytes.
// This is a method with a pointer receiver (*User).
//
// Why pointer receiver?
// - More efficient (doesn't copy the entire struct)
// - Allows modification of the original struct if needed
// - Convention: use pointers for methods on structs
func (u *User) ToJSON() ([]byte, error) {
	// json.Marshal converts Go structs to JSON bytes
	// It uses the struct tags to determine field names
	return json.Marshal(u)
}

// Validate checks if the User has all required fields.
// This is a simple validation - production code would be more thorough.
//
// Key concept: Methods should have single responsibility
// This method only validates, doesn't modify or save.
func (u *User) Validate() error {
	if u.Name == "" {
		return ErrInvalidInput("name is required")
	}
	if u.Email == "" {
		return ErrInvalidInput("email is required")
	}
	return nil
}

// ============================================================================
// API RESPONSE STRUCTURES
// ============================================================================

// APIResponse is a standardized response envelope.
// Using a consistent response format makes your API easier to consume.
//
// Benefits of response envelopes:
// - Consistent structure for all responses
// - Easy to add metadata (pagination, timestamps, etc.)
// - Clear separation of success/error cases
type APIResponse struct {
	Success bool        `json:"success"`           // Was the request successful?
	Message string      `json:"message,omitempty"` // Human-readable message
	Data    interface{} `json:"data,omitempty"`    // The actual response data
	Error   string      `json:"error,omitempty"`   // Error message if failed
}

// APIError represents a structured error response.
// This provides more context than just an error string.
type APIError struct {
	Code    int    `json:"code"`    // HTTP status code
	Message string `json:"message"` // Error description
	Details string `json:"details"` // Additional context (optional)
}

// ============================================================================
// CUSTOM ERROR TYPES
// ============================================================================

// AppError is a custom error type for our application.
// Implementing the error interface allows us to use it anywhere Go expects an error.
//
// Key concept: Custom error types
// - Give more context than string errors
// - Can carry additional data (error codes, stack traces, etc.)
// - Allow type assertions for specific error handling
type AppError struct {
	Message string
	Code    int
}

// Error implements the error interface.
// Any type with an Error() string method satisfies the error interface.
func (e AppError) Error() string {
	return e.Message
}

// ErrNotFound creates a "not found" error.
// These helper functions make it easy to create consistent errors.
func ErrNotFound(resource string) AppError {
	return AppError{
		Message: resource + " not found",
		Code:    http.StatusNotFound,
	}
}

// ErrInvalidInput creates a "bad request" error.
func ErrInvalidInput(detail string) AppError {
	return AppError{
		Message: "invalid input: " + detail,
		Code:    http.StatusBadRequest,
	}
}

// ErrUnauthorized creates an "unauthorized" error.
func ErrUnauthorized() AppError {
	return AppError{
		Message: "unauthorized access",
		Code:    http.StatusUnauthorized,
	}
}

// ErrInternal creates an "internal server error".
func ErrInternal(detail string) AppError {
	return AppError{
		Message: "internal error: " + detail,
		Code:    http.StatusInternalServerError,
	}
}
