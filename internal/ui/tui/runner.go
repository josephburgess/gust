package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
)

// Run opens the dashboard for city on the given tab.
func Run(city string, cfg *config.Config, auth *config.AuthConfig, tab Tab) error {
	fetch := func(city, units string) (*api.WeatherResponse, error) {
		return api.NewClient(cfg.ApiUrl, auth.APIKey, units).GetWeather(city)
	}

	model := NewModel(city, cfg.Units, cfg.ShowTips, tab, fetch)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		return fmt.Errorf("error running dashboard: %w", err)
	}
	return nil
}
