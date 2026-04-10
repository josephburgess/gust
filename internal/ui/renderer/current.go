package renderer

import (
	"fmt"
	"time"

	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
)

func (r *TerminalRenderer) RenderCurrentWeather(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	current := weather.Current

	fmt.Print(cityHeader("WEATHER FOR", city.Name))

	if len(current.Weather) > 0 {
		weatherCond := current.Weather[0]

		fmt.Printf("Current Conditions: %s %s\n\n",
			styles.HighlightStyleF(weatherCond.Description),
			models.GetWeatherEmoji(weatherCond.ID, &current))

		tempUnit := r.GetTemperatureUnit()

		fmt.Printf("Temperature: %s %s (Feels like: %.1f%s)\n",
			styles.TempStyle(fmt.Sprintf("%.1f%s", current.Temp, tempUnit)),
			"🌡️",
			current.FeelsLike, tempUnit)

		fmt.Printf("Humidity: %d%% %s\n", current.Humidity, "💧")
		if current.UVI > 0 {
			fmt.Printf("UV Index: %.1f ☀️\n", current.UVI)
		}

		r.displayWindInfo(current.WindSpeed, current.WindDeg, current.WindGust)

		if current.Clouds > 0 {
			fmt.Printf("Cloud coverage: %d%% ☁️\n", current.Clouds)
		}

		r.displayPrecipitation(current.Rain, current.Snow)
		fmt.Printf("Visibility: %s\n", models.VisibilityToString(current.Visibility))

		fmt.Printf("Sunrise: %s %s  Sunset: %s %s\n",
			time.Unix(current.Sunrise, 0).Format("15:04"),
			"🌅",
			time.Unix(current.Sunset, 0).Format("15:04"),
			"🌇")
		r.displayWeatherTip(weather, cfg)
		fmt.Printf("\n")
	}
	r.displayAlertSummary(weather.Alerts, city.Name)
}
