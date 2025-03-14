package setup

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func ClearViewport(width, height int) string {
	var sb strings.Builder

	blankLine := strings.Repeat(" ", width)
	for i := 0; i < height; i++ {
		sb.WriteString(blankLine + "\n")
	}

	sb.WriteString("\033[H")

	return sb.String()
}

// rendr current state of setup wizard
func (m Model) View() string {
	var result strings.Builder

	if m.Width > 0 && m.Height > 0 {
		result.WriteString(ClearViewport(m.Width, m.Height))
	}

	content := m.buildContent()
	result.WriteString(m.centerContent(content))

	return result.String()
}

// create correct content for currnet state
func (m Model) buildContent() string {
	var sb strings.Builder

	sb.WriteString(logoBoxStyle.Render(asciiLogo) + "\n\n")
	sb.WriteString(subtitleStyle.Render("Simple terminal weather 🌤️") + "\n\n")

	switch m.State {
	case StateCity:
		sb.WriteString(highlightStyle.Render("Enter a default city 🏙️") + "\n\n")
		sb.WriteString(m.CityInput.View() + "\n\n")
		sb.WriteString(hintStyle.Render("You can enter a country code too, but use a comma! (e.g. London,GB)"))

	case StateCitySearch:
		sb.WriteString(highlightStyle.Render("Searching for cities...") + "\n\n")
		sb.WriteString(fmt.Sprintf("%s Looking for \"%s\"", m.Spinner.View(), m.CitySearchQuery))
		sb.WriteString("\n\n")

	case StateCitySelect:
		sb.WriteString(highlightStyle.Render("Select your town or city: 🏙️") + "\n\n")

		if len(m.CityOptions) == 0 {
			sb.WriteString("No cities found. Please try a different search term.\n\n")
		} else {
			for i, city := range m.CityOptions {
				var line string
				var locationInfo string

				if city.State != "" && city.Country != "" {
					flag := getCountryEmoji(city.Country)
					locationInfo = fmt.Sprintf("%s, %s %s", city.State, city.Country, flag)
				} else if city.Country != "" {
					flag := getCountryEmoji(city.Country)
					locationInfo = fmt.Sprintf("%s %s", flag, locationInfo)
				} else {
					locationInfo = fmt.Sprintf("(%.4f, %.4f)", city.Lat, city.Lon)
				}

				displayName := fmt.Sprintf("%s - %s", city.Name, locationInfo)

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
		if m.Config.ShowTips {
			sb.WriteString("Tips enabled 💡\n")
		} else {
			sb.WriteString("Tips disabled 💡\n")
		}

	case StateTips:
		sb.WriteString(highlightStyle.Render("Would you like tips shown on daily forecasts? 💡") + "\n\n")
		sb.WriteString(m.renderOptions(m.TipOptions, m.TipCursor))
		sb.WriteString("\n" + hintStyle.Render("Press Enter to confirm"))

		if m.NeedsAuth {
			authStatus := "Authenticated ✅"
			if m.AuthCursor == 1 {
				authStatus = "Not authenticated ❌"
			}
			sb.WriteString(fmt.Sprintf("GitHub: %s\n", authStatus))
		}

	case StateApiKeyOption:
		sb.WriteString(highlightStyle.Render("Choose auth method: 🔑") + "\n\n")
		sb.WriteString(m.renderOptions(m.ApiKeyOptions, m.ApiKeyCursor))
		sb.WriteString("\n" + hintStyle.Render("Press Enter to confirm"))

	case StateApiKeyInput:
		sb.WriteString(highlightStyle.Render("Enter your OpenWeatherMap API key: 🔑") + "\n\n")
		sb.WriteString(m.ApiKeyInput.View() + "\n\n")
		sb.WriteString(hintStyle.Render("Get your API key from https://home.openweathermap.org/subscriptions/unauth_subscribe/onecall_30/base"))

	}

	// footer
	sb.WriteString("\n" + hintStyle.Render("↓j/↑k Navigate • Enter: Select • Ctrl + C: Quit"))

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

func getCountryEmoji(countryCode string) string {
	if countryCode == "" {
		return "🌍"
	}

	if len(countryCode) != 2 {
		return "🌍"
	}

	cc := strings.ToUpper(countryCode)
	const offset = 127397
	firstLetter := rune(cc[0]) + offset
	secondLetter := rune(cc[1]) + offset
	flag := string(firstLetter) + string(secondLetter)

	return flag
}

// Add helper functions to work with country names and emojis

// GetCountryEmojiByName returns the flag emoji for a given country name
// It converts the name to lowercase for case-insensitive matching
