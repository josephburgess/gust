package setup

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
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
				State:    StateCity,
				Client:   client,
				NeedsAuth: false,
				Quitting: false,
			},
		},
		{
			name:      "initializes with auth required",
			cfg:       cfg,
			needsAuth: true,
			want: Model{
				Config:    cfg,
				State:    StateCity,
				Client:   client,
				NeedsAuth: true,
				Quitting: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewModel(tt.cfg, tt.needsAuth, client)
			assert.Equal(t, tt.want.Config, got.Config)
			assert.Equal(t, tt.want.State, got.State)
			assert.Equal(t, tt.want.NeedsAuth, got.NeedsAuth)
			assert.Equal(t, tt.want.Quitting, got.Quitting)
			assert.NotNil(t, got.CityInput)
			assert.Len(t, got.UnitOptions, 3)
			assert.Len(t, got.ViewOptions, 5)
		})
	}
}

func TestStateTransitions(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})

	tests := []struct {
		name          string
		initialState  SetupState
		msg          tea.Msg
		expectedState SetupState
	}{
		{
			name:         "city search to city select",
			initialState: StateCitySearch,
			msg: CitiesSearchResult{
				cities: []models.City{{Name: "London"}},
				err:    nil,
			},
			expectedState: StateCitySelect,
		},
		{
			name:         "empty city search results return to city input",
			initialState: StateCitySearch,
			msg: CitiesSearchResult{
				cities: []models.City{},
				err:    nil,
			},
			expectedState: StateCity,
		},
		{
			name:         "error in city search returns to city input",
			initialState: StateCitySearch,
			msg: CitiesSearchResult{
				cities: nil,
				err:    assert.AnError,
			},
			expectedState: StateCity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.State = tt.initialState
			updatedModel, _ := m.Update(tt.msg)
			assert.Equal(t, tt.expectedState, updatedModel.(Model).State)
		})
	}
}

func TestKeyHandling(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})

	tests := []struct {
		name          string
		state         SetupState
		key          string
		expectedState SetupState
		expectedQuit  bool
	}{
		{
			name:          "ctrl+c quits from any state",
			state:         StateCity,
			key:          "ctrl+c",
			expectedState: StateCity,
			expectedQuit:  true,
		},
		{
			name:          "esc from city select returns to city input",
			state:         StateCitySelect,
			key:          "esc",
			expectedState: StateCity,
			expectedQuit:  false,
		},
		{
			name:          "down key in units state increments cursor",
			state:         StateUnits,
			key:          "down",
			expectedState: StateUnits,
			expectedQuit:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.State = tt.state
			initialCursor := m.UnitCursor
			updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)})
			updated := updatedModel.(Model)

			assert.Equal(t, tt.expectedState, updated.State)
			assert.Equal(t, tt.expectedQuit, updated.Quitting)

			if tt.key == "down" && tt.state == StateUnits {
				assert.Equal(t, initialCursor+1, updated.UnitCursor)
			}

			if tt.expectedQuit {
				assert.NotNil(t, cmd)
			}
		})
	}
}

func TestCitySearch(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})
	m.CitySearchQuery = "London"

	cmd := m.searchCities()
	assert.NotNil(t, cmd)

	// Test with nil client
	m.Client = nil
	cmd = m.searchCities()
	result := cmd().(CitiesSearchResult)
	assert.Error(t, result.err)
	assert.Len(t, result.cities, 0)
} 