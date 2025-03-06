package setup

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewModel(t *testing.T) {
	cfg := &config.Config{
		Units:       "metric",
		DefaultView: "default",
	}
	client := &api.Client{}

	tests := []struct {
		name      string
		cfg       *config.Config
		needsAuth bool
		want      Model
	}{
		{
			name:      "initializes with default values",
			cfg:       cfg,
			needsAuth: false,
			want: Model{
				Config:    cfg,
				State:     StateCity,
				Client:    client,
				NeedsAuth: false,
				Quitting:  false,
			},
		},
		{
			name:      "initializes with auth required",
			cfg:       cfg,
			needsAuth: true,
			want: Model{
				Config:    cfg,
				State:     StateCity,
				Client:    client,
				NeedsAuth: true,
				Quitting:  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewModel(tt.cfg, tt.needsAuth, client)

			// Check basic properties
			assert.Equal(t, tt.want.Config, got.Config)
			assert.Equal(t, tt.want.State, got.State)
			assert.Equal(t, tt.want.NeedsAuth, got.NeedsAuth)
			assert.Equal(t, tt.want.Quitting, got.Quitting)

			// Check initialized components
			assert.NotNil(t, got.CityInput)
			assert.Len(t, got.UnitOptions, 3)
			assert.Len(t, got.ViewOptions, 5)
			assert.Len(t, got.TipOptions, 2)
			assert.Len(t, got.AuthOptions, 2)

			// Verify textinput is properly configured
			assert.True(t, got.CityInput.Focused())
			assert.Equal(t, "Wherever the wind blows...", got.CityInput.Placeholder)
		})
	}
}

func TestNewModelWithDifferentConfigs(t *testing.T) {
	tests := []struct {
		name       string
		units      string
		view       string
		unitCursor int
		viewCursor int
		needsAuth  bool
	}{
		{"metric", "metric", "default", 0, 0, false},
		{"imperial", "imperial", "compact", 1, 1, false},
		{"standard", "standard", "daily", 2, 2, false},
		{"unknown units", "unknown", "hourly", 0, 3, true},
		{"unknown view", "metric", "unknown", 0, 0, true},
		{"full view", "imperial", "full", 1, 4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Units:       tt.units,
				DefaultView: tt.view,
			}
			model := NewModel(cfg, tt.needsAuth, nil)

			assert.Equal(t, tt.unitCursor, model.UnitCursor, "Unit cursor should match expected value")
			assert.Equal(t, tt.viewCursor, model.ViewCursor, "View cursor should match expected value")
			assert.Equal(t, tt.needsAuth, model.NeedsAuth, "Auth flag should match expected value")
		})
	}
}

func TestModelInit(t *testing.T) {
	model := NewModel(&config.Config{}, false, nil)
	cmd := model.Init()

	// We can't directly test the returned commands, but we can ensure one is returned
	assert.NotNil(t, cmd, "Init should return a command")

	// Try initializing with different configurations
	configs := []*config.Config{
		{Units: "metric", DefaultView: "default"},
		{Units: "imperial", DefaultView: "compact"},
		{DefaultCity: "London"},
	}

	for _, cfg := range configs {
		model = NewModel(cfg, true, nil)
		cmd = model.Init()
		assert.NotNil(t, cmd, "Init should return a command for all configurations")
	}
}

func TestModelMessageTypes(t *testing.T) {
	// Test message types have the correct structure
	authMsg := AuthenticateMsg{}
	setupCompleteMsg := SetupCompleteMsg{}

	// These are empty structs, so we can only verify they exist
	assert.IsType(t, AuthenticateMsg{}, authMsg)
	assert.IsType(t, SetupCompleteMsg{}, setupCompleteMsg)
}

func TestWindowSizeHandling(t *testing.T) {
	model := NewModel(&config.Config{}, false, nil)

	// Initial size should be 0,0
	assert.Equal(t, 0, model.Width)
	assert.Equal(t, 0, model.Height)

	// Test various window sizes
	sizes := []struct{ width, height int }{
		{80, 24},  // Standard terminal
		{120, 40}, // Large terminal
		{40, 10},  // Small terminal
		{0, 0},    // Zero dimensions
	}

	for _, size := range sizes {
		updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: size.width, Height: size.height})
		updated := updatedModel.(Model)
		assert.Equal(t, size.width, updated.Width)
		assert.Equal(t, size.height, updated.Height)
	}
}
