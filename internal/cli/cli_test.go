package cli

import (
	"testing"

	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockWeatherRenderer struct {
	mock.Mock
}

func (m *MockWeatherRenderer) RenderCurrentWeather(city *models.City, weather *models.OneCallResponse) {
	m.Called(city, weather)
}

func (m *MockWeatherRenderer) RenderDailyForecast(city *models.City, weather *models.OneCallResponse) {
	m.Called(city, weather)
}

func (m *MockWeatherRenderer) RenderHourlyForecast(city *models.City, weather *models.OneCallResponse) {
	m.Called(city, weather)
}

func (m *MockWeatherRenderer) RenderAlerts(city *models.City, weather *models.OneCallResponse) {
	m.Called(city, weather)
}

func (m *MockWeatherRenderer) RenderFullWeather(city *models.City, weather *models.OneCallResponse) {
	m.Called(city, weather)
}

func (m *MockWeatherRenderer) RenderCompactWeather(city *models.City, weather *models.OneCallResponse) {
	m.Called(city, weather)
}

func TestRenderWeatherView_CLIFlags(t *testing.T) {
	mockCity := &models.City{Name: "TestCity"}
	mockWeather := &models.OneCallResponse{}
	mockConfig := &config.Config{DefaultView: "compact"}

	testCases := []struct {
		name           string
		cli            *CLI
		expectedMethod string
	}{
		{
			name:           "Alerts flag",
			cli:            &CLI{Alerts: true},
			expectedMethod: "RenderAlerts",
		},
		{
			name:           "Hourly flag",
			cli:            &CLI{Hourly: true},
			expectedMethod: "RenderHourlyForecast",
		},
		{
			name:           "Daily flag",
			cli:            &CLI{Daily: true},
			expectedMethod: "RenderDailyForecast",
		},
		{
			name:           "Full flag",
			cli:            &CLI{Full: true},
			expectedMethod: "RenderFullWeather",
		},
		{
			name:           "Compact flag",
			cli:            &CLI{Compact: true},
			expectedMethod: "RenderCompactWeather",
		},
		{
			name:           "Detailed flag",
			cli:            &CLI{Detailed: true},
			expectedMethod: "RenderCurrentWeather",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRenderer := new(MockWeatherRenderer)

			switch tc.expectedMethod {
			case "RenderAlerts":
				mockRenderer.On("RenderAlerts", mockCity, mockWeather).Return()
			case "RenderHourlyForecast":
				mockRenderer.On("RenderHourlyForecast", mockCity, mockWeather).Return()
			case "RenderDailyForecast":
				mockRenderer.On("RenderDailyForecast", mockCity, mockWeather).Return()
			case "RenderFullWeather":
				mockRenderer.On("RenderFullWeather", mockCity, mockWeather).Return()
			case "RenderCompactWeather":
				mockRenderer.On("RenderCompactWeather", mockCity, mockWeather).Return()
			case "RenderCurrentWeather":
				mockRenderer.On("RenderCurrentWeather", mockCity, mockWeather).Return()
			}

			renderWeatherView(tc.cli, mockRenderer, mockCity, mockWeather, mockConfig)

			mockRenderer.AssertExpectations(t)
		})
	}
}

func TestRenderWeatherView_DefaultView(t *testing.T) {
	mockCity := &models.City{Name: "TestCity"}
	mockWeather := &models.OneCallResponse{}
	cli := &CLI{}

	testCases := []struct {
		name        string
		defaultView string
		expectedFn  string
	}{
		{
			name:        "Default view is compact",
			defaultView: "compact",
			expectedFn:  "RenderCompactWeather",
		},
		{
			name:        "Default view is daily",
			defaultView: "daily",
			expectedFn:  "RenderDailyForecast",
		},
		{
			name:        "Default view is hourly",
			defaultView: "hourly",
			expectedFn:  "RenderHourlyForecast",
		},
		{
			name:        "Default view is full",
			defaultView: "full",
			expectedFn:  "RenderFullWeather",
		},
		{
			name:        "Default view is unknown",
			defaultView: "unknown",
			expectedFn:  "RenderCurrentWeather",
		},
		{
			name:        "Default view is empty",
			defaultView: "",
			expectedFn:  "RenderCurrentWeather",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRenderer := new(MockWeatherRenderer)
			config := &config.Config{DefaultView: tc.defaultView}

			switch tc.expectedFn {
			case "RenderCompactWeather":
				mockRenderer.On("RenderCompactWeather", mockCity, mockWeather).Return()
			case "RenderDailyForecast":
				mockRenderer.On("RenderDailyForecast", mockCity, mockWeather).Return()
			case "RenderHourlyForecast":
				mockRenderer.On("RenderHourlyForecast", mockCity, mockWeather).Return()
			case "RenderFullWeather":
				mockRenderer.On("RenderFullWeather", mockCity, mockWeather).Return()
			case "RenderCurrentWeather":
				mockRenderer.On("RenderCurrentWeather", mockCity, mockWeather).Return()
			}

			renderWeatherView(cli, mockRenderer, mockCity, mockWeather, config)

			mockRenderer.AssertExpectations(t)
		})
	}
}

func TestRenderWeatherView_FlagPrecedence(t *testing.T) {
	mockRenderer := new(MockWeatherRenderer)
	mockCity := &models.City{Name: "TestCity"}
	mockWeather := &models.OneCallResponse{}

	config := &config.Config{DefaultView: "compact"}

	cli := &CLI{Daily: true}

	mockRenderer.On("RenderDailyForecast", mockCity, mockWeather).Return()

	renderWeatherView(cli, mockRenderer, mockCity, mockWeather, config)

	mockRenderer.AssertExpectations(t)
}

func TestRenderWeatherView_MultipleFlagsPriority(t *testing.T) {
	mockRenderer := new(MockWeatherRenderer)
	mockCity := &models.City{Name: "TestCity"}
	mockWeather := &models.OneCallResponse{}
	config := &config.Config{DefaultView: "compact"}

	cli := &CLI{
		Alerts:   true,
		Hourly:   true,
		Daily:    true,
		Full:     true,
		Compact:  true,
		Detailed: true,
	}

	mockRenderer.On("RenderAlerts", mockCity, mockWeather).Return()

	renderWeatherView(cli, mockRenderer, mockCity, mockWeather, config)

	mockRenderer.AssertExpectations(t)
}
