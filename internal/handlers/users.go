package handlers

import (
	"log/slog"
	"net/http"

	"github.com/xanderbilla/my-project/internal/constants"
	"github.com/xanderbilla/my-project/internal/types"
	"github.com/xanderbilla/my-project/internal/utils/response"
)

// This file handles all user-related HTTP endpoints.
// Each function corresponds to a specific endpoint and HTTP method.

// UserRoutes registers all user-related routes with the HTTP server.
// This is called during server startup to configure which URLs map to which handlers.
//
// Route format: "METHOD /path"
// - METHOD: HTTP verb (GET, POST, PATCH, DELETE)
// - /path: URL path with optional parameters in curly braces {param}
//
// Example usage in main.go:
//
//	mux := http.NewServeMux()
//	handlers.UserRoutes(mux)
func UserRoutes(mux *http.ServeMux) {
	// GET /health - Health check endpoint for Docker and monitoring
	mux.HandleFunc("GET /health", handleHealth)

	// GET /favicon.ico - Handle browser favicon requests
	// Browsers automatically request this, so we handle it to avoid 404 warnings
	mux.HandleFunc("GET /favicon.ico", handleFavicon)

	// POST /api/users - Create a new user
	mux.HandleFunc("POST /api/users", createUser)

	// GET /api/users - Get all users (with optional pagination)
	mux.HandleFunc("GET /api/users", getAllUsers)

	// GET /api/users/123 - Get a specific user by ID
	// {user_id} is a path parameter that can be accessed with r.PathValue("user_id")
	mux.HandleFunc("GET /api/users/{user_id}", getUserByID)

	// PATCH /api/users/123 - Update a user
	// PATCH is used for partial updates (not replacing entire resource)
	mux.HandleFunc("PATCH /api/users/{user_id}", updateUser)

	// DELETE /api/users/123 - Delete a user
	mux.HandleFunc("DELETE /api/users/{user_id}", deleteUser)
}

// handleHealth handles GET /health
// Health check endpoint used by Docker, Kubernetes, load balancers, and monitoring systems.
// Always returns 200 OK when the application is running.
//
// In production, you might want to:
// 1. Check database connectivity
// 2. Verify external service dependencies
// 3. Check disk space or memory
// 4. Return different status codes based on health (200, 503, etc.)
//
// Example response:
//
//	{
//	  "status": "healthy",
//	  "timestamp": "2026-03-04T10:30:00Z"
//	}
func handleHealth(w http.ResponseWriter, r *http.Request) {
	slog.Debug("health check requested")

	// Add caching headers (cache for 30 seconds)
	w.Header().Set("Cache-Control", "public, max-age=30")

	response.OK(w, r, "Service is healthy", map[string]string{
		"status": "healthy",
	})
}

// handleFavicon handles GET /favicon.ico
// Browsers automatically request favicon.ico when you visit a site in a browser.
// This prevents unnecessary 404 WARN logs in production.
//
// Returns 204 No Content (success with no body)
//
// In production, you would typically:
// 1. Serve an actual favicon.ico file
// 2. Or redirect to a static file server
// 3. Or return 204 No Content if you don't have a favicon
func handleFavicon(w http.ResponseWriter, r *http.Request) {
	// Get request ID for logging
	requestID := "unknown"
	if id, ok := r.Context().Value(constants.RequestIDKey).(string); ok {
		requestID = id
	}

	slog.Debug("handling favicon request", "requestId", requestID)

	// Return 204 No Content - successful request with no response body
	// This prevents browsers from showing an error icon in the tab
	w.WriteHeader(http.StatusNoContent)
}

