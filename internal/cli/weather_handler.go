package cli

import (
	"errors"
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
			weatherRenderer := renderer.NewWeatherRenderer(cfg.Units)
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
			if isRateLimitError(err) {
				return nil, fmt.Errorf("rate limit reached: %w", err)
			}
			return nil, err
		}
		return weather, nil
	}

	message := fmt.Sprintf("Fetching weather for %s...", city)
	weather, err := components.RunWithSpinner(message, components.WeatherEmojis, styles.Foam, fetchFunc)

	if err != nil {
		// city not found — show suggestions
		var notFound *api.CityNotFoundError
		if errors.As(err, &notFound) {
			return handleCityNotFound(client, city)
		}

		// network/server error — fall back to stale cache if available
		if weatherCache != nil && isNetworkError(err) {
			if stale, age, ok := weatherCache.GetStale(city, cfg.Units); ok {
				output.PrintStaleWarning(age)
				weatherRenderer := renderer.NewWeatherRenderer(cfg.Units)
				renderWeatherView(cli, weatherRenderer, stale.City, stale.Weather, cfg)
				return nil
			}
		}

		// rate limit — show friendly message
		if client.RateLimitInfo != nil && client.RateLimitInfo.Limit > 0 && isRateLimitError(err) {
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

		return err
	}

	if client.RateLimitInfo != nil && client.RateLimitInfo.Remaining <= 5 && client.RateLimitInfo.Remaining > 0 {
		output.PrintRateLimitWarning(
			client.RateLimitInfo.Remaining,
			client.RateLimitInfo.Limit,
			client.RateLimitInfo.ResetTime,
		)
	}

	if weatherCache != nil {
		weatherCache.Set(city, cfg.Units, weather)
	}

	weatherRenderer := renderer.NewWeatherRenderer(cfg.Units)
	renderWeatherView(cli, weatherRenderer, weather.City, weather.Weather, cfg)

	return nil
}

func handleCityNotFound(client *api.Client, city string) error {
	suggestions, err := client.SearchCities(city)
	if err != nil || len(suggestions) == 0 {
		output.PrintCityNotFound(city, nil)
		return fmt.Errorf("city %q not found", city)
	}

	formatted := make([]string, 0, len(suggestions))
	for _, s := range suggestions {
		var parts []string
		if s.Name != "" {
			parts = append(parts, s.Name)
		}
		if s.State != "" {
			parts = append(parts, s.State)
		}
		if s.Country != "" {
			parts = append(parts, s.Country)
		}
		label := strings.Join(parts, ", ")
		if flag := countryFlag(s.Country); flag != "" {
			label += " " + flag
		}
		formatted = append(formatted, label)
	}

	output.PrintCityNotFound(city, formatted)
	return fmt.Errorf("city %q not found", city)
}

func countryFlag(code string) string {
	if len(code) != 2 {
		return ""
	}
	code = strings.ToUpper(code)
	const offset = 127397
	return string(rune(code[0])+offset) + string(rune(code[1])+offset)
}

func isRateLimitError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "rate limit")
}

func isNetworkError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "failed to connect") ||
		strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "deadline exceeded") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "no such host")
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
