// Package server configures and initializes HTTP servers.
package server

import (
	"net/http"
	"time"

	"github.com/xanderbilla/my-project/internal/config"
	"github.com/xanderbilla/my-project/internal/handlers"
)

// NewHTTPServer creates a configured HTTP server instance.
func NewHTTPServer(cfg *config.Config) *http.Server {
	mux := http.NewServeMux()

	// Register routes.
	handlers.GetAllUsers(mux)

	return &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
