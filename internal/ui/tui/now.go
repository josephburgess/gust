package tui

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"charm.land/lipgloss/v2"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
	"github.com/josephburgess/gust/internal/ui/widgets"
)

// gustScale is the km/h that fills the wind bar.
const gustScale = 100

// Breakpoints, in terminal cells.
const (
	panelsGridWidth = 100 // under this, the Now panels go 2x2
	narrowWidth     = 80  // under this, drop the icon, big temp and hourly bars
	shortHeight     = 24  // under this, the Now tab hides the 12h panel
)

func (m Model) nowView(width int) string {
	sections := []string{
		"",
		m.hero(),
		"",
		m.nowPanels(width),
		m.statLine(),
	}
	if m.height >= shortHeight {
		sections = append(sections, "", m.next12(width))
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

	desc := lipgloss.NewStyle().Bold(true).Render(capitalise(cond.Description))
	if m.width < narrowWidth {
		desc = fg(styles.Gold).Bold(true).Render(fmt.Sprintf("%d°", deg(c.Temp))) + "  " + desc
	}
	lines := []string{
		desc,
		fg(styles.Subtle).Render("feels like ") + fmt.Sprintf("%d°", deg(c.FeelsLike)),
		"",
	}
	if len(m.weather.Daily) > 0 {
		today := m.weather.Daily[0].Temp
		lines = append(lines,
			fg(styles.Subtle).Render("high ")+fg(styles.Gold).Render(fmt.Sprintf("%d°", deg(today.Max)))+
				fg(styles.Subtle).Render("  low ")+fg(styles.Foam).Render(fmt.Sprintf("%d°", deg(today.Min))))
	}
	if m.width < narrowWidth {
		return strings.Join(lines, "\n")
	}
	text := lipgloss.NewStyle().PaddingTop(1).Render(strings.Join(lines, "\n"))
	icon := strings.Join(widgets.Icon(cond.ID, m.isNight()), "\n")
	temp := fg(styles.Gold).Render(strings.Join(widgets.BigTemp(deg(c.Temp)), "\n"))

	return lipgloss.JoinHorizontal(lipgloss.Top, icon, "   ", temp, "    ", text)
}

func (m Model) nowPanels(width int) string {
	perRow := 4
	if m.width < panelsGridWidth {
		perRow = 2
	}
	pw := (width - (perRow - 1)) / perRow
	last := width - (perRow - 1) - (perRow-1)*pw

	bodies := []struct {
		title string
		body  func(int) string
	}{{"wind", m.windBody}, {"rain", m.rainBody}, {"humidity", m.humidityBody}, {"sun", m.sunBody}}

	var rows []string
	for start := 0; start < len(bodies); start += perRow {
		var row []string
		for i := start; i < start+perRow; i++ {
			w := pw
			if i == start+perRow-1 {
				w = last
			}
			if i > start {
				row = append(row, " ")
			}
			row = append(row, widgets.Panel(bodies[i].title, bodies[i].body(w-4), w, false))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m Model) windBody(inner int) string {
	c := m.weather.Current
	speed := fmt.Sprintf("%d %s ", deg(windSpeed(c.WindSpeed, m.units)), windUnit(m.units)) +
		fg(styles.Iris).Render(models.GetWindArrow(c.WindDeg)) + " " + models.GetWindDirection(c.WindDeg)
	gusts := fg(styles.Love).Render(fmt.Sprintf("gusts %d", deg(windSpeed(c.WindGust, m.units))))
	bar := fg(styles.Love).Render(widgets.Bar(kmh(c.WindGust, m.units), gustScale, inner))
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
	bar := fg(styles.Foam).Render(widgets.Bar(pop, 1, inner))
	return strings.Join([]string{amount, chance, bar}, "\n")
}

func (m Model) humidityBody(inner int) string {
	c := m.weather.Current
	dew := fg(styles.Subtle).Render("dew point ") + fmt.Sprintf("%d°", deg(c.DewPoint))
	bar := fg(styles.Pine).Render(widgets.Bar(float64(c.Humidity), 100, inner))
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
	spark := []rune(widgets.Spark(temps))

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
	return widgets.Panel("next 12 hours", lipgloss.JoinHorizontal(lipgloss.Top, cols...), width, false)
}

// tipLine is the weather tip, or "" when tips are off.
func (m Model) tipLine() string {
	if !m.cfg.ShowTips {
		return ""
	}
	tip := models.GetWeatherTip(m.weather, m.units)
	return fg(styles.Iris).Italic(true).Render("› " + strings.ToLower(tip))
}
