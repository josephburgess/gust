package renderer

import (
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/pretty"
	"github.com/josephburgess/gust/internal/ui/styles"
)

type PrettyRenderer struct {
	Units string
}

func NewPrettyRenderer(units string) WeatherRenderer {
	return &PrettyRenderer{Units: units}
}

func (r *PrettyRenderer) RenderCurrentWeather(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	r.run(city, weather, cfg, pretty.TabCurrent)
}

func (r *PrettyRenderer) RenderHourlyForecast(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	r.run(city, weather, cfg, pretty.TabHourly)
}

func (r *PrettyRenderer) RenderDailyForecast(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	r.run(city, weather, cfg, pretty.TabDaily)
}

func (r *PrettyRenderer) RenderAlerts(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	r.run(city, weather, cfg, pretty.TabAlerts)
}

func (r *PrettyRenderer) RenderFullWeather(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	r.run(city, weather, cfg, pretty.TabCurrent)
}

func (r *PrettyRenderer) RenderCompactWeather(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	r.run(city, weather, cfg, pretty.TabCurrent)
}

func (r *PrettyRenderer) run(city *models.City, weather *models.OneCallResponse, cfg *config.Config, startTab int) {
	if err := pretty.Run(city, weather, cfg, r.Units, startTab); err != nil {
		styles.ExitWithError("pretty mode error", err)
	}
}
