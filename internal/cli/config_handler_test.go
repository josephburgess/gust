package cli

import (
	"path/filepath"
	"testing"

	"github.com/josephburgess/gust/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withTempConfig(t *testing.T) {
	t.Helper()
	orig := config.GetConfigPath
	t.Cleanup(func() { config.GetConfigPath = orig })
	config.GetConfigPath = func() (string, error) {
		return filepath.Join(t.TempDir(), "config.json"), nil
	}
}

func TestIsValidUnit(t *testing.T) {
	for _, unit := range []string{"metric", "imperial", "standard"} {
		assert.True(t, isValidUnit(unit), "expected %q to be valid", unit)
	}
	for _, unit := range []string{"celsius", "", "Metric"} {
		assert.False(t, isValidUnit(unit), "expected %q to be invalid", unit)
	}
}

func TestHandleConfigUpdates(t *testing.T) {
	withTempConfig(t)

	newConfig := func() *config.Config {
		return &config.Config{ApiUrl: "https://api.example.com", Units: "metric", DefaultCity: "London"}
	}

	tests := []struct {
		name        string
		cli         *CLI
		wantUpdated bool
		check       func(*testing.T, *config.Config)
	}{
		{
			name:        "update api url",
			cli:         &CLI{ApiUrl: "https://new.example.com"},
			wantUpdated: true,
			check: func(t *testing.T, c *config.Config) {
				assert.Equal(t, "https://new.example.com", c.ApiUrl)
			},
		},
		{
			name:        "update units",
			cli:         &CLI{Units: "imperial"},
			wantUpdated: true,
			check: func(t *testing.T, c *config.Config) {
				assert.Equal(t, "imperial", c.Units)
			},
		},
		{
			name:        "update default city",
			cli:         &CLI{Default: "Paris"},
			wantUpdated: true,
			check: func(t *testing.T, c *config.Config) {
				assert.Equal(t, "Paris", c.DefaultCity)
			},
		},
		{
			name:        "invalid units returns error",
			cli:         &CLI{Units: "kelvin"},
			wantUpdated: false,
			check:       nil,
		},
		{
			name:        "no flags — no update",
			cli:         &CLI{},
			wantUpdated: false,
			check:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newConfig()
			updated, err := handleConfigUpdates(tt.cli, cfg)

			if tt.cli.Units == "kelvin" {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantUpdated, updated)
			if tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}
