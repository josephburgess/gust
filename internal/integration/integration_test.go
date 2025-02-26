package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josephburgess/gust/internal/api"
)

// test the full flow of getting coordinates and weather
func TestSimpleIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/geo/1.0/direct":
			if _, err := w.Write([]byte(`[{"name":"London","lat":51.5074,"lon":-0.1278}]`)); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		case "/data/2.5/weather":
			if _, err := w.Write([]byte(`{
				"main": {"temp": 283.15},
				"weather": [{"description": "cloudy"}]
			}`)); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	originalBaseURL := api.GetBaseURL()
	api.SetBaseURL(server.URL + "/")
	defer api.SetBaseURL(originalBaseURL)

	city, err := api.GetCoordinates("London", "test-api-key")
	if err != nil {
		t.Fatalf("GetCoordinates failed: %v", err)
	}

	weather, err := api.GetWeather(city.Lat, city.Lon, "test-api-key")
	if err != nil {
		t.Fatalf("GetWeather failed: %v", err)
	}

	if city.Name != "London" {
		t.Errorf("Expected city name 'London', got '%s'", city.Name)
	}

	if weather.Temp != 283.15 {
		t.Errorf("Expected temperature 283.15, got %f", weather.Temp)
	}

	if weather.Description != "cloudy" {
		t.Errorf("Expected description 'cloudy', got '%s'", weather.Description)
	}
}
