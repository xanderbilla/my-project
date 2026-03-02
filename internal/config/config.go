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
		return nil, fmt.Errorf("config path is required")
	}

	var cfg Config

	// Read YAML and apply environment overrides.
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Perform additional logical validation.
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

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
