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
	header := m.header()
	tabs := m.tabs()
	rule := fg(styles.Overlay).Render(strings.Repeat("─", m.width))
	keys := m.keyBar()

	bodyHeight := max(m.height-lipgloss.Height(header)-lipgloss.Height(tabs)-2, 0)
	body := lipgloss.NewStyle().
		Padding(0, pad).
		Width(m.width).
		Height(bodyHeight).
		MaxHeight(bodyHeight).
		Render(m.body())

	return lipgloss.JoinVertical(lipgloss.Left, header, tabs, rule, body, keys)
}

func (m Model) body() string {
	switch {
	case m.err != nil:
		return "\n" + fg(styles.Love).Render(m.err.Error())
	case m.weather == nil:
		return "\n" + m.spinner.View() + " " + fg(styles.Subtle).Render("fetching "+m.cityLabel()+"…")
	}
	switch m.tab {
	case TabNow:
		return m.nowView(m.width - 2*pad)
	default:
		return "\n" + fg(styles.Subtle).Render(tabNames[m.tab]+" is on the way")
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
	return []keyHint{{"1-4", "tabs"}, {"q", "quit"}}
}

// keyBarInfo is the text on the right of the key bar.
func (m Model) keyBarInfo() string {
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
	hints := make([]string, 0, len(m.keyHints()))
	for _, h := range m.keyHints() {
		hints = append(hints, bar.Foreground(styles.Gold).Render(h.key)+bar.Render(" ")+bar.Foreground(styles.Subtle).Render(h.label))
	}
	left := bar.Render(strings.Repeat(" ", pad)) + strings.Join(hints, bar.Render("  "))
	right := bar.Foreground(styles.Subtle).Render(m.keyBarInfo()) + bar.Render(strings.Repeat(" ", pad))

	gap := max(m.width-ansi.StringWidth(left)-ansi.StringWidth(right), 0)
	return left + bar.Render(strings.Repeat(" ", gap)) + right
}
