package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/josephburgess/gust/internal/models"
)

type Renderer struct{}

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (r *Renderer) DisplayCurrentWeather(city *models.City, weather *models.OneCallResponse) {
	current := weather.Current

	fmt.Print(FormatHeader(fmt.Sprintf("WEATHER FOR %s", strings.ToUpper(city.Name))))

	if len(current.Weather) > 0 {
		weatherCond := current.Weather[0]
		fmt.Printf("Current Conditions: %s %s\n\n",
			HighlightStyle(weatherCond.Description),
			models.GetWeatherEmoji(weatherCond.ID))

		fmt.Printf("Temperature: %s %s (Feels like: %.1f°C)\n",
			TempStyle(fmt.Sprintf("%.1f°C", models.KelvinToCelsius(current.Temp))),
			"🌡️",
			models.KelvinToCelsius(current.FeelsLike))

		fmt.Printf("Humidity: %d%% %s\n", current.Humidity, "💧")
		if current.UVI > 0 {
			fmt.Printf("UV Index: %.1f ☀️\n", current.UVI)
		}

		displayWindInfo(current.WindSpeed, current.WindDeg, current.WindGust)

		if current.Clouds > 0 {
			fmt.Printf("Cloud coverage: %d%% ☁️\n", current.Clouds)
		}

		displayPrecipitation(current.Rain, current.Snow)
		fmt.Printf("Visibility: %s\n", models.VisibilityToString(current.Visibility))

		fmt.Printf("Sunrise: %s %s  Sunset: %s %s\n\n",
			time.Unix(current.Sunrise, 0).Format("15:04"),
			"🌅",
			time.Unix(current.Sunset, 0).Format("15:04"),
			"🌇")
	}

	displayAlertSummary(weather.Alerts, city.Name)
}

func (r *Renderer) DisplayDailyForecast(city *models.City, weather *models.OneCallResponse) {
	fmt.Print(FormatHeader(fmt.Sprintf("7-DAY FORECAST FOR %s", strings.ToUpper(city.Name))))

	if len(weather.Daily) > 0 {
		for i, day := range weather.Daily {
			if i >= 5 {
				break
			}

			date := time.Unix(day.Dt, 0).Format("Mon Jan 2")

			if i > 0 {
				fmt.Println()
			}

			fmt.Printf("%s: %s\n",
				HighlightStyle(date),
				day.Summary)

			fmt.Printf("  High/Low: %s/%s %s\n",
				TempStyle(fmt.Sprintf("%.1f°C", models.KelvinToCelsius(day.Temp.Max))),
				TempStyle(fmt.Sprintf("%.1f°C", models.KelvinToCelsius(day.Temp.Min))),
				"🌡️")

			fmt.Printf("  Morning: %.1f°C  Day: %.1f°C  Evening: %.1f°C  Night: %.1f°C\n",
				models.KelvinToCelsius(day.Temp.Morn),
				models.KelvinToCelsius(day.Temp.Day),
				models.KelvinToCelsius(day.Temp.Eve),
				models.KelvinToCelsius(day.Temp.Night))

			if len(day.Weather) > 0 {
				weather := day.Weather[0]
				condition := fmt.Sprintf("%s %s", weather.Description, models.GetWeatherEmoji(weather.ID))
				fmt.Printf("  Conditions: %s\n", InfoStyle(condition))
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

			fmt.Printf("  Wind: %.1f km/h %s\n",
				day.WindSpeed*3.6,
				models.GetWindDirection(day.WindDeg))

			fmt.Printf("  UV Index: %.1f\n", day.UVI)
		}
		fmt.Println()
	}
}

func (r *Renderer) DisplayHourlyForecast(city *models.City, weather *models.OneCallResponse) {
	fmt.Print(FormatHeader(fmt.Sprintf("HOURLY FORECAST FOR %s", strings.ToUpper(city.Name))))

	if len(weather.Hourly) > 0 {
		hourLimit := int(math.Min(24, float64(len(weather.Hourly))))
		currentDay := ""

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
				fmt.Printf("%s:\n", HighlightStyle(day))
				currentDay = day
			}

			weatherCond := hour.Weather[0]
			temp := TempStyle(fmt.Sprintf("%.1f°C", models.KelvinToCelsius(hour.Temp)))

			popStr := ""
			if hour.Pop > 0 {
				popStr = fmt.Sprintf(" (%.0f%% chance of precipitation)", hour.Pop*100)
			}

			fmt.Printf("  %s: %s %s  %s%s\n",
				hourStr,
				temp,
				models.GetWeatherEmoji(weatherCond.ID),
				InfoStyle(weatherCond.Description),
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
	fmt.Print(FormatHeader(fmt.Sprintf("WEATHER ALERTS FOR %s", strings.ToUpper(city.Name))))

	if len(weather.Alerts) == 0 {
		fmt.Println("No weather alerts for this area.")
		return
	}

	for i, alert := range weather.Alerts {
		if i > 0 {
			fmt.Println(Divider())
		}

		fmt.Printf("%s\n", AlertStyle(fmt.Sprintf("⚠️  %s", alert.Event)))
		fmt.Printf("Issued by: %s\n", alert.SenderName)
		fmt.Printf("Valid: %s to %s\n\n",
			TimeStyle(time.Unix(alert.Start, 0).Format("Mon Jan 2 15:04")),
			TimeStyle(time.Unix(alert.End, 0).Format("Mon Jan 2 15:04")))

		fmt.Println(alert.Description)
		fmt.Println()
	}
}

func (r *Renderer) DisplayFullWeather(city *models.City, weather *models.OneCallResponse) {
	r.DisplayCurrentWeather(city, weather)
	fmt.Println()

	if len(weather.Alerts) > 0 {
		r.DisplayAlerts(city, weather)
		fmt.Println()
	}

	r.DisplayDailyForecast(city, weather)
	fmt.Println()
}

// Helper functions
func displayWindInfo(speed float64, deg int, gust float64) {
	if gust > 0 {
		fmt.Printf("Wind: %.1f km/h %s %s (Gusts: %.1f km/h)\n",
			speed*3.6, // m/s to km/h
			models.GetWindDirection(deg),
			"💨",
			gust*3.6)
	} else {
		fmt.Printf("Wind: %.1f km/h %s %s\n",
			speed*3.6, // m/s to km/h
			models.GetWindDirection(deg),
			"💨")
	}
}

func displayPrecipitation(rain *models.RainData, snow *models.SnowData) {
	if rain != nil && rain.OneHour > 0 {
		fmt.Printf("Rain: %.1f mm (last hour) 🌧️\n", rain.OneHour)
	}

	if snow != nil && snow.OneHour > 0 {
		fmt.Printf("Snow: %.1f mm (last hour) ❄️\n", snow.OneHour)
	}
}

func displayAlertSummary(alerts []models.Alert, cityName string) {
	if len(alerts) > 0 {
		fmt.Printf("%s Use 'gust --alerts %s' to view them.\n",
			AlertStyle(fmt.Sprintf("⚠️  There are %d weather alerts for this area.", len(alerts))),
			cityName)
	}
}
