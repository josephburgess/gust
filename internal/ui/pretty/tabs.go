package pretty

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
)

var (
	labelStyle = lipgloss.NewStyle().Foreground(styles.Subtle)
	valueStyle = lipgloss.NewStyle().Foreground(styles.Text)
	tempStyle  = lipgloss.NewStyle().Bold(true).Foreground(styles.Gold)
	accentStyle = lipgloss.NewStyle().Foreground(styles.Foam)
	alertStyle  = lipgloss.NewStyle().Bold(true).Foreground(styles.Love)
	dimStyle    = lipgloss.NewStyle().Foreground(styles.Muted)
	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.Iris).
			MarginTop(1)
)

func kv(label, value string) string {
	return labelStyle.Render(label+":") + " " + valueStyle.Render(value)
}

// padTo truncates or space-pads s to exactly width runes.
func padTo(s string, width int) string {
	r := []rune(s)
	if len(r) >= width {
		return string(r[:width])
	}
	return s + strings.Repeat(" ", width-len(r))
}

func (m *Model) tempUnit() string {
	switch m.units {
	case "imperial":
		return "°F"
	case "metric":
		return "°C"
	default:
		return "K"
	}
}

func (m *Model) windUnit() string {
	if m.units == "imperial" {
		return "mph"
	}
	return "km/h"
}

func (m *Model) formatWindSpeed(speed float64) float64 {
	if m.units == "imperial" {
		return speed
	}
	return speed * 3.6
}

