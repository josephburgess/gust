package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/ui/components"
	"github.com/josephburgess/gust/internal/ui/output"
	"github.com/josephburgess/gust/internal/ui/styles"
	"github.com/josephburgess/gust/internal/ui/tui"
)

// printWidth is the layout width when stdout is not a terminal.
const printWidth = 100

type runMode int

const (
	modeDashboard runMode = iota
	modePrint
	modeCompact
)

// chooseMode picks the dashboard, print or compact, and the tab.
// A tab flag wins over opens-on. Print is the default when stdout is not a terminal.
func chooseMode(cli *CLI, cfg *config.Config, tty bool) (runMode, tui.Tab) {
	tab, tabFlag := tui.TabFromView(cfg.DefaultView), true
	switch {
	case cli.Alerts:
		tab = tui.TabAlerts
	case cli.Hourly:
		tab = tui.TabHourly
	case cli.Daily:
		tab = tui.TabDaily
	case cli.Now:
		tab = tui.TabNow
	default:
		tabFlag = false
	}

	switch {
	case cli.Compact, !tabFlag && cfg.DefaultView == "compact":
		return modeCompact, tab
	case cli.Print, !tty:
		return modePrint, tab
	}
	return modeDashboard, tab
}

func fetchAndRenderWeather(city string, cfg *config.Config, authConfig *config.AuthConfig, cli *CLI) error {
	tty := term.IsTerminal(os.Stdout.Fd())
	mode, tab := chooseMode(cli, cfg, tty)
	if mode == modeDashboard {
		return tui.Run(city, cfg, authConfig, tab)
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

	// the spinner draws on stdout, so it only runs on a terminal
	var weather *api.WeatherResponse
	var err error
	if tty {
		weather, err = components.RunWithSpinner(fmt.Sprintf("fetching %s…", city), spinner.Dot, styles.Gold, fetchFunc)
	} else {
		weather, err = fetchFunc()
	}

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
				} else {
					return fmt.Errorf("please try again in about %d minute(s) when your rate limit resets",
						minutesRemaining)
				}
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

	if mode == modeCompact {
		lipgloss.Println(tui.Compact(weather.City, weather.Weather, cfg.Units))
		return nil
	}

	width := printWidth
	if w, _, err := term.GetSize(os.Stdout.Fd()); tty && err == nil {
		width = w
	}
	lipgloss.Println(tui.Print(weather.City, weather.Weather, cfg, tab, width))
	return nil
}
