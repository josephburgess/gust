package renderer

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
)

func (r *TerminalRenderer) RenderHourlyForecast(city *models.City, weather *models.OneCallResponse, config *config.Config) {
	fmt.Print(styles.FormatHeader(fmt.Sprintf("24H FORECAST FOR %s", strings.ToUpper(city.Name))))

	if len(weather.Hourly) > 0 {
		hourLimit := int(math.Min(24, float64(len(weather.Hourly))))
		currentDay := ""

		tempUnit := r.GetTemperatureUnit()

		for i := 0; i < hourLimit; i++ {
			hour := weather.Hourly[i]
			if len(hour.Weather) == 0 {
				continue
			}

			t := time.Unix(hour.Dt, 0)
			day := t.Format("Mon Jan 2")
			hourStr := t.Format("15:04")

			if day != currentDay {
				if currentDay != "" {
					fmt.Println()
				}
				fmt.Printf("%s:\n", styles.HighlightStyleF(day))
				currentDay = day
			}

			weatherCond := hour.Weather[0]
			temp := styles.TempStyle(fmt.Sprintf("%.1f%s", hour.Temp, tempUnit))

			popStr := ""
			if hour.Pop > 0 {
				popStr = fmt.Sprintf(" (%.0f%% chance of precipitation)", hour.Pop*100)
			}

			extraSpace := ""
			if hour.Temp < 10 {
				extraSpace = " "
			}
			fmt.Printf("  %s:   %s  %s%s  %s%s\n",
				hourStr,
				temp,
				extraSpace,
				models.GetWeatherEmoji(weatherCond.ID, nil),
				styles.InfoStyle(weatherCond.Description),
				popStr)

			if hour.Rain != nil && hour.Rain.OneHour > 0 {
				fmt.Printf("       Rain: %.1f mm/h\n", hour.Rain.OneHour)
			}

			if hour.Snow != nil && hour.Snow.OneHour > 0 {
				fmt.Printf("       Snow: %.1f mm/h\n", hour.Snow.OneHour)
			}
		}
		fmt.Println()
	}
}
