package renderer

import (
	"fmt"

	"github.com/josephburgess/gust/internal/models"
)

func (r *TerminalRenderer) RenderFullWeather(city *models.City, weather *models.OneCallResponse) {
	r.RenderCurrentWeather(city, weather)
	fmt.Println()

	if len(weather.Alerts) > 0 {
		r.RenderAlerts(city, weather)
		fmt.Println()
	}

	r.RenderDailyForecast(city, weather)
	fmt.Println()
}
