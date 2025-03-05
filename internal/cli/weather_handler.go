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
	renderWeatherView(cli, weatherRenderer, weather.City, weather.Weather)

	return nil
}

func renderWeatherView(cli *CLI, weatherRenderer renderer.WeatherRenderer, city *models.City, weather *models.OneCallResponse) {
	if cli.Alerts {
		weatherRenderer.RenderAlerts(city, weather)
	} else if cli.Hourly {
		weatherRenderer.RenderHourlyForecast(city, weather)
	} else if cli.Daily {
		weatherRenderer.RenderDailyForecast(city, weather)
	} else if cli.Full {
		weatherRenderer.RenderFullWeather(city, weather)
	} else if cli.Compact {
		weatherRenderer.RenderCompactWeather(city, weather)
	} else if cli.Detailed {
		weatherRenderer.RenderCurrentWeather(city, weather)
	} else {
		weatherRenderer.RenderCurrentWeather(city, weather)
	}
}
