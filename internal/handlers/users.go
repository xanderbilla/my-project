// Package handlers contains all HTTP route handlers.
//
// A handler is simply a function that receives an HTTP request
// and writes back an HTTP response.
//
// This file contains user-related HTTP endpoints.
package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/xanderbilla/my-project/internal/types"
	"github.com/xanderbilla/my-project/internal/utils/response"
)

// RegisterUserRoutes registers all routes related to users.
//
// We pass the application's ServeMux and attach endpoints to it.
// This keeps route registration organized in one place.
func UserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/users", getAllUsers)
	mux.HandleFunc("POST /api/users", createUser)
}

// getAllUsers handles GET /api/users.
//
// For now, it just returns a simple message.
// Later, this is where you would fetch users from a database.
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "list users",
	})
}

// createUser handles POST /api/users.
//
// It expects a JSON request body containing user details.
// Example JSON:
//
//	{
//	  "name": "Aman",
//	  "email": "aman@example.com",
//	  "age": 22
//	}
//
// The function:
//   1. Validates the Content-Type
//   2. Decodes JSON into a User struct
//   3. Performs basic validation
//   4. Returns a JSON response
func createUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() // Always close request body to free resources

	// Ensure the client sends JSON
	if r.Header.Get("Content-Type") != "application/json" {
		response.JSON(w, http.StatusUnsupportedMediaType, map[string]string{
			"error": "Content-Type must be application/json",
		})
		return
	}

	var user types.User

	// Decode JSON request body into Go struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	// Basic validation (for learning purpose)
	if user.Name == "" || user.Email == "" {
		response.JSON(w, http.StatusBadRequest, map[string]string{
			"error": "name and email are required",
		})
		return
	}

	// Log useful information using structured logging
	slog.Info("User created", "email", user.Email)

	// Return created user as response
	response.JSON(w, http.StatusCreated, user)
}