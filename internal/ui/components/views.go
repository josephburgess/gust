package components

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
)

type Renderer struct {
	Units string
}

func NewRenderer(units string) *Renderer {
	return &Renderer{
		Units: units,
	}
}

func (r *Renderer) DisplayDefaultWeather(city *models.City, weather *models.OneCallResponse) {
	current := weather.Current

	fmt.Print(styles.FormatHeader(fmt.Sprintf("WEATHER FOR %s", strings.ToUpper(city.Name))))

	if len(current.Weather) > 0 {
		weatherCond := current.Weather[0]
		fmt.Printf("Current Conditions: %s %s\n\n",
			styles.HighlightStyleF(weatherCond.Description),
			models.GetWeatherEmoji(weatherCond.ID))

		tempUnit := r.getTemperatureUnit()

		fmt.Printf("Temperature: %s %s (F/L: %.1f%s)\n",
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

		fmt.Printf("Sunrise: %s %s  Sunset: %s %s\n\n",
			time.Unix(current.Sunrise, 0).Format("15:04"),
			"🌅",
			time.Unix(current.Sunset, 0).Format("15:04"),
			"🌇")
	}

	r.displayAlertSummary(weather.Alerts, city.Name)
}

func (r *Renderer) DisplayDailyForecast(city *models.City, weather *models.OneCallResponse) {
	fmt.Print(styles.FormatHeader(fmt.Sprintf("5-DAY FORECAST FOR %s", strings.ToUpper(city.Name))))

	if len(weather.Daily) > 0 {
		tempUnit := r.getTemperatureUnit()

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
				condition := fmt.Sprintf("%s %s", weather.Description, models.GetWeatherEmoji(weather.ID))
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

			windUnit := r.getWindSpeedUnit()
			windSpeed := r.formatWindSpeed(day.WindSpeed)

			fmt.Printf("  Wind: %.1f %s %s\n",
				windSpeed,
				windUnit,
				models.GetWindDirection(day.WindDeg))

			fmt.Printf("  UV Index: %.1f\n", day.UVI)
		}
		fmt.Println()
	}
}

func (r *Renderer) DisplayHourlyForecast(city *models.City, weather *models.OneCallResponse) {
	fmt.Print(styles.FormatHeader(fmt.Sprintf("24H FORECAST FOR %s", strings.ToUpper(city.Name))))

	if len(weather.Hourly) > 0 {
		hourLimit := int(math.Min(24, float64(len(weather.Hourly))))
		currentDay := ""

		tempUnit := r.getTemperatureUnit()

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
				models.GetWeatherEmoji(weatherCond.ID),
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

func (r *Renderer) DisplayAlerts(city *models.City, weather *models.OneCallResponse) {
	fmt.Print(styles.FormatHeader(fmt.Sprintf("WEATHER ALERTS FOR %s", strings.ToUpper(city.Name))))

	if len(weather.Alerts) == 0 {
		fmt.Println("No weather alerts for this area.")
		return
	}

	for i, alert := range weather.Alerts {
		if i > 0 {
			fmt.Println(styles.Divider(30))
		}

		fmt.Printf("%s\n", styles.AlertStyle(fmt.Sprintf("⚠️  %s", alert.Event)))
		fmt.Printf("Issued by: %s\n", alert.SenderName)
		fmt.Printf("Valid: %s to %s\n\n",
			styles.TimeStyle(time.Unix(alert.Start, 0).Format("Mon Jan 2 15:04")),
			styles.TimeStyle(time.Unix(alert.End, 0).Format("Mon Jan 2 15:04")))

		fmt.Println(alert.Description)
		fmt.Println()
	}
}

func (r *Renderer) DisplayFullWeather(city *models.City, weather *models.OneCallResponse) {
	r.DisplayDefaultWeather(city, weather)
	fmt.Println()

	if len(weather.Alerts) > 0 {
		r.DisplayAlerts(city, weather)
		fmt.Println()
	}

	r.DisplayDailyForecast(city, weather)
	fmt.Println()
}

func (r *Renderer) displayWindInfo(speed float64, deg int, gust float64) {
	windUnit := r.getWindSpeedUnit()
	windSpeed := r.formatWindSpeed(speed)

	if gust > 0 {
		gustSpeed := r.formatWindSpeed(gust)
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

func (r *Renderer) displayPrecipitation(rain *models.RainData, snow *models.SnowData) {
	if rain != nil && rain.OneHour > 0 {
		fmt.Printf("Rain: %.1f mm (last hour) 🌧️\n", rain.OneHour)
	}

	if snow != nil && snow.OneHour > 0 {
		fmt.Printf("Snow: %.1f mm (last hour) ❄️\n", snow.OneHour)
	}
}

func (r *Renderer) displayAlertSummary(alerts []models.Alert, cityName string) {
	if len(alerts) > 0 {
		fmt.Printf("%s Use 'gust --alerts %s' to view them.\n",
			styles.AlertStyle(fmt.Sprintf("⚠️  There are %d weather alerts for this area.", len(alerts))),
			cityName)
	}
}

func (r *Renderer) getTemperatureUnit() string {
	switch r.Units {
	case "imperial":
		return "°F"
	case "metric":
		return "°C"
	default:
		return "K"
	}
}

func (r *Renderer) getWindSpeedUnit() string {
	switch r.Units {
	case "imperial":
		return "mph"
	default:
		return "km/h"
	}
}

func (r *Renderer) formatWindSpeed(speed float64) float64 {
	switch r.Units {
	case "imperial":
		return speed
	default:
		return speed * 3.6
	}
}

func (r *Renderer) DisplayCompactWeather(city *models.City, weather *models.OneCallResponse) {
	current := weather.Current

	fmt.Print(styles.FormatHeader(fmt.Sprintf("%s WEATHER", strings.ToUpper(city.Name))))

	if len(current.Weather) > 0 {
		weatherCond := current.Weather[0]
		tempUnit := r.getTemperatureUnit()
		emoji := models.GetWeatherEmoji(weatherCond.ID)
		temp := styles.TempStyle(fmt.Sprintf("%.1f%s", current.Temp, tempUnit))
		feels := fmt.Sprintf("(%.1f%s)", current.FeelsLike, tempUnit)

		fmt.Printf("%s %-16s    %s %s\n",
			emoji,
			styles.HighlightStyleF(weatherCond.Description),
			temp,
			feels)

		windUnit := r.getWindSpeedUnit()
		windSpeed := r.formatWindSpeed(current.WindSpeed)
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
	}
}
