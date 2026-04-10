package cli

import (
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/stretchr/testify/mock"
)

// for renderWeatjerView
type MockWeatherRenderer struct {
	mock.Mock
}

func (m *MockWeatherRenderer) RenderCurrentWeather(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	m.Called(city, weather, cfg)
}

func (m *MockWeatherRenderer) RenderDailyForecast(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	m.Called(city, weather, cfg)
}

func (m *MockWeatherRenderer) RenderHourlyForecast(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	m.Called(city, weather, cfg)
}

func (m *MockWeatherRenderer) RenderAlerts(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	m.Called(city, weather, cfg)
}

func (m *MockWeatherRenderer) RenderFullWeather(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	m.Called(city, weather, cfg)
}

func (m *MockWeatherRenderer) RenderCompactWeather(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	m.Called(city, weather, cfg)
}


type MockWeatherClient struct {
	mock.Mock
}

func (m *MockWeatherClient) GetWeather(cityName string) (*api.WeatherResponse, error) {
	args := m.Called(cityName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.WeatherResponse), args.Error(1)
}

func (m *MockWeatherClient) SearchCities(query string) ([]models.City, error) {
	args := m.Called(query)
	cities, _ := args.Get(0).([]models.City)
	return cities, args.Error(1)
}

// create a test models.City
func createTestCity() *models.City {
	return &models.City{
		Name:    "TestCity",
		Lat:     51,
		Lon:     0,
		Country: "GB",
	}
}

// create a test models.OneCallResponse
func createTestWeather() *models.OneCallResponse {
	return &models.OneCallResponse{
		Current: models.CurrentWeather{
			Temp:      20.5,
			FeelsLike: 21.0,
			Humidity:  65,
		},
	}
}
