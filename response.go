package main

import (
	"encoding/json"
	"net/http"
)

// ============================================================================
// RESPONSE HELPERS
// ============================================================================
// These helper functions standardize how we send responses.
// Key concept: DRY (Don't Repeat Yourself)
// Instead of setting headers and encoding JSON in every handler,
// we centralize that logic here.

// JSON sends a JSON response with the given status code.
// This is the core helper that all other response helpers use.
//
// Important patterns demonstrated:
// 1. Setting Content-Type header BEFORE writing status code
// 2. Using json.NewEncoder for efficient streaming
// 3. Handling encoding errors gracefully
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	// Set Content-Type header FIRST
	// Once you call WriteHeader or Write, you can't modify headers
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code (200, 404, 500, etc.)
	w.WriteHeader(statusCode)

	// Encode the data as JSON and write to response
	// json.NewEncoder streams directly to the writer (more efficient than Marshal)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails, we've already written headers, so just log it
		// In production, you'd use a proper logger here
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Success sends a successful response with data.
// Use this for 200 OK responses with payload.
//
// Example usage in handler:
//
//	user := User{Name: "John", Email: "john@example.com"}
//	Success(w, "User created", user)
func Success(w http.ResponseWriter, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	JSON(w, http.StatusOK, response)
}

// Created sends a 201 Created response.
// Use this when a new resource is successfully created.
//
// HTTP 201 is specifically for successful creation operations.
// It's more semantic than just returning 200.
func Created(w http.ResponseWriter, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	JSON(w, http.StatusCreated, response)
}

// Error sends an error response with appropriate status code.
// This handles both custom AppErrors and standard Go errors.
//
// Key concept: Type assertion
// We check if the error is our custom AppError type to get the status code.
// If not, we default to 500 Internal Server Error.
func Error(w http.ResponseWriter, err error) {
	var statusCode int
	var message string

	// Type assertion: check if err is an AppError
	// This is a powerful Go pattern for working with interfaces
	if appErr, ok := err.(AppError); ok {
		// It's our custom error type, use its code
		statusCode = appErr.Code
		message = appErr.Message
	} else {
		// Generic error, use 500
		statusCode = http.StatusInternalServerError
		message = err.Error()
	}

	response := APIResponse{
		Success: false,
		Error:   message,
	}
	JSON(w, statusCode, response)
}

// BadRequest sends a 400 Bad Request response.
// Use this when the client sends invalid data.
func BadRequest(w http.ResponseWriter, message string) {
	response := APIResponse{
		Success: false,
		Error:   message,
	}
	JSON(w, http.StatusBadRequest, response)
}

// NotFound sends a 404 Not Found response.
// Use this when the requested resource doesn't exist.
func NotFound(w http.ResponseWriter, message string) {
	response := APIResponse{
		Success: false,
		Error:   message,
	}
	JSON(w, http.StatusNotFound, response)
}

// Unauthorized sends a 401 Unauthorized response.
// Use this when authentication is required but not provided or invalid.
func Unauthorized(w http.ResponseWriter, message string) {
	response := APIResponse{
		Success: false,
		Error:   message,
	}
	JSON(w, http.StatusUnauthorized, response)
}

// InternalError sends a 500 Internal Server Error response.
// Use this for unexpected server-side errors.
//
// Important: In production, don't expose internal error details to clients.
// Log the full error server-side, but send a generic message to the client.
func InternalError(w http.ResponseWriter, message string) {
	response := APIResponse{
		Success: false,
		Error:   message,
	}
	JSON(w, http.StatusInternalServerError, response)
}

// NoContent sends a 204 No Content response.
// Use this for successful operations that don't return data (like DELETE).
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
