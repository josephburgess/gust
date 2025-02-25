package main

import (
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

	cityName := "London"

	city, err := api.GetCoordinates(cityName, apiKey)
	errhandler.CheckFatal(err, "Failed to get coordinates")

	weather, err := api.GetWeather(city.Lat, city.Lon, apiKey)
	errhandler.CheckFatal(err, "Failed to get weather")

	fmt.Printf("Weather in %s:\n", city.Name)
	fmt.Printf("Temperature: %.1f°C\n", weather.Temp-273.15) // Convert from Kelvin
	fmt.Printf("Conditions: %s\n", weather.Description)
}
