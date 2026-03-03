package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/xanderbilla/my-project/internal/constants"
	"github.com/xanderbilla/my-project/internal/utils/response"
)

// This file implements panic recovery middleware.
// If your code panics (crashes), this middleware catches it and returns a proper error response.

// What is a panic?
// In Go, a panic is like an exception in other languages - it stops normal program execution.
// Panics can be caused by:
// - Accessing nil pointers
// - Out of bounds array access
// - Type assertion failures
// - Explicitly calling panic()
//
// Without recovery, a panic would crash your entire server!
// This middleware catches panics and keeps your server running.

// Recovery is a middleware function that recovers from panics.
//
// How it works:
//  1. Set up a deferred function (runs even if code panics)
//  2. Process the request normally
//  3. If a panic occurs, catch it and send an error response
//  4. Log the panic with a stack trace for debugging
//
// The defer/recover pattern is Go's way of handling panics:
// - defer: schedules a function to run after the current function returns
// - recover(): returns the panic value (or nil if no panic occurred)
//
// Example panic scenario:
//
//	var user *User = nil
//	name := user.Name  // This would panic (nil pointer dereference)
//	// Without Recovery middleware → entire server crashes
//	// With Recovery middleware → client gets 500 error, server keeps running
//
// Usage in your server setup:
//
//	handler := middleware.Recovery(yourMux)
//	http.ListenAndServe(":8080", handler)
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defer means this function will run after ServeHTTP completes
		// (or when a panic occurs)
		defer func() {
			// recover() returns nil if there was no panic
			// If there was a panic, it returns the panic value
			if panicValue := recover(); panicValue != nil {
				// Try to get request ID for better logging
				requestID := "unknown"
				if id, ok := r.Context().Value(constants.RequestIDKey).(string); ok {
					requestID = id
				}

				// A panic occurred! Log detailed information for debugging
				slog.Error("panic recovered",
					"requestId", requestID,
					"error", panicValue, // The panic message
					"path", r.URL.Path, // Which endpoint panicked
					"method", r.Method, // Which HTTP method
					"stack", string(debug.Stack()), // Full stack trace showing where panic occurred
				)

				// Send a proper error response to the client
				// Convert panic value to error (could be string, error, or other type)
				var err error
				switch v := panicValue.(type) {
				case error:
					err = v
				case string:
					err = fmt.Errorf("%s", v)
				default:
					err = fmt.Errorf("panic: %v", v)
				}

				response.InternalError(w, r, err)
			}
		}()

		// Process the request normally
		// If this panics, the defer function above will catch it
		next.ServeHTTP(w, r)
	})
}
