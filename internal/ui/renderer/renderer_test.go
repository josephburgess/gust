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
)

func TestNewWeatherRenderer(t *testing.T) {
	renderer := NewWeatherRenderer("terminal", "")
	if renderer == nil {
		t.Fatal("NewWeatherRenderer() should return a non-nil renderer")
	}
}

func TestNewTerminalRenderer(t *testing.T) {
	renderer := NewTerminalRenderer("")
	if renderer == nil {
		t.Fatal("NewTerminalRenderer() should return a non-nil renderer")
	}
}

func TestRenderCurrentWeather(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	city := &models.City{
		Name: "Test City",
		Lat:  51.5074,
		Lon:  -0.1278,
	}

	weather := &models.OneCallResponse{
		Current: models.CurrentWeather{
			Dt:         time.Now().Unix(),
			Temp:       10,
			FeelsLike:  8,
			Humidity:   65,
			UVI:        2.5,
			WindSpeed:  5.1,
			WindDeg:    180,
			Visibility: 10000,
			Sunrise:    time.Now().Add(-6 * time.Hour).Unix(),
			Sunset:     time.Now().Add(6 * time.Hour).Unix(),
			Weather: []models.WeatherCondition{
				{
					ID:          800,
					Main:        "Clear",
					Description: "clear sky",
					Icon:        "01d",
				},
			},
		},
	}

	cfg := &config.Config{
		Units:    "metric",
		ShowTips: true,
	}

	renderer := NewTerminalRenderer("metric")
	renderer.RenderCurrentWeather(city, weather, cfg)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expectedPhrases := []string{
		"WEATHER FOR TEST CITY",
		"clear sky",
		"10.0°C",
		"Humidity: 65%",
		"UV Index: 2.5",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Expected output to contain '%s', but it didn't", phrase)
		}
	}
}

func TestRenderAlerts(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	city := &models.City{Name: "Alert City"}
	weather := &models.OneCallResponse{
		Alerts: []models.Alert{
			{
				SenderName:  "Weather Service",
				Event:       "Severe Thunderstorm Warning",
				Start:       time.Now().Unix(),
				End:         time.Now().Add(3 * time.Hour).Unix(),
				Description: "A severe thunderstorm is expected in the area.",
			},
		},
	}

	cfg := &config.Config{
		Units:    "metric",
		ShowTips: true,
	}

	renderer := NewTerminalRenderer("")
	renderer.RenderAlerts(city, weather, cfg)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expectedPhrases := []string{
		"WEATHER ALERTS FOR ALERT CITY",
		"Severe Thunderstorm Warning",
		"Issued by: Weather Service",
		"A severe thunderstorm is expected in the area.",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Expected output to contain '%s', but it didn't", phrase)
		}
	}
}

func TestRenderAlertsNoAlerts(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	city := &models.City{Name: "Calm City"}
	weather := &models.OneCallResponse{
		Alerts: []models.Alert{},
	}

	cfg := &config.Config{
		Units:    "metric",
		ShowTips: true,
	}

	renderer := NewTerminalRenderer("")
	renderer.RenderAlerts(city, weather, cfg)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "No weather alerts for this area") {
		t.Error("Expected output to indicate no alerts")
	}
}

func TestBaseRendererHelpers(t *testing.T) {
	testCases := []struct {
		units                 string
		expectedTempUnit      string
		expectedWindSpeedUnit string
	}{
		{"metric", "°C", "km/h"},
		{"imperial", "°F", "mph"},
		{"standard", "K", "km/h"},
	}

	for _, tc := range testCases {
		t.Run("Units: "+tc.units, func(t *testing.T) {
			renderer := BaseRenderer{Units: tc.units}

			tempUnit := renderer.GetTemperatureUnit()
			if tempUnit != tc.expectedTempUnit {
				t.Errorf("GetTemperatureUnit() = %v, want %v", tempUnit, tc.expectedTempUnit)
			}

			windUnit := renderer.GetWindSpeedUnit()
			if windUnit != tc.expectedWindSpeedUnit {
				t.Errorf("GetWindSpeedUnit() = %v, want %v", windUnit, tc.expectedWindSpeedUnit)
			}

			// Test wind speed conversion
			windSpeed := 10.0
			convertedSpeed := renderer.FormatWindSpeed(windSpeed)

			if tc.units == "imperial" {
				if convertedSpeed != windSpeed {
					t.Errorf("FormatWindSpeed() = %v, want %v", convertedSpeed, windSpeed)
				}
			} else {
				if convertedSpeed != windSpeed*3.6 {
					t.Errorf("FormatWindSpeed() = %v, want %v", convertedSpeed, windSpeed*3.6)
				}
			}
		})
	}
}

func TestFormatDateTime(t *testing.T) {
	timestamp := int64(1609459200) // 2021-01-01 00:00:00 UTC
	format := "2006-01-02 15:04:05"

	result := FormatDateTime(timestamp, format)
	expected := "2021-01-01 00:00:00"

	if result != expected {
		t.Errorf("FormatDateTime() = %v, want %v", result, expected)
	}
}

func TestFormatHeader(t *testing.T) {
	title := "TEST HEADER"
	header := styles.FormatHeader(title)

	if !strings.Contains(header, title) {
		t.Errorf("Header doesn't contain title: %s", header)
	}

	if !strings.Contains(header, styles.Divider(len(title)*2)) {
		t.Errorf("Header doesn't contain divider: %s", header)
	}
}
