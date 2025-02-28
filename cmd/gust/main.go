package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/ui"
)

func main() {
	_ = godotenv.Load()

	// flags
	cityPtr := flag.String("city", "", "Name of the city")
	setDefaultCityPtr := flag.String("default", "", "Set a new default city")
	loginPtr := flag.Bool("login", false, "Authenticate with GitHub")
	logoutPtr := flag.Bool("logout", false, "Log out and remove authentication")
	apiURLPtr := flag.String("api", "", "Set custom API server URL")
	fullPtr := flag.Bool("full", false, "Show full weather report including daily and hourly forecasts")
	dailyPtr := flag.Bool("daily", false, "Show daily forecast only")
	hourlyPtr := flag.Bool("hourly", false, "Show hourly forecast only")
	alertsPtr := flag.Bool("alerts", false, "Show weather alerts only")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		ui.ExitWithError("Failed to load configuration", err)
	}

	if *apiURLPtr != "" {
		cfg.APIURL = *apiURLPtr
		if err := cfg.Save(); err != nil {
			ui.ExitWithError("Failed to save configuration", err)
		}
		fmt.Printf("API server URL set to: %s\n", *apiURLPtr)
		return
	}

	if cfg.APIURL == "" {
		cfg.APIURL = "https://gust.ngrok.io"
		if err := cfg.Save(); err != nil {
			ui.ExitWithError("Failed to save configuration", err)
		}
	}

	if *logoutPtr {
		handleLogout()
		return
	}

	if *loginPtr {
		handleLogin(cfg.APIURL)
		return
	}

	if *setDefaultCityPtr != "" {
		cfg.DefaultCity = *setDefaultCityPtr
		if err := cfg.Save(); err != nil {
			ui.ExitWithError("Failed to save configuration", err)
		}
		fmt.Printf("Default city set to: %s\n", *setDefaultCityPtr)
		return
	}

	authConfig, err := config.LoadAuthConfig()
	if err != nil {
		ui.ExitWithError("Failed to load authentication", err)
	}

	if authConfig == nil {
		fmt.Println("You need to authenticate with GitHub before using Gust.")
		fmt.Println("Run 'gust --login' to authenticate.")
		os.Exit(1)
	}

	cityName := determineCityName(*cityPtr, flag.Args(), cfg.DefaultCity)
	if cityName == "" {
		fmt.Println("No city specified and no default city set.")
		fmt.Println("Specify a city: gust [city name]")
		fmt.Println("Or set a default city: gust --default \"London\"")
		os.Exit(1)
	}

	client := api.NewClient(cfg.APIURL, authConfig.APIKey)

	weather, err := client.GetWeather(cityName)
	if err != nil {
		ui.ExitWithError("Failed to get weather data", err)
	}

	renderer := ui.NewRenderer()
	if *alertsPtr {
		renderer.DisplayAlerts(weather.City, weather.Weather)
	} else if *hourlyPtr {
		renderer.DisplayHourlyForecast(weather.City, weather.Weather)
	} else if *dailyPtr {
		renderer.DisplayDailyForecast(weather.City, weather.Weather)
	} else if *fullPtr {
		renderer.DisplayFullWeather(weather.City, weather.Weather)
	} else {
		renderer.DisplayCurrentWeather(weather.City, weather.Weather)
	}
}

func handleLogout() {
	authConfig, _ := config.LoadAuthConfig()
	if authConfig == nil {
		fmt.Println("Not currently logged in")
		return
	}

	authConfigPath, err := config.GetAuthConfigPath()
	if err != nil {
		ui.ExitWithError("Failed to get auth config path", err)
	}

	if err := os.Remove(authConfigPath); err != nil {
		ui.ExitWithError("Failed to remove auth config file", err)
	}

	fmt.Println("Logged out successfully")
}

func handleLogin(apiURL string) {
	fmt.Println("Starting GitHub authentication...")
	authConfig, err := config.Authenticate(apiURL)
	if err != nil {
		ui.ExitWithError("Authentication failed", err)
	}

	if err := config.SaveAuthConfig(authConfig); err != nil {
		ui.ExitWithError("Failed to save authentication", err)
	}

	fmt.Printf("Successfully authenticated as %s\n", authConfig.GithubUser)
}

func determineCityName(cityFlag string, args []string, defaultCity string) string {
	if cityFlag != "" {
		return cityFlag
	}

	if len(args) > 0 {
		return strings.Join(args, " ")
	}

	return defaultCity
}
