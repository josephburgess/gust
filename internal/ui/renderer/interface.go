package renderer

import (
	"github.com/josephburgess/gust/internal/models"
)

type WeatherRenderer interface {
	RenderCurrentWeather(city *models.City, weather *models.OneCallResponse)
	RenderDailyForecast(city *models.City, weather *models.OneCallResponse)
	RenderHourlyForecast(city *models.City, weather *models.OneCallResponse)
	RenderAlerts(city *models.City, weather *models.OneCallResponse)
	RenderFullWeather(city *models.City, weather *models.OneCallResponse)
	RenderCompactWeather(city *models.City, weather *models.OneCallResponse)
}

func NewWeatherRenderer(rendererType string, units string) WeatherRenderer {
	return NewTerminalRenderer(units)
}
