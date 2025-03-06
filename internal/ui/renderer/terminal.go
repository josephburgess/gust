package renderer

import (
	"fmt"

	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
)

type TerminalRenderer struct {
	BaseRenderer
}

func NewTerminalRenderer(units string) *TerminalRenderer {
	return &TerminalRenderer{
		BaseRenderer: BaseRenderer{
			Units: units,
		},
	}
}

func (r *TerminalRenderer) displayWindInfo(speed float64, deg int, gust float64) {
	windUnit := r.GetWindSpeedUnit()
	windSpeed := r.FormatWindSpeed(speed)

	if gust > 0 {
		gustSpeed := r.FormatWindSpeed(gust)
		fmt.Printf("Wind: %.1f %s %s %s (Gusts: %.1f %s)\n",
			windSpeed,
			windUnit,
			models.GetWindDirection(deg),
			"💨",
			gustSpeed,
			windUnit)
	} else {
		fmt.Printf("Wind: %.1f %s %s %s\n",
			windSpeed,
			windUnit,
			models.GetWindDirection(deg),
			"💨")
	}
}

func (r *TerminalRenderer) displayPrecipitation(rain *models.RainData, snow *models.SnowData) {
	if rain != nil && rain.OneHour > 0 {
		fmt.Printf("Rain: %.1f mm (last hour) 🌧️\n", rain.OneHour)
	}

	if snow != nil && snow.OneHour > 0 {
		fmt.Printf("Snow: %.1f mm (last hour) ❄️\n", snow.OneHour)
	}
}

func (r *TerminalRenderer) displayAlertSummary(alerts []models.Alert, cityName string) {
	if len(alerts) > 0 {
		fmt.Printf("%s Use 'gust --alerts %s' to view them.\n",
			styles.AlertStyle(fmt.Sprintf("⚠️  There are %d weather alerts for this area.", len(alerts))),
			cityName)
	}
}

func (r *TerminalRenderer) displayWeatherTip(weather *models.OneCallResponse, cfg *config.Config) {
	if !cfg.ShowTips {
		return
	}
	tip := models.GetWeatherTip(weather, r.Units)
	fmt.Printf("\n%s\n", styles.TipStyle(fmt.Sprintf("💡 %s", tip)))
}
