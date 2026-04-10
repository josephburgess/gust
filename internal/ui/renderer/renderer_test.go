package renderer

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
	"github.com/stretchr/testify/assert"
)

// captureOutput redirects stdout to a buffer for the duration of fn.
func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func testCity(name string) *models.City { return &models.City{Name: name, Lat: 51.5, Lon: -0.1} }
func noCfg() *config.Config             { return &config.Config{Units: "metric", ShowTips: false} }

func TestNewWeatherRenderer(t *testing.T) {
	assert.NotNil(t, NewWeatherRenderer("metric"))
}

func TestNewTerminalRenderer(t *testing.T) {
	assert.NotNil(t, NewTerminalRenderer("metric"))
}

func TestRenderCurrentWeather(t *testing.T) {
	weather := &models.OneCallResponse{
		Current: models.CurrentWeather{
			Dt:         time.Now().Unix(),
			Temp:       10, FeelsLike: 8, Humidity: 65, UVI: 2.5,
			WindSpeed: 5.1, WindDeg: 180, Visibility: 10000,
			Sunrise: time.Now().Add(-6 * time.Hour).Unix(),
			Sunset:  time.Now().Add(6 * time.Hour).Unix(),
			Weather: []models.WeatherCondition{{ID: 800, Main: "Clear", Description: "clear sky"}},
		},
	}

	out := captureOutput(t, func() {
		NewTerminalRenderer("metric").RenderCurrentWeather(testCity("Test City"), weather, noCfg())
	})

	assert.Contains(t, out, "WEATHER FOR TEST CITY")
	assert.Contains(t, out, "clear sky")
	assert.Contains(t, out, "10.0°C")
	assert.Contains(t, out, "Humidity: 65%")
	assert.Contains(t, out, "UV Index: 2.5")
}

func TestRenderCompactWeather(t *testing.T) {
	weather := &models.OneCallResponse{
		Current: models.CurrentWeather{
			Temp: 15, Humidity: 70, WindSpeed: 3, WindDeg: 90,
			Sunrise: time.Now().Add(-5 * time.Hour).Unix(),
			Sunset:  time.Now().Add(5 * time.Hour).Unix(),
			Weather: []models.WeatherCondition{{ID: 800, Description: "sunny"}},
		},
	}

	out := captureOutput(t, func() {
		NewTerminalRenderer("metric").RenderCompactWeather(testCity("London"), weather, noCfg())
	})

	assert.Contains(t, out, "WEATHER FOR LONDON")
	assert.Contains(t, out, "15.0°C")
	assert.Contains(t, out, "70")  // humidity
}

func TestRenderDailyForecast(t *testing.T) {
	weather := &models.OneCallResponse{
		Daily: []models.DayData{
			{
				Dt:      time.Now().Unix(),
				Summary: "Cloudy with a chance of rain",
				Temp:    models.TempData{Max: 18, Min: 10, Morn: 11, Day: 17, Eve: 15, Night: 10},
				Weather: []models.WeatherCondition{{ID: 500, Description: "light rain"}},
				Pop:     0.6, Rain: 2.5, WindSpeed: 4, WindDeg: 270, UVI: 1.5,
			},
		},
	}

	out := captureOutput(t, func() {
		NewTerminalRenderer("metric").RenderDailyForecast(testCity("Paris"), weather, noCfg())
	})

	assert.Contains(t, out, "5-DAY FORECAST FOR PARIS")
	assert.Contains(t, out, "18.0°C")
	assert.Contains(t, out, "light rain")
	assert.Contains(t, out, "60% chance")
	assert.Contains(t, out, "Rain: 2.5 mm")
}

func TestRenderHourlyForecast(t *testing.T) {
	weather := &models.OneCallResponse{
		Hourly: []models.HourData{
			{
				Dt:      time.Now().Unix(),
				Temp:    12.5,
				Pop:     0.3,
				Weather: []models.WeatherCondition{{ID: 803, Description: "broken clouds"}},
			},
		},
	}

	out := captureOutput(t, func() {
		NewTerminalRenderer("metric").RenderHourlyForecast(testCity("Berlin"), weather, noCfg())
	})

	assert.Contains(t, out, "24H FORECAST FOR BERLIN")
	assert.Contains(t, out, "12.5°C")
	assert.Contains(t, out, "broken clouds")
}

func TestRenderAlerts(t *testing.T) {
	weather := &models.OneCallResponse{
		Alerts: []models.Alert{
			{
				SenderName:  "Weather Service",
				Event:       "Severe Thunderstorm Warning",
				Start:       time.Now().Unix(),
				End:         time.Now().Add(3 * time.Hour).Unix(),
				Description: "A severe thunderstorm is expected.",
			},
		},
	}

	out := captureOutput(t, func() {
		NewTerminalRenderer("").RenderAlerts(testCity("Alert City"), weather, noCfg())
	})

	assert.Contains(t, out, "WEATHER ALERTS FOR ALERT CITY")
	assert.Contains(t, out, "Severe Thunderstorm Warning")
	assert.Contains(t, out, "Weather Service")
	assert.Contains(t, out, "A severe thunderstorm is expected.")
}

func TestRenderAlertsNone(t *testing.T) {
	out := captureOutput(t, func() {
		NewTerminalRenderer("").RenderAlerts(testCity("Calm City"), &models.OneCallResponse{}, noCfg())
	})
	assert.Contains(t, out, "No weather alerts")
}

func TestBaseRendererHelpers(t *testing.T) {
	tests := []struct {
		units    string
		tempUnit string
		windUnit string
	}{
		{"metric", "°C", "km/h"},
		{"imperial", "°F", "mph"},
		{"standard", "K", "km/h"},
	}
	for _, tt := range tests {
		t.Run(tt.units, func(t *testing.T) {
			r := BaseRenderer{Units: tt.units}
			assert.Equal(t, tt.tempUnit, r.GetTemperatureUnit())
			assert.Equal(t, tt.windUnit, r.GetWindSpeedUnit())
			if tt.units == "imperial" {
				assert.Equal(t, 10.0, r.FormatWindSpeed(10.0))
			} else {
				assert.Equal(t, 36.0, r.FormatWindSpeed(10.0))
			}
		})
	}
}

func TestFormatDateTime(t *testing.T) {
	result := FormatDateTime(1609459200, "2006-01-02")
	assert.True(t, strings.HasPrefix(result, "2021-01-0")) // local-time-safe
}

func TestFormatHeader(t *testing.T) {
	header := styles.FormatHeader("TEST HEADER")
	assert.Contains(t, header, "TEST HEADER")
}