// getAllUsers handles GET /api/users
// Returns a list of all users in the system.
//
// Query parameters (optional):
//   - page: Page number (default: 1)
//   - pageSize: Items per page (default: 20, max: 100)
//
// Example requests:
//
//	GET /api/users               -> First page with 20 items
//	GET /api/users?page=2        -> Second page with 20 items
//	GET /api/users?page=1&pageSize=50  -> First page with 50 items
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	// Get request ID for logging
	requestID := "unknown"
	if id, ok := r.Context().Value(constants.RequestIDKey).(string); ok {
		requestID = id
	}

	slog.Debug("processing getAllUsers request", "requestId", requestID, "handler", "getAllUsers")

	// In a real application, you would:
	// 1. Get pagination parameters from the request
	// 2. Query the database for the requested page of users
	// 3. Get the total count of users
	// 4. Return paginated response
	//
	// For now, we're using mock data for demonstration

	slog.Debug("creating mock user data", "requestId", requestID, "count", 2)

	// Mock data - in production, this would come from a database
	users := []types.User{
		{Name: "John Doe", Email: "john@example.com", Age: 30},
		{Name: "Jane Smith", Email: "jane@example.com", Age: 25},
	}

	// Add caching headers (cache for 60 seconds)
	w.Header().Set("Cache-Control", "public, max-age=60")

	slog.Debug("sending success response", "requestId", requestID, "user_count", len(users))

	// Send success response with the user data
	response.OK(w, r, "Users retrieved successfully", users)
}

// getUserByID handles GET /api/users/{user_id}
// Returns a single user by their ID.
//
// Path parameters:
//   - user_id: The unique identifier of the user
//
// Responses:
//   - 200 OK: User found and returned
//   - 404 Not Found: User with given ID doesn't exist
//
// Example requests:
//
//	GET /api/users/1   -> Returns user with ID 1
//	GET /api/users/999 -> Returns 404 if user doesn't exist
func getUserByID(w http.ResponseWriter, r *http.Request) {
	// Extract user_id from the URL path
	// PathValue gets the value of a path parameter defined in the route
	userID := r.PathValue("user_id")

	// Get request ID for logging
	requestID := "unknown"
	if id, ok := r.Context().Value(constants.RequestIDKey).(string); ok {
		requestID = id
	}

	slog.Debug("processing getUserByID request", "requestId", requestID, "handler", "getUserByID", "user_id", userID)

	// Mock validation: In real code, you would query the database
	// Here we're pretending only user ID "1" exists
	if userID != "1" {
		slog.Warn("user not found", "requestId", requestID, "user_id", userID, "reason", "invalid_user_id")
		// User not found - send 404 response
		response.NotFound(w, r, "User", userID)
		return
	}

	slog.Debug("user found in database", "requestId", requestID, "user_id", userID)

	// Mock user data - in production, this would come from database
	user := types.User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	// Add caching headers (cache for 60 seconds)
	w.Header().Set("Cache-Control", "public, max-age=60")

	slog.Debug("sending user response", "requestId", requestID, "user_name", user.Name)

	// Send success response with the user data
	response.OK(w, r, "User retrieved successfully", user)
}

// createUser handles POST /api/users
// Creates a new user in the system.
//
// Request body (JSON):
//
//	{
//	  "name": "John Doe",
//	  "email": "john@example.com",
//	  "age": 30
//	}
//
// Validations:
//   - Content-Type must be "application/json"
//   - JSON must be valid
//   - All required fields must be present
//   - Email must be unique (no duplicates)
//   - Age must be at least 18
//
// Responses:
//   - 201 Created: User successfully created
//   - 400 Bad Request: Invalid JSON or wrong Content-Type
//   - 409 Conflict: Email already exists
//   - 422 Unprocessable Entity: Validation errors
func createUser(w http.ResponseWriter, r *http.Request) {
	// Always close the request body when done
	// This prevents resource leaks
	defer func() {
		if err := r.Body.Close(); err != nil {
			slog.Warn("failed to close request body", "error", err)
		}
	}()

	// Get request ID for logging
	requestID := "unknown"
	if id, ok := r.Context().Value(constants.RequestIDKey).(string); ok {
		requestID = id
	}

	slog.Debug("processing createUser request", "requestId", requestID, "handler", "createUser", "content_type", r.Header.Get("Content-Type"))

	// Use the helper function to decode and validate JSON in one step
	// This function:
	// 1. Checks Content-Type
	// 2. Decodes the JSON
	// 3. Runs validation
	// 4. Sends error responses automatically if anything fails
	//
	// If it returns false, an error response was already sent, so we just return
	var user types.User

	// Create a wrapper function that matches the expected signature
	// The validation function expects interface{} but validateUser expects types.User
	validateFn := func(data interface{}) []types.ValidationError {
		// Convert interface{} back to *types.User
		u := data.(*types.User)
		return validateUser(*u)
	}

	if !response.DecodeJSONWithValidation(w, r, &user, validateFn) {
		slog.Warn("request validation failed", "requestId", requestID, "handler", "createUser")
		return
	}

	slog.Debug("request validated successfully", "requestId", requestID, "email", user.Email, "name", user.Name)

	// Check for business rule violations
	// In this mock example, we're checking if the email already exists
	// In production, you would query the database to check this
	if user.Email == "existing@example.com" {
		slog.Warn("email conflict detected", "requestId", requestID, "email", user.Email, "reason", "duplicate_email")
		response.Conflict(w, r,
			constants.ErrCodeUserAlreadyExists,
			"A user with this email already exists",
			"User with email "+user.Email+" already exists in database",
		)
		return
	}

	slog.Debug("creating user in database", "requestId", requestID, "email", user.Email)

	// If we got here, all validations passed!
	// In production, you would:
	// 1. Insert the user into the database
	// 2. Get the generated user ID
	// 3. Return the complete user object with its new ID
	//
	// For now, we just send back what was received
	response.Created(w, r, "User created successfully", user)
}

