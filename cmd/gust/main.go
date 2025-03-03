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
	defaulPtr := flag.String("default", "", "Set a new default city")
	loginPtr := flag.Bool("login", false, "Authenticate with GitHub")
	apiURLPtr := flag.String("api", "", "Set custom API server URL")
	fullPtr := flag.Bool("full", false, "Show full weather report including daily and hourly forecasts")
	dailyPtr := flag.Bool("daily", false, "Show daily forecast only")
	hourlyPtr := flag.Bool("hourly", false, "Show hourly forecast only")
	alertsPtr := flag.Bool("alerts", false, "Show weather alerts only")
	unitsPtr := flag.String("units", "", "Temperature units (metric, imperial, standard). Metric is default")
	setupPtr := flag.Bool("setup", false, "Run the interactive setup wizard")
	prettyPtr := flag.Bool("pretty", false, "Use the interactive UI to display weather")
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
		cfg.APIURL = "http://breeze.joeburgess.dev"
		if err := cfg.Save(); err != nil {
			ui.ExitWithError("Failed to save configuration", err)
		}
	}

	if *loginPtr {
		handleLogin(cfg.APIURL)
		return
	}

	authConfig, err := config.LoadAuthConfig()
	if err != nil {
		ui.ExitWithError("Failed to load authentication", err)
	}

	needsSetup := (cfg.DefaultCity == "" || *setupPtr)
	needsAuth := authConfig == nil

	if needsSetup {
		if err := ui.RunSetup(cfg, needsAuth); err != nil {
			ui.ExitWithError("Setup failed", err)
		}

		cfg, err = config.Load()
		if err != nil {
			ui.ExitWithError("Failed to load configuration after setup", err)
		}

		authConfig, err = config.LoadAuthConfig()
		if err != nil {
			ui.ExitWithError("Failed to load authentication after setup", err)
		}

		if *cityPtr == "" && len(flag.Args()) == 0 && !*fullPtr && !*dailyPtr && !*hourlyPtr && !*alertsPtr && !*prettyPtr {
			fmt.Println("Setup complete! Run 'gust' to check the weather for your default city.")
			return
		}
	}

	if *unitsPtr != "" {
		if !isValidUnit(*unitsPtr) {
			fmt.Println("Invalid units value. Must be one of: metric, imperial, standard")
			os.Exit(1)
		}

		cfg.Units = *unitsPtr
		if err := cfg.Save(); err != nil {
			ui.ExitWithError("Failed to save config", err)
		}
		fmt.Printf("Units set to: %s\n", *unitsPtr)
		return
	}

	if *defaulPtr != "" {
		cfg.DefaultCity = *defaulPtr
		if err := cfg.Save(); err != nil {
			ui.ExitWithError("Failed to save configuration", err)
		}
		fmt.Printf("Default city set to: %s\n", *defaulPtr)
		return
	}

	if authConfig == nil {
		fmt.Println("You need to authenticate with GitHub before using Gust.")
		fmt.Println("Run 'gust --login' to authenticate or 'gust --setup' to run the setup wizard.")
		os.Exit(1)
	}

	cityName := determineCityName(*cityPtr, flag.Args(), cfg.DefaultCity)
	if cityName == "" {
		fmt.Println("No city specified and no default city set.")
		fmt.Println("Specify a city: gust [city name]")
		fmt.Println("Or set a default city: gust --default \"London\"")
		fmt.Println("Or run the setup wizard: gust --setup")
		os.Exit(1)
	}

	client := api.NewClient(cfg.APIURL, authConfig.APIKey, cfg.Units)

	weather, err := client.GetWeather(cityName)
	if err != nil {
		ui.ExitWithError("Failed to get weather data", err)
	}

	renderer := ui.NewRenderer(cfg.Units)
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

func isValidUnit(unit string) bool {
	validUnits := map[string]bool{
		"metric":   true,
		"imperial": true,
		"standard": true,
	}

	return validUnits[unit]
}
