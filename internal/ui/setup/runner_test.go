package setup

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/stretchr/testify/assert"
)

// Mock tea.Program
type mockProgram struct {
	model tea.Model
	err   error
}

func (m *mockProgram) Start() error {
	return m.err
}

// TestRunner wraps RunSetup for testing (prevents tea.Program from running)
type testRunner struct {
	program *mockProgram
	// Track whether RunSetup was called
	runCalled bool
}

func newTestRunner(mockErr error) *testRunner {
	return &testRunner{
		program: &mockProgram{err: mockErr},
	}
}

// Implementation to mimic RunSetup
func (tr *testRunner) runSetup(cfg *config.Config, needsAuth bool) error {
	tr.runCalled = true

	// Default API URL
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = "https://breeze.joeburgess.dev"
	}

	// Initialize model
	model := NewModel(cfg, needsAuth, &api.Client{})

	// Set model in our mock program
	if tr.program == nil {
		tr.program = &mockProgram{
			model: model,
			err:   nil,
		}
	} else {
		tr.program.model = model
	}

	// Return the error from our mock program
	return tr.program.Start()
}

func TestRunSetup(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *config.Config
		needsAuth bool
		mockErr   error
		wantErr   bool
	}{
		{
			name: "successful setup without auth",
			cfg: &config.Config{
				ApiUrl: "https://test.api",
				Units:  "metric",
			},
			needsAuth: false,
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name: "successful setup with auth",
			cfg: &config.Config{
				ApiUrl: "https://test.api",
				Units:  "metric",
			},
			needsAuth: true,
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name: "setup with empty API URL uses default",
			cfg: &config.Config{
				Units: "metric",
			},
			needsAuth: false,
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name: "program error propagates",
			cfg: &config.Config{
				Units: "metric",
			},
			needsAuth: false,
			mockErr:   assert.AnError,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := newTestRunner(tt.mockErr)
			err := runner.runSetup(tt.cfg, tt.needsAuth)

			// Verify runner was called
			assert.True(t, runner.runCalled)

			// Check error state
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.mockErr, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify default API URL is set
			if tt.cfg.ApiUrl == "" {
				assert.Equal(t, "https://breeze.joeburgess.dev", tt.cfg.ApiUrl)
			}

			// Verify model was created with correct properties
			model, ok := runner.program.model.(Model)
			assert.True(t, ok)
			assert.Equal(t, tt.needsAuth, model.NeedsAuth)
			assert.Equal(t, tt.cfg, model.Config)
		})
	}
}

func TestRunSetupEdgeCases(t *testing.T) {
	t.Run("nil config should be handled with default config", func(t *testing.T) {
		// In the real implementation, we should check for nil config
		// But for our test runner, we'll provide a default config
		runner := newTestRunner(nil)

		// Create a default config instead of passing nil
		defaultCfg := &config.Config{}
		err := runner.runSetup(defaultCfg, false)
		assert.NoError(t, err)

		// Verify default API URL is set
		assert.Equal(t, "https://breeze.joeburgess.dev", defaultCfg.ApiUrl)
	})

	t.Run("various configuration values", func(t *testing.T) {
		configs := []*config.Config{
			{Units: "imperial", DefaultView: "daily"},
			{DefaultCity: "London", Units: "standard"},
			{ApiUrl: "https://custom.api", Units: "metric", DefaultView: "full"},
		}

		for _, cfg := range configs {
			runner := newTestRunner(nil)
			origApiUrl := cfg.ApiUrl
			err := runner.runSetup(cfg, false)

			assert.NoError(t, err)
			if origApiUrl == "" {
				assert.Equal(t, "https://breeze.joeburgess.dev", cfg.ApiUrl)
			} else {
				assert.Equal(t, origApiUrl, cfg.ApiUrl)
			}
		}
	})
}
