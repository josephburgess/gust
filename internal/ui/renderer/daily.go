package renderer

import (
	"fmt"
	"strings"
	"time"

	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
)

func (r *TerminalRenderer) RenderDailyForecast(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	fmt.Print(styles.FormatHeader(fmt.Sprintf("5-DAY FORECAST FOR %s", strings.ToUpper(city.Name))))

	if len(weather.Daily) > 0 {
		tempUnit := r.GetTemperatureUnit()

		for i, day := range weather.Daily {
			if i >= 5 {
				break
			}

			date := time.Unix(day.Dt, 0).Format("Mon Jan 2")

			if i > 0 {
				fmt.Println()
			}

			fmt.Printf("%s: %s\n",
				styles.HighlightStyleF(date),
				day.Summary)

			fmt.Printf("  High/Low: %s/%s %s\n",
				styles.TempStyle(fmt.Sprintf("%.1f%s", day.Temp.Max, tempUnit)),
				styles.TempStyle(fmt.Sprintf("%.1f%s", day.Temp.Min, tempUnit)),
				"🌡️")

			fmt.Printf("  Morning: %.1f%s  Day: %.1f%s  Evening: %.1f%s  Night: %.1f%s\n",
				day.Temp.Morn, tempUnit,
				day.Temp.Day, tempUnit,
				day.Temp.Eve, tempUnit,
				day.Temp.Night, tempUnit)

			if len(day.Weather) > 0 {
				weather := day.Weather[0]
				condition := fmt.Sprintf("%s %s", weather.Description, models.GetWeatherEmoji(weather.ID, nil))
				fmt.Printf("  Conditions: %s\n", styles.InfoStyle(condition))
			}

			if day.Pop > 0 {
				fmt.Printf("  Precipitation: %d%% chance\n", int(day.Pop*100))
			}

			if day.Rain > 0 {
				fmt.Printf("  Rain: %.1f mm 🌧️\n", day.Rain)
			}

			if day.Snow > 0 {
				fmt.Printf("  Snow: %.1f mm ❄️\n", day.Snow)
			}

			windUnit := r.GetWindSpeedUnit()
			windSpeed := r.FormatWindSpeed(day.WindSpeed)

			fmt.Printf("  Wind: %.1f %s %s\n",
				windSpeed,
				windUnit,
				models.GetWindDirection(day.WindDeg))

			fmt.Printf("  UV Index: %.1f\n", day.UVI)
		}
		fmt.Println()
	}
}
