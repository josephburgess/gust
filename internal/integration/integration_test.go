package integration

import (
	"os"
	"testing"

	"github.com/josephburgess/gust/internal/api"
)

func TestWeatherIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test: OPENWEATHER_API_KEY environment variable not set")
	}

	cities := []string{"London", "New York", "Tokyo"}

	for _, city := range cities {
		t.Run(city, func(t *testing.T) {
			cityData, err := api.GetCoordinates(city, apiKey)
			if err != nil {
				t.Fatalf("Failed to get coordinates for %s: %v", city, err)
			}

			if cityData.Name == "" || cityData.Lat == 0 || cityData.Lon == 0 {
				t.Errorf("Invalid city data returned: %+v", cityData)
			}

			weather, err := api.GetWeather(cityData.Lat, cityData.Lon, apiKey)
			if err != nil {
				t.Fatalf("Failed to get weather for %s: %v", city, err)
			}

			if weather.Temp == 0 {
				t.Errorf("Temperature should not be 0 Kelvin")
			}

			if weather.Description == "" {
				t.Errorf("Weather description should not be empty")
			}

			t.Logf("Weather for %s: %.1f°C, %s",
				city, weather.Temp-273.15, weather.Description)
		})
	}
}
