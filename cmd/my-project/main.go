// Package main is the application entry point.
// It initializes configuration, dependencies, and starts the HTTP server
// with graceful shutdown support.
package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/xanderbilla/my-project/internal/app"
	"github.com/xanderbilla/my-project/internal/config"
)

func main() {
	// Load application configuration.
	cfg, err := config.LoadConfig("config/dev.yaml") //CONFIG_PATH should not be hardcoded
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Initialize and run the application.
	application := app.New(cfg)

	if err := application.Run(); err != nil {
		slog.Error("application terminated unexpectedly", slog.Any("error", err))
		os.Exit(1)
	}
}