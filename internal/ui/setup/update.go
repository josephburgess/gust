package setup

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

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
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit
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
			m.Config.DefaultCity = m.CityInput.Value()
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
			m.Config.DefaultCity = m.CityInput.Value()
			m.State = StateUnits
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.CityInput, cmd = m.CityInput.Update(msg)
	return m, cmd
}
