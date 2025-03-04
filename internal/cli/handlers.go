package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/ui/renderer"
	"github.com/josephburgess/gust/internal/ui/setup"
)

func isValidUnit(unit string) bool {
	validUnits := map[string]bool{
		"metric":   true,
		"imperial": true,
		"standard": true,
	}

	return validUnits[unit]
}

func Run(ctx *kong.Context, cli *CLI) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if cli.APIURL != "" {
		cfg.APIURL = cli.APIURL
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
		fmt.Printf("API server URL set to: %s\n", cli.APIURL)
		return nil
	}

	if cfg.APIURL == "" {
		cfg.APIURL = "https://breeze.joeburgess.dev"
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
	}

	if cli.Login {
		return handleLogin(cfg.APIURL)
	}

	authConfig, err := config.LoadAuthConfig()
	if err != nil {
		return fmt.Errorf("failed to load authentication: %w", err)
	}

	needsSetup := (cfg.DefaultCity == "" || cli.Setup)
	needsAuth := authConfig == nil

	if needsSetup {
		if err := setup.RunSetup(cfg, needsAuth); err != nil {
			return fmt.Errorf("setup failed: %w", err)
		}

		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration after setup: %w", err)
		}

		authConfig, err = config.LoadAuthConfig()
		if err != nil {
			return fmt.Errorf("failed to load authentication after setup: %w", err)
		}

		if cli.City == "" && len(cli.Args) == 0 && !cli.Full && !cli.Daily &&
			!cli.Hourly && !cli.Alerts && !cli.Pretty {
			fmt.Println("Setup complete! Run 'gust' to check the weather for your default city.")
			return nil
		}
	}

	if cli.Units != "" {
		if !isValidUnit(cli.Units) {
			fmt.Println("Invalid units value. Must be one of: metric, imperial, standard")
			os.Exit(1)
		}

		cfg.Units = cli.Units
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Printf("Units set to: %s\n", cli.Units)
		return nil
	}

	if cli.Default != "" {
		cfg.DefaultCity = cli.Default
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
		fmt.Printf("Default city set to: %s\n", cli.Default)
		return nil
	}

	if authConfig == nil {
		fmt.Println("You need to authenticate with GitHub before using Gust.")
		fmt.Println("Run 'gust --login' to authenticate or 'gust --setup' to run the setup wizard.")
		os.Exit(1)
	}

	cityName := determineCityName(cli.City, cli.Args, cfg.DefaultCity)
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
		return fmt.Errorf("failed to get weather data: %w", err)
	}

	weatherRenderer := renderer.NewWeatherRenderer("terminal", cfg.Units)

	if cli.Alerts {
		weatherRenderer.RenderAlerts(weather.City, weather.Weather)
	} else if cli.Hourly {
		weatherRenderer.RenderHourlyForecast(weather.City, weather.Weather)
	} else if cli.Daily {
		weatherRenderer.RenderDailyForecast(weather.City, weather.Weather)
	} else if cli.Full {
		weatherRenderer.RenderFullWeather(weather.City, weather.Weather)
	} else if cli.Compact {
		weatherRenderer.RenderCompactWeather(weather.City, weather.Weather)
	} else if cli.Detailed {
		weatherRenderer.RenderCurrentWeather(weather.City, weather.Weather)
	} else {
		switch cfg.DefaultView {
		case "compact":
			weatherRenderer.RenderCompactWeather(weather.City, weather.Weather)
		case "daily":
			weatherRenderer.RenderDailyForecast(weather.City, weather.Weather)
		case "hourly":
			weatherRenderer.RenderHourlyForecast(weather.City, weather.Weather)
		case "full":
			weatherRenderer.RenderFullWeather(weather.City, weather.Weather)
		case "default", "":
			weatherRenderer.RenderCurrentWeather(weather.City, weather.Weather)
		}
	}

	return nil
}

func handleLogin(apiURL string) error {
	fmt.Println("Starting GitHub authentication...")
	authConfig, err := config.Authenticate(apiURL)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if err := config.SaveAuthConfig(authConfig); err != nil {
		return fmt.Errorf("failed to save authentication: %w", err)
	}

	fmt.Printf("Successfully authenticated as %s\n", authConfig.GithubUser)
	return nil
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
