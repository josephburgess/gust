package setup

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
)

// entry point for setup wizard
func RunSetup(cfg *config.Config, needsAuth bool) error {
	var apiClient *api.Client

	if cfg.ApiUrl == "" {
		cfg.ApiUrl = "https://breeze.joeburgess.dev"
	}

	authConfig, err := config.LoadAuthConfig()
	if err == nil && authConfig != nil {
		apiClient = api.NewClient(cfg.ApiUrl, authConfig.APIKey, cfg.Units)
	} else {
		// empty api key - dont need for setup
		apiClient = api.NewClient(cfg.ApiUrl, "", cfg.Units)
	}

	model := NewModel(cfg, needsAuth, apiClient)
	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running setup UI: %w", err)
	}

	// protection against nil model
	if finalModel == nil {
		return fmt.Errorf("unexpected nil model after running UI")
	}

	// get the final model
	finalSetupModel, ok := finalModel.(Model)
	if !ok {
		return fmt.Errorf("unexpected model type: %T", finalModel)
	}

	// dont save if quit with ctrl+c
	if finalSetupModel.Quitting {
		return nil
	}

	// save config if not saved earlier
	if err := finalSetupModel.Config.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	// handle auth if chosen
	if finalSetupModel.State == StateAuth && finalSetupModel.AuthCursor == 0 {
		fmt.Println("Starting GitHub authentication...")
		auth, err := config.Authenticate(cfg.ApiUrl)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
		if err := config.SaveAuthConfig(auth); err != nil {
			return fmt.Errorf("failed to save authentication: %w", err)
		}
		fmt.Printf("Successfully authenticated as %s\n", auth.GithubUser)
	}

	return nil
}
