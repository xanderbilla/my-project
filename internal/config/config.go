// Package config handles loading and validating application configuration.
//
// Configuration precedence:
//  1. YAML file values
//  2. Environment variables (override YAML)
//  3. env-default tag values
//
// Required environment variables are enforced using `env-required:"true"`.
package config

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

// HTTPServer contains HTTP server configuration.
type HTTPServer struct {
	// Address is the bind address for the HTTP server.
	// Example: ":8080" or "0.0.0.0:8080".
	Address string `yaml:"address" env:"HTTP_ADDRESS" env-required:"true"`
}

// Config represents the root application configuration.
type Config struct {
	// Env defines the runtime environment (e.g., dev, prod).
	// Defaults to "prod" if not provided.
	Env string `yaml:"env" env:"ENV" env-default:"prod"`

	// StoragePath is the filesystem path for persistent storage.
	// Must be provided either in YAML or via STORAGE_PATH.
	StoragePath string `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`

	// HTTPServer groups HTTP server-related settings.
	HTTPServer HTTPServer `yaml:"http-server"`
}

// Load reads configuration from the provided YAML file path.
// It applies environment variable overrides automatically.
//
// Returns an error if:
//   - path is empty
//   - YAML parsing fails
//   - required fields are missing
//   - validation fails
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		slog.Error("config path is required but not provided")
		return nil, fmt.Errorf("config path is required")
	}

	slog.Debug("reading YAML configuration file", "path", path)

	var cfg Config

	// Read YAML and apply environment overrides.
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		slog.Error("failed to read YAML configuration", "path", path, "error", err)
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	slog.Debug("YAML configuration parsed successfully", "env", cfg.Env, "http_address", cfg.HTTPServer.Address)

	// Perform additional logical validation.
	if err := cfg.Validate(); err != nil {
		slog.Error("configuration validation failed", "error", err)
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Warn if using default env
	if cfg.Env == "prod" {
		slog.Warn("using default environment", "env", "prod", "recommendation", "explicitly set ENV variable")
	}

	slog.Debug("configuration validation passed")
	return &cfg, nil
}

// LoadFromEnv loads configuration from environment variables only.
// No YAML file is required when using this method.
//
// This is useful for:
//   - Production deployments (using platform environment variables)
//   - Docker containers (using docker-compose or Kubernetes secrets)
//   - CI/CD pipelines
//
// All required environment variables must be set:
//   - HTTP_ADDRESS
//   - STORAGE_PATH
//
// Optional environment variables:
//   - ENV (defaults to "prod")
//
// Returns an error if:
//   - required environment variables are missing
//   - validation fails
func LoadFromEnv() (*Config, error) {
	slog.Debug("reading configuration from environment variables")

	var cfg Config

	// Read environment variables and apply defaults.
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		slog.Error("failed to read environment variables", "error", err)
		return nil, fmt.Errorf("failed to read environment variables: %w", err)
	}

	slog.Debug("environment variables parsed successfully", "http_address", cfg.HTTPServer.Address, "storage_path", cfg.StoragePath)

	// Perform additional logical validation.
	if err := cfg.Validate(); err != nil {
		slog.Error("configuration validation failed", "error", err)
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Warn if using default environment
	if cfg.Env == "prod" {
		slog.Warn("using default environment", "env", "prod", "recommendation", "explicitly set ENV variable")
	}

	slog.Debug("configuration validation passed")
	return &cfg, nil
}

// Validate performs additional logical validation
// that cannot be expressed via struct tags.
func (c *Config) Validate() error {
	if c.HTTPServer.Address == "" {
		return fmt.Errorf("http-server.address cannot be empty")
	}

	return nil
}
