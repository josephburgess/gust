package pretty

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/josephburgess/gust/internal/ui/styles"
)

func (m Model) View() string {
	if !m.ready {
		return "\n  Loading..."
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		m.tabBar(),
		m.viewport.View(),
		m.footer(),
	)
}

func (m Model) tabBar() string {
	tabs := make([]string, len(tabNames))
	for i, name := range tabNames {
		label := name
		if i == TabAlerts && len(m.weather.Alerts) > 0 {
			label = fmt.Sprintf("%s (%d)", name, len(m.weather.Alerts))
		}
		switch {
		case i == m.activeTab:
			tabs[i] = activeTabStyle.Render(label)
		case i == TabAlerts && len(m.weather.Alerts) > 0:
			tabs[i] = lipgloss.NewStyle().
				Bold(true).
				Foreground(styles.Love).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(styles.Love).
				Padding(0, 1).
				Render(label)
		default:
			tabs[i] = inactiveTabStyle.Render(label)
		}
	}

	bar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	// fill the rest of the row
	gap := m.width - lipgloss.Width(bar)
	if gap > 0 {
		bar += dimTabStyle.Render(strings.Repeat(" ", gap))
	}
	return bar
}

func (m Model) footer() string {
	hint := footerStyle.Render("  tab/→ next  ←/shift+tab prev  1-4 jump  j/k scroll  q quit")
	city := ""
	if m.city != nil {
		city = footerStyle.Render(m.city.Name)
	}
	gap := m.width - lipgloss.Width(hint) - lipgloss.Width(city) - 2
	if gap < 0 {
		gap = 0
	}
	return hint + strings.Repeat(" ", gap) + city
}
