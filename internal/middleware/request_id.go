package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/xanderbilla/my-project/internal/constants"
)

// This file implements request ID middleware.
// Every HTTP request gets a unique identifier (request ID) that can be used for tracking.

// Request IDs are incredibly useful for:
// - Tracking a single request across multiple services (in microservices)
// - Correlating log entries (all logs for one request share the same ID)
// - Debugging issues (users can quote the request ID when reporting problems)
// - Performance monitoring (track how long specific requests take)

// RequestID is a middleware function that adds a unique ID to every request.
//
// How it works:
//  1. Check if the client already sent a request ID in the X-Request-ID header
//  2. If not, generate a new UUID (Universally Unique Identifier)
//  3. Store the ID in the request context (so handlers can access it)
//  4. Send the ID back in the response headers (so clients can see it)
//
// Example: If a mobile app encounters an error, it can show the user:
//
//	"An error occurred (ID: abc-123-def). Please contact support."
//
// The support team can then search logs for "abc-123-def" to find the issue.
//
// Usage in your server setup:
//
//	handler := middleware.RequestID(yourMux)
//	http.ListenAndServe(":8080", handler)
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get request ID from client's header
		// Some clients (like API testing tools or frontend apps) might send their own request ID
		requestID := r.Header.Get("X-Request-ID")

		// If client didn't provide a request ID, generate a new one
		if requestID == "" {
			// Generate a UUID (Universally Unique Identifier)
			// Example UUID: "f7c9d01a-12b3-4d5e-a8c9-8a3d1f0f1234"
			// UUIDs are practically guaranteed to be unique (collision probability is extremely low)
			requestID = uuid.New().String()
			slog.Debug("generated new request ID", "requestId", requestID, "method", r.Method, "path", r.URL.Path)
		} else {
			slog.Debug("using client-provided request ID", "requestId", requestID, "method", r.Method, "path", r.URL.Path)
		}

		// Add the request ID to the request context
		// Context is like a key-value store that travels with the request
		// Other code can retrieve the requestID using: ctx.Value(constants.RequestIDKey)
		ctx := context.WithValue(r.Context(), constants.RequestIDKey, requestID)

		// Send the request ID back to the client in the response header
		// This allows clients to log the ID on their side for debugging
		w.Header().Set("X-Request-ID", requestID)

		// Continue processing the request with the updated context
		// r.WithContext(ctx) creates a copy of the request with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
