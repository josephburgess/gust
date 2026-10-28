package tui

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"charm.land/lipgloss/v2"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
)

// gustScale is the km/h that fills the wind bar.
const gustScale = 100

func (m Model) nowView(width int) string {
	sections := []string{
		"",
		m.hero(),
		"",
		m.nowPanels(width),
		m.statLine(),
		"",
		m.next12(width),
	}
	if tip := m.tipLine(); tip != "" {
		sections = append(sections, "", tip)
	}
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) isNight() bool {
	c := m.weather.Current
	return c.Dt > c.Sunset || c.Dt < c.Sunrise
}

func conditionOf(weather []models.WeatherCondition) models.WeatherCondition {
	if len(weather) == 0 {
		return models.WeatherCondition{}
	}
	return weather[0]
}

func capitalise(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func (m Model) hero() string {
	c := m.weather.Current
	cond := conditionOf(c.Weather)

	icon := strings.Join(Icon(cond.ID, m.isNight()), "\n")
	temp := fg(styles.Gold).Render(strings.Join(BigTemp(deg(c.Temp)), "\n"))

	lines := []string{
		lipgloss.NewStyle().Bold(true).Render(capitalise(cond.Description)),
		fg(styles.Subtle).Render("feels like ") + fmt.Sprintf("%d°", deg(c.FeelsLike)),
		"",
	}
	if len(m.weather.Daily) > 0 {
		today := m.weather.Daily[0].Temp
		lines = append(lines,
			fg(styles.Subtle).Render("high ")+fg(styles.Gold).Render(fmt.Sprintf("%d°", deg(today.Max)))+
				fg(styles.Subtle).Render("  low ")+fg(styles.Foam).Render(fmt.Sprintf("%d°", deg(today.Min))))
	}
	text := lipgloss.NewStyle().PaddingTop(1).Render(strings.Join(lines, "\n"))

	return lipgloss.JoinHorizontal(lipgloss.Top, icon, "   ", temp, "    ", text)
}

func (m Model) nowPanels(width int) string {
	const gap = 1
	pw := (width - 3*gap) / 4
	last := width - 3*gap - 3*pw

	panels := []string{
		Panel("wind", m.windBody(pw-4), pw, false),
		Panel("rain", m.rainBody(pw-4), pw, false),
		Panel("humidity", m.humidityBody(pw-4), pw, false),
		Panel("sun", m.sunBody(last-4), last, false),
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, panels[0], " ", panels[1], " ", panels[2], " ", panels[3])
}

func (m Model) windBody(inner int) string {
	c := m.weather.Current
	speed := fmt.Sprintf("%d %s ", deg(windSpeed(c.WindSpeed, m.units)), windUnit(m.units)) +
		fg(styles.Iris).Render(models.GetWindArrow(c.WindDeg)) + " " + models.GetWindDirection(c.WindDeg)
	gusts := fg(styles.Love).Render(fmt.Sprintf("gusts %d", deg(windSpeed(c.WindGust, m.units))))
	bar := fg(styles.Love).Render(Bar(kmh(c.WindGust, m.units), gustScale, inner))
	return strings.Join([]string{speed, gusts, bar}, "\n")
}

func (m Model) rainBody(inner int) string {
	c := m.weather.Current
	amount := fmt.Sprintf("%.1f mm/h", 0.0)
	switch {
	case c.Rain != nil && c.Rain.OneHour > 0:
		amount = fmt.Sprintf("%.1f mm/h", c.Rain.OneHour)
	case c.Snow != nil && c.Snow.OneHour > 0:
		amount = fmt.Sprintf("%.1f mm/h snow", c.Snow.OneHour)
	}
	pop := 0.0
	if len(m.weather.Hourly) > 1 {
		pop = m.weather.Hourly[1].Pop
	}
	chance := fg(styles.Foam).Render(fmt.Sprintf("%d%%", deg(pop*100))) + fg(styles.Subtle).Render(" next hour")
	bar := fg(styles.Foam).Render(Bar(pop, 1, inner))
	return strings.Join([]string{amount, chance, bar}, "\n")
}

func (m Model) humidityBody(inner int) string {
	c := m.weather.Current
	dew := fg(styles.Subtle).Render("dew point ") + fmt.Sprintf("%d°", deg(c.DewPoint))
	bar := fg(styles.Pine).Render(Bar(float64(c.Humidity), 100, inner))
	return strings.Join([]string{fmt.Sprintf("%d%%", c.Humidity), dew, bar}, "\n")
}

func (m Model) sunBody(inner int) string {
	c := m.weather.Current
	times := fg(styles.Gold).Render("↑ ") + m.weather.LocalTime(c.Sunrise).Format("15:04") + "  " +
		fg(styles.Rose).Render("↓ ") + m.weather.LocalTime(c.Sunset).Format("15:04")

	now := c.Dt
	var rel string
	switch {
	case now < c.Sunrise:
		rel = "rises in " + shortDuration(c.Sunrise-now)
	case now < c.Sunset:
		rel = "sets in " + shortDuration(c.Sunset-now)
	default:
		rel = "set " + shortDuration(now-c.Sunset) + " ago"
	}

	pos := 0
	if span := c.Sunset - c.Sunrise; span > 0 {
		pos = clamp(int(float64(now-c.Sunrise)/float64(span)*float64(inner-1)), 0, inner-1)
	}
	progress := fg(styles.Gold).Render(strings.Repeat("━", pos)+"●") +
		fg(styles.HighlightMed).Render(strings.Repeat("─", inner-1-pos))

	return strings.Join([]string{times, fg(styles.Subtle).Render(rel), progress}, "\n")
}

func shortDuration(secs int64) string {
	d := time.Duration(secs) * time.Second
	h, mins := int(d.Hours()), int(d.Minutes())%60
	switch {
	case h == 0:
		return fmt.Sprintf("%dm", mins)
	case mins == 0:
		return fmt.Sprintf("%dh", h)
	default:
		return fmt.Sprintf("%dh %dm", h, mins)
	}
}

func (m Model) statLine() string {
	c := m.weather.Current
	stat := func(label, value string) string {
		return fg(styles.Subtle).Render(label+" ") + fg(styles.Text).Render(value)
	}
	dot := fg(styles.Subtle).Render(" · ")
	return strings.Join([]string{
		stat("uv", fmt.Sprintf("%.1f", c.UVI)),
		stat("visibility", fmt.Sprintf("%.0f km", float64(c.Visibility)/1000)),
		stat("cloud", fmt.Sprintf("%d%%", c.Clouds)),
		stat("pressure", fmt.Sprintf("%d hPa", c.Pressure)),
	}, dot)
}

func (m Model) next12(width int) string {
	hours := m.weather.Hourly[:min(12, len(m.weather.Hourly))]
	if len(hours) == 0 {
		return ""
	}
	inner := width - 4
	colWidth := inner / len(hours)

	temps := make([]float64, len(hours))
	for i, h := range hours {
		temps[i] = h.Temp
	}
	spark := []rune(Spark(temps))

	cols := make([]string, len(hours))
	for i, h := range hours {
		label := m.weather.LocalTime(h.Dt).Format("15")
		if i == 0 {
			label = "now"
		}
		cell := strings.Join([]string{
			fg(styles.Subtle).Render(label),
			fg(styles.Gold).Render(string(spark[i])),
			lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%d°", deg(h.Temp))),
			fg(styles.Foam).Render(fmt.Sprintf("%d%%", deg(h.Pop*100))),
		}, "\n")
		cols[i] = lipgloss.NewStyle().Width(colWidth).Align(lipgloss.Center).Render(cell)
	}
	return Panel("next 12 hours", lipgloss.JoinHorizontal(lipgloss.Top, cols...), width, false)
}

// tipLine is the weather tip without its emoji, or "" when tips are off.
func (m Model) tipLine() string {
	if !m.showTips {
		return ""
	}
	tip := models.GetWeatherTip(m.weather, m.units)
	tip = strings.TrimRightFunc(tip, func(r rune) bool { return r > unicode.MaxASCII || unicode.IsSpace(r) })
	return fg(styles.Iris).Italic(true).Render("› " + strings.ToLower(tip))
}