// updateUser handles PATCH /api/users/{user_id}
// Updates an existing user's information.
//
// Path parameters:
//   - user_id: The ID of the user to update
//
// Request body (JSON) - all fields optional:
//
//	{
//	  "name": "New Name",
//	  "email": "newemail@example.com",
//	  "age": 31
//	}
//
// Responses:
//   - 200 OK: User successfully updated
//   - 400 Bad Request: Invalid JSON
//   - 404 Not Found: User doesn't exist
func updateUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the URL path
	userID := r.PathValue("user_id")

	// Always close request body
	defer func() {
		if err := r.Body.Close(); err != nil {
			slog.Warn("failed to close request body", "error", err)
		}
	}()

	// Check if user exists
	// In production, you would check the database
	if userID != "1" {
		response.NotFound(w, r, "User", userID)
		return
	}

	// Decode the update data
	// Note: We're not validating here because PATCH allows partial updates
	// You could add custom validation for PATCH if needed
	var user types.User
	if !response.DecodeJSONWithValidation(w, r, &user, nil) {
		return
	}

	// In production, you would:
	// 1. Load the existing user from database
	// 2. Update only the fields that were provided
	// 3. Save back to database
	// 4. Return the updated user
	//
	// For this mock, we just return what was sent
	response.OK(w, r, "User updated successfully", user)
}

// deleteUser handles DELETE /api/users/{user_id}
// Deletes a user from the system.
//
// Path parameters:
//   - user_id: The ID of the user to delete
//
// Responses:
//   - 204 No Content: User successfully deleted
//   - 404 Not Found: User doesn't exist
//
// Note: 204 responses don't include a body (just status code)
func deleteUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL path
	userID := r.PathValue("user_id")

	// Check if user exists
	// In production, you would check the database
	if userID != "1" {
		response.NotFound(w, r, "User", userID)
		return
	}

	// In production, you would delete from database here
	// db.Delete("users", userID)

	// Send 204 No Content response
	// This means success, but no data to return
	response.NoContent(w, r)
}

// validateUser performs validation checks on a User struct.
// Returns a slice of validation errors (empty if valid).
//
// Validation rules:
//   - name: required (cannot be empty)
//   - email: required (cannot be empty) and must be valid format
//   - age: must be between 18 and 120
//
// This function is passed to DecodeJSONWithValidation to automatically
// validate user input after decoding JSON.
func validateUser(user types.User) []types.ValidationError {
	var errors []types.ValidationError

	// Use helper functions to build validation errors
	// These append to the errors slice if validation fails
	response.ValidateRequired("name", user.Name, &errors)
	response.ValidateRequired("email", user.Email, &errors)
	response.ValidateEmail("email", user.Email, &errors)
	response.ValidateRange("age", user.Age, 18, 120, &errors)

	return errors
}
