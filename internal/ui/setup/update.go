package setup

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
)

type CitiesSearchResult struct {
	cities []models.City
	err    error
}

func (m Model) searchCities() tea.Cmd {
	return func() tea.Msg {
		if m.Client == nil {
			return CitiesSearchResult{[]models.City{}, fmt.Errorf("API client not initialized")}
		}

		cities, err := m.Client.SearchCities(m.CitySearchQuery)
		if err != nil {
			return CitiesSearchResult{[]models.City{}, err}
		}
		return CitiesSearchResult{cities, nil}
	}
}

// updates the model based on messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
		switch m.State {
		case StateCity, StateCitySearch, StateCitySelect:
			return m.handleCityKey(msg)
		case StatePreferences:
			return m.handlePreferencesKey(msg)
		case StateApiKeyOption:
			return m.handleSignInKey(msg)
		case StateApiKeyInput:
			return m.handleApiKeyInput(msg)
		case StateComplete:
			return m.handleCompleteKey(msg)
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case SetupCompleteMsg:
		m.State = StateComplete
	case CitiesSearchResult:
		if m.State != StateCitySearch {
			return m, nil
		}
		m.CityCursor = 0
		m.CityOptions = msg.cities
		switch {
		case msg.err != nil:
			m.SearchErr = msg.err.Error()
			m.State = StateCity
		case len(msg.cities) == 0:
			m.SearchErr = fmt.Sprintf("no matches for %q", m.CitySearchQuery)
			m.State = StateCity
		default:
			m.State = StateCitySelect
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	default:
		var cmd tea.Cmd
		switch m.State {
		case StateCity, StateCitySearch, StateCitySelect:
			m.CityInput, cmd = m.CityInput.Update(msg)
		case StateApiKeyInput:
			m.ApiKeyInput, cmd = m.ApiKeyInput.Update(msg)
		}
		return m, cmd
	}
	return m, nil
}

// handleCityKey runs the city step. The input keeps focus, so typing again
// and pressing enter searches again.
func (m Model) handleCityKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		query := strings.TrimSpace(m.CityInput.Value())
		if m.State == StateCitySelect && query == m.CitySearchQuery {
			m.Config.DefaultCity = m.CityOptions[m.CityCursor].Name
			m.State = StatePreferences
			return m, nil
		}
		if query == "" || m.State == StateCitySearch {
			return m, nil
		}
		m.CitySearchQuery = query
		m.SearchErr = ""
		m.State = StateCitySearch
		return m, m.searchCities()
	case "up":
		if m.State == StateCitySelect {
			m.CityCursor = max(m.CityCursor-1, 0)
		}
		return m, nil
	case "down":
		if m.State == StateCitySelect {
			m.CityCursor = min(m.CityCursor+1, len(m.CityOptions)-1)
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.CityInput, cmd = m.CityInput.Update(msg)
	return m, cmd
}

func (m Model) handlePreferencesKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	cursors := []*int{&m.UnitCursor, &m.ViewCursor, &m.TipCursor}
	sizes := []int{len(unitOptions), len(viewOptions), len(tipOptions)}
	cursor, size := cursors[m.PrefGroup], sizes[m.PrefGroup]

	switch msg.String() {
	case "tab":
		m.PrefGroup = (m.PrefGroup + 1) % 3
	case "shift+tab":
		m.PrefGroup = (m.PrefGroup + 2) % 3
	case "up", "k":
		*cursor = max(*cursor-1, 0)
	case "down", "j":
		*cursor = min(*cursor+1, size-1)
	case "esc":
		m.State = StateCitySelect
		if len(m.CityOptions) == 0 {
			m.State = StateCity
		}
	case "enter":
		m.Config.Units = unitOptions[m.UnitCursor].value
		m.Config.DefaultView = viewOptions[m.ViewCursor].value
		m.Config.ShowTips = tipOptions[m.TipCursor].value == "on"
		m.State = StateApiKeyOption
	}
	return m, nil
}

func (m Model) handleSignInKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.ApiKeyCursor = 0
	case "down", "j":
		m.ApiKeyCursor = 1
	case "esc":
		m.State = StatePreferences
	case "enter":
		if m.ApiKeyCursor == 1 {
			m.State = StateApiKeyInput
			return m, m.ApiKeyInput.Focus()
		}
		m.WantsGitHub = m.NeedsAuth
		m.State = StateComplete
	}
	return m, nil
}

func (m Model) handleApiKeyInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.ApiKeyInput.Blur()
		m.State = StateApiKeyOption
		return m, nil
	case "enter":
		key := strings.TrimSpace(m.ApiKeyInput.Value())
		if key == "" {
			return m, nil
		}
		authConfig := &config.AuthConfig{
			APIKey:     key,
			ServerURL:  m.Config.ApiUrl,
			LastAuth:   time.Now(),
			GithubUser: "OpenWeather API User",
		}
		if err := config.SaveAuthConfig(authConfig); err != nil {
			m.SearchErr = err.Error()
			return m, nil
		}
		m.NeedsAuth = false
		m.WantsGitHub = false
		m.ApiKeyInput.Blur()
		m.State = StateComplete
		return m, nil
	}

	var cmd tea.Cmd
	m.ApiKeyInput, cmd = m.ApiKeyInput.Update(msg)
	return m, cmd
}

func (m Model) handleCompleteKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		return m, tea.Quit
	case "esc":
		m.WantsGitHub = false
		m.State = StateApiKeyOption
	}
	return m, nil
}
