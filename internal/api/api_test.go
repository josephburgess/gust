package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCoordinates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/geo/1.0/direct" {
			t.Errorf("Expected request to '/geo/1.0/direct', got: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`[{"name":"London","lat":51.5074,"lon":-0.1278}]`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	originalBaseURL := baseURL
	originalGeoCodeURL := geoCodeURL
	originalWeatherURL := weatherURL

	baseURL = server.URL + "/"
	geoCodeURL = baseURL + "geo/1.0/direct?q=%s&limit=1&appid=%s"
	weatherURL = baseURL + "data/2.5/weather?lat=%f&lon=%f&appid=%s"

	defer func() {
		baseURL = originalBaseURL
		geoCodeURL = originalGeoCodeURL
		weatherURL = originalWeatherURL
	}()

	city, err := GetCoordinates("London", "test-api-key")
	if err != nil {
		t.Fatalf("GetCoordinates returned an error: %v", err)
	}

	if city.Name != "London" {
		t.Errorf("Expected city name 'London', got '%s'", city.Name)
	}
	if city.Lat != 51.5074 {
		t.Errorf("Expected latitude 51.5074, got %f", city.Lat)
	}
	if city.Lon != -0.1278 {
		t.Errorf("Expected longitude -0.1278, got %f", city.Lon)
	}
}

func TestGetCoordinatesNoResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	originalBaseURL := baseURL
	originalGeoCodeURL := geoCodeURL
	baseURL = server.URL + "/"
	geoCodeURL = baseURL + "geo/1.0/direct?q=%s&limit=1&appid=%s"

	defer func() {
		baseURL = originalBaseURL
		geoCodeURL = originalGeoCodeURL
	}()

	_, err := GetCoordinates("NonExistentCity", "test-api-key")

	if err == nil {
		t.Fatal("Expected an error for non-existent city, but got nil")
	}
}

func TestGetWeather(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data/2.5/weather" {
			t.Errorf("Expected request to '/data/2.5/weather', got: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"main": {"temp": 283.15},
			"weather": [{"description": "cloudy"}]
		}`))
	}))
	defer server.Close()

	originalBaseURL := baseURL
	originalWeatherURL := weatherURL
	baseURL = server.URL + "/"
	weatherURL = baseURL + "data/2.5/weather?lat=%f&lon=%f&appid=%s"

	defer func() {
		baseURL = originalBaseURL
		weatherURL = originalWeatherURL
	}()

	weather, err := GetWeather(51.5074, -0.1278, "test-api-key")
	if err != nil {
		t.Fatalf("GetWeather returned an error: %v", err)
	}
	if weather.Temp != 283.15 {
		t.Errorf("Expected temperature 283.15, got %f", weather.Temp)
	}
	if weather.Description != "cloudy" {
		t.Errorf("Expected description 'cloudy', got '%s'", weather.Description)
	}
}
