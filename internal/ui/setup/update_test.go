package setup

import (
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestStateTransitions(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})

	tests := []struct {
		name          string
		initialState  SetupState
		msg           tea.Msg
		expectedState SetupState
		expectedCmd   bool // Whether we expect a non-nil command
	}{
		{
			name:         "city search to city select",
			initialState: StateCitySearch,
			msg: CitiesSearchResult{
				cities: []models.City{{Name: "London"}},
				err:    nil,
			},
			expectedState: StateCitySelect,
			expectedCmd:   false,
		},
		{
			name:         "empty city search results return to city input",
			initialState: StateCitySearch,
			msg: CitiesSearchResult{
				cities: []models.City{},
				err:    nil,
			},
			expectedState: StateCity,
			expectedCmd:   false,
		},
		{
			name:         "error in city search returns to city input",
			initialState: StateCitySearch,
			msg: CitiesSearchResult{
				cities: nil,
				err:    assert.AnError,
			},
			expectedState: StateCity,
			expectedCmd:   false,
		},
		{
			name:          "setup complete transitions to complete state",
			initialState:  StateAuth,
			msg:           SetupCompleteMsg{},
			expectedState: StateComplete,
			expectedCmd:   false,
		},
		{
			name:          "window size msg updates dimensions",
			initialState:  StateCity,
			msg:           tea.WindowSizeMsg{Width: 100, Height: 50},
			expectedState: StateCity,
			expectedCmd:   false,
		},
		{
			name:          "spinner tick updates spinner",
			initialState:  StateCitySearch,
			msg:           spinner.TickMsg{},
			expectedState: StateCitySearch,
			expectedCmd:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.State = tt.initialState
			updatedModel, cmd := m.Update(tt.msg)
			updated := updatedModel.(Model)

			assert.Equal(t, tt.expectedState, updated.State)
			if tt.expectedCmd {
				assert.NotNil(t, cmd)
			}
		})
	}
}

func TestKeyHandling(t *testing.T) {
	tests := []struct {
		name          string
		state         SetupState
		key           string
		keyType       tea.KeyType
		expectedState SetupState
		expectedQuit  bool
	}{
		{
			name:          "ctrl+c quits from any state",
			state:         StateCity,
			key:           "ctrl+c",
			keyType:       tea.KeyCtrlC,
			expectedState: StateCity,
			expectedQuit:  true,
		},
		{
			name:          "esc from city select returns to city input",
			state:         StateCitySelect,
			key:           "esc",
			keyType:       tea.KeyEsc,
			expectedState: StateCity,
			expectedQuit:  false,
		},
		{
			name:          "enter in city select transitions to units state",
			state:         StateCitySelect,
			key:           "enter",
			keyType:       tea.KeyEnter,
			expectedState: StateUnits,
			expectedQuit:  false,
		},
		{
			name:          "down key in units state increments cursor",
			state:         StateUnits,
			key:           "down",
			keyType:       tea.KeyDown,
			expectedState: StateUnits,
			expectedQuit:  false,
		},
		{
			name:          "up key in view state decrements cursor",
			state:         StateView,
			key:           "up",
			keyType:       tea.KeyUp,
			expectedState: StateView,
			expectedQuit:  false,
		},
		{
			name:          "j key acts as down key",
			state:         StateTips,
			key:           "j",
			keyType:       tea.KeyRunes,
			expectedState: StateTips,
			expectedQuit:  false,
		},
		{
			name:          "k key acts as up key",
			state:         StateAuth,
			key:           "k",
			keyType:       tea.KeyRunes,
			expectedState: StateAuth,
			expectedQuit:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup model for test
			m := NewModel(&config.Config{}, false, &api.Client{})
			m.State = tt.state

			// Add required test data
			if tt.state == StateCitySelect {
				m.CityOptions = []models.City{{Name: "London"}}
			}

			// Track initial cursor value for cursor movement tests
			var initialCursor int
			switch tt.state {
			case StateUnits:
				initialCursor = m.UnitCursor
			case StateView:
				initialCursor = m.ViewCursor
			case StateTips:
				initialCursor = m.TipCursor
			case StateAuth:
				initialCursor = m.AuthCursor
			}

			// Create key message
			var keyMsg tea.KeyMsg
			if tt.keyType == tea.KeyRunes {
				keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			} else {
				keyMsg = tea.KeyMsg{Type: tt.keyType}
			}

			// Call Update
			updatedModel, cmd := m.Update(keyMsg)
			updated := updatedModel.(Model)

			// Verify expected state
			assert.Equal(t, tt.expectedState, updated.State)
			assert.Equal(t, tt.expectedQuit, updated.Quitting)

			// Verify cursor movement for movement tests
			if tt.key == "down" || tt.key == "j" {
				switch tt.state {
				case StateUnits:
					assert.Equal(t, initialCursor+1, updated.UnitCursor)
				case StateView:
					assert.Equal(t, initialCursor+1, updated.ViewCursor)
				case StateTips:
					assert.Equal(t, initialCursor+1, updated.TipCursor)
				case StateAuth:
					assert.Equal(t, initialCursor+1, updated.AuthCursor)
				}
			} else if tt.key == "up" || tt.key == "k" {
				switch tt.state {
				case StateUnits:
					// If initialCursor is already 0, it should remain 0
					if initialCursor > 0 {
						assert.Equal(t, initialCursor-1, updated.UnitCursor)
					} else {
						assert.Equal(t, 0, updated.UnitCursor)
					}
				case StateView:
					// If initialCursor is already 0, it should remain 0
					if initialCursor > 0 {
						assert.Equal(t, initialCursor-1, updated.ViewCursor)
					} else {
						assert.Equal(t, 0, updated.ViewCursor)
					}
				case StateTips:
					// If initialCursor is already 0, it should remain 0
					if initialCursor > 0 {
						assert.Equal(t, initialCursor-1, updated.TipCursor)
					} else {
						assert.Equal(t, 0, updated.TipCursor)
					}
				case StateAuth:
					// If initialCursor is already 0, it should remain 0
					if initialCursor > 0 {
						assert.Equal(t, initialCursor-1, updated.AuthCursor)
					} else {
						assert.Equal(t, 0, updated.AuthCursor)
					}
				}
			}

			// Verify quit command
			if tt.expectedQuit {
				assert.NotNil(t, cmd)
			}
		})
	}
}

