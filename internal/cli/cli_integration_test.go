package cli

import (
	"testing"

	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/renderer"
	"github.com/stretchr/testify/assert"
)

func TestWeatherFlowIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode - skipping int tests")
	}

	testCity := "London"

	cityData := &models.City{
		Name:    testCity,
		Country: "GB",
		Lat:     51,
		Lon:     0,
	}

	weatherData := &models.OneCallResponse{
		Current: models.CurrentWeather{
			Temp:      20.5,
			FeelsLike: 21.0,
			Humidity:  65,
		},
	}

	weatherResponse := &api.WeatherResponse{
		City:    cityData,
		Weather: weatherData,
	}

	testCases := []struct {
		name         string
		cityFlag     string
		args         []string
		defaultCity  string
		expectedCity string
	}{
		{
			name:         "Using city flag",
			cityFlag:     testCity,
			args:         []string{},
			defaultCity:  "Berlin",
			expectedCity: testCity,
		},
		{
			name:         "Using positional args",
			cityFlag:     "",
			args:         []string{testCity},
			defaultCity:  "Berlin",
			expectedCity: testCity,
		},
		{
			name:         "Using default city",
			cityFlag:     "",
			args:         []string{},
			defaultCity:  "Berlin",
			expectedCity: "Berlin",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(MockWeatherClient)
			mockRenderer := new(MockWeatherRenderer)
			mockConfig := &config.Config{ShowTips: true}

			mockClient.On("GetWeather", tc.expectedCity).Return(weatherResponse, nil)
			mockRenderer.On("RenderCompactWeather", cityData, weatherData, mockConfig).Return()

			cli := &CLI{
				City: tc.cityFlag,
				Args: tc.args,
			}

			city := determineCityName(cli.City, cli.Args, tc.defaultCity)

			assert.Equal(t, tc.expectedCity, city)

			weather, err := mockClient.GetWeather(city)
			assert.NoError(t, err)
			assert.Equal(t, weatherResponse, weather)

			testRenderWeatherView := func(cli *CLI, renderer renderer.WeatherRenderer, city *models.City, weather *models.OneCallResponse, defaultView string, cfg *config.Config) {
				switch {
				case cli.Alerts:
					renderer.RenderAlerts(city, weather, cfg)
				case cli.Hourly:
					renderer.RenderHourlyForecast(city, weather, cfg)
				case cli.Daily:
					renderer.RenderDailyForecast(city, weather, cfg)
				case cli.Full:
					renderer.RenderFullWeather(city, weather, cfg)
				case cli.Compact:
					renderer.RenderCompactWeather(city, weather, cfg)
				case cli.Detailed:
					renderer.RenderCurrentWeather(city, weather, cfg)
				default:
					switch defaultView {
					case "compact":
						renderer.RenderCompactWeather(city, weather, cfg)
					case "daily":
						renderer.RenderDailyForecast(city, weather, cfg)
					case "hourly":
						renderer.RenderHourlyForecast(city, weather, cfg)
					case "full":
						renderer.RenderFullWeather(city, weather, cfg)
					default:
						renderer.RenderCurrentWeather(city, weather, cfg)
					}
				}
			}

			testRenderWeatherView(cli, mockRenderer, weather.City, weather.Weather, "compact", mockConfig)

			mockClient.AssertExpectations(t)
			mockRenderer.AssertExpectations(t)
		})
	}
}

