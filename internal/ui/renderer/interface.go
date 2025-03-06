package renderer

import (
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
)

type WeatherRenderer interface {
	RenderCurrentWeather(city *models.City, weather *models.OneCallResponse, cfg *config.Config)
	RenderDailyForecast(city *models.City, weather *models.OneCallResponse, cfg *config.Config)
	RenderHourlyForecast(city *models.City, weather *models.OneCallResponse, cfg *config.Config)
	RenderAlerts(city *models.City, weather *models.OneCallResponse, cfg *config.Config)
	RenderFullWeather(city *models.City, weather *models.OneCallResponse, cfg *config.Config)
	RenderCompactWeather(city *models.City, weather *models.OneCallResponse, cfg *config.Config)
}

func NewWeatherRenderer(rendererType string, units string) WeatherRenderer {
	return NewTerminalRenderer(units)
}
