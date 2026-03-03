package response

import (
	"log/slog"
	"net/http"
	"time"

	json "github.com/goccy/go-json"

	"github.com/xanderbilla/my-project/internal/constants"
	"github.com/xanderbilla/my-project/internal/types"
)

// Success sends a standardized successful response to the client.
// This is the base function that all successful responses use.
//
// Parameters:
//   - w: ResponseWriter to send the response to
//   - r: The original HTTP request (needed for metadata like path, method)
//   - status: HTTP status code (200, 201, etc.)
//   - message: Human-readable success message
//   - data: The actual data to return (can be any type)
//
// The response will automatically include:
//   - Timestamp when the response was generated
//   - Request ID for tracking
//   - API version
func Success(w http.ResponseWriter, r *http.Request, status int, message string, data interface{}) {
	// Extract the request ID from the request context
	// This was added by the RequestID middleware earlier in the request chain
	requestID := getRequestID(r.Context())

	// Build the standard response structure
	response := types.APIResponse{
		Success:   true,             // Indicates the request was successful
		Status:    status,           // HTTP status code (200, 201, etc.)
		Timestamp: time.Now().UTC(), // Current time in UTC (Universal Time Coordinated)
		RequestID: requestID,        // Unique ID for tracking this specific request
		Path:      r.URL.Path,       // The endpoint that was called (e.g., "/api/users")
		Method:    r.Method,         // HTTP method used (GET, POST, etc.)
		Message:   message,          // Success message to display
		Data:      data,             // The actual response data
		Meta: &types.Meta{
			APIVersion: constants.APIVersion, // Include the API version for client reference
		},
	}

	// Set the response header to indicate we're sending JSON
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code
	w.WriteHeader(status)

	// Convert the response struct to JSON and send it
	// json.NewEncoder(w).Encode() automatically handles the conversion
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to encode success response", "error", err)
	}
}

// OK sends a 200 OK response.
// Use this for successful GET, PUT, or PATCH requests.
//
// HTTP 200 OK means: "The request succeeded and returned the expected data"
//
// Example usage:
//
//	response.OK(w, r, "User retrieved successfully", user)
func OK(w http.ResponseWriter, r *http.Request, message string, data interface{}) {
	Success(w, r, http.StatusOK, message, data)
}

// Created sends a 201 Created response.
// Use this for successful POST requests that create a new resource.
//
// HTTP 201 Created means: "A new resource was successfully created"
//
// Example usage:
//
//	response.Created(w, r, "User created successfully", newUser)
func Created(w http.ResponseWriter, r *http.Request, message string, data interface{}) {
	Success(w, r, http.StatusCreated, message, data)
}

// NoContent sends a 204 No Content response.
// Use this for successful operations that don't return any data.
//
// HTTP 204 No Content means: "The request succeeded but there's no data to return"
// Common use cases:
//   - DELETE requests (record deleted, nothing to return)
//   - Update operations where you don't need to return the updated resource
//
// Note: A 204 response must NOT include a response body
//
// Example usage:
//
//	response.NoContent(w, r)  // For DELETE /api/users/123
func NoContent(w http.ResponseWriter, r *http.Request) {
	// Only set the status code, no body
	w.WriteHeader(http.StatusNoContent)
}
