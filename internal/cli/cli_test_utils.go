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

// auth config handling
type MockAuthConfig struct {
	mock.Mock
	APIKey     string
	GithubUser string
}

func (m *MockAuthConfig) AuthFunc(apiURL string) (*config.AuthConfig, error) {
	args := m.Called(apiURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*config.AuthConfig), args.Error(1)
}

func (m *MockAuthConfig) SaveAuthFunc(config *config.AuthConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockAuthConfig) LoadAuthFunc() (*config.AuthConfig, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*config.AuthConfig), args.Error(1)
}

// mocks setup wizard
type MockSetup struct {
	mock.Mock
}

func (m *MockSetup) RunSetupFunc(cfg *config.Config, needsAuth bool) error {
	args := m.Called(cfg, needsAuth)
	return args.Error(0)
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
