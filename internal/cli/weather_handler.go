package cli

import (
	"fmt"

	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/components"
	"github.com/josephburgess/gust/internal/ui/renderer"
	"github.com/josephburgess/gust/internal/ui/styles"
)

func fetchAndRenderWeather(city string, cfg *config.Config, authConfig *config.AuthConfig, cli *CLI) error {
	client := api.NewClient(cfg.ApiUrl, authConfig.APIKey, cfg.Units)

	fetchFunc := func() (*api.WeatherResponse, error) {
		weather, err := client.GetWeather(city)
		if err != nil {
			return nil, fmt.Errorf("failed to get weather data: %w", err)
		}
		return weather, nil
	}

	message := fmt.Sprintf("Fetching weather for %s...", city)
	weather, err := components.RunWithSpinner(message, components.WeatherEmojis, styles.Foam, fetchFunc)
	if err != nil {
		return err
	}

	weatherRenderer := renderer.NewWeatherRenderer("terminal", cfg.Units)
	renderWeatherView(cli, weatherRenderer, weather.City, weather.Weather, cfg)

	return nil
}

func renderWeatherView(cli *CLI, weatherRenderer renderer.WeatherRenderer, city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	switch {
	case cli.Alerts:
		weatherRenderer.RenderAlerts(city, weather, cfg)
	case cli.Hourly:
		weatherRenderer.RenderHourlyForecast(city, weather, cfg)
	case cli.Daily:
		weatherRenderer.RenderDailyForecast(city, weather, cfg)
	case cli.Full:
		weatherRenderer.RenderFullWeather(city, weather, cfg)
	case cli.Compact:
		weatherRenderer.RenderCompactWeather(city, weather, cfg)
	case cli.Detailed:
		weatherRenderer.RenderCurrentWeather(city, weather, cfg)
	default:
		switch cfg.DefaultView {
		case "compact":
			weatherRenderer.RenderCompactWeather(city, weather, cfg)
		case "daily":
			weatherRenderer.RenderDailyForecast(city, weather, cfg)
		case "hourly":
			weatherRenderer.RenderHourlyForecast(city, weather, cfg)
		case "full":
			weatherRenderer.RenderFullWeather(city, weather, cfg)
		default:
			weatherRenderer.RenderCurrentWeather(city, weather, cfg)
		}
	}
}
