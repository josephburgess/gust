package cli

import (
	"fmt"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWeatherRenderer struct {
	mock.Mock
}

func (m *MockWeatherRenderer) RenderCurrentWeather(city *models.City, weather *models.OneCallResponse) {
	m.Called(city, weather)
}

func (m *MockWeatherRenderer) RenderDailyForecast(city *models.City, weather *models.OneCallResponse) {
	m.Called(city, weather)
}

func (m *MockWeatherRenderer) RenderHourlyForecast(city *models.City, weather *models.OneCallResponse) {
	m.Called(city, weather)
}

func (m *MockWeatherRenderer) RenderAlerts(city *models.City, weather *models.OneCallResponse) {
	m.Called(city, weather)
}

func (m *MockWeatherRenderer) RenderFullWeather(city *models.City, weather *models.OneCallResponse) {
	m.Called(city, weather)
}

func (m *MockWeatherRenderer) RenderCompactWeather(city *models.City, weather *models.OneCallResponse) {
	m.Called(city, weather)
}

func TestRenderWeatherView_CLIFlags(t *testing.T) {
	mockCity := &models.City{Name: "TestCity"}
	mockWeather := &models.OneCallResponse{}
	mockConfig := &config.Config{DefaultView: "compact"}

	testCases := []struct {
		name           string
		cli            *CLI
		expectedMethod string
	}{
		{
			name:           "Alerts flag",
			cli:            &CLI{Alerts: true},
			expectedMethod: "RenderAlerts",
		},
		{
			name:           "Hourly flag",
			cli:            &CLI{Hourly: true},
			expectedMethod: "RenderHourlyForecast",
		},
		{
			name:           "Daily flag",
			cli:            &CLI{Daily: true},
			expectedMethod: "RenderDailyForecast",
		},
		{
			name:           "Full flag",
			cli:            &CLI{Full: true},
			expectedMethod: "RenderFullWeather",
		},
		{
			name:           "Compact flag",
			cli:            &CLI{Compact: true},
			expectedMethod: "RenderCompactWeather",
		},
		{
			name:           "Detailed flag",
			cli:            &CLI{Detailed: true},
			expectedMethod: "RenderCurrentWeather",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRenderer := new(MockWeatherRenderer)
			mockRenderer.On(tc.expectedMethod, mockCity, mockWeather).Return()

			renderWeatherView(tc.cli, mockRenderer, mockCity, mockWeather, mockConfig)

			mockRenderer.AssertExpectations(t)
		})
	}
}

func TestRenderWeatherView_DefaultView(t *testing.T) {
	mockCity := &models.City{Name: "TestCity"}
	mockWeather := &models.OneCallResponse{}
	cli := &CLI{}

	testCases := []struct {
		name        string
		defaultView string
		expectedFn  string
	}{
		{
			name:        "Default view is compact",
			defaultView: "compact",
			expectedFn:  "RenderCompactWeather",
		},
		{
			name:        "Default view is daily",
			defaultView: "daily",
			expectedFn:  "RenderDailyForecast",
		},
		{
			name:        "Default view is hourly",
			defaultView: "hourly",
			expectedFn:  "RenderHourlyForecast",
		},
		{
			name:        "Default view is full",
			defaultView: "full",
			expectedFn:  "RenderFullWeather",
		},
		{
			name:        "Default view is unknown",
			defaultView: "unknown",
			expectedFn:  "RenderCurrentWeather",
		},
		{
			name:        "Default view is empty",
			defaultView: "",
			expectedFn:  "RenderCurrentWeather",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRenderer := new(MockWeatherRenderer)
			mockRenderer.On(tc.expectedFn, mockCity, mockWeather).Return()

			config := &config.Config{DefaultView: tc.defaultView}
			renderWeatherView(cli, mockRenderer, mockCity, mockWeather, config)

			mockRenderer.AssertExpectations(t)
		})
	}
}

