// Package server configures and initializes HTTP servers.
package server

import (
	"net/http"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/xanderbilla/my-project/internal/config"
	"github.com/xanderbilla/my-project/internal/handlers"
	"github.com/xanderbilla/my-project/internal/middleware"
)

// NewHTTPServer creates a configured HTTP server instance with all middleware.
//
// Middleware is applied in the following order (from outer to inner):
//  1. Recovery - Catches panics and returns proper error responses
//  2. Logger - Logs all requests and responses
//  3. RequestID - Adds unique ID to each request
//  4. Routes - Your actual endpoint handlers
//
// Why this order matters:
//   - Recovery must be outermost to catch panics from all other middleware
//   - Logger should be early to log all requests (including failed ones)
//   - RequestID should be early so all logs include the request ID
//
// Example request flow:
//
//	Request → Recovery → Logger → RequestID → Handler → Response
func NewHTTPServer(cfg *config.Config) *http.Server {
	// Create a new HTTP multiplexer (router)
	// This will map URLs to handler functions
	mux := http.NewServeMux()

	// Register all user-related routes
	// This sets up endpoints like GET /api/users, POST /api/users, etc.
	handlers.UserRoutes(mux)

	// Wrap the mux with middleware layers
	// Think of this like wrapping gifts - each layer wraps the previous one
	//
	// CRITICAL: The order matters!
	// RequestID MUST be applied BEFORE Logger so the logger can access the request ID from context.
	//
	// Correct order:
	//   Request
	//     → Recovery (outermost - catches any panics)
	//       → RequestID (adds unique ID to context)
	//         → Logger (reads ID from context and logs)
	//           → Your Handler (processes request)
	//         → Logger (logs completion with same ID)
	//       → RequestID (cleanup if needed)
	//     → Recovery (handles any panics)
	//   Response
	//
	// If Logger comes before RequestID, the logger will see requestID="unknown"
	// because the ID hasn't been added to context yet!
	handler := middleware.Recovery( // 1. Catch panics (outermost)
		middleware.RequestID( // 2. Add request ID first!
			middleware.Logger(mux), // 3. Log with the ID available
		),
	)

	// Add gzip compression (60-80% smaller responses)
	handler = gziphandler.GzipHandler(handler)

	// Create and configure the HTTP server
	return &http.Server{
		// Server address from configuration (e.g., ":8080")
		Addr: cfg.HTTPServer.Address,

		// Handler is the root handler (with all middleware applied)
		Handler: handler,

		// ReadTimeout is how long to wait for client to send request
		// Prevents slow clients from holding connections open forever
		ReadTimeout: 10 * time.Second,

		// ReadHeaderTimeout prevents Slowloris attacks
		ReadHeaderTimeout: 5 * time.Second,

		// WriteTimeout is how long to wait while sending response to client
		// Prevents slow clients from holding connections open forever
		WriteTimeout: 10 * time.Second,

		// IdleTimeout is how long to keep idle connections open
		// Allows connection reuse while preventing resource exhaustion
		IdleTimeout: 60 * time.Second,

		// MaxHeaderBytes limits request header size (prevent large header attacks)
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
}
