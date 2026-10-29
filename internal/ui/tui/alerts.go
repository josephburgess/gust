package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/josephburgess/gust/internal/ui/styles"
)

func (m Model) alertsView(width int) string {
	if len(m.weather.Alerts) == 0 {
		return "\n" + fg(styles.Subtle).Render("no alerts for "+m.cityLabel())
	}

	panels := []string{""}
	for _, a := range m.weather.Alerts {
		when := m.weather.LocalTime(a.Start).Format("Mon 2 Jan 15:04") + " – " +
			m.weather.LocalTime(a.End).Format("Mon 2 Jan 15:04")
		body := strings.Join([]string{
			lipgloss.NewStyle().Bold(true).Render(a.Event),
			fg(styles.Subtle).Render(a.SenderName + " · " + when),
			"",
			lipgloss.NewStyle().Width(width - 4).Render(strings.TrimSpace(a.Description)),
		}, "\n")
		panels = append(panels, lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Love).
			Padding(0, 1).
			Width(width).
			Render(body))
	}
	return strings.Join(panels, "\n")
}
