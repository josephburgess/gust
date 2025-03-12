package setup

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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

func (m Model) handleApiKeyInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit
	case "enter":
		if m.ApiKeyInput.Value() != "" {
			authConfig := &config.AuthConfig{
				APIKey:     m.ApiKeyInput.Value(),
				ServerURL:  m.Config.ApiUrl,
				LastAuth:   time.Now(),
				GithubUser: "OpenWeather API User",
			}

			if err := config.SaveAuthConfig(authConfig); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				m.NeedsAuth = false
			}
		}
		m.State = StateComplete
		return m, nil
	case "esc":
		m.State = StateApiKeyOption
		return m, nil
	}

	var cmd tea.Cmd
	m.ApiKeyInput, cmd = m.ApiKeyInput.Update(msg)
	return m, cmd
}

// updates the model based on messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.State {
		case StateCity:
			return m.handleTextInput(msg)
		case StateApiKeyInput:
			return m.handleApiKeyInput(msg)
		}
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case SetupCompleteMsg:
		m.State = StateComplete
		return m, nil
	case CitiesSearchResult:
		if msg.err != nil {
			fmt.Printf("Error searching cities: %v\n", msg.err)
			m.State = StateCity
			return m, nil
		}

		m.CityOptions = msg.cities
		m.CityCursor = 0

		if len(m.CityOptions) == 0 {
			fmt.Println("No cities found. Please try a different search.")
			m.State = StateCity
			return m, nil
		}

		m.State = StateCitySelect
		return m, nil
	case spinner.TickMsg:
		var spinnerCmd tea.Cmd
		m.Spinner, spinnerCmd = m.Spinner.Update(msg)
		cmds = append(cmds, spinnerCmd)
	default:
		switch m.State {
		case StateCity:
			var tiCmd tea.Cmd
			m.CityInput, tiCmd = m.CityInput.Update(msg)
			if tiCmd != nil {
				cmds = append(cmds, tiCmd)
			}
		case StateApiKeyInput:
			var tiCmd tea.Cmd
			m.ApiKeyInput, tiCmd = m.ApiKeyInput.Update(msg)
			if tiCmd != nil {
				cmds = append(cmds, tiCmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit
	case "esc":
		if m.State == StateCitySelect {
			m.State = StateCity
		}
		return m, nil
	case "enter":
		return m.handleEnterKey()
	case "up", "k":
		return m.handleUpKey()
	case "down", "j":
		return m.handleDownKey()
	}
	return m, nil
}

func (m Model) handleEnterKey() (tea.Model, tea.Cmd) {
	switch m.State {
	case StateCity:
		if m.CityInput.Value() != "" {
			m.CitySearchQuery = m.CityInput.Value()
			m.State = StateCitySearch
			return m, m.searchCities()
		}

	case StateCitySelect:
		if len(m.CityOptions) > 0 {
			selectedCity := m.CityOptions[m.CityCursor]
			m.Config.DefaultCity = selectedCity.Name
			m.State = StateUnits
		}

	case StateUnits:
		unitValues := []string{"metric", "imperial", "standard"}
		m.Config.Units = unitValues[m.UnitCursor]
		m.State = StateView

	case StateView:
		viewValues := []string{"default", "compact", "daily", "hourly", "full"}
		m.Config.DefaultView = viewValues[m.ViewCursor]
		m.State = StateTips

	case StateTips:
		m.Config.ShowTips = (m.TipCursor == 0)
		m.State = StateApiKeyOption

	case StateApiKeyOption:
		if m.ApiKeyCursor == 0 {
			if m.NeedsAuth {
				m.State = StateAuth
			} else {
				m.State = StateComplete
			}
		} else {
			m.State = StateApiKeyInput
			m.ApiKeyInput.Focus()
		}

	case StateApiKeyInput:
		if m.ApiKeyInput.Value() != "" {
			authConfig := &config.AuthConfig{
				APIKey:     m.ApiKeyInput.Value(),
				ServerURL:  m.Config.ApiUrl,
				LastAuth:   time.Now(),
				GithubUser: "",
			}

			if err := config.SaveAuthConfig(authConfig); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				m.NeedsAuth = false
			}
		}
		m.State = StateComplete

	case StateAuth:
		if err := m.Config.Save(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		if m.AuthCursor == 0 {
			// user chose to auth
			return m, tea.Quit
		} else {
			// user skipped auth
			m.State = StateComplete
		}

	case StateComplete:
		return m, tea.Quit
	}

	return m, nil
}

// process up n down keypresses
func (m Model) handleUpKey() (tea.Model, tea.Cmd) {
	switch m.State {
	case StateUnits:
		if m.UnitCursor > 0 {
			m.UnitCursor--
		}
	case StateView:
		if m.ViewCursor > 0 {
			m.ViewCursor--
		}
	case StateAuth:
		if m.AuthCursor > 0 {
			m.AuthCursor--
		}
	case StateCitySelect:
		if m.CityCursor > 0 {
			m.CityCursor--
		}
	case StateTips:
		if m.TipCursor > 0 {
			m.TipCursor--
		}
	case StateApiKeyOption:
		if m.ApiKeyCursor > 0 {
			m.ApiKeyCursor--
		}
	}
	return m, nil
}

func (m Model) handleDownKey() (tea.Model, tea.Cmd) {
	switch m.State {
	case StateUnits:
		if m.UnitCursor < len(m.UnitOptions)-1 {
			m.UnitCursor++
		}
	case StateView:
		if m.ViewCursor < len(m.ViewOptions)-1 {
			m.ViewCursor++
		}
	case StateAuth:
		if m.AuthCursor < len(m.AuthOptions)-1 {
			m.AuthCursor++
		}
	case StateCitySelect:
		if m.CityCursor < len(m.CityOptions)-1 {
			m.CityCursor++
		}
	case StateTips:
		if m.TipCursor < len(m.TipOptions)-1 {
			m.TipCursor++
		}
	case StateApiKeyOption:
		if m.ApiKeyCursor < len(m.ApiKeyOptions)-1 {
			m.ApiKeyCursor++
		}
	}
	return m, nil
}

func (m Model) handleTextInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit
	case "enter":
		if m.CityInput.Value() != "" {
			m.CitySearchQuery = m.CityInput.Value()
			m.State = StateCitySearch
			return m, m.searchCities()
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.CityInput, cmd = m.CityInput.Update(msg)
	return m, cmd
}