func TestRenderWeatherView_FlagsAndPriority(t *testing.T) {
	mockCity := &models.City{Name: "TestCity"}
	mockWeather := &models.OneCallResponse{}
	mockConfig := &config.Config{DefaultView: "compact"}

	testCases := []struct {
		name           string
		cli            *CLI
		expectedMethod string
	}{
		{
			name:           "Default falls back to config",
			cli:            &CLI{},
			expectedMethod: "RenderCompactWeather",
		},
		{
			name:           "Flag overrides default",
			cli:            &CLI{Daily: true},
			expectedMethod: "RenderDailyForecast",
		},
		{
			name: "Multiple flags respect priority",
			cli: &CLI{
				Alerts:   true,
				Hourly:   true,
				Daily:    true,
				Full:     true,
				Compact:  true,
				Detailed: true,
			},
			expectedMethod: "RenderAlerts",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRenderer := new(MockWeatherRenderer)
			mockRenderer.On(tc.expectedMethod, mockCity, mockWeather).Return()

			renderWeatherView(tc.cli, mockRenderer, mockCity, mockWeather, mockConfig)

			mockRenderer.AssertExpectations(t)
		})
	}
}

func TestDetermineCityName(t *testing.T) {
	testCases := []struct {
		name        string
		cityFlag    string
		args        []string
		defaultCity string
		expected    string
	}{
		{
			name:        "Use city flag when provided",
			cityFlag:    "London",
			args:        []string{"Paris"},
			defaultCity: "Berlin",
			expected:    "London",
		},
		{
			name:        "Use args when no flag but args provided",
			cityFlag:    "",
			args:        []string{"New", "York"},
			defaultCity: "Berlin",
			expected:    "New York",
		},
		{
			name:        "Use default city when no flag or args",
			cityFlag:    "",
			args:        []string{},
			defaultCity: "Berlin",
			expected:    "Berlin",
		},
		{
			name:        "Return empty when no values provided",
			cityFlag:    "",
			args:        []string{},
			defaultCity: "",
			expected:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := determineCityName(tc.cityFlag, tc.args, tc.defaultCity)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsValidUnit(t *testing.T) {
	validUnits := []string{"metric", "imperial", "standard"}
	invalidUnits := []string{"celsius", "", "Metric"}

	for _, unit := range validUnits {
		t.Run("Valid: "+unit, func(t *testing.T) {
			assert.True(t, isValidUnit(unit))
		})
	}

	for _, unit := range invalidUnits {
		t.Run("Invalid: "+unit, func(t *testing.T) {
			assert.False(t, isValidUnit(unit))
		})
	}
}

func TestHandleConfigUpdates(t *testing.T) {
	createConfig := func() *config.Config {
		return &config.Config{
			ApiUrl:      "https://api.example.com",
			Units:       "metric",
			DefaultCity: "London",
		}
	}

	testCases := []struct {
		name            string
		cli             *CLI
		expectedUpdated bool
		configMutator   func(*config.Config)
	}{
		{
			name: "Update API URL",
			cli: &CLI{
				ApiUrl: "https://test-api.example.com",
			},
			expectedUpdated: true,
			configMutator: func(c *config.Config) {
				c.ApiUrl = "https://test-api.example.com"
			},
		},
		{
			name: "Update Units",
			cli: &CLI{
				Units: "imperial",
			},
			expectedUpdated: true,
			configMutator: func(c *config.Config) {
				c.Units = "imperial"
			},
		},
		{
			name: "Update Default City",
			cli: &CLI{
				Default: "Paris",
			},
			expectedUpdated: true,
			configMutator: func(c *config.Config) {
				c.DefaultCity = "Paris"
			},
		},
		{
			name:            "No Updates",
			cli:             &CLI{},
			expectedUpdated: false,
			configMutator:   func(c *config.Config) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initialConfig := createConfig()
			expectedConfig := createConfig()
			tc.configMutator(expectedConfig)

			testHandleConfigUpdates := func(cli *CLI, cfg *config.Config) (bool, error) {
				updated := false

				if cli.ApiUrl != "" {
					cfg.ApiUrl = cli.ApiUrl
					updated = true
				}

				if cli.Units != "" {
					if !isValidUnit(cli.Units) {
						return false, fmt.Errorf("invalid units")
					}
					cfg.Units = cli.Units
					updated = true
				}

				if cli.Default != "" {
					cfg.DefaultCity = cli.Default
					updated = true
				}

				return updated, nil
			}

			updated, err := testHandleConfigUpdates(tc.cli, initialConfig)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedUpdated, updated)
			assert.Equal(t, expectedConfig, initialConfig)
		})
	}
}

func TestNeedsSetup(t *testing.T) {
	testCases := []struct {
		name     string
		cli      *CLI
		cfg      *config.Config
		expected bool
	}{
		{
			name:     "Empty default city",
			cli:      &CLI{},
			cfg:      &config.Config{DefaultCity: ""},
			expected: true,
		},
		{
			name:     "Setup flag set",
			cli:      &CLI{Setup: true},
			cfg:      &config.Config{DefaultCity: "London"},
			expected: true,
		},
		{
			name:     "Default city set and no setup flag",
			cli:      &CLI{},
			cfg:      &config.Config{DefaultCity: "London"},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := needsSetup(tc.cli, tc.cfg)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNewApp(t *testing.T) {
	app, cli := NewApp()

	assert.NotNil(t, app)
	assert.NotNil(t, cli)

	fields := []struct {
		name  string
		value any
	}{
		{"City", cli.City},
		{"Default", cli.Default},
		{"ApiUrl", cli.ApiUrl},
		{"Login", cli.Login},
		{"Setup", cli.Setup},
		{"Compact", cli.Compact},
		{"Detailed", cli.Detailed},
		{"Full", cli.Full},
		{"Daily", cli.Daily},
		{"Hourly", cli.Hourly},
		{"Alerts", cli.Alerts},
		{"Units", cli.Units},
		{"Pretty", cli.Pretty},
	}

	for _, field := range fields {
		switch v := field.value.(type) {
		case string:
			assert.Empty(t, v, "Field %s should be empty", field.name)
		case bool:
			assert.False(t, v, "Field %s should be false", field.name)
		}
	}

	assert.Empty(t, cli.Args)
}

func TestErrorHandlers(t *testing.T) {
	t.Run("Missing city error", func(t *testing.T) {
		err := handleMissingCity()
		assert.EqualError(t, err, "no city provided")
	})

	t.Run("Missing auth error", func(t *testing.T) {
		err := handleMissingAuth()
		assert.EqualError(t, err, "authentication required")
	})
}

type MockAuthConfig struct {
	mock.Mock
	APIKey     string
	GithubUser string
}

func (m *MockAuthConfig) AuthFunc(apiURL string) (*config.AuthConfig, error) {
	args := m.Called(apiURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*config.AuthConfig), args.Error(1)
}

func (m *MockAuthConfig) SaveAuthFunc(config *config.AuthConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockAuthConfig) LoadAuthFunc() (*config.AuthConfig, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*config.AuthConfig), args.Error(1)
}

type MockSetup struct {
	mock.Mock
}

func (m *MockSetup) RunSetupFunc(cfg *config.Config, needsAuth bool) error {
	args := m.Called(cfg, needsAuth)
	return args.Error(0)
}

type MockWeatherClient struct {
	mock.Mock
}

func (m *MockWeatherClient) GetWeather(cityName string) (*api.WeatherResponse, error) {
	args := m.Called(cityName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.WeatherResponse), args.Error(1)
}

func (m *MockWeatherClient) SearchCities(query string) ([]models.City, error) {
	args := m.Called(query)
	cities, _ := args.Get(0).([]models.City)
	return cities, args.Error(1)
}

func TestRunLogin(t *testing.T) {
	runWithDeps := createRunWithDeps()
	ctx := &kong.Context{}
	cli := &CLI{Login: true}
	mockAuthConfig := &config.AuthConfig{
		APIKey:     "test-api-key",
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
			name:        "No city provided",
			cli:         &CLI{},
			defaultCity: "",
			expectError: true,
		},
		{
			name:        "City provided through flag",
			cli:         &CLI{City: "London"},
			defaultCity: "",
			expectError: false,
		},
		{
			name:        "City provided through args",
			cli:         &CLI{Args: []string{"New", "York"}},
			defaultCity: "",
			expectError: false,
		},
		{
			name:        "Default city used",
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
			name:        "Auth required",
			needsAuth:   true,
			expectError: true,
		},
		{
			name:        "Auth not required",
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
				return fmt.Errorf("failed to save authentication: %w", err)
			}

			fmt.Printf("Successfully authenticated as %s\n", authConfig.GithubUser)
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
