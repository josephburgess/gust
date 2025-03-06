package cli

import (
	"testing"

	"github.com/josephburgess/gust/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNeedsSetup(t *testing.T) {
	testCases := []struct {
		name     string
		cli      *CLI
		cfg      *config.Config
		expected bool
	}{
		{
			name:     "empty default city",
			cli:      &CLI{},
			cfg:      &config.Config{DefaultCity: ""},
			expected: true,
		},
		{
			name:     "setup flag set",
			cli:      &CLI{Setup: true},
			cfg:      &config.Config{DefaultCity: "London"},
			expected: true,
		},
		{
			name:     "city set w/ no setup flag",
			cli:      &CLI{},
			cfg:      &config.Config{DefaultCity: "London"},
			expected: false,
		},
		{
			name:     "empty default and setup flag",
			cli:      &CLI{Setup: true},
			cfg:      &config.Config{DefaultCity: ""},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := needsSetup(tc.cli, tc.cfg)
			assert.Equal(t, tc.expected, result)
		})
	}
}
