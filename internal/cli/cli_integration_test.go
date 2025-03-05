package cli

import (
	"testing"

	"github.com/josephburgess/gust/internal/api"
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

			mockClient.On("GetWeather", tc.expectedCity).Return(weatherResponse, nil)
			mockRenderer.On("RenderCompactWeather", cityData, weatherData).Return()

			cli := &CLI{
				City: tc.cityFlag,
				Args: tc.args,
			}

			city := determineCityName(cli.City, cli.Args, tc.defaultCity)

			assert.Equal(t, tc.expectedCity, city)

			weather, err := mockClient.GetWeather(city)
			assert.NoError(t, err)
			assert.Equal(t, weatherResponse, weather)

			testRenderWeatherView := func(cli *CLI, renderer renderer.WeatherRenderer, city *models.City, weather *models.OneCallResponse, defaultView string) {
				switch {
				case cli.Alerts:
					renderer.RenderAlerts(city, weather)
				case cli.Hourly:
					renderer.RenderHourlyForecast(city, weather)
				case cli.Daily:
					renderer.RenderDailyForecast(city, weather)
				case cli.Full:
					renderer.RenderFullWeather(city, weather)
				case cli.Compact:
					renderer.RenderCompactWeather(city, weather)
				case cli.Detailed:
					renderer.RenderCurrentWeather(city, weather)
				default:
					switch defaultView {
					case "compact":
						renderer.RenderCompactWeather(city, weather)
					case "daily":
						renderer.RenderDailyForecast(city, weather)
					case "hourly":
						renderer.RenderHourlyForecast(city, weather)
					case "full":
						renderer.RenderFullWeather(city, weather)
					default:
						renderer.RenderCurrentWeather(city, weather)
					}
				}
			}

			testRenderWeatherView(cli, mockRenderer, weather.City, weather.Weather, "compact")

			mockClient.AssertExpectations(t)
			mockRenderer.AssertExpectations(t)
		})
	}
}

func TestViewSelectionIntegration(t *testing.T) {
	mockCity := &models.City{Name: "TestCity"}
	mockWeather := &models.OneCallResponse{}

	type fakeConfig struct {
		DefaultView string
	}

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
			config := &fakeConfig{DefaultView: tc.defaultView}

			switch tc.expectedMethod {
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
			case "RenderAlerts":
				mockRenderer.On("RenderAlerts", mockCity, mockWeather).Return()
			}

			testRenderWeatherView := func(cli *CLI, renderer renderer.WeatherRenderer, city *models.City, weather *models.OneCallResponse, cfg any) {
				configObj, ok := cfg.(*fakeConfig)
				var defaultView string
				if ok {
					defaultView = configObj.DefaultView
				}

				switch {
				case cli.Alerts:
					renderer.RenderAlerts(city, weather)
				case cli.Hourly:
					renderer.RenderHourlyForecast(city, weather)
				case cli.Daily:
					renderer.RenderDailyForecast(city, weather)
				case cli.Full:
					renderer.RenderFullWeather(city, weather)
				case cli.Compact:
					renderer.RenderCompactWeather(city, weather)
				case cli.Detailed:
					renderer.RenderCurrentWeather(city, weather)
				default:
					switch defaultView {
					case "compact":
						renderer.RenderCompactWeather(city, weather)
					case "daily":
						renderer.RenderDailyForecast(city, weather)
					case "hourly":
						renderer.RenderHourlyForecast(city, weather)
					case "full":
						renderer.RenderFullWeather(city, weather)
					default:
						renderer.RenderCurrentWeather(city, weather)
					}
				}
			}

			testRenderWeatherView(tc.cli, mockRenderer, mockCity, mockWeather, config)

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
