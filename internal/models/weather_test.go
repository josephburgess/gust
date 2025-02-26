package models

import (
	"encoding/json"
	"testing"
)

// TestCityJSON tests JSON marshaling and unmarshaling for City struct
func TestCityJSON(t *testing.T) {
	// Create a test city
	city := City{
		Name: "London",
		Lat:  51.5074,
		Lon:  -0.1278,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(city)
	if err != nil {
		t.Fatalf("Failed to marshal City to JSON: %v", err)
	}

	// Convert back from JSON to a new struct
	var parsedCity City
	err = json.Unmarshal(jsonData, &parsedCity)
	if err != nil {
		t.Fatalf("Failed to unmarshal City from JSON: %v", err)
	}

	// Check that the values match
	if parsedCity.Name != city.Name {
		t.Errorf("City name doesn't match: expected %s, got %s", city.Name, parsedCity.Name)
	}
	if parsedCity.Lat != city.Lat {
		t.Errorf("City latitude doesn't match: expected %f, got %f", city.Lat, parsedCity.Lat)
	}
	if parsedCity.Lon != city.Lon {
		t.Errorf("City longitude doesn't match: expected %f, got %f", city.Lon, parsedCity.Lon)
	}
}

// TestWeatherJSON tests JSON marshaling and unmarshaling for Weather struct
func TestWeatherJSON(t *testing.T) {
	// Create a test weather
	weather := Weather{
		Temp:        283.15,
		Description: "cloudy",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(weather)
	if err != nil {
		t.Fatalf("Failed to marshal Weather to JSON: %v", err)
	}

	// Convert back from JSON to a new struct
	var parsedWeather Weather
	err = json.Unmarshal(jsonData, &parsedWeather)
	if err != nil {
		t.Fatalf("Failed to unmarshal Weather from JSON: %v", err)
	}

	// Check that the values match
	if parsedWeather.Temp != weather.Temp {
		t.Errorf("Weather temperature doesn't match: expected %f, got %f", weather.Temp, parsedWeather.Temp)
	}
	if parsedWeather.Description != weather.Description {
		t.Errorf("Weather description doesn't match: expected %s, got %s", weather.Description, parsedWeather.Description)
	}
}

// TestParseWeatherJSON tests parsing a sample API response
func TestParseWeatherJSON(t *testing.T) {
	// Sample JSON that might be returned by the API
	jsonData := `{
		"temp": 283.15,
		"description": "cloudy"
	}`

	// Parse the JSON
	var weather Weather
	err := json.Unmarshal([]byte(jsonData), &weather)
	if err != nil {
		t.Fatalf("Failed to parse Weather JSON: %v", err)
	}

	// Check the values
	if weather.Temp != 283.15 {
		t.Errorf("Expected temperature 283.15, got %f", weather.Temp)
	}
	if weather.Description != "cloudy" {
		t.Errorf("Expected description 'cloudy', got '%s'", weather.Description)
	}
}
