package main

import (
	"encoding/json"
	"net/http"
)

// ============================================================================
// HTTP HANDLERS
// ============================================================================
// Handlers are the "controllers" of your API - they process requests and send responses.
// Each handler is responsible for:
// 1. Reading and validating input
// 2. Performing business logic (or calling services that do)
// 3. Sending an appropriate response

// HealthCheck handles GET /health requests.
// Health check endpoints are standard practice for:
// - Load balancers to verify the service is alive
// - Monitoring systems to track uptime
// - Kubernetes liveness/readiness probes
//
// Key concept: Handlers always have the same signature:
// func(http.ResponseWriter, *http.Request)
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Simple health check response
	// In production, you might check database connectivity, etc.
	Success(w, "Service is healthy", map[string]string{
		"status":  "ok",
		"version": "1.0.0",
	})
}

// GetAPIInfo handles GET /api requests.
// This provides general information about the API.
func GetAPIInfo(w http.ResponseWriter, r *http.Request) {
	apiInfo := map[string]interface{}{
		"name":        "Go Learning API",
		"version":     "1.0.0",
		"description": "A simple REST API for learning Go",
		"endpoints": []string{
			"GET  /health - Health check",
			"GET  /api - API information",
			"GET  /api/users - List all users",
			"POST /api/users - Create a user",
			"GET  /api/users/{id} - Get user by ID",
		},
	}

	Success(w, "API information", apiInfo)
}

// ============================================================================
// USER HANDLERS (CRUD Operations)
// ============================================================================
// These handlers demonstrate typical CRUD operations:
// Create, Read, Update, Delete

// In-memory storage for demo purposes.
// In production, you'd use a database.
// Key concept: Package-level variables
// These are initialized when the package loads.
var (
	users     = make(map[int]User) // Map of user ID to User
	nextID    = 1                  // Auto-incrementing ID
)

// ListUsers handles GET /api/users requests.
// Returns all users in the system.
//
// This demonstrates:
// - Converting a map to a slice
// - Iterating over maps with range
func ListUsers(w http.ResponseWriter, r *http.Request) {
	// Convert map to slice for JSON response
	// Maps don't have a guaranteed order, so we convert to slice
	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}

	Success(w, "Users retrieved successfully", userList)
}

// CreateUser handles POST /api/users requests.
// Creates a new user from JSON body.
//
// This demonstrates:
// - Reading request body
// - JSON decoding
// - Input validation
// - Error handling
// - Proper HTTP status codes (201 Created)
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Step 1: Decode JSON request body
	// json.NewDecoder streams from the request body (more efficient than ReadAll)
	var user User
	decoder := json.NewDecoder(r.Body)

	// Decode into our User struct
	if err := decoder.Decode(&user); err != nil {
		// If JSON is malformed, send 400 Bad Request
		BadRequest(w, "Invalid JSON: "+err.Error())
		return
	}

	// Step 2: Validate the input
	// Always validate user input before processing
	if err := user.Validate(); err != nil {
		Error(w, err) // Our custom Error helper handles AppError
		return
	}

	// Step 3: Assign ID and store
	// In production, the database would generate the ID
	user.ID = nextID
	nextID++
	users[user.ID] = user

	// Step 4: Return the created user with 201 status
	// 201 Created is the proper status for successful creation
	Created(w, "User created successfully", user)
}

// GetUser handles GET /api/users/{id} requests.
// Retrieves a specific user by ID.
//
// Note: Our simple router doesn't support path parameters.
// In a real app, you'd use a router like chi, gorilla/mux, or gin.
// For learning, this shows how to handle query parameters instead.
func GetUser(w http.ResponseWriter, r *http.Request) {
	// Get ID from query parameters (?id=1)
	// In production with a proper router: r.PathValue("id") or mux.Vars(r)["id"]
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		BadRequest(w, "Missing 'id' query parameter")
		return
	}

	// Convert string to int
	// strconv.Atoi converts ASCII to Integer
	var id int
	if _, err := parseID(idStr, &id); err != nil {
		BadRequest(w, "Invalid ID format")
		return
	}

	// Look up the user
	user, exists := users[id]
	if !exists {
		Error(w, ErrNotFound("user"))
		return
	}

	Success(w, "User found", user)
}

// parseID is a helper to parse ID strings.
// This shows how to create small, focused helper functions.
func parseID(s string, id *int) (bool, error) {
	n := 0
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false, ErrInvalidInput("ID must be numeric")
		}
		n = n*10 + int(ch-'0')
	}
	*id = n
	return true, nil
}

// DeleteUser handles DELETE /api/users requests.
// Removes a user by ID.
//
// This demonstrates:
// - DELETE operations
// - 204 No Content response (successful delete, no body)
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		BadRequest(w, "Missing 'id' query parameter")
		return
	}

	var id int
	if _, err := parseID(idStr, &id); err != nil {
		BadRequest(w, "Invalid ID format")
		return
	}

	// Check if user exists
	if _, exists := users[id]; !exists {
		Error(w, ErrNotFound("user"))
		return
	}

	// Delete from map
	delete(users, id)

	// 204 No Content - successful deletion, no response body
	NoContent(w)
}

// UpdateUser handles PUT /api/users requests.
// Updates an existing user.
//
// This demonstrates:
// - PUT semantics (full replacement)
// - Combining read and write operations
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID to update
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		BadRequest(w, "Missing 'id' query parameter")
		return
	}

	var id int
	if _, err := parseID(idStr, &id); err != nil {
		BadRequest(w, "Invalid ID format")
		return
	}

	// Check if user exists
	if _, exists := users[id]; !exists {
		Error(w, ErrNotFound("user"))
		return
	}

	// Decode the new user data
	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		BadRequest(w, "Invalid JSON: "+err.Error())
		return
	}

	// Validate the updated data
	if err := updatedUser.Validate(); err != nil {
		Error(w, err)
		return
	}

	// Keep the same ID
	updatedUser.ID = id
	users[id] = updatedUser

	Success(w, "User updated successfully", updatedUser)
}
