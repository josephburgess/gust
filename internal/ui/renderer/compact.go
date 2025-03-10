package renderer

import (
	"fmt"
	"strings"
	"time"

	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
)

func (r *TerminalRenderer) RenderCompactWeather(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	current := weather.Current

	fmt.Print(styles.FormatHeader(fmt.Sprintf("%s WEATHER", strings.ToUpper(city.Name))))

	if len(current.Weather) > 0 {
		weatherCond := current.Weather[0]
		tempUnit := r.GetTemperatureUnit()
		emoji := models.GetWeatherEmoji(weatherCond.ID, &current)
		temp := styles.TempStyle(fmt.Sprintf("%.1f%s", current.Temp, tempUnit))
		feels := fmt.Sprintf("(%.1f%s)", current.FeelsLike, tempUnit)

		fmt.Printf("%s %-16s    %s %s\n",
			emoji,
			styles.HighlightStyleF(weatherCond.Description),
			temp,
			feels)

		windUnit := r.GetWindSpeedUnit()
		windSpeed := r.FormatWindSpeed(current.WindSpeed)
		windDir := models.GetWindDirection(current.WindDeg)

		fmt.Printf("💧 %-3d%%         💨 %-4.1f %-3s %-2s",
			current.Humidity,
			windSpeed,
			windUnit,
			windDir)

		if current.Rain != nil && current.Rain.OneHour > 0 {
			fmt.Printf("     🌧️ %.1f mm", current.Rain.OneHour)
		}
		if current.Snow != nil && current.Snow.OneHour > 0 {
			fmt.Printf("     ❄️ %.1f mm", current.Snow.OneHour)
		}
		fmt.Println()

		sunrise := time.Unix(current.Sunrise, 0).Format("15:04")
		sunset := time.Unix(current.Sunset, 0).Format("15:04")
		fmt.Printf("🌅 %-8s     🌇 %-8s", sunrise, sunset)

		if len(weather.Alerts) > 0 {
			fmt.Printf("     %s",
				styles.AlertStyle(fmt.Sprintf("⚠️ %d alerts", len(weather.Alerts))))
		}
		fmt.Println()
		r.displayWeatherTip(weather, cfg)
	}
}
