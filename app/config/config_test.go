package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Default(t *testing.T) {
	cfg := LoadConfig("")

	if cfg == nil {
		t.Fatal("expected config to be loaded, got nil")
	}

	if cfg.ServerAddress != ":8080" {
		t.Errorf("expected ServerAddress to be :8080, got %s", cfg.ServerAddress)
	}
}

func TestLoadConfig_EnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "9090")

	defer func() {
		os.Unsetenv("PORT")
	}()

	cfg := LoadConfig("")

	if cfg.ServerAddress != ":9090" {
		t.Errorf("expected ServerAddress to be :9090, got %s", cfg.ServerAddress)
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Test with non-existent config file
	cfg := LoadConfig("/non/existent/file.yaml")

	if cfg == nil {
		t.Fatal("expected config to be loaded with defaults, got nil")
	}

	// Should fall back to defaults
	if cfg.ServerAddress != ":8080" {
		t.Errorf("expected ServerAddress to be :8080, got %s", cfg.ServerAddress)
	}
}

func TestEnvironment_Struct(t *testing.T) {
	env := Environment{
		Port: "8080",
	}

	if env.Port != "8080" {
		t.Errorf("expected Port to be 8080, got %s", env.Port)
	}
}
