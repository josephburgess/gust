package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/ui/styles"
)

func Run(ctx *kong.Context, cli *CLI) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if updated, err := handleConfigUpdates(cli, cfg); updated || err != nil {
		return err
	}

	if cli.Login {
		return handleLogin(cfg.ApiUrl)
	}

	authConfig, err := config.LoadAuthConfig()
	if err != nil {
		return fmt.Errorf("failed to reload configuration after setup: %w", err)
	}

	needsAuth := authConfig == nil

	if needsSetup(cli, cfg) {
		styles.PrintInfo("Defaults not set, running setup...")
		if err := handleSetup(cfg, &needsAuth); err != nil {
			return err
		}
	}

	if needsAuth {
		return handleMissingAuth()
	}

	city := determineCityName(cli.City, cli.Args, cfg.DefaultCity)
	if city == "" {
		return handleMissingCity()
	}

	return fetchAndRenderWeather(city, cfg, authConfig, cli)
}
