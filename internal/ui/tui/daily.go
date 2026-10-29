package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
)

const (
	dayCount   = 5
	rangeWidth = 30
)

func (m Model) days() []models.DayData {
	return m.weather.Daily[:min(dayCount, len(m.weather.Daily))]
}

// weekRange is the lowest low and highest high of the shown days.
func (m Model) weekRange() (float64, float64) {
	days := m.days()
	lo, hi := days[0].Temp.Min, days[0].Temp.Max
	for _, d := range days[1:] {
		lo = min(lo, d.Temp.Min)
		hi = max(hi, d.Temp.Max)
	}
	return lo, hi
}

func (m Model) dayLabel(i int, d models.DayData) string {
	if i == 0 {
		return "Today"
	}
	return m.weather.LocalTime(d.Dt).Format("Mon 2")
}

func (m Model) dailyView(width int) string {
	days := m.days()
	if len(days) == 0 {
		return ""
	}
	weekMin, weekMax := m.weekRange()

	rows := []string{""}
	for i, d := range days {
		selected := i == m.dayCursor
		st := rowStyle(selected)
		cond := conditionOf(d.Weather)
		pre, fill, post := RangeBar(d.Temp.Min, d.Temp.Max, weekMin, weekMax, rangeWidth)

		row := strings.Join([]string{
			cursorCell(selected),
			st(lipgloss.NewStyle().Width(8)).Render(m.dayLabel(i, d)),
			st(fg(conditionColor(cond.ID)).Width(conditionCols)).Render(cond.Description),
			st(fg(styles.Foam).Width(5).Align(lipgloss.Right)).Render(fmt.Sprintf("%d°", deg(d.Temp.Min))),
			st(lipgloss.NewStyle()).Render(" "),
			st(fg(styles.HighlightMed)).Render(pre),
			st(fg(styles.Gold)).Render(fill),
			st(fg(styles.HighlightMed)).Render(post),
			st(fg(styles.Gold).Bold(true).Width(6)).Render(fmt.Sprintf(" %d°", deg(d.Temp.Max))),
			st(fg(styles.Foam).Width(8)).Render(fmt.Sprintf("☂ %d%%", deg(d.Pop*100))),
			st(lipgloss.NewStyle()).Render(fmt.Sprintf("%d %s ", deg(windSpeed(d.WindSpeed, m.units)), windUnit(m.units))),
			st(fg(styles.Iris)).Render(models.GetWindArrow(d.WindDeg)),
			st(lipgloss.NewStyle()).Render(" " + models.GetWindDirection(d.WindDeg)),
		}, "")
		if selected {
			row += lipgloss.NewStyle().Background(styles.Overlay).Render(strings.Repeat(" ", max(width-lipgloss.Width(row), 0)))
		}
		rows = append(rows, row)
	}

	rows = append(rows, "", m.dayDetail(days[m.dayCursor], width))
	if tip := m.dailyTip(); tip != "" {
		rows = append(rows, "", tip)
	}
	return strings.Join(rows, "\n")
}

func (m Model) dayDetail(d models.DayData, width int) string {
	date := m.weather.LocalTime(d.Dt).Format("Monday 2 January")
	sun := fg(styles.Gold).Render("↑ ") + m.weather.LocalTime(d.Sunrise).Format("15:04") + "  " +
		fg(styles.Rose).Render("↓ ") + m.weather.LocalTime(d.Sunset).Format("15:04")
	head := spread(lipgloss.NewStyle().Bold(true).Render(date), sun, width-4)

	part := func(label string, t float64) string {
		return lipgloss.NewStyle().Width((width - 4) / 4).Render(
			fg(styles.Subtle).Render(label) + "\n" + lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%d°", deg(t))))
	}
	temps := lipgloss.JoinHorizontal(lipgloss.Top,
		part("morning", d.Temp.Morn), part("day", d.Temp.Day), part("evening", d.Temp.Eve), part("night", d.Temp.Night))

	stat := func(label, value string) string {
		return fg(styles.Subtle).Render(label+" ") + value
	}
	dot := fg(styles.Subtle).Render(" · ")
	stats := strings.Join([]string{
		stat("rain", fmt.Sprintf("%.1f mm", d.Rain)),
		stat("wind", fmt.Sprintf("%d %s", deg(windSpeed(d.WindSpeed, m.units)), windUnit(m.units))),
		stat("gusts", fg(styles.Love).Render(fmt.Sprintf("%d", deg(windSpeed(d.WindGust, m.units))))),
		stat("uv", fmt.Sprintf("%.1f", d.UVI)),
		stat("humidity", fmt.Sprintf("%d%%", d.Humidity)),
	}, dot)

	body := []string{head}
	if d.Summary != "" {
		body = append(body, fg(styles.Subtle).Render(d.Summary))
	}
	body = append(body, "", temps, "", stats)
	return Panel("", strings.Join(body, "\n"), width, true)
}

// coldNight is the night low worth a mention, per unit system.
func coldNight(units string) float64 {
	switch units {
	case "imperial":
		return 37 // 3C
	case "standard":
		return 276 // 3C
	default:
		return 3
	}
}

// dailyTip names the driest day, plus the coldest night when it is near freezing.
func (m Model) dailyTip() string {
	if !m.showTips {
		return ""
	}
	days := m.days()
	driest, coldest := 0, 0
	for i, d := range days {
		if d.Pop < days[driest].Pop {
			driest = i
		}
		if d.Temp.Night < days[coldest].Temp.Night {
			coldest = i
		}
	}
	tip := fmt.Sprintf("driest day is %s, %d%% chance of rain", strings.ToLower(m.dayLabel(driest, days[driest])), deg(days[driest].Pop*100))
	if night := days[coldest].Temp.Night; night <= coldNight(m.units) {
		tip += fmt.Sprintf(". coldest night is %s at %d°", strings.ToLower(m.dayLabel(coldest, days[coldest])), deg(night))
	}
	return fg(styles.Iris).Italic(true).Render("› " + tip)
}