func TestViewSelectionIntegration(t *testing.T) {
	mockCity := &models.City{Name: "TestCity"}
	mockWeather := &models.OneCallResponse{}
	mockConfig := &config.Config{ShowTips: true}

	testCases := []struct {
		name           string
		cli            *CLI
		defaultView    string
		expectedMethod string
	}{
		{
			name:           "cli flag overrides defaults",
			cli:            &CLI{Hourly: true},
			defaultView:    "compact",
			expectedMethod: "RenderHourlyForecast",
		},
		{
			name:           "multiple flags follow prio",
			cli:            &CLI{Hourly: true, Daily: true, Compact: true},
			defaultView:    "full",
			expectedMethod: "RenderHourlyForecast",
		},
		{
			name:           "default used when no flags",
			cli:            &CLI{},
			defaultView:    "daily",
			expectedMethod: "RenderDailyForecast",
		},
		{
			name:           "fallback to current when no flag or valid config",
			cli:            &CLI{},
			defaultView:    "invalid",
			expectedMethod: "RenderCurrentWeather",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRenderer := new(MockWeatherRenderer)
			fakeConfig := &struct {
				DefaultView string
				ShowTips    bool
			}{
				DefaultView: tc.defaultView,
				ShowTips:    true,
			}

			switch tc.expectedMethod {
			case "RenderHourlyForecast":
				mockRenderer.On("RenderHourlyForecast", mockCity, mockWeather, mockConfig).Return()
			case "RenderDailyForecast":
				mockRenderer.On("RenderDailyForecast", mockCity, mockWeather, mockConfig).Return()
			case "RenderFullWeather":
				mockRenderer.On("RenderFullWeather", mockCity, mockWeather, mockConfig).Return()
			case "RenderCompactWeather":
				mockRenderer.On("RenderCompactWeather", mockCity, mockWeather, mockConfig).Return()
			case "RenderCurrentWeather":
				mockRenderer.On("RenderCurrentWeather", mockCity, mockWeather, mockConfig).Return()
			case "RenderAlerts":
				mockRenderer.On("RenderAlerts", mockCity, mockWeather, mockConfig).Return()
			}

			testRenderWeatherView := func(cli *CLI, renderer renderer.WeatherRenderer, city *models.City, weather *models.OneCallResponse, cfg any) {
				configObj, _ := cfg.(*struct {
					DefaultView string
					ShowTips    bool
				})
				var defaultView string
				if configObj != nil {
					defaultView = configObj.DefaultView
				}

				realConfig := &config.Config{ShowTips: true}

				switch {
				case cli.Alerts:
					renderer.RenderAlerts(city, weather, realConfig)
				case cli.Hourly:
					renderer.RenderHourlyForecast(city, weather, realConfig)
				case cli.Daily:
					renderer.RenderDailyForecast(city, weather, realConfig)
				case cli.Full:
					renderer.RenderFullWeather(city, weather, realConfig)
				case cli.Compact:
					renderer.RenderCompactWeather(city, weather, realConfig)
				case cli.Detailed:
					renderer.RenderCurrentWeather(city, weather, realConfig)
				default:
					switch defaultView {
					case "compact":
						renderer.RenderCompactWeather(city, weather, realConfig)
					case "daily":
						renderer.RenderDailyForecast(city, weather, realConfig)
					case "hourly":
						renderer.RenderHourlyForecast(city, weather, realConfig)
					case "full":
						renderer.RenderFullWeather(city, weather, realConfig)
					default:
						renderer.RenderCurrentWeather(city, weather, realConfig)
					}
				}
			}

			testRenderWeatherView(tc.cli, mockRenderer, mockCity, mockWeather, fakeConfig)

			mockRenderer.AssertExpectations(t)
		})
	}
}

func TestCityDeterminationIntegration(t *testing.T) {
	testCases := []struct {
		name         string
		cityFlag     string
		args         []string
		defaultCity  string
		expectedCity string
	}{
		{
			name:         "city flag takes prio",
			cityFlag:     "London",
			args:         []string{"Paris"},
			defaultCity:  "Berlin",
			expectedCity: "London",
		},
		{
			name:         "args used when no flag",
			cityFlag:     "",
			args:         []string{"Paris"},
			defaultCity:  "Berlin",
			expectedCity: "Paris",
		},
		{
			name:         "multi-word city from args",
			cityFlag:     "",
			args:         []string{"New", "York"},
			defaultCity:  "Berlin",
			expectedCity: "New York",
		},
		{
			name:         "default city when no flag or args",
			cityFlag:     "",
			args:         []string{},
			defaultCity:  "Berlin",
			expectedCity: "Berlin",
		},
		{
			name:         "empty when no sources available",
			cityFlag:     "",
			args:         []string{},
			defaultCity:  "",
			expectedCity: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cli := &CLI{
				City: tc.cityFlag,
				Args: tc.args,
			}

			city := determineCityName(cli.City, cli.Args, tc.defaultCity)

			assert.Equal(t, tc.expectedCity, city)
		})
	}
}
