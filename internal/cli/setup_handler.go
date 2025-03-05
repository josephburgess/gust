package cli

import (
	"fmt"

	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/ui/output"
	"github.com/josephburgess/gust/internal/ui/setup"
)

func needsSetup(cli *CLI, cfg *config.Config) bool {
	return cfg.DefaultCity == "" || cli.Setup
}

func handleSetup(cfg *config.Config) (bool, error) {
	output.PrintInfo("Running setup wizard...")
	authConfig, _ := config.LoadAuthConfig()
	needsAuth := authConfig == nil

	if err := setup.RunSetup(cfg, needsAuth); err != nil {
		return needsAuth, fmt.Errorf("setup failed: %w", err)
	}

	newCfg, err := config.Load()
	if err != nil {
		return needsAuth, fmt.Errorf("failed to reload configuration after setup: %w", err)
	}

	cfg.DefaultCity = newCfg.DefaultCity

	authConfig, err = config.LoadAuthConfig()
	if err != nil {
		return true, fmt.Errorf("failed to load auth config after setup: %w", err)
	}

	needsAuth = authConfig == nil
	output.PrintSuccess("Setup complete! Run 'gust' to check the weather for your default city.")

	return needsAuth, nil
}
