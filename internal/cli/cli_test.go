package cli

import (
	"testing"
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
			name:        "Use city flag",
			cityFlag:    "London",
			args:        []string{},
			defaultCity: "Paris",
			expected:    "London",
		},
		{
			name:        "Use args when no flag",
			cityFlag:    "",
			args:        []string{"New", "York"},
			defaultCity: "Paris",
			expected:    "New York",
		},
		{
			name:        "Use default when no flag or args",
			cityFlag:    "",
			args:        []string{},
			defaultCity: "Paris",
			expected:    "Paris",
		},
		{
			name:        "Handle multi-word args",
			cityFlag:    "",
			args:        []string{"San", "Francisco", "CA"},
			defaultCity: "Paris",
			expected:    "San Francisco CA",
		},
		{
			name:        "Handle empty everything",
			cityFlag:    "",
			args:        []string{},
			defaultCity: "",
			expected:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := determineCityName(tc.cityFlag, tc.args, tc.defaultCity)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}
