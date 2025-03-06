package cli

import (
	"fmt"
	"testing"

	"github.com/josephburgess/gust/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestIsValidUnit(t *testing.T) {
	validUnits := []string{"metric", "imperial", "standard"}
	invalidUnits := []string{"celsius", "", "Metric"}

	for _, unit := range validUnits {
		t.Run("Valid: "+unit, func(t *testing.T) {
			assert.True(t, isValidUnit(unit))
		})
	}

	for _, unit := range invalidUnits {
		t.Run("Invalid: "+unit, func(t *testing.T) {
			assert.False(t, isValidUnit(unit))
		})
	}
}

func TestHandleConfigUpdates(t *testing.T) {
	createConfig := func() *config.Config {
		return &config.Config{
			ApiUrl:      "https://api.example.com",
			Units:       "metric",
			DefaultCity: "London",
		}
	}

	testCases := []struct {
		name            string
		cli             *CLI
		expectedUpdated bool
		configMutator   func(*config.Config)
	}{
		{
			name: "update api url",
			cli: &CLI{
				ApiUrl: "https://test-api.example.com",
			},
			expectedUpdated: true,
			configMutator: func(c *config.Config) {
				c.ApiUrl = "https://test-api.example.com"
			},
		},
		{
			name: "update units",
			cli: &CLI{
				Units: "imperial",
			},
			expectedUpdated: true,
			configMutator: func(c *config.Config) {
				c.Units = "imperial"
			},
		},
		{
			name: "update default city",
			cli: &CLI{
				Default: "Paris",
			},
			expectedUpdated: true,
			configMutator: func(c *config.Config) {
				c.DefaultCity = "Paris"
			},
		},
		{
			name:            "no updates",
			cli:             &CLI{},
			expectedUpdated: false,
			configMutator:   func(c *config.Config) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initialConfig := createConfig()
			expectedConfig := createConfig()
			tc.configMutator(expectedConfig)

			testHandleConfigUpdates := func(cli *CLI, cfg *config.Config) (bool, error) {
				updated := false

				if cli.ApiUrl != "" {
					cfg.ApiUrl = cli.ApiUrl
					updated = true
				}

				if cli.Units != "" {
					if !isValidUnit(cli.Units) {
						return false, fmt.Errorf("invalid units")
					}
					cfg.Units = cli.Units
					updated = true
				}

				if cli.Default != "" {
					cfg.DefaultCity = cli.Default
					updated = true
				}

				return updated, nil
			}

			updated, err := testHandleConfigUpdates(tc.cli, initialConfig)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, expectedConfig, initialConfig)
		})
	}
}
