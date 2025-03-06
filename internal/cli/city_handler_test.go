package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineCityName(t *testing.T) {
	testCases := []struct {
		name        string
		cityFlag    string
		args        []string
		defaultCity string
		expected    string
	}{
		{
			name:        "Use city flag when provided",
			cityFlag:    "London",
			args:        []string{"Paris"},
			defaultCity: "Berlin",
			expected:    "London",
		},
		{
			name:        "Use args when no flag but args provided",
			cityFlag:    "",
			args:        []string{"New", "York"},
			defaultCity: "Berlin",
			expected:    "New York",
		},
		{
			name:        "Use default city when no flag or args",
			cityFlag:    "",
			args:        []string{},
			defaultCity: "Berlin",
			expected:    "Berlin",
		},
		{
			name:        "Return empty when no values provided",
			cityFlag:    "",
			args:        []string{},
			defaultCity: "",
			expected:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := determineCityName(tc.cityFlag, tc.args, tc.defaultCity)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestHandleMissingCity(t *testing.T) {
	err := handleMissingCity()

	assert.Error(t, err)
	assert.Equal(t, "no city provided", err.Error())
}
