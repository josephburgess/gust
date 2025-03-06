package cli

import (
	"fmt"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/josephburgess/gust/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRunMissingCity(t *testing.T) {
	testFunc := func(cli *CLI, cfg *config.Config) error {
		city := determineCityName(cli.City, cli.Args, cfg.DefaultCity)
		if city == "" {
			return handleMissingCity()
		}
		return nil
	}

	testCases := []struct {
		name        string
		cli         *CLI
		defaultCity string
		expectError bool
	}{
		{
			name:        "no city provided",
			cli:         &CLI{},
			defaultCity: "",
			expectError: true,
		},
		{
			name:        "city provided through flag",
			cli:         &CLI{City: "London"},
			defaultCity: "",
			expectError: false,
		},
		{
			name:        "city provided through args",
			cli:         &CLI{Args: []string{"New", "York"}},
			defaultCity: "",
			expectError: false,
		},
		{
			name:        "default city used",
			cli:         &CLI{},
			defaultCity: "Paris",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{DefaultCity: tc.defaultCity}
			err := testFunc(tc.cli, cfg)

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "no city provided")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRunAuthRequired(t *testing.T) {
	testFunc := func(needsAuth bool) error {
		if needsAuth {
			return handleMissingAuth()
		}
		return nil
	}

	testCases := []struct {
		name        string
		needsAuth   bool
		expectError bool
	}{
		{
			name:        "auth required",
			needsAuth:   true,
			expectError: true,
		},
		{
			name:        "auth not required",
			needsAuth:   false,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := testFunc(tc.needsAuth)

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "authentication required")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// testable version of run with deps injected
func createRunWithDeps() func(
	ctx *kong.Context,
	cli *CLI,
	loadConfig func() (*config.Config, error),
	authenticate func(string) (*config.AuthConfig, error),
	saveAuthConfig func(*config.AuthConfig) error,
	loadAuthConfig func() (*config.AuthConfig, error),
	runSetup func(*config.Config, bool) error,
	fetchWeather func(string, *config.Config, *config.AuthConfig, *CLI) error,
) error {
	return func(
		ctx *kong.Context,
		cli *CLI,
		loadConfig func() (*config.Config, error),
		authenticate func(string) (*config.AuthConfig, error),
		saveAuthConfig func(*config.AuthConfig) error,
		loadAuthConfig func() (*config.AuthConfig, error),
		runSetup func(*config.Config, bool) error,
		fetchWeather func(string, *config.Config, *config.AuthConfig, *CLI) error,
	) error {
		cfg, err := loadConfig()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		if cli.Login {
			authConfig, err := authenticate(cfg.ApiUrl)
			if err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}

			if err := saveAuthConfig(authConfig); err != nil {
				return fmt.Errorf("failed to save auth: %w", err)
			}

			fmt.Printf("authed as %s\n", authConfig.GithubUser)
			return nil
		}

		authConfig, _ := loadAuthConfig()
		needsAuth := authConfig == nil

		if needsSetup(cli, cfg) {
			err = runSetup(cfg, needsAuth)
			if err != nil {
				return err
			}

			authConfig, _ = loadAuthConfig()
		}

		if needsAuth {
			return handleMissingAuth()
		}

		city := determineCityName(cli.City, cli.Args, cfg.DefaultCity)
		if city == "" {
			return handleMissingCity()
		}

		return fetchWeather(city, cfg, authConfig, cli)
	}
}

func TestRunLogin(t *testing.T) {
	runWithDeps := createRunWithDeps()
	ctx := &kong.Context{}
	cli := &CLI{Login: true}
	mockAuthConfig := &config.AuthConfig{
		APIKey:     "testkey",
		GithubUser: "testuser",
	}

	t.Run("Successful login", func(t *testing.T) {
		var loadConfigCalled, authenticateCalled, saveAuthCalled bool

		err := runWithDeps(
			ctx,
			cli,
			func() (*config.Config, error) {
				loadConfigCalled = true
				return &config.Config{ApiUrl: "https://api.example.com"}, nil
			},
			func(apiURL string) (*config.AuthConfig, error) {
				authenticateCalled = true
				assert.Equal(t, "https://api.example.com", apiURL)
				return mockAuthConfig, nil
			},
			func(config *config.AuthConfig) error {
				saveAuthCalled = true
				assert.Equal(t, mockAuthConfig, config)
				return nil
			},
			nil, nil, nil,
		)

		assert.NoError(t, err)
		assert.True(t, loadConfigCalled)
		assert.True(t, authenticateCalled)
		assert.True(t, saveAuthCalled)
	})

	t.Run("Authentication failure", func(t *testing.T) {
		err := runWithDeps(
			ctx,
			cli,
			func() (*config.Config, error) {
				return &config.Config{ApiUrl: "https://api.example.com"}, nil
			},
			func(apiURL string) (*config.AuthConfig, error) {
				return nil, fmt.Errorf("authentication failed")
			},
			nil, nil, nil, nil,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "authentication failed")
	})
}
