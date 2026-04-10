package cli

import (
	"testing"

	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeatherFlowIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode - skipping int tests")
	}

	cityData := &models.City{Name: "London", Country: "GB", Lat: 51, Lon: 0}
	weatherData := &models.OneCallResponse{Current: models.CurrentWeather{Temp: 20.5, FeelsLike: 21.0, Humidity: 65}}
	weatherResponse := &api.WeatherResponse{City: cityData, Weather: weatherData}
	cfg := &config.Config{ShowTips: false, DefaultView: "compact"}

	tests := []struct {
		name         string
		cityFlag     string
		args         []string
		defaultCity  string
		expectedCity string
	}{
		{"city flag", "London", []string{}, "Berlin", "London"},
		{"positional args", "", []string{"London"}, "Berlin", "London"},
		{"default city", "", []string{}, "Berlin", "Berlin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockWeatherClient)
			mockRenderer := new(MockWeatherRenderer)

			mockClient.On("GetWeather", tt.expectedCity).Return(weatherResponse, nil)
			mockRenderer.On("RenderCompactWeather", cityData, weatherData, cfg).Return()

			cli := &CLI{City: tt.cityFlag, Args: tt.args}
			city := determineCityName(cli.City, cli.Args, tt.defaultCity)
			assert.Equal(t, tt.expectedCity, city)

			weather, err := mockClient.GetWeather(city)
			require.NoError(t, err)

			renderWeatherView(cli, mockRenderer, weather.City, weather.Weather, cfg)

			mockClient.AssertExpectations(t)
			mockRenderer.AssertExpectations(t)
		})
	}
}

func TestViewSelectionIntegration(t *testing.T) {
	mockCity := &models.City{Name: "TestCity"}
	mockWeather := &models.OneCallResponse{}

	tests := []struct {
		name           string
		cli            *CLI
		defaultView    string
		expectedMethod string
	}{
		{"cli flag overrides default", &CLI{Hourly: true}, "compact", "RenderHourlyForecast"},
		{"multiple flags respect priority", &CLI{Hourly: true, Daily: true, Compact: true}, "full", "RenderHourlyForecast"},
		{"default view used when no flags", &CLI{}, "daily", "RenderDailyForecast"},
		{"fallback to current for unknown default", &CLI{}, "invalid", "RenderCurrentWeather"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{DefaultView: tt.defaultView, ShowTips: false}
			mockRenderer := new(MockWeatherRenderer)
			mockRenderer.On(tt.expectedMethod, mockCity, mockWeather, cfg).Return()

			renderWeatherView(tt.cli, mockRenderer, mockCity, mockWeather, cfg)

			mockRenderer.AssertExpectations(t)
		})
	}
}

func TestCityDeterminationIntegration(t *testing.T) {
	tests := []struct {
		name         string
		cityFlag     string
		args         []string
		defaultCity  string
		expectedCity string
	}{
		{"flag takes priority", "London", []string{"Paris"}, "Berlin", "London"},
		{"args used when no flag", "", []string{"Paris"}, "Berlin", "Paris"},
		{"multi-word city from args", "", []string{"New", "York"}, "Berlin", "New York"},
		{"default when no flag or args", "", []string{}, "Berlin", "Berlin"},
		{"empty when no sources", "", []string{}, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &CLI{City: tt.cityFlag, Args: tt.args}
			assert.Equal(t, tt.expectedCity, determineCityName(cli.City, cli.Args, tt.defaultCity))
		})
	}
}
