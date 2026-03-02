// Package app wires together application dependencies
// and manages application lifecycle.
package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xanderbilla/my-project/internal/config"
	"github.com/xanderbilla/my-project/internal/server"
)

// App represents the application.
type App struct {
	cfg    *config.Config
	server *http.Server
}

// New creates a new application instance.
func New(cfg *config.Config) *App {
	httpServer := server.NewHTTPServer(cfg)

	return &App{
		cfg:    cfg,
		server: httpServer,
	}
}

// Run starts the application and handles graceful shutdown.
func (a *App) Run() error {
	// Start server in background.
	go func() {
		slog.Info("starting HTTP server", slog.String("address", a.server.Addr))

		if err := a.server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Wait for termination signal.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return a.server.Shutdown(ctx)
}