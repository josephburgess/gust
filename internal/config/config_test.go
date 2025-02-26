package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigSaveAndLoad(t *testing.T) {
	// temp dir
	tmpDir, err := os.MkdirTemp("", "gust-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// save/restore original function b4 override
	origGetConfigPath := GetConfigPath
	defer func() { GetConfigPath = origGetConfigPath }()

	// override
	GetConfigPath = func() (string, error) {
		return filepath.Join(tmpDir, "config.json"), nil
	}

	// test configuraion
	testConfig := &Config{
		APIKey:      "test-api-key",
		DefaultCity: "Test City",
	}

	// save the test cfg
	if err := testConfig.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// load testcfg
	loadedConfig, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// check loaded matches saved
	if loadedConfig.APIKey != testConfig.APIKey {
		t.Errorf("Expected API key %s, got %s", testConfig.APIKey, loadedConfig.APIKey)
	}

	if loadedConfig.DefaultCity != testConfig.DefaultCity {
		t.Errorf("Expected default city %s, got %s", testConfig.DefaultCity, loadedConfig.DefaultCity)
	}
}

func TestLoadNonExistentConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gust-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origGetConfigPath := GetConfigPath
	defer func() { GetConfigPath = origGetConfigPath }()

	GetConfigPath = func() (string, error) {
		return filepath.Join(tmpDir, "nonexistent.json"), nil
	}

	config, err := Load()
	if err != nil {
		t.Fatalf("Failed to load non-existent config: %v", err)
	}

	if config.APIKey != "" {
		t.Errorf("Expected empty API key, got %s", config.APIKey)
	}

	if config.DefaultCity != "" {
		t.Errorf("Expected empty default city, got %s", config.DefaultCity)
	}
}
