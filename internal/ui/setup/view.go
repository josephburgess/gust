package setup

import (
	"fmt"
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/ui/styles"
	"github.com/josephburgess/gust/internal/ui/widgets"
)

// Column widths for the centred content, in cells.
const (
	columnWidth      = 80
	preferencesWidth = 110
)

var steps = []string{"city", "preferences", "sign in"}

type keyHint struct{ key, label string }

func fg(c color.Color) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(c)
}

// rendr current state of setup wizard
func (m Model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	return v
}

func (m Model) step() int {
	switch m.State {
	case StatePreferences:
		return 1
	case StateApiKeyOption, StateApiKeyInput:
		return 2
	case StateComplete:
		return 3
	default:
		return 0
	}
}

func (m Model) width() int {
	want := columnWidth
	if m.State == StatePreferences {
		want = preferencesWidth
	}
	if m.Width == 0 {
		return want
	}
	return min(want, m.Width-4)
}

func (m Model) render() string {
	width := m.width()
	var heading, sub, content string
	var hints []keyHint

	switch m.State {
	case StatePreferences:
		heading, sub = "How do you like it?", "Pick your units, the tab gust opens on, and tips."
		content = m.preferencesView(width)
		hints = []keyHint{{"tab", "group"}, {"↑↓", "choose"}, {"enter", "next"}, {"esc", "back"}}
	case StateApiKeyOption, StateApiKeyInput:
		heading, sub = "Last bit: how should gust get weather?", "gust talks to breeze, a small proxy, so most people never need a key."
		content = m.signInView(width)
		hints = []keyHint{{"↑↓", "choose"}, {"enter", "select"}, {"esc", "back"}}
	case StateComplete:
		content = m.completeView(width)
		hints = []keyHint{{"enter", "show me the weather"}, {"esc", "back"}}
	default:
		heading, sub = "Where's home?", "Your default city. Change it any time with c or s."
		content = m.cityView(width)
		hints = []keyHint{{"enter", "search / pick"}, {"↑↓", "choose"}, {"ctrl+c", "quit"}}
	}

	header := lipgloss.NewStyle().Width(width).Render(spread(
		fg(styles.Rose).Bold(true).Render("gust")+fg(styles.Subtle).Render(" setup"),
		widgets.StepRail(steps, m.step()),
		width,
	))
	parts := []string{header, ""}
	if heading != "" {
		parts = append(parts, lipgloss.NewStyle().Bold(true).Render(heading), fg(styles.Subtle).Render(sub), "")
	}
	column := lipgloss.JoinVertical(lipgloss.Left, append(parts, content)...)

	if m.Width == 0 || m.Height == 0 {
		return column
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.Place(m.Width, m.Height-1, lipgloss.Center, lipgloss.Center, column),
		m.keyBar(hints),
	)
}

// spread puts left and right at the two ends of a width-cell row.
func spread(left, right string, width int) string {
	gap := max(width-ansi.StringWidth(left)-ansi.StringWidth(right), 1)
	return left + strings.Repeat(" ", gap) + right
}

func (m Model) keyBar(hints []keyHint) string {
	bar := lipgloss.NewStyle().Background(styles.Surface)
	parts := make([]string, len(hints))
	for i, h := range hints {
		parts[i] = bar.Foreground(styles.Gold).Render(h.key) + bar.Render(" ") + bar.Foreground(styles.Subtle).Render(h.label)
	}
	left := bar.Render(" ") + strings.Join(parts, bar.Render("  "))

	info := "done"
	if m.step() < len(steps) {
		info = fmt.Sprintf("step %d of %d", m.step()+1, len(steps))
	}
	right := bar.Foreground(styles.Subtle).Render(info) + bar.Render(" ")

	gap := max(m.Width-ansi.StringWidth(left)-ansi.StringWidth(right), 0)
	return left + bar.Render(strings.Repeat(" ", gap)) + right
}

func (m Model) cityView(width int) string {
	m.CityInput.SetWidth(width - 6)
	input := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Iris).
		Padding(0, 1).
		Width(width).
		Render(m.CityInput.View())

	parts := []string{input, ""}
	switch {
	case m.State == StateCitySearch:
		parts = append(parts, m.Spinner.View()+" "+fg(styles.Subtle).Render(fmt.Sprintf("searching for %q…", m.CitySearchQuery)))
	case m.SearchErr != "":
		parts = append(parts, fg(styles.Love).Render(m.SearchErr))
	case m.State == StateCitySelect:
		parts = append(parts, fg(styles.Subtle).Render(fmt.Sprintf("%d MATCHES", len(m.CityOptions))))
		for i := range m.CityOptions {
			parts = append(parts, m.cityRow(i, width))
		}
	}
	parts = append(parts, "", fg(styles.Subtle).Italic(true).Render("not there? add a country code, like london,gb"))
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m Model) cityRow(i, width int) string {
	city := m.CityOptions[i]
	selected := i == m.CityCursor
	bg := func(s lipgloss.Style) lipgloss.Style {
		if selected {
			return s.Background(styles.Overlay)
		}
		return s
	}

	cursor, nameColor := "  ", styles.Text
	if selected {
		cursor, nameColor = "▸ ", styles.Foam
	}
	place := city.Country
	if city.State != "" {
		place = city.State + ", " + city.Country
	}

	left := bg(fg(styles.Love)).Render(cursor) +
		bg(fg(nameColor).Bold(true)).Render(city.Name) +
		bg(lipgloss.NewStyle()).Render("  ") +
		bg(fg(styles.Subtle)).Render(place)
	right := bg(fg(styles.Subtle)).Render(fmt.Sprintf("%.2f, %.2f", city.Lat, city.Lon))
	gap := max(width-ansi.StringWidth(left)-ansi.StringWidth(right), 1)
	return left + bg(lipgloss.NewStyle()).Render(strings.Repeat(" ", gap)) + right
}

