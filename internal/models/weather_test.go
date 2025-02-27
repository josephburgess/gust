package models

import (
	"encoding/json"
	"testing"
)

func TestKelvinToCelsius(t *testing.T) {
	tests := []struct {
		kelvin   float64
		expected float64
	}{
		{273.15, 0.0},
		{283.15, 10.0},
		{293.15, 20.0},
		{0, -273.15},
	}

	for _, test := range tests {
		result := KelvinToCelsius(test.kelvin)
		if result != test.expected {
			t.Errorf("KelvinToCelsius(%f) = %f, expected %f", test.kelvin, result, test.expected)
		}
	}
}

func TestGetWeatherEmoji(t *testing.T) {
	tests := []struct {
		id       int
		expected string
	}{
		{200, "⚡"},  // storm
		{300, "🌦"},  // drizzle
		{500, "☔"},  // rain
		{600, "⛄"},  // snow
		{700, "🌫"},  // fog
		{800, "🔆"},  // clear
		{801, "🌥️"}, // cloudy
		{900, "🌡"},  // default
	}

	for _, test := range tests {
		result := GetWeatherEmoji(test.id)
		if result != test.expected {
			t.Errorf("GetWeatherEmoji(%d) = %s, expected %s", test.id, result, test.expected)
		}
	}
}

func TestGetWindDirection(t *testing.T) {
	tests := []struct {
		degrees  int
		expected string
	}{
		{0, "N"},
		{45, "NE"},
		{90, "E"},
		{135, "SE"},
		{180, "S"},
		{225, "SW"},
		{270, "W"},
		{315, "NW"},
		{360, "N"},  // full circle
		{400, "NE"}, // > 360
		{-45, "NW"}, // neg
	}

	for _, test := range tests {
		result := GetWindDirection(test.degrees)
		if result != test.expected {
			t.Errorf("GetWindDirection(%d) = %s, expected %s", test.degrees, result, test.expected)
		}
	}
}

func TestVisibilityToString(t *testing.T) {
	tests := []struct {
		meters   int
		contains string
	}{
		{12000, "Excellent"},
		{10000, "Excellent"},
		{7500, "Good"},
		{5000, "Good"},
		{3000, "Moderate"},
		{2000, "Moderate"},
		{1000, "Poor"},
		{500, "Poor"},
	}

	for _, test := range tests {
		result := VisibilityToString(test.meters)
		if result == "" || !contains(result, test.contains) {
			t.Errorf("VisibilityToString(%d) = %s, expected to contain %s", test.meters, result, test.contains)
		}
	}
}

func TestCityJSON(t *testing.T) {
	city := City{
		Name: "London",
		Lat:  51.5074,
		Lon:  -0.1278,
	}

	jsonData, err := json.Marshal(city)
	if err != nil {
		t.Fatalf("Failed to marshal City to JSON: %v", err)
	}

	var parsedCity City
	err = json.Unmarshal(jsonData, &parsedCity)
	if err != nil {
		t.Fatalf("Failed to unmarshal City from JSON: %v", err)
	}

	if parsedCity.Name != city.Name {
		t.Errorf("Expected Name to be %s, got %s", city.Name, parsedCity.Name)
	}

	if parsedCity.Lat != city.Lat {
		t.Errorf("Expected Lat to be %f, got %f", city.Lat, parsedCity.Lat)
	}

	if parsedCity.Lon != city.Lon {
		t.Errorf("Expected Lon to be %f, got %f", city.Lon, parsedCity.Lon)
	}
}

func TestWeatherConditionJSON(t *testing.T) {
	weatherCond := WeatherCondition{
		ID:          800,
		Main:        "Clear",
		Description: "clear sky",
		Icon:        "01d",
	}

	jsonData, err := json.Marshal(weatherCond)
	if err != nil {
		t.Fatalf("Failed to marshal WeatherCondition to JSON: %v", err)
	}

	var parsedWeather WeatherCondition
	err = json.Unmarshal(jsonData, &parsedWeather)
	if err != nil {
		t.Fatalf("Failed to unmarshal WeatherCondition from JSON: %v", err)
	}

	if parsedWeather.ID != weatherCond.ID {
		t.Errorf("Expected ID to be %d, got %d", weatherCond.ID, parsedWeather.ID)
	}

	if parsedWeather.Main != weatherCond.Main {
		t.Errorf("Expected Main to be %s, got %s", weatherCond.Main, parsedWeather.Main)
	}

	if parsedWeather.Description != weatherCond.Description {
		t.Errorf("Expected Description to be %s, got %s",
			weatherCond.Description, parsedWeather.Description)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
