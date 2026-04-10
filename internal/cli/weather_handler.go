package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/cache"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/components"
	"github.com/josephburgess/gust/internal/ui/output"
	"github.com/josephburgess/gust/internal/ui/renderer"
	"github.com/josephburgess/gust/internal/ui/styles"
)

func fetchAndRenderWeather(city string, cfg *config.Config, authConfig *config.AuthConfig, cli *CLI) error {
	weatherCache, err := cache.New()
	if err != nil {
		// cache failure is non-fatal — just skip caching
		weatherCache = nil
	}

	// serve from cache if available and refresh not requested
	if !cli.Refresh && weatherCache != nil {
		if cached, age, ok := weatherCache.Get(city, cfg.Units); ok {
			ttlRemaining := cache.TTL - age
			weatherRenderer := renderer.NewWeatherRenderer("terminal", cfg.Units)
			renderWeatherView(cli, weatherRenderer, cached.City, cached.Weather, cfg)
			fmt.Printf("%s\n", styles.HintStyle.Render(
				fmt.Sprintf("↩ cached %s · refreshes in %dm · gust -R to force refresh",
					cache.FormatAge(age),
					int(ttlRemaining.Minutes())+1,
				),
			))
			return nil
		}
	}

	client := api.NewClient(cfg.ApiUrl, authConfig.APIKey, cfg.Units)

	fetchFunc := func() (*api.WeatherResponse, error) {
		weather, err := client.GetWeather(city)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "rate limit") {
				return nil, fmt.Errorf("rate limit reached: %w", err)
			}
			return nil, fmt.Errorf("failed to get weather data: %w", err)
		}
		return weather, nil
	}

	message := fmt.Sprintf("Fetching weather for %s...", city)
	weather, err := components.RunWithSpinner(message, components.WeatherEmojis, styles.Foam, fetchFunc)

	if client.RateLimitInfo != nil && client.RateLimitInfo.Limit > 0 {
		if err != nil && strings.Contains(strings.ToLower(err.Error()), "rate limit") {
			output.PrintRateLimitError(client.RateLimitInfo.Limit, client.RateLimitInfo.ResetTime)

			timeUntilReset := time.Until(client.RateLimitInfo.ResetTime)
			if timeUntilReset > 0 {
				minutesRemaining := int(timeUntilReset.Minutes()) + 1
				hoursRemaining := minutesRemaining / 60

				if hoursRemaining > 0 {
					remainingMinutes := minutesRemaining % 60
					return fmt.Errorf("please try again in about %d hour(s) and %d minute(s) when your rate limit resets",
						hoursRemaining, remainingMinutes)
				}
				return fmt.Errorf("please try again in about %d minute(s) when your rate limit resets",
					minutesRemaining)
			}
			return fmt.Errorf("rate limit reached, please try again later")
		}

		if client.RateLimitInfo.Remaining <= 5 && client.RateLimitInfo.Remaining > 0 {
			output.PrintRateLimitWarning(
				client.RateLimitInfo.Remaining,
				client.RateLimitInfo.Limit,
				client.RateLimitInfo.ResetTime,
			)
		}
	}

	if err != nil {
		return err
	}

	if weatherCache != nil {
		// non-fatal if store fails
		weatherCache.Set(city, cfg.Units, weather)
	}

	weatherRenderer := renderer.NewWeatherRenderer("terminal", cfg.Units)
	renderWeatherView(cli, weatherRenderer, weather.City, weather.Weather, cfg)

	return nil
}

func handleStatus(cfg *config.Config, authConfig *config.AuthConfig) error {
	client := api.NewClient(cfg.ApiUrl, authConfig.APIKey, cfg.Units)
	quota, err := client.GetQuota()
	if err != nil {
		return fmt.Errorf("failed to fetch quota: %w", err)
	}
	output.PrintQuotaStatus(quota.DailyLimit, quota.DailyUsed, quota.ResetAt, quota.Unlimited)
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
