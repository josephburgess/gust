package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/josephburgess/gust/internal/api"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	cityName := "London"

	city, err := api.GetCoordinates(cityName, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	weather, err := api.GetWeather(city.Lat, city.Lon, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Weather in %s:\n", city.Name)
	fmt.Printf("Temperature: %.1f°C\n", weather.Temp-273.15)
	fmt.Printf("Conditions: %s\n", weather.Description)
}
