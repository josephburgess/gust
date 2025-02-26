package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/errhandler"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var weatherIcons = map[string]string{
	"clear sky":            "☀️  ",
	"few clouds":           "🌤️  ",
	"scattered clouds":     "☁️  ",
	"broken clouds":        "☁️  ",
	"shower rain":          "🌧️  ",
	"rain":                 "🌧️  ",
	"thunderstorm":         "⛈️  ",
	"snow":                 "❄️  ",
	"mist":                 "🌫️  ",
	"overcast clouds":      "☁️  ",
	"light rain":           "🌦️  ",
	"moderate rain":        "🌧️  ",
	"heavy intensity rain": "🌧️  ",
	"cloudy":               "☁️  ",
	"partly cloudy":        "⛅ ",
	"fog":                  "🌫️  ",
	"haze":                 "🌫️  ",
	"dust":                 "🌫️  ",
	"smoke":                "🌫️  ",
}

func getWindDirection(degrees int) string {
	directions := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
	index := int((float64(degrees)+11.25)/22.5) % 16
	return directions[index]
}

func getWeatherIcon(description string) string {
	desc := strings.ToLower(description)

	for pattern, icon := range weatherIcons {
		if strings.ToLower(pattern) == desc {
			return icon
		}
	}

	for pattern, icon := range weatherIcons {
		if strings.Contains(desc, strings.ToLower(pattern)) {
			return icon
		}
	}

	return "🌡️  "
}

func kelvinToCelsius(kelvin float64) float64 {
	return kelvin - 273.15
}

func formatVisibility(meters int) string {
	if meters >= 10000 {
		return "10+ km"
	}
	return fmt.Sprintf("%.1f km", float64(meters)/1000)
}

func formatTime(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("15:04")
}

func main() {
	err := godotenv.Load()
	errhandler.CheckFatal(err, "Error loading .env file")

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		errhandler.CheckFatal(fmt.Errorf("API key not found"), "Missing OpenWeather API key")
	}

	cityName := "London"
	if len(os.Args) > 1 {
		cityName = os.Args[1]
	}

	boldWhite := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc()

	now := time.Now()
	dateStr := cyan(now.Format("Monday, January 2, 2006"))
	timeStr := cyan(now.Format("15:04"))

	fmt.Printf("Fetching weather data for %s...\n", boldWhite(cityName))
	city, err := api.GetCoordinates(cityName, apiKey)
	errhandler.CheckFatal(err, "Failed to get coordinates")

	weather, err := api.GetWeather(city.Lat, city.Lon, apiKey)
	errhandler.CheckFatal(err, "Failed to get weather")

	tempC := kelvinToCelsius(weather.Temp)
	feelsLikeC := kelvinToCelsius(weather.FeelsLike)
	tempMinC := kelvinToCelsius(weather.TempMin)
	tempMaxC := kelvinToCelsius(weather.TempMax)

	icon := getWeatherIcon(weather.Description)

	fmt.Print("\033[H\033[2J")

	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("   %s Weather - %s at %s\n", boldWhite("Gust"), dateStr, timeStr)
	fmt.Println(strings.Repeat("─", 60))

	locationStr := city.Name
	if city.Country != "" {
		locationStr += ", " + city.Country
	}
	fmt.Printf("  Location: %s\n", boldWhite(locationStr))
	fmt.Printf("  Coordinates: %s\n", gray(fmt.Sprintf("%.4f, %.4f", city.Lat, city.Lon)))
	fmt.Println()

	caser := cases.Title(language.English)
	fmt.Printf("  %s %s\n", icon, boldWhite(caser.String(weather.Description)))
	fmt.Printf("  Temperature: %s", yellow(fmt.Sprintf("%.1f°C", tempC)))
	fmt.Printf("  (Feels like: %s)\n", yellow(fmt.Sprintf("%.1f°C", feelsLikeC)))

	fmt.Printf("  High: %s  Low: %s\n",
		yellow(fmt.Sprintf("%.1f°C", tempMaxC)),
		blue(fmt.Sprintf("%.1f°C", tempMinC)))
	fmt.Println()

	fmt.Printf("  Wind: %s %s %s\n",
		green(fmt.Sprintf("%.1f m/s", weather.WindSpeed)),
		green(getWindDirection(weather.WindDeg)),
		gray(fmt.Sprintf("(%d°)", weather.WindDeg)))

	fmt.Printf("  Humidity: %s  Pressure: %s\n",
		cyan(fmt.Sprintf("%d%%", weather.Humidity)),
		magenta(fmt.Sprintf("%d hPa", weather.Pressure)))

	fmt.Printf("  Visibility: %s  Clouds: %s\n",
		blue(formatVisibility(weather.Visibility)),
		gray(fmt.Sprintf("%d%%", weather.CloudsAll)))

	if weather.Sunrise > 0 && weather.Sunset > 0 {
		fmt.Printf("  Sunrise: %s  Sunset: %s\n",
			yellow(formatTime(weather.Sunrise)),
			blue(formatTime(weather.Sunset)))
	}

	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("  Data from OpenWeather API | %s\n", gray("github.com/josephburgess/gust"))
	fmt.Println(strings.Repeat("─", 60))
}
