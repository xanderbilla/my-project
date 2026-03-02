// Package response provides helper functions
// for sending HTTP responses in JSON format.
package response

import (
	"encoding/json"
	"net/http"
)

// JSON writes a JSON response to the client.
//
// It:
//   1. Sets the Content-Type header
//   2. Writes the HTTP status code
//   3. Encodes the given data as JSON
//
// This function helps avoid repeating the same
// response-writing code in every handler.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Encode the response data as JSON
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails, send a generic internal server error
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}