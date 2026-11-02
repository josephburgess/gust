package tui

import (
	"fmt"
	"strings"

	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
	"github.com/josephburgess/gust/internal/ui/widgets"
)

// rainLikely is the pop at which "rain back" is worth a mention.
const rainLikely = 0.4

func raining(rain *models.RainData) bool {
	return rain != nil && rain.OneHour > 0
}

// Compact is the two-line print layout: conditions now, then the next 24h.
func Compact(city *models.City, w *models.OneCallResponse, units string) string {
	c := w.Current
	subtle := fg(styles.Subtle)
	dot := subtle.Render(" · ")
	hours := w.Hourly[:min(hourCount, len(w.Hourly))]

	first := []string{
		fg(styles.Rose).Bold(true).Render(city.Name) + "  " +
			fg(styles.Gold).Bold(true).Render(fmt.Sprintf("%d°", deg(c.Temp))) + " " +
			fg(styles.Foam).Render(conditionOf(c.Weather).Description) + "  " +
			subtle.Render(fmt.Sprintf("feels %d°", deg(c.FeelsLike))),
	}
	if len(w.Daily) > 0 {
		first = append(first, subtle.Render(fmt.Sprintf("↑%d° ↓%d°", deg(w.Daily[0].Temp.Max), deg(w.Daily[0].Temp.Min))))
	}
	wind := subtle.Render(fmt.Sprintf("%s %d %s", models.GetWindArrow(c.WindDeg), deg(windSpeed(c.WindSpeed, units)), windUnit(units)))
	if c.WindGust > 0 {
		wind += fg(styles.Love).Render(fmt.Sprintf(" gusts %d", deg(windSpeed(c.WindGust, units))))
	}
	first = append(first, wind)

	if len(hours) > 0 {
		rain := subtle.Render(fmt.Sprintf("☂ %d%%", deg(hours[0].Pop*100)))
		if raining(c.Rain) {
			for _, h := range hours[1:] {
				if !raining(h.Rain) {
					rain += subtle.Render(" till " + w.LocalTime(h.Dt).Format("15:04"))
					break
				}
			}
		}
		first = append(first, rain)
	}
	lines := []string{strings.Join(first, dot)}

	if len(hours) > 0 {
		temps := make([]float64, len(hours))
		for i, h := range hours {
			temps[i] = h.Temp
		}
		lo, hi := minMax(temps)
		second := fg(styles.Gold).Render(widgets.Spark(temps)) + subtle.Render(fmt.Sprintf("  next 24h %d°–%d°", deg(lo), deg(hi)))
		if !raining(c.Rain) {
			for _, h := range hours[1:] {
				if h.Pop > rainLikely {
					second += subtle.Render(", rain back ~" + w.LocalTime(h.Dt).Format("15:04"))
					break
				}
			}
		}
		lines = append(lines, "  "+second)
	}
	return strings.Join(lines, "\n")
}
