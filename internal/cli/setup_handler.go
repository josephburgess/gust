package cli

import (
	"fmt"

	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/ui/setup"
)

func needsSetup(cli *CLI, cfg *config.Config) bool {
	return cfg.DefaultCity == "" || cli.Setup
}

func handleSetup(cfg *config.Config, needsAuth *bool) error {
	fmt.Println("Running setup wizard...")

	if needsAuth == nil {
		localNeedsAuth := true
		needsAuth = &localNeedsAuth
	}

	if err := setup.RunSetup(cfg, *needsAuth); err != nil {
		return fmt.Errorf("setup failed: %w", err)
	}

	newCfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to reload configuration after setup: %w", err)
	}

	cfg.DefaultCity = newCfg.DefaultCity

	authConfig, err := config.LoadAuthConfig()
	if err != nil {
		return fmt.Errorf("failed to load auth config after setup: %w", err)
	}
	*needsAuth = authConfig == nil

	fmt.Println("Setup complete! Run 'gust' to check the weather for your default city.")

	return nil
}
