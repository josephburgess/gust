package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/josephburgess/gust/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	tests := []struct {
		name      string
		tab       Tab
		want      []string
		notWanted []string
	}{
		{"now", TabNow, []string{"gust › London, GB", "gusts 41", "NEXT 12 HOURS"}, []string{"1 Now", "q quit"}},
		{"hourly shows all 24 hours", TabHourly, []string{"20:42", "19:42", "── Tue 3 Nov"}, []string{"j/k move"}},
		{"daily", TabDaily, []string{"Today", "Fri 6", "Monday 2 November"}, nil},
		{"alerts", TabAlerts, []string{"no alerts for London, GB"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ansi.Strip(Print(fixtureCity, fixtureWeather, &config.Config{Units: "metric", ShowTips: true}, tt.tab, 120))

			for _, want := range tt.want {
				assert.Contains(t, got, want)
			}
			for _, notWanted := range tt.notWanted {
				assert.NotContains(t, got, notWanted)
			}
		})
	}
}

func TestPrintHasNoCursorOrTrailingSpaces(t *testing.T) {
	got := ansi.Strip(Print(fixtureCity, fixtureWeather, &config.Config{Units: "metric"}, TabHourly, 120))

	assert.NotContains(t, got, "▸")
	for _, line := range strings.Split(got, "\n") {
		assert.Equal(t, strings.TrimRight(line, " "), line)
	}
}
