package tui

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/josephburgess/gust/internal/ui/styles"
)

// pad is the left and right page margin.
const pad = 1

func fg(c color.Color) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(c)
}

func (m Model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	return v
}

func (m Model) render() string {
	if m.width == 0 {
		return ""
	}
	if m.screen == screenSplash {
		return lipgloss.JoinVertical(lipgloss.Left, m.splashView(), m.keyBar())
	}
	header := m.header()
	tabs := m.tabs()
	rule := fg(styles.Overlay).Render(strings.Repeat("─", m.width))
	keys := m.keyBar()

	body := lipgloss.NewStyle().
		Padding(0, pad).
		Width(m.width).
		Height(m.bodyHeight()).
		MaxHeight(m.bodyHeight()).
		Render(m.body())

	return lipgloss.JoinVertical(lipgloss.Left, header, tabs, rule, body, keys)
}

// bodyHeight is the rows left for the tab after the header, tabs, rule and key bar.
func (m Model) bodyHeight() int {
	return max(m.height-4, 0)
}

func (m Model) body() string {
	switch {
	case m.showHelp:
		return "\n" + helpView()
	}
	width := m.width - 2*pad
	switch m.tab {
	case TabHourly:
		return m.hourlyView(width)
	case TabDaily:
		return m.dailyView(width)
	case TabAlerts:
		return m.alertsView(width)
	default:
		return m.nowView(width)
	}
}

// spread puts left and right at the two ends of a width-cell row.
func spread(left, right string, width int) string {
	gap := max(width-ansi.StringWidth(left)-ansi.StringWidth(right), 1)
	return left + strings.Repeat(" ", gap) + right
}

func (m Model) header() string {
	left := fg(styles.Rose).Bold(true).Render("gust") +
		fg(styles.Subtle).Render(" › ") +
		lipgloss.NewStyle().Bold(true).Render(m.cityLabel())

	now := m.now()
	if m.weather != nil {
		now = m.weather.LocalTime(now.Unix())
	}
	right := fg(styles.Subtle).Render(now.Format("Mon 2 Jan · 15:04") + " · " + m.units)

	return lipgloss.NewStyle().Padding(0, pad).Render(spread(left, right, m.width-2*pad))
}

func (m Model) tabs() string {
	parts := make([]string, len(tabNames))
	for i, name := range tabNames {
		label := fmt.Sprintf("%d %s", i+1, name)
		style := lipgloss.NewStyle().Padding(0, 1).Foreground(styles.Subtle)
		if Tab(i) == m.tab {
			style = style.Background(styles.Overlay).Foreground(styles.Rose).Bold(true)
		}
		parts[i] = style.Render(label)
		if Tab(i) == TabAlerts && m.weather != nil && len(m.weather.Alerts) > 0 {
			parts[i] += fg(styles.Love).Render(fmt.Sprintf("●%d", len(m.weather.Alerts)))
		}
	}
	return lipgloss.NewStyle().Padding(0, pad-1).Render(strings.Join(parts, ""))
}

type keyHint struct{ key, label string }

func (m Model) keyHints() []keyHint {
	if m.screen == screenSplash {
		return m.splashKeyHints()
	}
	hints := []keyHint{{"1-4", "tabs"}}
	if m.tab == TabHourly || m.tab == TabDaily {
		hints = append(hints, keyHint{"j/k", "move"})
	}
	return append(hints, keyHint{"c", "city"}, keyHint{"u", "units"}, keyHint{"r", "refresh"}, keyHint{"?", "help"}, keyHint{"q", "quit"})
}

// keyBarInfo is the text on the right of the key bar.
func (m Model) keyBarInfo() string {
	switch {
	case m.screen == screenSplash:
		return "via " + m.via
	case m.loading:
		return m.spinner.View() + " refreshing"
	case m.err != nil:
		return fg(styles.Love).Render(m.err.Error())
	}
	switch m.tab {
	case TabHourly:
		if len(m.hours()) > 0 {
			return fmt.Sprintf("%s · %d of %d", m.weather.LocalTime(m.hours()[m.hourCursor].Dt).Format("15:04"), m.hourCursor+1, len(m.hours()))
		}
	case TabDaily:
		if len(m.days()) > 0 {
			lo, hi := m.weekRange()
			return fmt.Sprintf("week %d°–%d°", deg(lo), deg(hi))
		}
	}
	if m.tab == TabNow && !m.fetchedAt.IsZero() {
		mins := int(m.now().Sub(m.fetchedAt) / time.Minute)
		if mins < 1 {
			return "updated just now"
		}
		return fmt.Sprintf("updated %dm ago", mins)
	}
	return ""
}

func (m Model) keyBar() string {
	bar := lipgloss.NewStyle().Background(styles.Surface)
	if m.prompting {
		m.cityInput.SetWidth(m.width - 2*pad - 4)
		row := bar.Render(strings.Repeat(" ", pad)) + m.cityInput.View()
		return row + bar.Render(strings.Repeat(" ", max(m.width-ansi.StringWidth(row), 0)))
	}
	hints := make([]string, 0, len(m.keyHints()))
	for _, h := range m.keyHints() {
		hints = append(hints, bar.Foreground(styles.Gold).Render(h.key)+bar.Render(" ")+bar.Foreground(styles.Subtle).Render(h.label))
	}
	left := bar.Render(strings.Repeat(" ", pad)) + strings.Join(hints, bar.Render("  "))
	right := bar.Foreground(styles.Subtle).Render(m.keyBarInfo()) + bar.Render(strings.Repeat(" ", pad))

	gap := max(m.width-ansi.StringWidth(left)-ansi.StringWidth(right), 0)
	return left + bar.Render(strings.Repeat(" ", gap)) + right
}
