package server

import (
	"testing"
	"time"

	"github.com/xanderbilla/my-project/internal/config"
)

// TestNewHTTPServer validates the server creation and configuration
func TestNewHTTPServer(t *testing.T) {
	cfg := &config.Config{
		HTTPServer: config.HTTPServer{
			Address: ":8080",
		},
	}

	server := NewHTTPServer(cfg)

	if server == nil {
		t.Fatal("expected non-nil server")
	}

	if server.Addr != ":8080" {
		t.Errorf("expected address :8080, got %s", server.Addr)
	}

	if server.Handler == nil {
		t.Error("expected non-nil handler")
	}
}

// TestNewHTTPServer_Timeouts validates server timeout configurations
func TestNewHTTPServer_Timeouts(t *testing.T) {
	cfg := &config.Config{
		HTTPServer: config.HTTPServer{
			Address: ":9000",
		},
	}

	server := NewHTTPServer(cfg)

	expectedReadTimeout := 10 * time.Second
	expectedWriteTimeout := 10 * time.Second
	expectedIdleTimeout := 60 * time.Second

	if server.ReadTimeout != expectedReadTimeout {
		t.Errorf("expected ReadTimeout %v, got %v", expectedReadTimeout, server.ReadTimeout)
	}

	if server.WriteTimeout != expectedWriteTimeout {
		t.Errorf("expected WriteTimeout %v, got %v", expectedWriteTimeout, server.WriteTimeout)
	}

	if server.IdleTimeout != expectedIdleTimeout {
		t.Errorf("expected IdleTimeout %v, got %v", expectedIdleTimeout, server.IdleTimeout)
	}
}

// TestNewHTTPServer_DifferentAddresses validates server works with various addresses
func TestNewHTTPServer_DifferentAddresses(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{
			name:    "Port8080",
			address: ":8080",
		},
		{
			name:    "Port3000",
			address: ":3000",
		},
		{
			name:    "PortWithHost",
			address: "localhost:8080",
		},
		{
			name:    "Port80",
			address: ":80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				HTTPServer: config.HTTPServer{
					Address: tt.address,
				},
			}

			server := NewHTTPServer(cfg)

			if server == nil {
				t.Fatal("expected non-nil server")
			}

			if server.Addr != tt.address {
				t.Errorf("expected address %s, got %s", tt.address, server.Addr)
			}

			if server.Handler == nil {
				t.Error("expected non-nil handler with middleware chain")
			}
		})
	}
}

// TestNewHTTPServer_HandlerNotNil validates that middleware is applied
func TestNewHTTPServer_HandlerNotNil(t *testing.T) {
	cfg := &config.Config{
		HTTPServer: config.HTTPServer{
			Address: ":8080",
		},
	}

	server := NewHTTPServer(cfg)

	// The handler should be wrapped with middleware (Recovery → RequestID → Logger)
	// We can't easily test the middleware order without making requests,
	// but we can verify that a handler exists
	if server.Handler == nil {
		t.Fatal("expected handler to be set with middleware chain")
	}

	// Verify it's an http.Handler (implements ServeHTTP)
	var _ = server.Handler
}
