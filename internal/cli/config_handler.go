package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/josephburgess/gust/internal/config"
)

func handleConfigUpdates(cli *CLI, cfg *config.Config) (bool, error) {
	updated := false

	if cli.ApiUrl != "" {
		cfg.ApiUrl = cli.ApiUrl
		updated = true
	}

	if cli.ApiKey != "" {
		authConfig, _ := config.LoadAuthConfig()

		newAuthConfig := &config.AuthConfig{
			APIKey:     cli.ApiKey,
			ServerURL:  cfg.ApiUrl,
			LastAuth:   time.Now(),
			GithubUser: "",
		}

		if authConfig != nil {
			newAuthConfig.GithubUser = authConfig.GithubUser
		}

		if err := config.SaveAuthConfig(newAuthConfig); err != nil {
			return false, fmt.Errorf("failed to save API key: %w", err)
		}

		fmt.Println("API key updated.")
		updated = true
	}

	if cli.Units != "" {
		if !isValidUnit(cli.Units) {
			fmt.Println("Invalid units. Must be one of: metric, imperial, standard")
			os.Exit(1)
		}
		cfg.Units = cli.Units
		updated = true
	}

	if cli.Default != "" {
		cfg.DefaultCity = cli.Default
		updated = true
	}

	if updated {
		if err := cfg.Save(); err != nil {
			return false, fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Println("Configuration updated.")
	}

	return updated, nil
}

func isValidUnit(unit string) bool {
	validUnits := map[string]bool{
		"metric":   true,
		"imperial": true,
		"standard": true,
	}

	return validUnits[unit]
}
