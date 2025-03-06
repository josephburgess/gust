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
		msg           tea.Msg
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
		key           string
		expectedState SetupState
		expectedQuit  bool
	}{
		{
			name:          "ctrl+c quits from any state",
			state:         StateCity,
			key:           "ctrl+c",
			expectedState: StateCity,
			expectedQuit:  true,
		},
		{
			name:          "esc from city select returns to city input",
			state:         StateCitySelect,
			key:           "esc",
			expectedState: StateCity,
			expectedQuit:  false,
		},
		{
			name:          "down key in units state increments cursor",
			state:         StateUnits,
			key:           "down",
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

	m.Client = nil
	cmd = m.searchCities()
	result := cmd().(CitiesSearchResult)
	assert.Error(t, result.err)
	assert.Len(t, result.cities, 0)
}

func TestView(t *testing.T) {
	tests := []struct {
		name          string
		setupModel    func() Model
		expectedParts []string
	}{
		{
			name: "city input state",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateCity
				return m
			},
			expectedParts: []string{
				"Enter a default city",
				"You can enter a country code",
			},
		},
		{
			name: "city search state",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateCitySearch
				m.CitySearchQuery = "London"
				return m
			},
			expectedParts: []string{
				"Searching for cities",
				"Looking for \"London\"",
			},
		},
		{
			name: "city select state with results",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateCitySelect
				m.CityOptions = []models.City{
					{Name: "London", Country: "GB"},
					{Name: "Paris", Country: "FR"},
				}
				return m
			},
			expectedParts: []string{
				"Select your town or city",
				"London",
				"Paris",
			},
		},
		{
			name: "complete state",
			setupModel: func() Model {
				m := NewModel(&config.Config{
					DefaultCity: "London",
					Units:       "metric",
					DefaultView: "detailed",
					ShowTips:    true,
				}, false, &api.Client{})
				m.State = StateComplete
				return m
			},
			expectedParts: []string{
				"Setup complete",
				"Default city: London",
				"Units: metric",
				"Default view: detailed",
				"Tips enabled",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.setupModel()
			view := model.View()
			for _, expectedPart := range tt.expectedParts {
				assert.Contains(t, view, expectedPart)
			}
		})
	}
}

func TestCountryEmoji(t *testing.T) {
	tests := []struct {
		name        string
		countryCode string
		want        string
	}{
		{
			name:        "valid country code",
			countryCode: "GB",
			want:        "🇬🇧",
		},
		{
			name:        "empty country code",
			countryCode: "",
			want:        "🌍",
		},
		{
			name:        "invalid length",
			countryCode: "GBR",
			want:        "🌍",
		},
		{
			name:        "lowercase code",
			countryCode: "fr",
			want:        "🇫🇷",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getCountryEmoji(tt.countryCode)
			assert.Equal(t, tt.want, got)
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
			name: "tips to complete (no auth needed)",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.State = StateTips
				return m
			},
			expectedState: StateComplete,
		},
		{
			name: "tips to auth (auth needed)",
			setupModel: func() Model {
				m := NewModel(&config.Config{}, true, &api.Client{})
				m.State = StateTips
				return m
			},
			expectedState: StateAuth,
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

func TestCursorMovement(t *testing.T) {
	tests := []struct {
		name       string
		state      SetupState
		setupModel func() Model
		key        string
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
			key: "up",
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
			key: "down",
			checkFunc: func(t *testing.T, m Model) {
				assert.Equal(t, 1, m.UnitCursor)
			},
		},
		{
			name:  "view cursor movement bounds",
			state: StateView,
			setupModel: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				m.ViewCursor = 0
				return m
			},
			key: "up",
			checkFunc: func(t *testing.T, m Model) {
				assert.Equal(t, 0, m.ViewCursor, "cursor should not go below 0")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.setupModel()
			model.State = tt.state
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)})
			updated := updatedModel.(Model)
			tt.checkFunc(t, updated)
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("empty city search query", func(t *testing.T) {
		m := NewModel(&config.Config{}, false, &api.Client{})
		m.State = StateCity
		updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		updated := updatedModel.(Model)
		assert.Equal(t, StateCity, updated.State, "shouldnt proceed with empty city")
	})

	t.Run("cursor bounds checking", func(t *testing.T) {
		m := NewModel(&config.Config{}, false, &api.Client{})
		m.State = StateUnits
		m.UnitCursor = len(m.UnitOptions) - 1

		updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("down")})
		updated := updatedModel.(Model)
		assert.Equal(t, len(m.UnitOptions)-1, updated.UnitCursor, "curor shouldnt go oob")
	})

	t.Run("window resize handling", func(t *testing.T) {
		m := NewModel(&config.Config{}, false, &api.Client{})
		newWidth, newHeight := 100, 50
		updatedModel, _ := m.Update(tea.WindowSizeMsg{Width: newWidth, Height: newHeight})
		updated := updatedModel.(Model)
		assert.Equal(t, newWidth, updated.Width)
		assert.Equal(t, newHeight, updated.Height)
	})
}

// mock tea.Program
type mockProgram struct {
	model Model
	err   error
}

func (m *mockProgram) Start() error {
	return m.err
}

func newMockProgram(model tea.Model) *mockProgram {
	return &mockProgram{
		model: model.(Model),
		err:   nil,
	}
}

// wraps runsetup for testing (otherwise it launches and softlocks tests)
type testRunner struct {
	program *mockProgram
}

func newTestRunner(mockErr error) *testRunner {
	return &testRunner{
		program: &mockProgram{err: mockErr},
	}
}

func (tr *testRunner) runSetup(cfg *config.Config, needsAuth bool) error {
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = "https://breeze.joeburgess.dev"
	}

	model := NewModel(cfg, needsAuth, &api.Client{})
	if tr.program == nil {
		tr.program = &mockProgram{
			model: model,
			err:   nil,
		}
	} else {
		tr.program.model = model
	}
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
			name: "setup with empty API URL",
			cfg: &config.Config{
				Units: "metric",
			},
			needsAuth: false,
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name: "program error",
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

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.mockErr, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.cfg.ApiUrl == "" {
				assert.Equal(t, "https://breeze.joeburgess.dev", tt.cfg.ApiUrl)
			}
		})
	}
}

func TestModelInit(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})
	cmd := m.Init()
	assert.NotNil(t, cmd, "Init should return a command")
}

func TestSpinnerUpdate(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})
	m.State = StateCitySearch

	updatedModel, cmd := m.Update(spinner.TickMsg{})
	assert.NotNil(t, cmd, "Spinner update should return a command")
	assert.Equal(t, StateCitySearch, updatedModel.(Model).State)
}

func TestViewRendering(t *testing.T) {
	tests := []struct {
		name     string
		model    Model
		validate func(*testing.T, string)
	}{
		{
			name: "renders content",
			model: func() Model {
				m := NewModel(&config.Config{}, false, &api.Client{})
				return m
			}(),
			validate: func(t *testing.T, view string) {
				assert.Contains(t, view, "____ ___  _______/ /_")
				assert.Contains(t, view, "Simple terminal weather")
				assert.Contains(t, view, "Enter a default city")
				assert.Contains(t, view, "╔═")
				assert.Contains(t, view, "║")
				assert.Contains(t, view, "╚═")
				assert.Contains(t, view, "╗")
				assert.Contains(t, view, "╝")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := tt.model.View()
			tt.validate(t, view)
		})
	}
}
