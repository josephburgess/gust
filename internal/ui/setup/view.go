package setup

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// rendr current state of setup wizard
func (m Model) View() string {
	content := m.buildContent()
	return m.centerContent(content)
}

// create correct content for currnet state
func (m Model) buildContent() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render(asciiLogo) + "\n\n")
	sb.WriteString(boxStyle.Render(subtitleStyle.Render("Simple terminal weather 🌤️")) + "\n\n")

	switch m.State {
	case StateCity:
		sb.WriteString(highlightStyle.Render("Enter a default city 🏙️") + "\n\n")
		sb.WriteString(m.CityInput.View() + "\n\n")
		sb.WriteString(hintStyle.Render("Press Enter to continue"))

	case StateCitySearch:
		sb.WriteString(highlightStyle.Render("Searching for cities...") + "\n\n")
		sb.WriteString(fmt.Sprintf("%s Looking for \"%s\"", m.Spinner.View(), m.CitySearchQuery))
		sb.WriteString("\n\n")

	case StateCitySelect:
		sb.WriteString(highlightStyle.Render("Select your city: 🏙️") + "\n\n")

		if len(m.CityOptions) == 0 {
			sb.WriteString("No cities found. Please try a different search term.\n\n")
		} else {
			for i, city := range m.CityOptions {
				var line string
				displayName := fmt.Sprintf("%s (%f, %f)", city.Name, city.Lat, city.Lon)

				if m.CityCursor == i {
					line = fmt.Sprintf("%s %s", cursorStyle.Render("→"), selectedItemStyle.Render(displayName))
				} else {
					line = fmt.Sprintf("  %s", displayName)
				}
				sb.WriteString(line + "\n")
			}
			sb.WriteString("\n")
		}

		sb.WriteString(hintStyle.Render("Press Enter to select or Esc to search again"))

	case StateUnits:
		sb.WriteString(highlightStyle.Render("Choose your preferred units: 🌡️") + "\n\n")
		sb.WriteString(m.renderOptions(m.UnitOptions, m.UnitCursor))
		sb.WriteString("\n" + hintStyle.Render("Press Enter to confirm"))

	case StateView:
		sb.WriteString(highlightStyle.Render("Choose your preferred view: 📊") + "\n\n")
		sb.WriteString(m.renderOptions(m.ViewOptions, m.ViewCursor))
		sb.WriteString("\n" + hintStyle.Render("Press Enter to confirm"))

	case StateAuth:
		sb.WriteString(highlightStyle.Render("GitHub Auth 🔒") + "\n\n")
		sb.WriteString("To get weather data you need to authenticate with GitHub (don't worry, no permissions requested!).\n\n")
		sb.WriteString(m.renderOptions(m.AuthOptions, m.AuthCursor))
		sb.WriteString("\n" + hintStyle.Render("Press Enter to confirm your selection"))

	case StateComplete:
		sb.WriteString(highlightStyle.Render("✓ Setup complete! 🎉") + "\n\n")
		sb.WriteString(fmt.Sprintf("Default city: %s 🏙️\n", m.Config.DefaultCity))
		sb.WriteString(fmt.Sprintf("Units: %s 🌡️\n", m.Config.Units))
		sb.WriteString(fmt.Sprintf("Default view: %s 📊\n", m.Config.DefaultView))

		if m.NeedsAuth {
			authStatus := "Authenticated ✅"
			if m.AuthCursor == 1 {
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

// renders a list of opts with current selection highlighted
func (m Model) renderOptions(options []string, cursor int) string {
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

// center in viewport
func (m Model) centerContent(content string) string {
	var sb strings.Builder
	lines := strings.Split(content, "\n")

	// only try to center if we have dimensions
	if m.Width > 0 && m.Height > 0 {
		// vert
		contentHeight := len(lines)
		verticalPadding := (m.Height - contentHeight) / 2

		if verticalPadding > 0 {
			sb.WriteString(strings.Repeat("\n", verticalPadding))
		}

		// horiz
		for _, line := range lines {
			visibleLen := lipgloss.Width(line)
			padding := (m.Width - visibleLen) / 2

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
