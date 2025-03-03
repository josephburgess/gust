package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/josephburgess/gust/internal/models"
)

func TestNewRenderer(t *testing.T) {
	renderer := NewRenderer("")
	if renderer == nil {
		t.Fatal("NewRenderer() should return a non-nil renderer")
	}
}

func TestDisplayCurrentWeather(t *testing.T) {
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

	renderer := NewRenderer("metric")
	renderer.DisplayCurrentWeather(city, weather)

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

func TestDisplayAlerts(t *testing.T) {
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

	renderer := NewRenderer("")
	renderer.DisplayAlerts(city, weather)

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

func TestDisplayAlertsNoAlerts(t *testing.T) {
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

	renderer := NewRenderer("")
	renderer.DisplayAlerts(city, weather)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "No weather alerts for this area") {
		t.Error("Expected output to indicate no alerts")
	}
}

func TestFormatHeader(t *testing.T) {
	title := "TEST HEADER"
	header := FormatHeader(title)

	if !strings.Contains(header, title) {
		t.Errorf("Header doesn't contain title: %s", header)
	}

	if !strings.Contains(header, Divider()) {
		t.Errorf("Header doesn't contain divider: %s", header)
	}
}
