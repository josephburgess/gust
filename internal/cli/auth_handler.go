package cli

import (
	"fmt"

	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/ui/output"
)

func handleLogin(apiURL string) error {
	output.PrintInfo("Starting GitHub authentication...")
	authConfig, err := config.Authenticate(apiURL)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if err := config.SaveAuthConfig(authConfig); err != nil {
		return fmt.Errorf("failed to save authentication: %w", err)
	}

	output.PrintSuccess(fmt.Sprintf("Successfully authenticated as %s\n", authConfig.GithubUser))
	return nil
}

func handleMissingAuth() error {
	output.PrintError("You need to authenticate with GitHub before using Gust.")
	output.PrintInfo("Run 'gust --login' to authenticate or 'gust --setup' to run the setup wizard.")
	return fmt.Errorf("authentication required")
}