func radioRows(options []option, cursor int) string {
	rows := make([]string, len(options))
	for i, o := range options {
		if i == cursor {
			rows[i] = fg(styles.Foam).Render("● "+o.label) + "  " + fg(styles.Subtle).Render(o.detail)
		} else {
			rows[i] = fg(styles.Subtle).Render("○ " + o.label + "  " + o.detail)
		}
	}
	return strings.Join(rows, "\n")
}

func (m Model) preferencesView(width int) string {
	left := 44
	groups := lipgloss.JoinVertical(lipgloss.Left,
		widgets.Panel("units", radioRows(unitOptions, m.UnitCursor), left, m.PrefGroup == GroupUnits),
		widgets.Panel("opens on", radioRows(viewOptions, m.ViewCursor), left, m.PrefGroup == GroupView),
		widgets.Panel("tips", radioRows(tipOptions, m.TipCursor), left, m.PrefGroup == GroupTips),
	)
	return lipgloss.JoinHorizontal(lipgloss.Top, groups, "  ", m.preview(width-left-2))
}

// preview is a fake Now line in the chosen units.
func (m Model) preview(width int) string {
	temp, wind := "9°C", "24 km/h"
	switch unitOptions[m.UnitCursor].value {
	case "imperial":
		temp, wind = "49°F", "15 mph"
	case "standard":
		temp, wind = "283K", "24 km/h"
	}
	city := m.Config.DefaultCity
	if city == "" {
		city = "London"
	}

	surface := lipgloss.NewStyle().Background(styles.Surface)
	lines := []string{
		fg(styles.Subtle).Render("PREVIEW"),
		"",
		fg(styles.Rose).Bold(true).Render(city) + "  " + fg(styles.Gold).Bold(true).Render(temp) + " " +
			fg(styles.Foam).Render("light rain") + fg(styles.Subtle).Render(" · ↗ "+wind+" SW"),
		fg(styles.Subtle).Render("opens on ") + viewOptions[m.ViewCursor].label,
	}
	if tipOptions[m.TipCursor].value == "on" {
		lines = append(lines, "", fg(styles.Iris).Italic(true).Render("› don't go out without an umbrella!"))
	}
	return surface.Padding(1, 2).Width(width).Render(strings.Join(lines, "\n"))
}

func (m Model) signInView(width int) string {
	github := lipgloss.NewStyle().Bold(true).Render("Sign in with GitHub") + "  " + fg(styles.Gold).Render("recommended") + "\n" +
		fg(styles.Subtle).Render("breeze gives gust a key. No GitHub permissions are asked for.")

	key := lipgloss.NewStyle().Bold(true).Render("Use my own OpenWeatherMap key") + "\n" +
		fg(styles.Subtle).Render("Needs a One Call 3.0 subscription.")
	if m.State == StateApiKeyInput {
		m.ApiKeyInput.SetWidth(width - 8)
		key += "\n\n" + fg(styles.Subtle).Render("API KEY") + "\n" + m.ApiKeyInput.View() + "\n\n" +
			fg(styles.Subtle).Render("https://openweathermap.org/api/one-call-3") + "\n" +
			fg(styles.Subtle).Render("saved to ~/.config/gust/auth.json")
		if m.SearchErr != "" {
			key += "\n" + fg(styles.Love).Render(m.SearchErr)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		widgets.Panel("", github, width, m.ApiKeyCursor == 0),
		widgets.Panel("", key, width, m.ApiKeyCursor == 1),
	)
}

func (m Model) weatherVia() string {
	if m.ApiKeyCursor == 1 {
		key := m.ApiKeyInput.Value()
		return "your key ••••" + key[max(len(key)-4, 0):]
	}
	return "GitHub"
}

func (m Model) completeView(width int) string {
	path, _ := config.GetConfigPath()
	tips := "off"
	if m.Config.ShowTips {
		tips = "on"
	}

	row := func(label, value string) string {
		return fg(styles.Subtle).Width(14).Render(label) + value
	}
	summary := strings.Join([]string{
		row("city", m.Config.DefaultCity),
		row("units", m.Config.Units),
		row("opens on", viewOptions[optionIndex(viewOptions, m.Config.DefaultView)].label),
		row("tips", tips),
		row("weather via", m.weatherVia()),
	}, "\n")

	cell := func(k, label string) string {
		return lipgloss.NewStyle().Width((width - 4) / 3).Render(fg(styles.Gold).Render(k) + " " + fg(styles.Subtle).Render(label))
	}
	handy := widgets.Panel("handy keys", lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, cell("1-4", "tabs"), cell("c", "city"), cell("u", "units")),
		lipgloss.JoinHorizontal(lipgloss.Top, cell("r", "refresh"), cell("s", "settings"), cell("?", "help")),
	), width, false)

	cta := lipgloss.NewStyle().Background(styles.Gold).Foreground(styles.Base).Bold(true).Padding(0, 2).
		Render("enter  show me the weather")

	return lipgloss.JoinVertical(lipgloss.Left,
		fg(styles.Foam).Bold(true).Render("✓ ")+lipgloss.NewStyle().Bold(true).Render("You're set."),
		fg(styles.Subtle).Render("saved to "+path),
		"",
		summary,
		"",
		handy,
		"",
		cta,
	)
}
