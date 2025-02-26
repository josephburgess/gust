// cmd/gust/main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/errhandler"
)

func main() {
	_ = godotenv.Load()

	cityPtr := flag.String("city", "", "Name of the city")
	setDefaultCityPtr := flag.String("set-default-city", "", "Set a new default city")
	setAPIKeyPtr := flag.String("set-api-key", "", "Set a new OpenWeather API key")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		errhandler.CheckFatal(err, "Failed to load configuration")
	}

	if *setDefaultCityPtr != "" {
		cfg.DefaultCity = *setDefaultCityPtr
		if err := cfg.Save(); err != nil {
			errhandler.CheckFatal(err, "Failed to save configuration")
		}
		fmt.Printf("Default city set to: %s\n", *setDefaultCityPtr)
		return
	}

	if *setAPIKeyPtr != "" {
		cfg.APIKey = *setAPIKeyPtr
		if err := cfg.Save(); err != nil {
			errhandler.CheckFatal(err, "Failed to save configuration")
		}
		fmt.Println("API key updated successfully")
		return
	}

	if cfg.APIKey == "" || cfg.DefaultCity == "" {
		fmt.Println("First-time setup required")
		newCfg, err := config.PromptForConfiguration()
		if err != nil {
			errhandler.CheckFatal(err, "Failed to configure")
		}
		*cfg = *newCfg
		if err := cfg.Save(); err != nil {
			errhandler.CheckFatal(err, "Failed to save configuration")
		}
	}

	apiKey := cfg.APIKey

	if apiKey == "" {
		apiKey = os.Getenv("OPENWEATHER_API_KEY")
	}

	if apiKey == "" {
		errhandler.CheckFatal(
			fmt.Errorf("no OpenWeather API key found"),
			"Please set an API key using --set-api-key or run the setup again",
		)
	}

	var cityName string

	args := flag.Args()
	if *cityPtr != "" {
		cityName = *cityPtr
	} else if len(args) > 0 {
		cityName = strings.Join(args, " ")
	} else {
		cityName = cfg.DefaultCity
	}

	city, err := api.GetCoordinates(cityName, apiKey)
	errhandler.CheckFatal(err, "Failed to get coordinates")

	weather, err := api.GetWeather(city.Lat, city.Lon, apiKey)
	errhandler.CheckFatal(err, "Failed to get weather")

	fmt.Printf("Current weather in %s:\n", city.Name)
	fmt.Printf("  Temperature: %.1f°C\n", weather.Temp-273.15)
	fmt.Printf("  Conditions : %s %s\n", weather.Description, weather.Emoji())
}
