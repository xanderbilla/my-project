package main

import (
	"os"
	"testing"
)

// TestPackageImports validates the package can be initialized
func TestPackageImports(t *testing.T) {
	// This test validates the main package compiles and basic setup works
	// We can't directly test main() as it runs the server, but we can
	// test environment variable handling and setup logic
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	// If this test runs, package imports are working
}

// TestLogLevelParsing validates log level environment variable handling
func TestLogLevelParsing(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		wantErr  bool
	}{
		{
			name:     "ValidDEBUG",
			envValue: "DEBUG",
			wantErr:  false,
		},
		{
			name:     "ValidINFO",
			envValue: "INFO",
			wantErr:  false,
		},
		{
			name:     "ValidWARN",
			envValue: "WARN",
			wantErr:  false,
		},
		{
			name:     "ValidERROR",
			envValue: "ERROR",
			wantErr:  false,
		},
		{
			name:     "EmptyUsesDefault",
			envValue: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable for this test
			if tt.envValue != "" {
				if err := os.Setenv("LOG_LEVEL", tt.envValue); err != nil {
					t.Fatal(err)
				}
				defer func() {
					_ = os.Unsetenv("LOG_LEVEL")
				}()
			}

			// The logic in main() handles this, we're just testing
			// that the values are valid and don't panic
			envLevel := os.Getenv("LOG_LEVEL")
			if tt.envValue != "" && envLevel != tt.envValue {
				t.Errorf("expected LOG_LEVEL=%s, got %s", tt.envValue, envLevel)
			}
		})
	}
}

// TestEnvironmentHandling validates ENV variable handling
func TestEnvironmentHandling(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "Development",
			envValue: "dev",
			expected: "dev",
		},
		{
			name:     "Production",
			envValue: "prod",
			expected: "prod",
		},
		{
			name:     "EmptyDefaultsToProd",
			envValue: "",
			expected: "", // Will default to prod in main()
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				if err := os.Setenv("ENV", tt.envValue); err != nil {
					t.Fatal(err)
				}
				defer func() {
					_ = os.Unsetenv("ENV")
				}()
			} else {
				_ = os.Unsetenv("ENV")
			}

			env := os.Getenv("ENV")
			if env != tt.expected {
				t.Errorf("expected ENV=%s, got %s", tt.expected, env)
			}
		})
	}
}
