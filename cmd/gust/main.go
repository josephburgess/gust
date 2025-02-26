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
	setDefaultCityPtr := flag.String("default", "", "Set a new default city")
	setAPIKeyPtr := flag.String("apikey", "", "Set a new OpenWeather API key")
	forecastPtr := flag.Bool("f", false, "Show output including forecast")
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

	if *forecastPtr {
		weather, forecast, err := api.GetWeatherAndForecast(city.Lat, city.Lon, apiKey)
		errhandler.CheckFatal(err, "Failed to get weather data")
		displayWeather(city, weather, forecast, *forecastPtr)
	} else {
		weather, err := api.GetCurrentWeather(city.Lat, city.Lon, apiKey)
		errhandler.CheckFatal(err, "Failed to get weather data")
		displayWeather(city, weather, nil, *forecastPtr)
	}
}

func displayWeather(city *models.City, weather *models.Weather, forecast []models.ForecastItem, showForecast bool) {
	// styled output funcs
	headerStyle := color.New(color.FgHiCyan, color.Bold).SprintFunc()
	tempStyle := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	highlightStyle := color.New(color.FgHiWhite).SprintFunc()
	// subtleStyle := color.New(color.FgWhite).SprintFunc()
	infoStyle := color.New(color.FgHiBlue).SprintFunc()
	// errorStyle := color.New(color.FgHiRed, color.Bold).SprintFunc()
	// successStyle := color.New(color.FgHiGreen, color.Bold).SprintFunc()

	// header
	fmt.Printf("\n%s\n", headerStyle(fmt.Sprintf("WEATHER FOR %s", strings.ToUpper(city.Name))))
	fmt.Println(strings.Repeat("─", 50))

	// emoji
	fmt.Printf("Current Conditions: %s %s\n\n",
		highlightStyle(weather.Description),
		weather.Emoji())

	// temp
	fmt.Printf("Temperature: %s %s (Feels like: %.1f°C)\n",
		tempStyle(fmt.Sprintf("%.1f°C", models.KelvinToCelsius(weather.Temp))),
		"🌡️",
		models.KelvinToCelsius(weather.FeelsLike))

	// temp range
	if weather.TempMin > 0 && weather.TempMax > 0 {
		fmt.Printf("Range: %.1f°C to %.1f°C\n",
			models.KelvinToCelsius(weather.TempMin),
			models.KelvinToCelsius(weather.TempMax))
	}

	// humidity
	fmt.Printf("Humidity: %d%% %s\n", weather.Humidity, "💧")

	// wind / gusts if avail
	if weather.WindGust > 0 {
		fmt.Printf("Wind: %.1f km/h %s %s (Gusts: %.1f km/h)\n",
			weather.WindSpeed*3.6, // m/s to km/h
			models.GetWindDirection(weather.WindDeg),
			"💨",
			weather.WindGust*3.6)
	} else {
		fmt.Printf("Wind: %.1f km/h %s %s\n",
			weather.WindSpeed*3.6, // m/s to km/h
			models.GetWindDirection(weather.WindDeg),
			"💨")
	}

	// cloud
	if weather.Clouds > 0 {
		fmt.Printf("Cloud coverage: %d%% ☁️\n", weather.Clouds)
	}

	// rain
	if weather.Rain1h > 0 {
		fmt.Printf("Rain: %.1f mm (last hour) 🌧️\n", weather.Rain1h)
	}

	// sunrise
	fmt.Printf("Sunrise: %s %s  Sunset: %s %s\n\n",
		models.FormatTimestamp(weather.Sunrise),
		"🌅",
		models.FormatTimestamp(weather.Sunset),
		"🌇")

	// forecast
	if showForecast && len(forecast) > 0 {
		fmt.Println(headerStyle("5-DAY FORECAST"))
		fmt.Println(strings.Repeat("─", 50))

		for i, day := range forecast {
			date := models.FormatDay(day.DateTime)

			// add spacer if more than 1
			if i > 0 {
				fmt.Println()
			}

			// Day header with date
			fmt.Printf("%s:\n", highlightStyle(date))

			// Temperature range with emoji
			fmt.Printf("  Highs: %s %s\n",
				tempStyle(fmt.Sprintf("%.1f°C", models.KelvinToCelsius(day.TempMax))),
				"🌡️")

			// Condition with emoji
			// fmt.Printf("  Conditions: %s %s\n",
			// 	day.Description,
			// 	day.Emoji())

			condition := fmt.Sprintf("%s %s", day.Description, day.Emoji())

			fmt.Printf("  Conditions: %s\n", infoStyle(condition))
			// // Show rain chance if > 0
			// if day.Pop > 0 {
			// 	fmt.Printf("  Rain Chance: %s\n",
			// 		infoStyle(fmt.Sprintf("%.0f%%", day.Pop*100)))
			// }
		}
		fmt.Println()
	}
}
