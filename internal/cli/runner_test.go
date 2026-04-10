package cli

import (
	"testing"

	"github.com/josephburgess/gust/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRunMissingCity(t *testing.T) {
	tests := []struct {
		name        string
		cli         *CLI
		defaultCity string
		wantErr     bool
	}{
		{"no city provided", &CLI{}, "", true},
		{"city via flag", &CLI{City: "London"}, "", false},
		{"city via args", &CLI{Args: []string{"New", "York"}}, "", false},
		{"default city", &CLI{}, "Paris", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{DefaultCity: tt.defaultCity}
			city := determineCityName(tt.cli.City, tt.cli.Args, cfg.DefaultCity)
			var err error
			if city == "" {
				err = handleMissingCity()
			}
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "no city provided")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRunAuthRequired(t *testing.T) {
	err := handleMissingAuth()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication required")
}
