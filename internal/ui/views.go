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

		tempUnit := getTemperatureUnit(current.Temp)

		fmt.Printf("Temperature: %s %s (Feels like: %.1f%s)\n",
			TempStyle(fmt.Sprintf("%.1f%s", current.Temp, tempUnit)),
			"🌡️",
			current.FeelsLike, tempUnit)

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
		tempUnit := getTemperatureUnit(weather.Daily[0].Temp.Day)

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
				TempStyle(fmt.Sprintf("%.1f%s", day.Temp.Max, tempUnit)),
				TempStyle(fmt.Sprintf("%.1f%s", day.Temp.Min, tempUnit)),
				"🌡️")

			fmt.Printf("  Morning: %.1f%s  Day: %.1f%s  Evening: %.1f%s  Night: %.1f%s\n",
				day.Temp.Morn, tempUnit,
				day.Temp.Day, tempUnit,
				day.Temp.Eve, tempUnit,
				day.Temp.Night, tempUnit)

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

			windUnit := getWindSpeedUnit(day.WindSpeed)
			windSpeed := formatWindSpeed(day.WindSpeed)

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
	fmt.Print(FormatHeader(fmt.Sprintf("HOURLY FORECAST FOR %s", strings.ToUpper(city.Name))))

	if len(weather.Hourly) > 0 {
		hourLimit := int(math.Min(24, float64(len(weather.Hourly))))
		currentDay := ""

		// Get temperature unit from the first hourly forecast
		tempUnit := getTemperatureUnit(weather.Hourly[0].Temp)

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
			temp := TempStyle(fmt.Sprintf("%.1f%s", hour.Temp, tempUnit))

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

func displayWindInfo(speed float64, deg int, gust float64) {
	windUnit := getWindSpeedUnit(speed)
	windSpeed := formatWindSpeed(speed)

	if gust > 0 {
		gustSpeed := formatWindSpeed(gust)
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

func getTemperatureUnit(temp float64) string {
	if temp > 100 {
		return "K" // Kelvin
	} else if temp > 50 || temp < -50 {
		return "°F" // Likely Fahrenheit
	} else {
		return "°C" // Celsius
	}
}

func getWindSpeedUnit(speed float64) string {
	if speed < 30 {
		return "km/h"
	} else {
		return "mph"
	}
}

func formatWindSpeed(speed float64) float64 {
	if speed < 30 {
		return speed * 3.6
	}
	return speed
}
