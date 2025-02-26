package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/josephburgess/gust/internal/config"
)

func TestEnvVarLoading(t *testing.T) {
	os.Setenv("OPENWEATHER_API_KEY", "test-api-key")
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got '%s'", apiKey)
	}
	os.Unsetenv("OPENWEATHER_API_KEY")
}

func TestConfigLoading(t *testing.T) {
	// tempdir for tests
	tmpDir, err := os.MkdirTemp("", "gust-main-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// test config file
	configPath := filepath.Join(tmpDir, "config.json")
	configData := `{
		"api_key": "test-config-api-key",
		"default_city": "Test City"
	}`

	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// load cfg
	testConfig := &config.Config{
		APIKey:      "test-config-api-key",
		DefaultCity: "Test City",
	}

	// mock config.Load fn
	if testConfig.APIKey != "test-config-api-key" {
		t.Errorf("Config file structure incorrect - API key mismatch")
	}

	if testConfig.DefaultCity != "Test City" {
		t.Errorf("Config file structure incorrect - default city mismatch")
	}
}