func (m *Model) currentTab() string {
	current := m.weather.Current
	if len(current.Weather) == 0 {
		return dimStyle.Render("No current weather data available.")
	}

	cond := current.Weather[0]
	emoji := models.GetWeatherEmoji(cond.ID, &current)
	tempUnit := m.tempUnit()
	windUnit := m.windUnit()
	windSpeed := m.formatWindSpeed(current.WindSpeed)

	var sb strings.Builder

	// City title
	sb.WriteString(titleStyle.Render(fmt.Sprintf("%s %s", strings.ToUpper(m.city.Name), emoji)))
	sb.WriteString("\n")
	if m.city.Country != "" {
		sb.WriteString(dimStyle.Render(m.city.Country))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Conditions
	sb.WriteString(sectionStyle.Render("CONDITIONS"))
	sb.WriteString("\n")
	sb.WriteString(kv("Description", cond.Description))
	sb.WriteString("\n")
	sb.WriteString(kv("Temperature", tempStyle.Render(fmt.Sprintf("%.1f%s", current.Temp, tempUnit))+
		dimStyle.Render(fmt.Sprintf("  (feels like %.1f%s)", current.FeelsLike, tempUnit))))
	sb.WriteString("\n")
	sb.WriteString(kv("Humidity", fmt.Sprintf("%d%%", current.Humidity)))
	sb.WriteString("\n")
	if current.UVI > 0 {
		sb.WriteString(kv("UV Index", fmt.Sprintf("%.1f", current.UVI)))
		sb.WriteString("\n")
	}
	sb.WriteString(kv("Cloud Cover", fmt.Sprintf("%d%%", current.Clouds)))
	sb.WriteString("\n")

	// Wind
	sb.WriteString("\n")
	sb.WriteString(sectionStyle.Render("WIND"))
	sb.WriteString("\n")
	sb.WriteString(kv("Speed", fmt.Sprintf("%.1f %s %s", windSpeed, windUnit, models.GetWindDirection(current.WindDeg))))
	sb.WriteString("\n")
	if current.WindGust > 0 {
		sb.WriteString(kv("Gusts", fmt.Sprintf("%.1f %s", m.formatWindSpeed(current.WindGust), windUnit)))
		sb.WriteString("\n")
	}

	// Precipitation
	if current.Rain != nil && current.Rain.OneHour > 0 {
		sb.WriteString("\n")
		sb.WriteString(sectionStyle.Render("PRECIPITATION"))
		sb.WriteString("\n")
		sb.WriteString(kv("Rain", fmt.Sprintf("%.1f mm/h", current.Rain.OneHour)))
		sb.WriteString("\n")
	}
	if current.Snow != nil && current.Snow.OneHour > 0 {
		if current.Rain == nil || current.Rain.OneHour == 0 {
			sb.WriteString("\n")
			sb.WriteString(sectionStyle.Render("PRECIPITATION"))
			sb.WriteString("\n")
		}
		sb.WriteString(kv("Snow", fmt.Sprintf("%.1f mm/h", current.Snow.OneHour)))
		sb.WriteString("\n")
	}

	// Visibility & Sun
	sb.WriteString("\n")
	sb.WriteString(sectionStyle.Render("VISIBILITY & SUN"))
	sb.WriteString("\n")
	sb.WriteString(kv("Visibility", models.VisibilityToString(current.Visibility)))
	sb.WriteString("\n")
	sb.WriteString(kv("Sunrise", time.Unix(current.Sunrise, 0).Format("15:04")))
	sb.WriteString("\n")
	sb.WriteString(kv("Sunset", time.Unix(current.Sunset, 0).Format("15:04")))
	sb.WriteString("\n")

	// Tip
	if m.cfg != nil && m.cfg.ShowTips {
		tip := models.GetWeatherTip(m.weather, m.units)
		sb.WriteString("\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(styles.Pine).Italic(true).Render("💡 "+tip))
		sb.WriteString("\n")
	}

	// Alerts summary
	if len(m.weather.Alerts) > 0 {
		sb.WriteString("\n")
		sb.WriteString(alertStyle.Render(fmt.Sprintf("⚠️  %d active weather alert(s) — see Alerts tab", len(m.weather.Alerts))))
		sb.WriteString("\n")
	}

	return sb.String()
}

func (m *Model) hourlyTab() string {
	if len(m.weather.Hourly) == 0 {
		return dimStyle.Render("No hourly forecast data available.")
	}

	tempUnit := m.tempUnit()
	limit := int(math.Min(24, float64(len(m.weather.Hourly))))

	// first pass: find the longest description
	maxDescLen := 0
	for i := 0; i < limit; i++ {
		if len(m.weather.Hourly[i].Weather) > 0 {
			if n := len([]rune(m.weather.Hourly[i].Weather[0].Description)); n > maxDescLen {
				maxDescLen = n
			}
		}
	}

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("24-HOUR FORECAST"))
	sb.WriteString("\n\n")

	currentDay := ""
	for i := 0; i < limit; i++ {
		hour := m.weather.Hourly[i]
		if len(hour.Weather) == 0 {
			continue
		}

		t := time.Unix(hour.Dt, 0)
		day := t.Format("Mon Jan 2")
		hourStr := t.Format("15:04")

		if day != currentDay {
			if currentDay != "" {
				sb.WriteString("\n")
			}
			sb.WriteString(accentStyle.Bold(true).Render(day))
			sb.WriteString("\n")
			currentDay = day
		}

		cond := hour.Weather[0]
		emoji := models.GetWeatherEmoji(cond.ID, nil)
		temp := tempStyle.Render(fmt.Sprintf("%5.1f%s", hour.Temp, tempUnit))
		desc := valueStyle.Render(padTo(cond.Description, maxDescLen))

		var extras []string
		if hour.Pop > 0 {
			extras = append(extras, fmt.Sprintf("%3.0f%% 💧", hour.Pop*100))
		}
		if hour.Rain != nil && hour.Rain.OneHour > 0 {
			extras = append(extras, fmt.Sprintf("🌧 %4.1fmm", hour.Rain.OneHour))
		}
		if hour.Snow != nil && hour.Snow.OneHour > 0 {
			extras = append(extras, fmt.Sprintf("❄ %4.1fmm", hour.Snow.OneHour))
		}
		extraStr := ""
		if len(extras) > 0 {
			extraStr = dimStyle.Render("  " + strings.Join(extras, "  "))
		}

		sb.WriteString(fmt.Sprintf("  %s   %s  %s  %s%s\n",
			labelStyle.Render(hourStr),
			temp,
			emoji,
			desc,
			extraStr,
		))
	}

	return sb.String()
}

func (m *Model) dailyTab() string {
	if len(m.weather.Daily) == 0 {
		return dimStyle.Render("No daily forecast data available.")
	}

	tempUnit := m.tempUnit()
	windUnit := m.windUnit()

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("5-DAY FORECAST"))
	sb.WriteString("\n\n")

	for i, day := range m.weather.Daily {
		if i >= 5 {
			break
		}
		if i > 0 {
			sb.WriteString(dimStyle.Render(strings.Repeat("─", 40)))
			sb.WriteString("\n")
		}

		date := time.Unix(day.Dt, 0).Format("Monday, Jan 2")
		sb.WriteString(accentStyle.Bold(true).Render(date))
		sb.WriteString("\n")

		if day.Summary != "" {
			sb.WriteString(dimStyle.Italic(true).Render("  "+day.Summary))
			sb.WriteString("\n")
		}

		if len(day.Weather) > 0 {
			cond := day.Weather[0]
			emoji := models.GetWeatherEmoji(cond.ID, nil)
			sb.WriteString(fmt.Sprintf("  %s  %s\n", emoji, valueStyle.Render(cond.Description)))
		}

		sb.WriteString(fmt.Sprintf("  %s  %s / %s\n",
			labelStyle.Render("High/Low"),
			tempStyle.Render(fmt.Sprintf("%.1f%s", day.Temp.Max, tempUnit)),
			tempStyle.Render(fmt.Sprintf("%.1f%s", day.Temp.Min, tempUnit)),
		))

		sb.WriteString(fmt.Sprintf("  %s\n",
			dimStyle.Render(fmt.Sprintf("Morn %.1f%s  Day %.1f%s  Eve %.1f%s  Night %.1f%s",
				day.Temp.Morn, tempUnit,
				day.Temp.Day, tempUnit,
				day.Temp.Eve, tempUnit,
				day.Temp.Night, tempUnit,
			)),
		))

		if day.Pop > 0 {
			sb.WriteString(fmt.Sprintf("  %s\n", kv("Precipitation chance", fmt.Sprintf("%d%%", int(day.Pop*100)))))
		}
		if day.Rain > 0 {
			sb.WriteString(fmt.Sprintf("  %s\n", kv("Rain", fmt.Sprintf("%.1f mm", day.Rain))))
		}
		if day.Snow > 0 {
			sb.WriteString(fmt.Sprintf("  %s\n", kv("Snow", fmt.Sprintf("%.1f mm", day.Snow))))
		}

		windSpeed := m.formatWindSpeed(day.WindSpeed)
		sb.WriteString(fmt.Sprintf("  %s\n", kv("Wind", fmt.Sprintf("%.1f %s %s", windSpeed, windUnit, models.GetWindDirection(day.WindDeg)))))
		sb.WriteString(fmt.Sprintf("  %s\n", kv("UV Index", fmt.Sprintf("%.1f", day.UVI))))
		sb.WriteString("\n")
	}

	return sb.String()
}

