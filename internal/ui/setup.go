package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josephburgess/gust/internal/config"
)

type setupState int

const (
	stateCity setupState = iota
	stateUnits
	stateAuth
	stateComplete
)

const asciiLogo = `
                        __
      ____ ___  _______/ /_
     / ** '/ / / / **_/ __/
    / /_/ / /_/ (__  ) /_   _
    \__, /\__,_/____/\__/  (_)
   /____/                      `

type setupModel struct {
	config        *config.Config
	state         setupState
	cityInput     textinput.Model
	unitOptions   []string
	unitCursor    int
	authOptions   []string
	authCursor    int
	needsAuth     bool
	width, height int
	quitting      bool
}

func NewSetupModel(cfg *config.Config, needsAuth bool) setupModel {
	ti := textinput.New()
	ti.Placeholder = "Wherever the wind blows..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = len(ti.Placeholder)
	ti.PromptStyle = lipgloss.NewStyle().Foreground(love)
	ti.TextStyle = lipgloss.NewStyle().Foreground(text)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(gold)

	return setupModel{
		config:      cfg,
		state:       stateCity,
		cityInput:   ti,
		unitOptions: []string{"metric (°C, km/h) 🌡️", "imperial (°F, mph) 🌡️", "standard (K, m/s) 🌡️"},
		unitCursor:  0, // Default to metric
		authOptions: []string{"Yes, authenticate with GitHub 🔑", "No, I'll do it later ⏱️"},
		authCursor:  0,
		needsAuth:   needsAuth,
		quitting:    false,
	}
}

// init
func (m setupModel) Init() tea.Cmd {
	return textinput.Blink
}

// update
func (m setupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == stateCity {
			return m.handleTextInput(msg)
		}
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case setupCompleteMsg:
		m.state = stateComplete
		return m, nil
	}

	return m, nil
}

// view
func (m setupModel) View() string {
	content := m.buildContent()
	return m.centerContent(content)
}

func (m setupModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
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

func (m setupModel) handleEnterKey() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateCity:
		if m.cityInput.Value() != "" {
			m.config.DefaultCity = m.cityInput.Value()
			m.state = stateUnits
		}

	case stateUnits:
		unitValues := []string{"metric", "imperial", "standard"}
		m.config.Units = unitValues[m.unitCursor]

		if m.needsAuth {
			m.state = stateAuth
		} else {
			m.state = stateComplete
		}

	case stateAuth:
		// save config in case auth fails
		if err := m.config.Save(); err != nil {
			fmt.Println("Warning: Failed to save configuration before authentication.")
		}

		if m.authCursor == 0 {
			// user chose to auth
			return m, tea.Quit
		} else {
			// user skipped auth
			m.state = stateComplete
		}

	case stateComplete:
		return m, tea.Quit
	}

	return m, nil
}

func (m setupModel) handleUpKey() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateUnits:
		if m.unitCursor > 0 {
			m.unitCursor--
		}
	case stateAuth:
		if m.authCursor > 0 {
			m.authCursor--
		}
	}
	return m, nil
}

func (m setupModel) handleDownKey() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateUnits:
		if m.unitCursor < len(m.unitOptions)-1 {
			m.unitCursor++
		}
	case stateAuth:
		if m.authCursor < len(m.authOptions)-1 {
			m.authCursor++
		}
	}
	return m, nil
}

func (m setupModel) handleTextInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "enter":
		if m.cityInput.Value() != "" {
			m.config.DefaultCity = m.cityInput.Value()
			m.state = stateUnits
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.cityInput, cmd = m.cityInput.Update(msg)
	return m, cmd
}

