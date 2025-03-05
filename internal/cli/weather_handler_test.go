package cli

import (
	"testing"

	"github.com/josephburgess/gust/internal/config"
)

func TestRenderWeatherView_CLIFlags(t *testing.T) {
	mockCity := createTestCity()
	mockWeather := createTestWeather()
	mockConfig := &config.Config{DefaultView: "compact"}

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
			mockRenderer.On(tc.expectedMethod, mockCity, mockWeather).Return()

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
			mockRenderer.On(tc.expectedFn, mockCity, mockWeather).Return()

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
			mockRenderer.On(tc.expectedMethod, mockCity, mockWeather).Return()

			renderWeatherView(tc.cli, mockRenderer, mockCity, mockWeather, mockConfig)

			mockRenderer.AssertExpectations(t)
		})
	}
}
