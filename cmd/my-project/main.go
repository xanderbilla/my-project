// Package main is the application entry point.
// It initializes configuration, dependencies, and starts the HTTP server
// with graceful shutdown support.
package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/xanderbilla/my-project/internal/app"
	"github.com/xanderbilla/my-project/internal/config"
)

func main() {
	// Load environment variables from .env file (if it exists)
	// This is optional and typically used for local development
	// In production, environment variables are set by the deployment platform
	// Errors are ignored if .env doesn't exist
	_ = godotenv.Load()
	// Configure structured logging with slog
	// This sets up the default logger that all your code will use
	//
	// Log levels work hierarchically - each level shows itself and all higher severity levels:
	//   DEBUG → Shows: DEBUG + INFO + WARN + ERROR (everything, most verbose)
	//   INFO  → Shows: INFO + WARN + ERROR (default, production-ready)
	//   WARN  → Shows: WARN + ERROR (warnings and errors only)
	//   ERROR → Shows: ERROR only (least verbose, errors only)
	//
	// Examples:
	//   LOG_LEVEL=DEBUG  # Development: see everything including debug details
	//   LOG_LEVEL=INFO   # Production: normal operation logs + warnings + errors
	//   LOG_LEVEL=WARN   # Production: only problems (warnings and errors)
	//   LOG_LEVEL=ERROR  # Production: only critical failures
	logLevel := slog.LevelInfo // Default to INFO (shows INFO, WARN, ERROR)

	// Check if user set a custom log level via environment variable
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		switch envLevel {
		case "DEBUG":
			logLevel = slog.LevelDebug
		case "INFO":
			logLevel = slog.LevelInfo
		case "WARN", "WARNING":
			logLevel = slog.LevelWarn
		case "ERROR":
			logLevel = slog.LevelError
		default:
			log.Printf("Invalid LOG_LEVEL '%s', using INFO. Valid values: DEBUG, INFO, WARN, ERROR", envLevel)
		}
	}

	// Create appropriate handler for structured logging
	// For development: Use TextHandler for human-readable output
	// For production: Use JSONHandler for log aggregation tools
	var handler slog.Handler

	// Check ENV to determine if we're in development or production
	// Default to production (JSON) for safety
	env := os.Getenv("ENV")
	if env == "" {
		env = "prod" // Default to production
	}

	if env == "dev" || env == "development" {
		// Development: Human-readable text logs
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	} else {
		// Production: JSON logs for aggregation
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	}

	// Set this as the default logger with base fields for consistency
	// This ensures ALL logs (including from other packages) have service/env fields
	logger := slog.New(handler).With(
		slog.String("service", "my-project"),
		slog.String("env", env),
	)
	slog.SetDefault(logger)

	// Log startup sequence with DEBUG details
	slog.Debug("initializing application", "go_version", "1.26", "log_format", fmt.Sprintf("%T", handler))
	slog.Info("logger initialized", "log_level", logLevel.String())
	slog.Debug("logger configuration complete")

	// Load application configuration
	// Configuration can be loaded from:
	//   1. YAML file (if CONFIG_PATH environment variable is set)
	//   2. Environment variables only (if CONFIG_PATH is not set)
	//
	// All YAML values can be overridden by environment variables:
	//   ENV, STORAGE_PATH, HTTP_ADDRESS
	//
	// See .env.example for all available environment variables
	configPath := os.Getenv("CONFIG_PATH")

	var cfg *config.Config
	var err error

	slog.Debug("starting configuration load", "config_path", configPath)

	if configPath != "" {
		// Load from YAML file with environment variable overrides
		slog.Info("loading configuration from file", "path", configPath)
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			slog.Error("failed to load configuration from file", "error", err, "path", configPath)
			log.Fatalf("failed to load configuration from %s: %v", configPath, err)
		}
		slog.Debug("configuration loaded successfully from file", "env", cfg.Env, "http_address", cfg.HTTPServer.Address)
	} else {
		// Load from environment variables only
		slog.Info("loading configuration from environment variables")
		cfg, err = config.LoadFromEnv()
		if err != nil {
			slog.Error("failed to load configuration from environment", "error", err)
			log.Fatalf("failed to load configuration from environment: %v", err)
		}
		slog.Debug("configuration loaded successfully", "http_address", cfg.HTTPServer.Address, "storage_path", cfg.StoragePath)
	}

	// Create the application with all dependencies
	// This sets up:
	//   - HTTP server with all middleware
	//   - Route handlers
	//   - Graceful shutdown handling
	slog.Debug("initializing application components", "config_env", cfg.Env)
	application := app.New(cfg)
	slog.Debug("application initialized successfully")

	// Start the application
	// This blocks until the server shuts down (on SIGINT or SIGTERM)
	// The server handles graceful shutdown automatically:
	//   - Stops accepting new connections
	//   - Waits for active requests to complete (up to 5 seconds)
	//   - Then exits
	slog.Info("starting application")
	if err := application.Run(); err != nil {
		slog.Error("application terminated unexpectedly", "error", err)
		os.Exit(1)
	}
	slog.Info("application shutdown complete")
}
