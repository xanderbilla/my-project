package config

import (
	"os"
	"testing"
)

// TestLoadFromEnv tests loading configuration from environment variables
func TestLoadFromEnv(t *testing.T) {
	if err := os.Setenv("ENV", "test"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("HTTP_ADDRESS", "localhost:9999"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("STORAGE_PATH", "test-storage.db"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Unsetenv("ENV")
		_ = os.Unsetenv("HTTP_ADDRESS")
		_ = os.Unsetenv("STORAGE_PATH")
	}()

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv should not return error: %v", err)
	}

	if cfg.Env != "test" {
		t.Errorf("Expected Env to be 'test', got '%s'", cfg.Env)
	}

	if cfg.HTTPServer.Address != "localhost:9999" {
		t.Errorf("Expected HTTP address to be 'localhost:9999', got '%s'", cfg.HTTPServer.Address)
	}

	if cfg.StoragePath != "test-storage.db" {
		t.Errorf("Expected storage path to be 'test-storage.db', got '%s'", cfg.StoragePath)
	}
}

// TestLoadFromEnv_MissingRequired tests error when required fields are missing
func TestLoadFromEnv_MissingRequired(t *testing.T) {
	_ = os.Unsetenv("ENV")
	_ = os.Unsetenv("HTTP_ADDRESS")
	_ = os.Unsetenv("STORAGE_PATH")

	_, err := LoadFromEnv()
	if err == nil {
		t.Error("LoadFromEnv should return error when required fields are missing")
	}
}

// TestLoadConfig_InvalidPath tests error handling for invalid paths
func TestLoadConfig_InvalidPath(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/to/config.yaml")
	if err == nil {
		t.Error("Expected error when loading from invalid path")
	}
}
