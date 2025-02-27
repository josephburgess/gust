package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndLoadAuthConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gust-auth-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalGetAuthConfigPath := GetAuthConfigPath
	defer func() { GetAuthConfigPath = originalGetAuthConfigPath }()

	GetAuthConfigPath = func() (string, error) {
		return filepath.Join(tempDir, "auth.json"), nil
	}

	testAuth := &AuthConfig{
		APIKey:     "test-api-key-123",
		ServerURL:  "https://test-server.example.com",
		LastAuth:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		GithubUser: "testuser",
	}

	err = SaveAuthConfig(testAuth)
	if err != nil {
		t.Fatalf("Failed to save auth config: %v", err)
	}

	loadedAuth, err := LoadAuthConfig()
	if err != nil {
		t.Fatalf("Failed to load auth config: %v", err)
	}

	if loadedAuth.APIKey != testAuth.APIKey {
		t.Errorf("Expected APIKey to be %s, got %s", testAuth.APIKey, loadedAuth.APIKey)
	}

	if loadedAuth.ServerURL != testAuth.ServerURL {
		t.Errorf("Expected ServerURL to be %s, got %s", testAuth.ServerURL, loadedAuth.ServerURL)
	}

	if !loadedAuth.LastAuth.Equal(testAuth.LastAuth) {
		t.Errorf("Expected LastAuth to be %v, got %v", testAuth.LastAuth, loadedAuth.LastAuth)
	}

	if loadedAuth.GithubUser != testAuth.GithubUser {
		t.Errorf("Expected GithubUser to be %s, got %s", testAuth.GithubUser, loadedAuth.GithubUser)
	}
}

func TestLoadNonExistentAuthConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gust-auth-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalGetAuthConfigPath := GetAuthConfigPath
	defer func() { GetAuthConfigPath = originalGetAuthConfigPath }()

	GetAuthConfigPath = func() (string, error) {
		return filepath.Join(tempDir, "nonexistent-auth.json"), nil
	}

	authConfig, err := LoadAuthConfig()
	if err != nil {
		t.Fatalf("Expected no error loading non-existent auth config, got: %v", err)
	}

	if authConfig != nil {
		t.Fatalf("Expected nil auth config when file doesn't exist, got: %+v", authConfig)
	}
}

func TestDefaultGetAuthConfigPath(t *testing.T) {
	path, err := defaultGetAuthConfigPath()
	if err != nil {
		t.Fatalf("Failed to get auth config path: %v", err)
	}

	if filepath.Base(path) != "auth.json" {
		t.Errorf("Expected auth filename to be auth.json, got %s", filepath.Base(path))
	}

	if !filepath.IsAbs(path) {
		t.Errorf("Expected absolute path, got %s", path)
	}

	if !filepath.IsAbs(path) || !contains(path, filepath.Join(".config", "gust")) {
		t.Errorf("Expected path to contain .config/gust, got %s", path)
	}
}

func contains(path, substr string) bool {
	return filepath.ToSlash(path) == filepath.ToSlash(substr) ||
		contains2(filepath.ToSlash(path), filepath.ToSlash(substr))
}

func contains2(path, substr string) bool {
	for i := 0; i <= len(path)-len(substr); i++ {
		if path[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