func (m *Model) alertsTab() string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render("WEATHER ALERTS"))
	sb.WriteString("\n\n")

	if len(m.weather.Alerts) == 0 {
		sb.WriteString(dimStyle.Render("No weather alerts for this area."))
		sb.WriteString("\n")
		return sb.String()
	}

	for i, alert := range m.weather.Alerts {
		if i > 0 {
			sb.WriteString(dimStyle.Render(strings.Repeat("─", 40)))
			sb.WriteString("\n\n")
		}

		sb.WriteString(alertStyle.Render(fmt.Sprintf("⚠️  %s", alert.Event)))
		sb.WriteString("\n")
		sb.WriteString(kv("Issued by", alert.SenderName))
		sb.WriteString("\n")
		sb.WriteString(kv("Valid",
			fmt.Sprintf("%s → %s",
				time.Unix(alert.Start, 0).Format("Mon Jan 2 15:04"),
				time.Unix(alert.End, 0).Format("Mon Jan 2 15:04"),
			),
		))
		sb.WriteString("\n\n")
		sb.WriteString(valueStyle.Render(alert.Description))
		sb.WriteString("\n\n")
	}

	return sb.String()
}

func (m *Model) tabContent() string {
	switch m.activeTab {
	case TabCurrent:
		return m.currentTab()
	case TabHourly:
		return m.hourlyTab()
	case TabDaily:
		return m.dailyTab()
	case TabAlerts:
		return m.alertsTab()
	default:
		return ""
	}
}
