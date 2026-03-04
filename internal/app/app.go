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
	slog.Debug("initializing server goroutine")

	// Start server in background.
	go func() {
		slog.Info("server starting", "address", a.server.Addr)
		slog.Debug("server listening for incoming connections", "address", a.server.Addr, "protocol", "http")

		if err := a.server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed to start", "error", err, "address", a.server.Addr)
			slog.Warn("💡 port might be in use", "tip", "run 'make run-quick' to auto-kill old server")
			os.Exit(1)
		}
		slog.Debug("server stopped accepting connections")
	}()

	slog.Debug("registering signal handlers", "signals", []string{"SIGINT", "SIGTERM"})

	// Wait for termination signal.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	slog.Debug("waiting for shutdown signal")
	<-stop

	slog.Info("shutdown signal received, shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Debug("waiting for active connections to close", "timeout", "5s")
	err := a.server.Shutdown(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Warn("graceful shutdown timeout exceeded, forcing shutdown", "timeout", "5s")
		} else {
			slog.Error("error during shutdown", "error", err)
		}
	} else {
		slog.Debug("shutdown completed successfully")
	}
	return err
}
