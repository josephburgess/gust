package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAndSave(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gust-config-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalGetConfigPath := GetConfigPath
	defer func() { GetConfigPath = originalGetConfigPath }()

	GetConfigPath = func() (string, error) {
		return filepath.Join(tempDir, "config.json"), nil
	}

	testConfig := &Config{
		DefaultCity: "Tokyo",
		ApiUrl:      "https://test-api.example.com",
	}

	err = testConfig.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	loadedConfig, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loadedConfig.DefaultCity != testConfig.DefaultCity {
		t.Errorf("Expected DefaultCity to be %s, got %s", testConfig.DefaultCity, loadedConfig.DefaultCity)
	}

	if loadedConfig.ApiUrl != testConfig.ApiUrl {
		t.Errorf("Expected ApiUrl to be %s, got %s", testConfig.ApiUrl, loadedConfig.ApiUrl)
	}
}

func TestLoadNonExistentConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gust-config-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalGetConfigPath := GetConfigPath
	defer func() { GetConfigPath = originalGetConfigPath }()

	GetConfigPath = func() (string, error) {
		return filepath.Join(tempDir, "nonexistent-config.json"), nil
	}

	config, err := Load()
	if err != nil {
		t.Fatalf("Expected no error loading non-existent config, got: %v", err)
	}

	if config == nil {
		t.Fatal("Expected empty config object, got nil")
	}

	if config.DefaultCity != "" {
		t.Errorf("Expected empty DefaultCity, got %s", config.DefaultCity)
	}

	if config.ApiUrl != "" {
		t.Errorf("Expected empty ApiUrl, got %s", config.ApiUrl)
	}
}

func TestDefaultGetConfigPath(t *testing.T) {
	path, err := defaultGetConfigPath()
	if err != nil {
		t.Fatalf("Failed to get default config path: %v", err)
	}

	if filepath.Base(path) != "config.json" {
		t.Errorf("Expected config filename to be config.json, got %s", filepath.Base(path))
	}

	if !filepath.IsAbs(path) {
		t.Errorf("Expected absolute path, got %s", path)
	}
}
