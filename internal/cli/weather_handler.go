package cli

import (
	"fmt"
	"strings"
	"time"

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
			displayRateLimitError(0, client.RateLimitInfo.Limit, client.RateLimitInfo.ResetTime)

			timeUntilReset := time.Until(client.RateLimitInfo.ResetTime)
			if timeUntilReset > 0 {
				minutesRemaining := int(timeUntilReset.Minutes()) + 1
				hoursRemaining := minutesRemaining / 60

				if hoursRemaining > 0 {
					remainingMinutes := minutesRemaining % 60
					return fmt.Errorf("please try again in about %d hour(s) and %d minute(s) when your rate limit resets",
						hoursRemaining, remainingMinutes)
				} else {
					return fmt.Errorf("please try again in about %d minute(s) when your rate limit resets",
						minutesRemaining)
				}
			}
			return fmt.Errorf("rate limit reached, please try again later")
		}

		if client.RateLimitInfo.Remaining <= 0 {
			displayRateLimitError(0, client.RateLimitInfo.Limit, client.RateLimitInfo.ResetTime)
		} else if client.RateLimitInfo.Remaining <= 5 {
			displayRateLimitWarning(
				client.RateLimitInfo.Remaining,
				client.RateLimitInfo.Limit,
				client.RateLimitInfo.ResetTime,
			)
		}

		defer func() {
			if client.RateLimitInfo.Remaining > 5 {
				fmt.Println()
				displayRateLimitStatus(client.RateLimitInfo.Remaining, client.RateLimitInfo.Limit)
			}
		}()
	}

	if err != nil {
		return err
	}

	weatherRenderer := renderer.NewWeatherRenderer("terminal", cfg.Units)
	renderWeatherView(cli, weatherRenderer, weather.City, weather.Weather, cfg)

	return nil
}

func displayRateLimitWarning(remaining, limit int, resetTime time.Time) {
	timeUntilReset := time.Until(resetTime)
	minutesUntilReset := int(timeUntilReset.Minutes())
	resetFormatted := resetTime.Format("15:04")

	fmt.Println()
	fmt.Println(styles.BoxStyle.Render(fmt.Sprintf(
		"%s API Rate Limit Warning\n\n"+
			"You have %s requests remaining out of %d.\n"+
			"Your rate limit will reset at %s (%d minutes from now).",
		styles.WarningStyle("⚠️"),
		styles.HighlightStyleF(fmt.Sprintf("%d", remaining)),
		limit,
		styles.TimeStyle(resetFormatted),
		minutesUntilReset,
	)))
	fmt.Println()
}

func displayRateLimitError(remaining, limit int, resetTime time.Time) {
	timeUntilReset := time.Until(resetTime)
	minutesUntilReset := int(timeUntilReset.Minutes())
	resetFormatted := resetTime.Format("15:04")

	fmt.Println()
	fmt.Println(styles.BoxStyle.BorderForeground(styles.Love).Render(fmt.Sprintf(
		"%s API Rate Limit Reached\n\n"+
			"You have used all %d available requests.\n"+
			"Your rate limit will reset at %s (%d minutes from now).\n\n"+
			"%s To get more data, please wait until the limit resets.",
		styles.ErrorStyle("❌"),
		limit,
		styles.TimeStyle(resetFormatted),
		minutesUntilReset,
		styles.InfoStyle("💡"),
	)))
	fmt.Println()
}

func displayRateLimitStatus(remaining, limit int) {
	if limit <= 0 {
		return
	}

	const barWidth = 20
	used := limit - remaining

	filledCount := min(int(float64(used)/float64(limit)*barWidth), barWidth)

	emptyCount := barWidth - filledCount

	filled := styles.HighlightStyleF(strings.Repeat("█", filledCount))
	empty := strings.Repeat("░", emptyCount)

	percentage := float64(used) / float64(limit) * 100

	var usageText string
	if percentage >= 90 {
		usageText = styles.ErrorStyle(fmt.Sprintf("%.0f%% used", percentage))
	} else if percentage >= 75 {
		usageText = styles.WarningStyle(fmt.Sprintf("%.0f%% used", percentage))
	} else {
		usageText = styles.InfoStyle(fmt.Sprintf("%.0f%% used", percentage))
	}

	fmt.Printf("API Usage: [%s%s] %s (%d/%d)\n", filled, empty, usageText, used, limit)
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