func (m setupModel) buildContent() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render(asciiLogo) + "\n\n")
	sb.WriteString(boxStyle.Render(subtitleStyle.Render("Simple terminal weather 🌤️")) + "\n\n")

	switch m.state {
	case stateCity:
		sb.WriteString(highlightStyle.Render("Enter a default city 🏙️") + "\n\n")
		sb.WriteString(m.cityInput.View() + "\n\n")
		sb.WriteString(hintStyle.Render("Press Enter to continue"))

	case stateUnits:
		sb.WriteString(highlightStyle.Render("Choose your preferred units: 🌡️") + "\n\n")
		sb.WriteString(m.renderOptions(m.unitOptions, m.unitCursor))
		sb.WriteString("\n" + hintStyle.Render("Press Enter to confirm"))

	case stateAuth:
		sb.WriteString(highlightStyle.Render("GitHub Auth 🔒") + "\n\n")
		sb.WriteString("To get weather data you need to authenticate with GitHub (don't worry, no permissions requested!).\n\n")
		sb.WriteString(m.renderOptions(m.authOptions, m.authCursor))
		sb.WriteString("\n" + hintStyle.Render("Press Enter to confirm your selection"))

	case stateComplete:
		sb.WriteString(highlightStyle.Render("✓ Setup complete! 🎉") + "\n\n")
		sb.WriteString(fmt.Sprintf("Default city: %s 🏙️\n", m.config.DefaultCity))
		sb.WriteString(fmt.Sprintf("Units: %s 🌡️\n", m.config.Units))

		if m.needsAuth {
			authStatus := "Authenticated ✅"
			if m.authCursor == 1 {
				authStatus = "Not authenticated ❌"
			}
			sb.WriteString(fmt.Sprintf("GitHub: %s\n", authStatus))
		}

		sb.WriteString("\n" + hintStyle.Render("Press any key to continue"))
	}

	// footer
	sb.WriteString("\n" + hintStyle.Render("↑/↓: Navigate • Enter: Select • Ctrl + C: Quit"))

	return sb.String()
}

func (m setupModel) renderOptions(options []string, cursor int) string {
	var sb strings.Builder

	for i, option := range options {
		var line string
		if cursor == i {
			line = fmt.Sprintf("%s %s", cursorStyle.Render("→"), selectedItemStyle.Render(option))
		} else {
			line = fmt.Sprintf("  %s", option)
		}
		sb.WriteString(line + "\n")
	}

	return sb.String()
}

func (m setupModel) centerContent(content string) string {
	var sb strings.Builder
	lines := strings.Split(content, "\n")

	// only try to center if we have dimensions
	if m.width > 0 && m.height > 0 {
		// vert
		contentHeight := len(lines)
		verticalPadding := (m.height - contentHeight) / 2

		if verticalPadding > 0 {
			sb.WriteString(strings.Repeat("\n", verticalPadding))
		}

		// horizontal
		for _, line := range lines {
			visibleLen := lipgloss.Width(line)
			padding := (m.width - visibleLen) / 2

			if padding > 0 {
				sb.WriteString(strings.Repeat(" ", padding))
			}
			sb.WriteString(line + "\n")
		}
	} else {
		// else just render
		sb.WriteString(content)
	}

	return sb.String()
}

type (
	authenticateMsg  struct{}
	setupCompleteMsg struct{}
)

func RunSetup(cfg *config.Config, needsAuth bool) error {
	model := NewSetupModel(cfg, needsAuth)
	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running setup UI: %w", err)
	}

	// protection against nil model
	if finalModel == nil {
		return fmt.Errorf("unexpected nil model after running UI")
	}

	// get the final model
	finalSetupModel, ok := finalModel.(setupModel)
	if !ok {
		return fmt.Errorf("unexpected model type: %T", finalModel)
	}

	// dont save if quit with q
	if finalSetupModel.quitting {
		return nil
	}

	// save config if not saved earlier
	if err := finalSetupModel.config.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	// handle auth if chosen
	if finalSetupModel.state == stateAuth && finalSetupModel.authCursor == 0 {
		fmt.Println("Starting GitHub authentication...")

		auth, err := config.Authenticate(cfg.APIURL)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		if err := config.SaveAuthConfig(auth); err != nil {
			return fmt.Errorf("failed to save authentication: %w", err)
		}

		fmt.Printf("Successfully authenticated as %s\n", auth.GithubUser)
	}

	return nil
}