func TestEnterKeyHandling(t *testing.T) {
	tests := []struct {
		name          string
		setupModel    func() Model
		expectedState SetupState
	}{
		{
			name: "city select to units",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateCitySelect
				m.CityOptions = []models.City{{Name: "London"}}
				return m
			},
			expectedState: StateUnits,
		},
		{
			name: "units to view",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateUnits
				return m
			},
			expectedState: StateView,
		},
		{
			name: "view to tips",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateView
				return m
			},
			expectedState: StateTips,
		},
		{
			name: "tips to auth choices",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, true, &api.Client{})
				m.State = StateTips
				return m
			},
			expectedState: StateApiKeyOption,
		},
		{
			name: "auth to complete (no auth selected)",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, true, &api.Client{})
				m.State = StateAuth
				m.AuthCursor = 1
				return m
			},
			expectedState: StateComplete,
		},
		{
			name: "empty city input doesn't advance",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateCity
				m.CityInput.SetValue("")
				return m
			},
			expectedState: StateCity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.setupModel()
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			updated := updatedModel.(Model)
			assert.Equal(t, tt.expectedState, updated.State)
		})
	}
}

func TestCitySearch(t *testing.T) {
	t.Run("successful search creates command", func(t *testing.T) {
		m := NewModel(&config.Config{}, false, &api.Client{})
		m.CitySearchQuery = "London"

		cmd := m.searchCities()
		assert.NotNil(t, cmd)
	})

	t.Run("nil client returns error result", func(t *testing.T) {
		m := NewModel(&config.Config{}, false, nil)
		m.CitySearchQuery = "London"

		cmd := m.searchCities()
		result := cmd().(CitiesSearchResult)
		assert.Error(t, result.err)
		assert.Contains(t, result.err.Error(), "API client not initialized")
		assert.Len(t, result.cities, 0)
	})

	t.Run("empty search query", func(t *testing.T) {
		m := NewModel(&config.Config{}, false, &api.Client{})
		m.CitySearchQuery = ""

		// Even with empty query, we should get a valid command
		cmd := m.searchCities()
		assert.NotNil(t, cmd)
	})
}

