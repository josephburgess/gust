package pretty

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
)

func Run(city *models.City, weather *models.OneCallResponse, cfg *config.Config, units string, startTab int) error {
	m := New(city, weather, cfg, units, startTab)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
