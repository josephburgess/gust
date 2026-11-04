package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineCityName(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		defaultCity string
		expected    string
	}{
		{
			name:        "Use args when provided",
			args:        []string{"New", "York"},
			defaultCity: "Berlin",
			expected:    "New York",
		},
		{
			name:        "Use default city when no args",
			args:        []string{},
			defaultCity: "Berlin",
			expected:    "Berlin",
		},
		{
			name:        "Return empty when no values provided",
			args:        []string{},
			defaultCity: "",
			expected:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := determineCityName(tc.args, tc.defaultCity)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestHandleMissingCity(t *testing.T) {
	err := handleMissingCity()

	assert.Error(t, err)
	assert.Equal(t, "no city provided", err.Error())
}
