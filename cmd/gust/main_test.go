package main

import (
	"os"
	"testing"
)

func TestEnvVarLoading(t *testing.T) {
	os.Setenv("OPENWEATHER_API_KEY", "test-api-key")
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got '%s'", apiKey)
	}
	os.Unsetenv("OPENWEATHER_API_KEY")
}
