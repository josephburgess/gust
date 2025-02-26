// File: /gust/cmd/gust/main.go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/errhandler"
)

func main() {
	err := godotenv.Load()
	errhandler.CheckFatal(err, "Error loading .env file")

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		errhandler.CheckFatal(fmt.Errorf("API key not found"), "Missing OpenWeather API key")
	}

	cityPtr := flag.String("city", "London", "City name to fetch weather for")
	flag.Parse()

	cityName := *cityPtr

	city, err := api.GetCoordinates(cityName, apiKey)
	errhandler.CheckFatal(err, "Failed to get coordinates")

	weather, err := api.GetWeather(city.Lat, city.Lon, apiKey)
	errhandler.CheckFatal(err, "Failed to get weather")

	fmt.Printf("Weather in %s:\n", city.Name)
	fmt.Printf("Temperature: %.1f°C\n", weather.Temp-273.15) // Kelvin to °C
	fmt.Printf("Conditions: %s\n", weather.Description)
}
