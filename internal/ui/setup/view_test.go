package setup

import (
	"strings"
	"testing"

	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestView(t *testing.T) {
	tests := []struct {
		name            string
		setupModel      func() Model
		expectedParts   []string
		unexpectedParts []string
	}{
		{
			name: "city input state",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateCity
				return m
			},
			expectedParts: []string{
				"Enter a default city",
				"You can enter a country code",
				"Simple terminal weather",
			},
		},
		{
			name: "city search state",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateCitySearch
				m.CitySearchQuery = "London"
				return m
			},
			expectedParts: []string{
				"Searching for cities",
				"Looking for \"London\"",
			},
		},
		{
			name: "city select state with results",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateCitySelect
				m.CityOptions = []models.City{
					{Name: "London", Country: "GB"},
					{Name: "Paris", Country: "FR"},
				}
				return m
			},
			expectedParts: []string{
				"Select your town or city",
				"London",
				"Paris",
				"🇬🇧",
				"🇫🇷",
				"Press Enter to select or Esc to search again",
			},
		},
		{
			name: "city select state with empty results",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateCitySelect
				m.CityOptions = []models.City{}
				return m
			},
			expectedParts: []string{
				"No cities found",
				"Please try a different search term",
			},
		},
		{
			name: "units select state",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateUnits
				return m
			},
			expectedParts: []string{
				"Choose your preferred units",
				"metric",
				"imperial",
				"standard",
			},
		},
		{
			name: "view select state",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateView
				return m
			},
			expectedParts: []string{
				"Choose your preferred view",
				"detailed",
				"compact",
				"5-day",
				"24-hour",
				"full",
			},
		},
		{
			name: "tips select state",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateTips
				return m
			},
			expectedParts: []string{
				"Would you like tips shown on daily forecasts",
				"Yes, show weather tips",
				"No, don't show tips",
			},
		},
		{
			name: "auth state",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, true, &api.Client{})
				m.State = StateAuth
				return m
			},
			expectedParts: []string{
				"GitHub Auth",
				"authenticate with GitHub",
				"no permissions requested",
			},
		},
		{
			name: "complete state",
			setupModel: func() Model {
				m := NewModel(&config.Config{
					DefaultCity: "London",
					Units:       "metric",
					DefaultView: "detailed",
					ShowTips:    true,
				}, false, &api.Client{})
				m.State = StateComplete
				return m
			},
			expectedParts: []string{
				"Setup complete",
				"Default city: London",
				"Units: metric",
				"Default view: detailed",
				"Tips enabled",
			},
			unexpectedParts: []string{
				"Tips disabled",
			},
		},
		{
			name: "complete state with tips disabled",
			setupModel: func() Model {
				m := NewModel(&config.Config{
					DefaultCity: "London",
					Units:       "metric",
					DefaultView: "detailed",
					ShowTips:    false,
				}, false, &api.Client{})
				m.State = StateComplete
				return m
			},
			expectedParts: []string{
				"Tips disabled",
			},
			unexpectedParts: []string{
				"Tips enabled",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.setupModel()
			view := model.View()

			for _, expectedPart := range tt.expectedParts {
				assert.Contains(t, view, expectedPart)
			}

			if tt.unexpectedParts != nil {
				for _, unexpectedPart := range tt.unexpectedParts {
					assert.NotContains(t, view, unexpectedPart)
				}
			}
		})
	}
}

func TestCountryEmoji(t *testing.T) {
	tests := []struct {
		name        string
		countryCode string
		want        string
	}{
		{
			name:        "valid country code",
			countryCode: "GB",
			want:        "🇬🇧",
		},
		{
			name:        "empty country code",
			countryCode: "",
			want:        "🌍",
		},
		{
			name:        "invalid length",
			countryCode: "GBR",
			want:        "🌍",
		},
		{
			name:        "lowercase code",
			countryCode: "fr",
			want:        "🇫🇷",
		},
		{
			name:        "single letter code",
			countryCode: "X",
			want:        "🌍",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getCountryEmoji(tt.countryCode)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRenderOptions(t *testing.T) {
	model := NewModel(&config.Config{}, false, nil)
	options := []string{"Option 1", "Option 2", "Option 3"}

	for cursor := 0; cursor < len(options); cursor++ {
		t.Run("cursor at position "+string(rune('0'+cursor)), func(t *testing.T) {
			result := model.renderOptions(options, cursor)

			for i, option := range options {
				if i == cursor {
					assert.Contains(t, result, "→ "+option)
				} else {
					assert.Contains(t, result, "  "+option)
				}
			}

			for _, option := range options {
				assert.Contains(t, result, option)
			}
		})
	}

	t.Run("empty options", func(t *testing.T) {
		result := model.renderOptions([]string{}, 0)
		assert.Empty(t, result)
	})
}

func TestCenterContent(t *testing.T) {
	testCases := []struct {
		name      string
		content   string
		width     int
		height    int
		checkFunc func(*testing.T, string, string)
	}{
		{
			name:    "zero dimensions",
			content: "test content\nmultiple lines\ndiff lengths",
			width:   0,
			height:  0,
			checkFunc: func(t *testing.T, original, result string) {
				assert.Equal(t, original, result, "should with zero dimensions")
			},
		},
		{
			name:    "normal dimensions",
			content: "test content\nmultiple lines\ndiff lengths",
			width:   80,
			height:  24,
			checkFunc: func(t *testing.T, original, result string) {
				assert.NotEqual(t, original, result, "should be padded/centered")

				for _, line := range strings.Split(original, "\n") {
					assert.Contains(t, result, line)
				}

				assert.Greater(t,
					strings.Count(result, "\n"),
					strings.Count(original, "\n"),
					"Result should have vertical padding")
			},
		},
		{
			name:    "single line content",
			content: "Just one line",
			width:   50,
			height:  10,
			checkFunc: func(t *testing.T, original, result string) {
				assert.Contains(t, result, original)
				assert.Greater(t, len(result), len(original), "Result should be padded")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := NewModel(&config.Config{}, false, nil)
			model.Width = tc.width
			model.Height = tc.height

			result := model.centerContent(tc.content)
			tc.checkFunc(t, tc.content, result)
		})
	}
}

func TestViewRendering(t *testing.T) {
	testCases := []struct {
		width  int
		height int
	}{
		{0, 0},
		{40, 10},
		{80, 24},
		{120, 40},
		{200, 100},
	}

	for _, tc := range testCases {
		t.Run("screen size "+string(rune('0'+tc.width))+"x"+string(rune('0'+tc.height)), func(t *testing.T) {
			model := NewModel(&config.Config{}, false, &api.Client{})
			model.Width = tc.width
			model.Height = tc.height

			view := model.View()

			assert.Contains(t, view, "____ ___  _______/ /_")
			assert.Contains(t, view, "Simple terminal weather")

			assert.NotEmpty(t, view)
		})
	}
}

func TestBuildContent(t *testing.T) {
	states := []SetupState{
		StateCity,
		StateCitySearch,
		StateCitySelect,
		StateUnits,
		StateView,
		StateTips,
		StateAuth,
		StateComplete,
	}

	seenContents := make(map[string]bool)

	for _, state := range states {
		model := NewModel(&config.Config{}, false, &api.Client{})
		model.State = state

		if state == StateCitySelect {
			model.CityOptions = []models.City{{Name: "London"}}
		}

		content := model.buildContent()

		_, exists := seenContents[content]
		assert.False(t, exists, "%v should be unique", state)

		seenContents[content] = true

		assert.NotEmpty(t, content)

		assert.Contains(t, content, "____/")
		assert.Contains(t, content, "Navigate")
	}
}
