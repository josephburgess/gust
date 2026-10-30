package renderer

import (
	"charm.land/lipgloss/v2"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/tui"
)

func (r *TerminalRenderer) RenderCompactWeather(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	lipgloss.Println(tui.Compact(city, weather, r.Units))
}