func TestCursorMovement(t *testing.T) {
	tests := []struct {
		name       string
		state      SetupState
		setupModel func() Model
		key        tea.KeyMsg
		checkFunc  func(*testing.T, Model)
	}{
		{
			name:  "units cursor up",
			state: StateUnits,
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.UnitCursor = 1
				return m
			},
			key: tea.KeyMsg{Type: tea.KeyUp},
			checkFunc: func(t *testing.T, m Model) {
				assert.Equal(t, 0, m.UnitCursor)
			},
		},
		{
			name:  "units cursor down",
			state: StateUnits,
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.UnitCursor = 0
				return m
			},
			key: tea.KeyMsg{Type: tea.KeyDown},
			checkFunc: func(t *testing.T, m Model) {
				assert.Equal(t, 1, m.UnitCursor)
			},
		},
		{
			name:  "units cursor upper bound",
			state: StateUnits,
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.UnitCursor = 0
				return m
			},
			key: tea.KeyMsg{Type: tea.KeyUp},
			checkFunc: func(t *testing.T, m Model) {
				assert.Equal(t, 0, m.UnitCursor, "cursor shouldn't go below 0")
			},
		},
		{
			name:  "units cursor lower bound",
			state: StateUnits,
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.UnitCursor = 2 // Last option
				return m
			},
			key: tea.KeyMsg{Type: tea.KeyDown},
			checkFunc: func(t *testing.T, m Model) {
				assert.Equal(t, 2, m.UnitCursor, "cursor shouldn't exceed options length")
			},
		},
		{
			name:  "view cursor movement",
			state: StateView,
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.ViewCursor = 1
				return m
			},
			key: tea.KeyMsg{Type: tea.KeyDown},
			checkFunc: func(t *testing.T, m Model) {
				assert.Equal(t, 2, m.ViewCursor)
			},
		},
		{
			name:  "city cursor movement",
			state: StateCitySelect,
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.CityOptions = []models.City{
					{Name: "London"},
					{Name: "Paris"},
					{Name: "Berlin"},
				}
				m.CityCursor = 0
				return m
			},
			key: tea.KeyMsg{Type: tea.KeyDown},
			checkFunc: func(t *testing.T, m Model) {
				assert.Equal(t, 1, m.CityCursor)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.setupModel()
			model.State = tt.state
			updatedModel, _ := model.Update(tt.key)
			updated := updatedModel.(Model)
			tt.checkFunc(t, updated)
		})
	}
}

func TestTextInputHandling(t *testing.T) {
	t.Run("enter with value transitions to search", func(t *testing.T) {
		m := NewModel(&config.Config{}, false, &api.Client{})
		m.State = StateCity
		m.CityInput.SetValue("London")

		updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		updated := updatedModel.(Model)

		assert.Equal(t, StateCitySearch, updated.State)
		assert.Equal(t, "London", updated.CitySearchQuery)
		assert.NotNil(t, cmd)
	})

	t.Run("enter with empty value stays in city state", func(t *testing.T) {
		m := NewModel(&config.Config{}, false, &api.Client{})
		m.State = StateCity
		m.CityInput.SetValue("")

		updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		updated := updatedModel.(Model)

		assert.Equal(t, StateCity, updated.State)
	})

	t.Run("text input updates correctly", func(t *testing.T) {
		m := NewModel(&config.Config{}, false, &api.Client{})
		m.State = StateCity

		// Type a character
		updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("L")})
		updated := updatedModel.(Model)

		assert.Equal(t, StateCity, updated.State)
		assert.Equal(t, "L", updated.CityInput.Value())
		assert.NotNil(t, cmd)
	})

	t.Run("ctrl+c quits", func(t *testing.T) {
		m := NewModel(&config.Config{}, false, &api.Client{})
		m.State = StateCity

		updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
		updated := updatedModel.(Model)

		assert.True(t, updated.Quitting)
		assert.NotNil(t, cmd)
	})
}

func TestMiscUpdateHandling(t *testing.T) {
	t.Run("unknown message passes through", func(t *testing.T) {
		m := NewModel(&config.Config{}, false, &api.Client{})

		// Create a custom message type
		type customMsg struct{}

		updatedModel, cmd := m.Update(customMsg{})
		updated := updatedModel.(Model)

		// Model should remain unchanged
		assert.Equal(t, m.State, updated.State)
		assert.Nil(t, cmd)
	})
}
