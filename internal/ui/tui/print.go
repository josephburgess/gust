package tui

import (
	"charm.land/lipgloss/v2"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
)

// printHeight is tall enough that print mode never scrolls or hides a panel.
const printHeight = 1000

// Print renders one tab for print mode: the header line and the tab, without
// the tabs row or the key bar.
func Print(city *models.City, w *models.OneCallResponse, cfg *config.Config, tab Tab, width int) string {
	m := NewModel(city.Name, cfg, tab, nil)
	m.screen = screenMain
	m.location, m.weather = city, w
	m.loading = false
	m.width, m.height = width, printHeight

	body := lipgloss.NewStyle().Padding(0, pad).Render(m.body())
	return lipgloss.JoinVertical(lipgloss.Left, m.header(), body)
}
