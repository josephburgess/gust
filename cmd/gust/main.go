package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/errhandler"
	"github.com/josephburgess/gust/internal/models"
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

	weather, forecast, err := api.GetWeatherAndForecast(city.Lat, city.Lon, apiKey)
	errhandler.CheckFatal(err, "Failed to get weather data")

	displayWeather(city, weather, forecast)
}

func displayWeather(city *models.City, weather *models.Weather, forecast []models.ForecastItem) {
	// styled output funcs
	headerStyle := color.New(color.FgHiCyan, color.Bold).SprintFunc()
	tempStyle := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	highlightStyle := color.New(color.FgHiWhite).SprintFunc()
	subtleStyle := color.New(color.FgWhite).SprintFunc()
	infoStyle := color.New(color.FgHiBlue).SprintFunc()

	// header
	fmt.Printf("\n%s\n", headerStyle(fmt.Sprintf("WEATHER FOR %s", strings.ToUpper(city.Name))))
	fmt.Println(strings.Repeat("─", 50))

	// emoji
	fmt.Printf("Current Conditions: %s %s\n\n",
		highlightStyle(weather.Description),
		weather.Emoji())

	// temp
	fmt.Printf("Temperature: %s (Feels like: %.1f°C)\n",
		tempStyle(fmt.Sprintf("%.1f°C", models.KelvinToCelsius(weather.Temp))),
		models.KelvinToCelsius(weather.FeelsLike))

	// temp range
	if weather.TempMin > 0 && weather.TempMax > 0 {
		fmt.Printf("Range: %.1f°C to %.1f°C\n",
			models.KelvinToCelsius(weather.TempMin),
			models.KelvinToCelsius(weather.TempMax))
	}

	// humidity
	fmt.Printf("Humidity: %d%%\n", weather.Humidity)

	// wind / gusts if avail
	if weather.WindGust > 0 {
		fmt.Printf("Wind: %.1f km/h %s (Gusts: %.1f km/h)\n",
			weather.WindSpeed*3.6, // Convert m/s to km/h
			models.GetWindDirection(weather.WindDeg),
			weather.WindGust*3.6)
	} else {
		fmt.Printf("Wind: %.1f km/h %s\n",
			weather.WindSpeed*3.6, // Convert m/s to km/h
			models.GetWindDirection(weather.WindDeg))
	}

	// pressure
	fmt.Printf("Pressure: %d hPa\n", weather.Pressure)

	// visibility
	if weather.Visibility > 0 {
		fmt.Printf("Visibility: %s\n", models.VisibilityToString(weather.Visibility))
	}

	// cloud
	if weather.Clouds > 0 {
		fmt.Printf("Cloud coverage: %d%%\n", weather.Clouds)
	}

	// rain
	if weather.Rain1h > 0 {
		fmt.Printf("Rain: %.1f mm (last hour)\n", weather.Rain1h)
	}

	// sunrise
	fmt.Printf("Sunrise: %s  Sunset: %s\n\n",
		models.FormatTimestamp(weather.Sunrise),
		models.FormatTimestamp(weather.Sunset))

	// forecast
	if len(forecast) > 0 {
		fmt.Println(headerStyle("5-DAY FORECAST"))
		fmt.Println(strings.Repeat("─", 50))

		for _, day := range forecast {
			date := models.FormatDay(day.DateTime)
			rainChance := ""
			if day.Pop > 0 {
				rainChance = fmt.Sprintf(" %s", infoStyle(fmt.Sprintf("(%.0f%% chance of rain)", day.Pop*100)))
			}

			fmt.Printf("%s: %s to %s  %s %s%s\n",
				highlightStyle(date),
				subtleStyle(fmt.Sprintf("%.1f°C", models.KelvinToCelsius(day.TempMin))),
				tempStyle(fmt.Sprintf("%.1f°C", models.KelvinToCelsius(day.TempMax))),
				day.Description,
				day.Emoji(),
				rainChance)
		}
		fmt.Println()
	}
}
