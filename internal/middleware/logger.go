package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/xanderbilla/my-project/internal/constants"
)

// This file implements HTTP request logging middleware.
// Middleware is code that runs before and after your request handlers.
// The Logger middleware logs information about every HTTP request that comes in.

// responseWriter is a custom wrapper around http.ResponseWriter.
// We need this because the standard http.ResponseWriter doesn't let us see
// what status code was written. This wrapper captures that information.
type responseWriter struct {
	http.ResponseWriter     // Embeds the original ResponseWriter
	statusCode          int // Stores the HTTP status code that was written
	bytesWritten        int // Stores how many bytes were sent
}

// WriteHeader captures the status code before passing it to the real ResponseWriter.
// This method is called when your handler sets the HTTP status code.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code                // Save the status code
	rw.ResponseWriter.WriteHeader(code) // Pass it to the original ResponseWriter
}

// Write captures the number of bytes written.
// This method is called when your handler sends response body data.
func (rw *responseWriter) Write(b []byte) (int, error) {
	// Call the original Write method
	n, err := rw.ResponseWriter.Write(b)

	// Track how many bytes were written
	rw.bytesWritten += n

	return n, err
}

// Logger is a middleware function that logs all HTTP requests.
// It logs two messages for each request:
//  1. When the request starts (request received)
//  2. When the request completes (response sent)
//
// Information logged:
//   - Request method (GET, POST, etc.)
//   - Request path and query parameters
//   - Response status code
//   - Time taken to process the request
//   - Number of bytes sent
//   - Client IP address and user agent
//
// Example usage in your server setup:
//
//	handler := middleware.Logger(yourMux)
//	http.ListenAndServe(":8080", handler)
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record when the request started
		start := time.Now()

		// Wrap the ResponseWriter so we can capture status code and bytes written
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default to 200 if handler doesn't set it
		}

		// Try to get the request ID from context
		// This will be "unknown" if RequestID middleware wasn't used
		requestID := "unknown"
		if id, ok := r.Context().Value(constants.RequestIDKey).(string); ok {
			requestID = id
		}

		// Log that the request has started
		slog.Debug("request received from client",
			"requestId", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"content_type", r.Header.Get("Content-Type"),
			"user_agent", r.UserAgent(),
		)

		slog.Info("Request started",
			"requestId", requestID, // Unique ID for tracking
			"method", r.Method, // HTTP method (GET, POST, etc.)
			"path", r.URL.Path, // Endpoint path (/api/users)
			"query", r.URL.RawQuery, // Query parameters (id=123&sort=name)
			"remoteAddr", r.RemoteAddr, // Client IP address
			"userAgent", r.UserAgent(), // Browser/client information
		)

		// Process the request by calling the next handler
		// This is where your actual handler code runs
		slog.Debug("invoking request handler", "requestId", requestID, "path", r.URL.Path)
		next.ServeHTTP(wrapped, r)
		slog.Debug("handler execution completed", "requestId", requestID, "status", wrapped.statusCode)

		// Calculate how long the request took
		duration := time.Since(start)

		// Determine the log level based on the HTTP status code
		// Errors get logged at higher severity levels for easier monitoring
		logLevel := slog.LevelInfo
		if wrapped.statusCode >= 500 {
			// Server errors (500-599) are logged as errors
			logLevel = slog.LevelError
		} else if wrapped.statusCode >= 400 {
			// Client errors (400-499) are logged as warnings
			logLevel = slog.LevelWarn
		}

		// Log that the request has completed
		// Calculate duration in milliseconds with decimal precision
		// duration.Milliseconds() returns int64 which truncates to 0 for fast requests
		// duration.Seconds() * 1000 gives us fractional milliseconds
		// Round to 4 decimal places for consistent formatting (e.g., 0.1234 ms)
		durationMs := float64(int(duration.Seconds()*1000*10000)) / 10000

		slog.Log(r.Context(), logLevel, "Request completed",
			"requestId", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode, // HTTP status code (200, 404, etc.)
			"statusText", http.StatusText(wrapped.statusCode), // Human-readable status (OK, Not Found, etc.)
			"durationMs", durationMs, // Time taken in milliseconds with 4 decimal places
			"bytesWritten", wrapped.bytesWritten, // Size of response sent
		)
	})
}

// LoggerConfig allows customizing the logger middleware behavior.
// This is useful when you want to skip logging for certain paths (like health checks)
// or when you have other special requirements.
type LoggerConfig struct {
	// SkipPaths is a list of paths to skip logging for.
	// Example: []string{"/health", "/metrics", "/ping"}
	// This is useful for health check endpoints that get called very frequently
	// and would otherwise spam your logs.
	SkipPaths []string

	// LogRequestBody controls whether to log request body content.
	// WARNING: Be very careful with this! Don't log sensitive data like passwords.
	// Usually you should leave this false.
	LogRequestBody bool

	// LogResponseBody controls whether to log response body content.
	// WARNING: This can make logs very large. Usually leave this false.
	LogResponseBody bool
}

// LoggerWithConfigFunc creates a customizable logger middleware.
// Use this when you need more control than the standard Logger provides.
//
// Example usage:
//
//	config := middleware.LoggerConfig{
//	  SkipPaths: []string{"/health", "/metrics"},  // Don't log these endpoints
//	}
//	handler := middleware.LoggerWithConfigFunc(config)(yourMux)
//	http.ListenAndServe(":8080", handler)
func LoggerWithConfigFunc(config LoggerConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if this path should skip logging
			for _, path := range config.SkipPaths {
				if r.URL.Path == path {
					// Skip logging, just pass through to next handler
					next.ServeHTTP(w, r)
					return
				}
			}

			// Same logging logic as the standard Logger
			start := time.Now()

			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			requestID := "unknown"
			if id, ok := r.Context().Value(constants.RequestIDKey).(string); ok {
				requestID = id
			}

			slog.Debug("request received from client (config mode)",
				"requestId", requestID,
				"method", r.Method,
				"path", r.URL.Path,
			)

			slog.Info("Request started",
				"requestId", requestID,
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"remoteAddr", r.RemoteAddr,
				"userAgent", r.UserAgent(),
			)

			slog.Debug("processing request in handler", "requestId", requestID)
			next.ServeHTTP(wrapped, r)
			slog.Debug("request processing complete", "requestId", requestID, "status", wrapped.statusCode)

			duration := time.Since(start)

			logLevel := slog.LevelInfo
			if wrapped.statusCode >= 500 {
				logLevel = slog.LevelError
			} else if wrapped.statusCode >= 400 {
				logLevel = slog.LevelWarn
			}

			// Calculate duration with decimal precision
			// Round to 4 decimal places for consistent formatting (e.g., 0.1234 ms)
			durationMs := float64(int(duration.Seconds()*1000*10000)) / 10000

			slog.Log(r.Context(), logLevel, "Request completed",
				"requestId", requestID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.statusCode,
				"statusText", http.StatusText(wrapped.statusCode),
				"durationMs", durationMs,
				"bytesWritten", wrapped.bytesWritten,
			)
		})
	}
}
