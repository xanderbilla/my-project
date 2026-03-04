package app

import (
	"testing"

	"github.com/xanderbilla/my-project/internal/config"
)

// TestNew validates the New constructor creates a valid App instance
func TestNew(t *testing.T) {
	cfg := &config.Config{
		HTTPServer: config.HTTPServer{
			Address: ":8080",
		},
		StoragePath: "./testdata",
	}

	app := New(cfg)

	if app == nil {
		t.Fatal("expected non-nil App instance")
	}

	if app.cfg == nil {
		t.Error("expected non-nil config")
	}

	if app.server == nil {
		t.Error("expected non-nil server")
	}

	if app.cfg.HTTPServer.Address != ":8080" {
		t.Errorf("expected address :8080, got %s", app.cfg.HTTPServer.Address)
	}
}

// TestAppStructure validates the App struct fields
func TestAppStructure(t *testing.T) {
	cfg := &config.Config{
		HTTPServer: config.HTTPServer{
			Address: ":9000",
		},
		Env:         "test",
		StoragePath: "/tmp/test",
	}

	app := New(cfg)

	// Validate cfg is properly stored
	if app.cfg.Env != "test" {
		t.Errorf("expected env 'test', got %s", app.cfg.Env)
	}

	if app.cfg.StoragePath != "/tmp/test" {
		t.Errorf("expected storage path '/tmp/test', got %s", app.cfg.StoragePath)
	}

	if app.cfg.HTTPServer.Address != ":9000" {
		t.Errorf("expected address :9000, got %s", app.cfg.HTTPServer.Address)
	}
}

// TestNewWithDifferentConfigs validates New works with various configs
func TestNewWithDifferentConfigs(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name: "ValidDevConfig",
			cfg: &config.Config{
				Env: "dev",
				HTTPServer: config.HTTPServer{
					Address: ":8080",
				},
				StoragePath: "./data",
			},
			wantErr: false,
		},
		{
			name: "ValidProdConfig",
			cfg: &config.Config{
				Env: "prod",
				HTTPServer: config.HTTPServer{
					Address: ":80",
				},
				StoragePath: "/var/lib/myapp",
			},
			wantErr: false,
		},
		{
			name: "CustomPort",
			cfg: &config.Config{
				Env: "test",
				HTTPServer: config.HTTPServer{
					Address: ":3000",
				},
				StoragePath: "./test-data",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New(tt.cfg)

			if app == nil {
				t.Fatal("expected non-nil App instance")
			}

			if app.cfg.Env != tt.cfg.Env {
				t.Errorf("expected env %s, got %s", tt.cfg.Env, app.cfg.Env)
			}

			if app.server.Addr != tt.cfg.HTTPServer.Address {
				t.Errorf("expected server address %s, got %s", tt.cfg.HTTPServer.Address, app.server.Addr)
			}
		})
	}
}
