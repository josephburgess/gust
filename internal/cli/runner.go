package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/ui/output"
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
		return fmt.Errorf("failed to load auth config: %w", err)
	}
	needsAuth := authConfig == nil

	if needsSetup(cli, cfg) {
		output.PrintInfo("Defaults not set, running setup...")
		needsAuth, err = handleSetup(cfg)
		if err != nil {
			return err
		}

		authConfig, err = config.LoadAuthConfig()
		if err != nil {
			return fmt.Errorf("failed to load auth config: %w", err)
		}
	}

	if needsAuth {
		return handleMissingAuth()
	}

	if cli.Status {
		return handleStatus(cfg, authConfig)
	}

	city := determineCityName(cli.City, cli.Args, cfg.DefaultCity)
	if city == "" {
		return handleMissingCity()
	}

	return fetchAndRenderWeather(city, cfg, authConfig, cli)
}
