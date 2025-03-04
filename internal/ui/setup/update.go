package setup

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/josephburgess/gust/internal/models"
)

type CitiesSearchResult struct {
	cities []models.City
	err    error
}

func (m Model) searchCities() tea.Cmd {
	return func() tea.Msg {
		fmt.Println("Searching for cities with query:", m.CitySearchQuery)
		if m.Client == nil {
			fmt.Println("ERROR: API client is nil!")
			return CitiesSearchResult{[]models.City{}, fmt.Errorf("API client not initialized")}
		}

		cities, err := m.Client.SearchCities(m.CitySearchQuery)
		if err != nil {
			fmt.Println("Error searching for cities:", err)
			return CitiesSearchResult{[]models.City{}, err}
		}
		fmt.Printf("Found %d cities\n", len(cities))
		return CitiesSearchResult{cities, nil}
	}
}

// updates the model based on messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.State == StateCity {
			return m.handleTextInput(msg)
		}
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case SetupCompleteMsg:
		m.State = StateComplete
		return m, nil
	case CitiesSearchResult:
		fmt.Printf("Received city search result: %d cities, error: %v\n",
			len(msg.cities), msg.err)

		if msg.err != nil {
			// Log the error and go back to input
			fmt.Println("Error searching for cities:", msg.err)
			m.State = StateCity
			return m, nil
		}

		m.CityOptions = msg.cities
		m.CityCursor = 0

		if len(m.CityOptions) == 0 {
			fmt.Println("No cities found, returning to input")
			// No results, go back to input
			m.State = StateCity
			return m, nil
		}

		fmt.Println("Moving to city selection state")
		m.State = StateCitySelect
		return m, nil
	}

	return m, nil
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

		if m.NeedsAuth {
			m.State = StateAuth
		} else {
			m.State = StateComplete
		}

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
			m.CitySearchQuery = m.CityInput.Value() // Set the search query
			m.State = StateCitySearch               // Change state to searching
			return m, m.searchCities()              // Start the search command
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.CityInput, cmd = m.CityInput.Update(msg)
	return m, cmd
}
