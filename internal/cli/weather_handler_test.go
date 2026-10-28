package cli

import (
	"testing"

	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/ui/tui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRenderWeatherView_CLIFlags(t *testing.T) {
	mockCity := createTestCity()
	mockWeather := createTestWeather()
	mockConfig := &config.Config{DefaultView: "compact", ShowTips: false}

	testCases := []struct {
		name           string
		cli            *CLI
		expectedMethod string
	}{
		{
			name:           "alerts flag",
			cli:            &CLI{Alerts: true},
			expectedMethod: "RenderAlerts",
		},
		{
			name:           "hourly flag",
			cli:            &CLI{Hourly: true},
			expectedMethod: "RenderHourlyForecast",
		},
		{
			name:           "daily flag",
			cli:            &CLI{Daily: true},
			expectedMethod: "RenderDailyForecast",
		},
		{
			name:           "full flag",
			cli:            &CLI{Full: true},
			expectedMethod: "RenderFullWeather",
		},
		{
			name:           "compact flag",
			cli:            &CLI{Compact: true},
			expectedMethod: "RenderCompactWeather",
		},
		{
			name:           "detailed flag",
			cli:            &CLI{Detailed: true},
			expectedMethod: "RenderCurrentWeather",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRenderer := new(MockWeatherRenderer)
			mockRenderer.On(tc.expectedMethod, mockCity, mockWeather, mockConfig).Return()

			renderWeatherView(tc.cli, mockRenderer, mockCity, mockWeather, mockConfig)

			mockRenderer.AssertExpectations(t)
		})
	}
}

func TestRenderWeatherView_DefaultView(t *testing.T) {
	mockCity := createTestCity()
	mockWeather := createTestWeather()
	cli := &CLI{}

	testCases := []struct {
		name        string
		defaultView string
		expectedFn  string
	}{
		{
			name:        "default compact",
			defaultView: "compact",
			expectedFn:  "RenderCompactWeather",
		},
		{
			name:        "default view daily",
			defaultView: "daily",
			expectedFn:  "RenderDailyForecast",
		},
		{
			name:        "default view hourly",
			defaultView: "hourly",
			expectedFn:  "RenderHourlyForecast",
		},
		{
			name:        "default view full",
			defaultView: "full",
			expectedFn:  "RenderFullWeather",
		},
		{
			name:        "default view unknown",
			defaultView: "unknown",
			expectedFn:  "RenderCurrentWeather",
		},
		{
			name:        "default view empty",
			defaultView: "",
			expectedFn:  "RenderCurrentWeather",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRenderer := new(MockWeatherRenderer)
			mockRenderer.On(tc.expectedFn, mockCity, mockWeather, mock.Anything).Return()

			config := &config.Config{DefaultView: tc.defaultView}
			renderWeatherView(cli, mockRenderer, mockCity, mockWeather, config)

			mockRenderer.AssertExpectations(t)
		})
	}
}

func TestRenderWeatherView_FlagsAndPriority(t *testing.T) {
	mockCity := createTestCity()
	mockWeather := createTestWeather()
	mockConfig := &config.Config{DefaultView: "compact"}

	testCases := []struct {
		name           string
		cli            *CLI
		expectedMethod string
	}{
		{
			name:           "default falls back",
			cli:            &CLI{},
			expectedMethod: "RenderCompactWeather",
		},
		{
			name:           "flag overrides default",
			cli:            &CLI{Daily: true},
			expectedMethod: "RenderDailyForecast",
		},
		{
			name: "multiple flags respect prio",
			cli: &CLI{
				Alerts:   true,
				Hourly:   true,
				Daily:    true,
				Full:     true,
				Compact:  true,
				Detailed: true,
			},
			expectedMethod: "RenderAlerts",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRenderer := new(MockWeatherRenderer)
			mockRenderer.On(tc.expectedMethod, mockCity, mockWeather, mock.Anything).Return()

			renderWeatherView(tc.cli, mockRenderer, mockCity, mockWeather, mockConfig)

			mockRenderer.AssertExpectations(t)
		})
	}
}

func TestDashboardTab(t *testing.T) {
	tests := []struct {
		name string
		cli  *CLI
		view string
		want tui.Tab
	}{
		{"alerts flag", &CLI{Alerts: true}, "hourly", tui.TabAlerts},
		{"hourly flag", &CLI{Hourly: true}, "daily", tui.TabHourly},
		{"daily flag", &CLI{Daily: true}, "default", tui.TabDaily},
		{"no flag uses the configured view", &CLI{}, "hourly", tui.TabHourly},
		{"old full view opens now", &CLI{}, "full", tui.TabNow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, dashboardTab(tt.cli, &config.Config{DefaultView: tt.view}))
		})
	}
}
