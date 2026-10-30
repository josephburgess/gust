package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/josephburgess/gust/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCompactWhileRaining(t *testing.T) {
	lines := strings.Split(ansi.Strip(Compact(fixtureCity, fixtureWeather, "metric")), "\n")

	assert.Len(t, lines, 2)
	assert.Equal(t, "London  9° light rain  feels 6° · ↑11° ↓5° · ↗ 24 km/h gusts 41 · ☂ 80% till 21:42", lines[0])
	assert.Equal(t, "  █▇▅▄▂▁█▇▅▄▂▁█▇▅▄▂▁█▇▅▄▂▁  next 24h 4°–9°", lines[1])
}

func TestCompactWhenDryWithRainLater(t *testing.T) {
	dry := *fixtureWeather
	dry.Current.Rain = nil
	dry.Hourly = append([]models.HourData{}, fixtureWeather.Hourly...)
	for i := range dry.Hourly[:5] {
		dry.Hourly[i].Pop = 0.1
	}

	got := ansi.Strip(Compact(fixtureCity, &dry, "metric"))

	assert.Contains(t, got, "☂ 10%")
	assert.NotContains(t, got, "till")
	assert.Contains(t, got, "next 24h 4°–9°, rain back ~01:42")
}

func TestCompactWhenDryAllDay(t *testing.T) {
	dry := *fixtureWeather
	dry.Current.Rain = nil
	dry.Hourly = append([]models.HourData{}, fixtureWeather.Hourly...)
	for i := range dry.Hourly {
		dry.Hourly[i].Pop = 0
	}

	assert.NotContains(t, ansi.Strip(Compact(fixtureCity, &dry, "metric")), "rain back")
}
