package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/errhandler"
	"github.com/josephburgess/gust/internal/models"
)

func main() {
	_ = godotenv.Load()

	// Command-line flags
	cityPtr := flag.String("city", "", "Name of the city")
	setDefaultCityPtr := flag.String("default", "", "Set a new default city")
	loginPtr := flag.Bool("login", false, "Authenticate with GitHub")
	logoutPtr := flag.Bool("logout", false, "Log out and remove authentication")
	apiURLPtr := flag.String("api", "", "Set custom API server URL")
	forecastPtr := flag.Bool("f", false, "Show output including forecast")
	flag.Parse()

	// Load regular config
	cfg, err := config.Load()
	if err != nil {
		errhandler.CheckFatal(err, "Failed to load configuration")
	}

	// Handle API URL setting
	if *apiURLPtr != "" {
		cfg.APIURL = *apiURLPtr
		if err := cfg.Save(); err != nil {
			errhandler.CheckFatal(err, "Failed to save configuration")
		}
		fmt.Printf("API server URL set to: %s\n", *apiURLPtr)
		return
	}

	// Set default API URL if not set
	if cfg.APIURL == "" {
		cfg.APIURL = "https://gust.ngrok.io" // Default API server
		if err := cfg.Save(); err != nil {
			errhandler.CheckFatal(err, "Failed to save configuration")
		}
	}

	// Handle logout
	if *logoutPtr {
		authConfig, _ := config.LoadAuthConfig()
		if authConfig != nil {
			authConfigPath, err := config.GetAuthConfigPath()
			if err != nil {
				log.Fatalf("Failed to get auth config path: %v", err)
			}

			if err := os.Remove(authConfigPath); err != nil {
				log.Fatalf("Failed to remove auth config file: %v", err)
			}
			fmt.Println("Logged out successfully")
		} else {
			fmt.Println("Not currently logged in")
		}
		return
	}

	// Handle login
	if *loginPtr {
		fmt.Println("Starting GitHub authentication...")
		authConfig, err := config.Authenticate(cfg.APIURL)
		if err != nil {
			errhandler.CheckFatal(err, "Authentication failed")
		}

		if err := config.SaveAuthConfig(authConfig); err != nil {
			errhandler.CheckFatal(err, "Failed to save authentication")
		}

		fmt.Printf("Successfully authenticated as %s\n", authConfig.GithubUser)
		return
	}

	// Handle default city setting
	if *setDefaultCityPtr != "" {
		cfg.DefaultCity = *setDefaultCityPtr
		if err := cfg.Save(); err != nil {
			errhandler.CheckFatal(err, "Failed to save configuration")
		}
		fmt.Printf("Default city set to: %s\n", *setDefaultCityPtr)
		return
	}

	// Load auth config
	authConfig, err := config.LoadAuthConfig()
	if err != nil {
		errhandler.CheckFatal(err, "Failed to load authentication")
	}

	// Check if auth is required
	if authConfig == nil {
		fmt.Println("You need to authenticate with GitHub before using Gust.")
		fmt.Println("Run 'gust --login' to authenticate.")
		os.Exit(1)
	}

	// Determine city
	var cityName string
	args := flag.Args()
	if *cityPtr != "" {
		cityName = *cityPtr
	} else if len(args) > 0 {
		cityName = strings.Join(args, " ")
	} else {
		cityName = cfg.DefaultCity
	}

	// Check if default city is set
	if cityName == "" {
		fmt.Println("No city specified and no default city set.")
		fmt.Println("Specify a city: gust [city name]")
		fmt.Println("Or set a default city: gust --default \"London\"")
		os.Exit(1)
	}

	// Make API request to get weather data
	if *forecastPtr {
		weather, forecast, err := getWeatherAndForecast(cfg.APIURL, authConfig.APIKey, cityName)
		errhandler.CheckFatal(err, "Failed to get weather data")
		displayWeather(weather.City, weather.Weather, forecast, *forecastPtr)
	} else {
		weather, err := getCurrentWeather(cfg.APIURL, authConfig.APIKey, cityName)
		errhandler.CheckFatal(err, "Failed to get weather data")
		displayWeather(weather.City, weather.Weather, nil, *forecastPtr)
	}
}

// API response structure
type WeatherResponse struct {
	City    *models.City    `json:"city"`
	Weather *models.Weather `json:"weather"`
}

type ForecastResponse struct {
	City     *models.City          `json:"city"`
	Forecast []models.ForecastItem `json:"forecast"`
}

// Get current weather from API
func getCurrentWeather(apiURL, apiKey, cityName string) (*WeatherResponse, error) {
	url := fmt.Sprintf("%s/api/weather/%s?api_key=%s",
		apiURL, url.QueryEscape(cityName), apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var response WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return &response, nil
}

// Get weather and forecast from API
func getWeatherAndForecast(apiURL, apiKey, cityName string) (*WeatherResponse, []models.ForecastItem, error) {
	// Get current weather
	weather, err := getCurrentWeather(apiURL, apiKey, cityName)
	if err != nil {
		return nil, nil, err
	}

	// Get forecast
	url := fmt.Sprintf("%s/api/forecast/%s?api_key=%s",
		apiURL, url.QueryEscape(cityName), apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return weather, nil, fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return weather, nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var forecastResp ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecastResp); err != nil {
		return weather, nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return weather, forecastResp.Forecast, nil
}

// The displayWeather function remains the same as in your original code
func displayWeather(city *models.City, weather *models.Weather, forecast []models.ForecastItem, showForecast bool) {
	// styled output funcs
	headerStyle := color.New(color.FgHiCyan, color.Bold).SprintFunc()
	tempStyle := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	highlightStyle := color.New(color.FgHiWhite).SprintFunc()
	infoStyle := color.New(color.FgHiBlue).SprintFunc()

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

			condition := fmt.Sprintf("%s %s", day.Description, day.Emoji())
			fmt.Printf("  Conditions: %s\n", infoStyle(condition))
		}
		fmt.Println()
	}
}
