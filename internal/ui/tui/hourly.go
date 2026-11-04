package tui

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
	"github.com/josephburgess/gust/internal/ui/widgets"
)

const (
	hourCount     = 24
	tempBarWidth  = 12
	rainBarWidth  = 10
	conditionCols = 18
)

func (m Model) hours() []models.HourData {
	return m.weather.Hourly[:min(hourCount, len(m.weather.Hourly))]
}

// conditionColor colours a condition by its ID: rain blue, clear gold, cloud grey.
func conditionColor(id int) color.Color {
	switch {
	case id >= 300 && id <= 531:
		return styles.Foam
	case id == 800:
		return styles.Gold
	case id >= 803 && id <= 804:
		return styles.Subtle
	default:
		return styles.Text
	}
}

// rowStyle adds the selected-row background to a cell style.
func rowStyle(selected bool) func(lipgloss.Style) lipgloss.Style {
	return func(s lipgloss.Style) lipgloss.Style {
		if selected {
			return s.Background(styles.Overlay)
		}
		return s
	}
}

func cursorCell(selected bool) string {
	if selected {
		return lipgloss.NewStyle().Background(styles.Overlay).Foreground(styles.Love).Render(" ▸ ")
	}
	return "   "
}

// hourlyLines returns the table lines and the line index of the cursor row.
func (m Model) hourlyLines(width int) ([]string, int) {
	hours := m.hours()
	temps := make([]float64, len(hours))
	for i, h := range hours {
		temps[i] = h.Temp
	}
	lo, _ := minMax(temps)

	var lines []string
	cursorLine := 0
	day := ""
	for i, h := range hours {
		t := m.weather.LocalTime(h.Dt)
		if d := t.Format("Mon 2 Jan"); d != day {
			if day != "" {
				label := "── " + d + " "
				lines = append(lines, fg(styles.Iris).Render(label+strings.Repeat("─", max(width-len([]rune(label)), 0))))
			}
			day = d
		}
		if i == m.hourCursor {
			cursorLine = len(lines)
		}
		lines = append(lines, m.hourRow(h, i == m.hourCursor, deg(lo), width))
	}
	return lines, cursorLine
}

func (m Model) hourRow(h models.HourData, selected bool, lo, width int) string {
	st := rowStyle(selected)
	cond := conditionOf(h.Weather)

	tempBar := widgets.HBlocks(min((deg(h.Temp)-lo+1)*2, tempBarWidth))
	filled := deg(h.Pop * rainBarWidth)
	rainBar := strings.Repeat("▮", filled) + strings.Repeat("▯", rainBarWidth-filled)
	mm := "–"
	if h.Rain != nil && h.Rain.OneHour > 0 {
		mm = fmt.Sprintf("%.1f", h.Rain.OneHour)
	}

	narrow := m.width < narrowWidth
	cells := []string{
		cursorCell(selected),
		st(lipgloss.NewStyle().Width(7)).Render(m.weather.LocalTime(h.Dt).Format("15:04")),
		st(lipgloss.NewStyle().Width(5).Bold(true)).Render(fmt.Sprintf("%d°", deg(h.Temp))),
	}
	if !narrow {
		cells = append(cells, st(fg(styles.Gold).Width(tempBarWidth+2)).Render(tempBar))
	}
	cells = append(cells,
		st(fg(conditionColor(cond.ID)).Width(conditionCols)).Render(cond.Description),
		st(fg(styles.Foam).Width(5).Align(lipgloss.Right)).Render(fmt.Sprintf("%d%%", deg(h.Pop*100))),
		st(lipgloss.NewStyle()).Render(" "),
	)
	if !narrow {
		cells = append(cells, st(fg(styles.Pine).Width(rainBarWidth+2)).Render(rainBar))
	}
	cells = append(cells,
		st(fg(styles.Foam).Width(6)).Render(mm),
		st(lipgloss.NewStyle()).Render(fmt.Sprintf("%d ", deg(windSpeed(h.WindSpeed, m.units)))),
		st(fg(styles.Iris)).Render(models.GetWindArrow(h.WindDeg)),
	)

	row := strings.Join(cells, "")
	if selected {
		row += lipgloss.NewStyle().Background(styles.Overlay).Render(strings.Repeat(" ", max(width-lipgloss.Width(row), 0)))
	}
	return row
}

func (m Model) hourlyView(width int) string {
	lines, _ := m.hourlyLines(width)
	end := min(m.hourScroll+m.bodyHeight()-1, len(lines))
	return "\n" + strings.Join(lines[min(m.hourScroll, end):end], "\n")
}

// scrollToCursor keeps the hourly cursor row on screen.
func (m Model) scrollToCursor() Model {
	if m.weather == nil {
		return m
	}
	_, cursorLine := m.hourlyLines(m.width - 2*pad)
	visible := max(m.bodyHeight()-1, 1)
	switch {
	case cursorLine < m.hourScroll:
		m.hourScroll = cursorLine
	case cursorLine >= m.hourScroll+visible:
		m.hourScroll = cursorLine - visible + 1
	}
	return m
}
