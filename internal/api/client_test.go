package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	baseURL := "https://example.com"
	apiKey := "test-api-key"
	units := "metric"

	client := NewClient(baseURL, apiKey, units)

	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL to be %s, got %s", baseURL, client.baseURL)
	}

	if client.apiKey != apiKey {
		t.Errorf("Expected apiKey to be %s, got %s", apiKey, client.apiKey)
	}

	if client.client == nil {
		t.Error("HTTP client should not be nil")
	}
}

func TestGetWeather(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/weather/London" {
			t.Errorf("Expected path /api/weather/London, got %s", r.URL.Path)
		}

		if apiKey := r.URL.Query().Get("api_key"); apiKey != "test-api-key" {
			t.Errorf("Expected api_key=test-api-key, got %s", apiKey)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"city": {
				"name": "London",
				"lat": 51.5074,
				"lon": -0.1278
			},
			"weather": {
				"lat": 51.5074,
				"lon": -0.1278,
				"timezone": "Europe/London",
				"timezone_offset": 0,
				"current": {
					"dt": 1613896743,
					"temp": 283.15,
					"weather": [{"id": 800, "main": "Clear", "description": "clear sky", "icon": "01d"}]
				}
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key", "metric")

	resp, err := client.GetWeather("London")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.City.Name != "London" {
		t.Errorf("Expected city name London, got %s", resp.City.Name)
	}

	if resp.Weather.Current.Temp != 283.15 {
		t.Errorf("Expected temp 283.15, got %f", resp.Weather.Current.Temp)
	}

	if len(resp.Weather.Current.Weather) == 0 || resp.Weather.Current.Weather[0].Description != "clear sky" {
		t.Error("Weather conditions not parsed correctly")
	}
}

func TestGetWeatherError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "City not found"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key", "metric")

	resp, err := client.GetWeather("NonExistentCity")

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if resp != nil {
		t.Errorf("Expected nil response, got %+v", resp)
	}
}
