package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/errhandler"
)

func main() {
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		errhandler.CheckFatal(
			fmt.Errorf("no OPENWEATHER_API_KEY found"),
			"Please set an API key in the environment",
		)
	}

	cityPtr := flag.String("city", "", "Name of the city")
	flag.Parse()

	args := flag.Args()
	var cityName string

	if *cityPtr != "" {
		cityName = *cityPtr
	} else if len(args) > 0 {
		cityName = strings.Join(args, " ")
	} else {
		cityName = "London"
	}

	city, err := api.GetCoordinates(cityName, apiKey)
	errhandler.CheckFatal(err, "Failed to get coordinates")

	weather, err := api.GetWeather(city.Lat, city.Lon, apiKey)
	errhandler.CheckFatal(err, "Failed to get weather")

	fmt.Printf("Current weather in %s:\n", city.Name)
	fmt.Printf("  Temperature: %.1f°C\n", weather.Temp-273.15)
	fmt.Printf("  Conditions : %s %s\n", weather.Description, weather.Emoji())
}
