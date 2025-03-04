package cli

import (
	"fmt"

	"github.com/josephburgess/gust/internal/config"
)

func handleLogin(apiURL string) error {
	fmt.Println("Starting GitHub authentication...")
	authConfig, err := config.Authenticate(apiURL)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if err := config.SaveAuthConfig(authConfig); err != nil {
		return fmt.Errorf("failed to save authentication: %w", err)
	}

	fmt.Printf("Successfully authenticated as %s\n", authConfig.GithubUser)
	return nil
}

func handleMissingAuth() error {
	fmt.Println("You need to authenticate with GitHub before using Gust.")
	fmt.Println("Run 'gust --login' to authenticate or 'gust --setup' to run the setup wizard.")
	return fmt.Errorf("authentication required")
}
