package tui

import (
	"errors"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/models"
	"github.com/stretchr/testify/assert"
)

var (
	fixtureNow = time.Date(2026, 11, 2, 20, 42, 0, 0, time.UTC)

	fixtureCity = &models.City{Name: "London", Country: "GB"}

	fixtureWeather = &models.OneCallResponse{
		Timezone:       "Europe/London",
		TimezoneOffset: 0,
		Current: models.CurrentWeather{
			Dt:         fixtureNow.Unix(),
			Sunrise:    time.Date(2026, 11, 2, 6, 57, 0, 0, time.UTC).Unix(),
			Sunset:     time.Date(2026, 11, 2, 16, 28, 0, 0, time.UTC).Unix(),
			Temp:       9.4,
			FeelsLike:  6.2,
			Pressure:   1012,
			Humidity:   78,
			DewPoint:   5.6,
			UVI:        0,
			Clouds:     75,
			Visibility: 10000,
			WindSpeed:  6.7,
			WindGust:   11.4,
			WindDeg:    225,
			Rain:       &models.RainData{OneHour: 0.6},
			Weather:    []models.WeatherCondition{{ID: 500, Description: "light rain"}},
		},
		Hourly: func() []models.HourData {
			hours := make([]models.HourData, 24)
			for i := range hours {
				hours[i] = models.HourData{
					Dt:      fixtureNow.Add(time.Duration(i) * time.Hour).Unix(),
					Temp:    9 - float64(i%6),
					Pop:     0.8,
					Weather: []models.WeatherCondition{{ID: 500, Description: "light rain"}},
				}
			}
			return hours
		}(),
		Daily: []models.DayData{
			{Dt: fixtureNow.Unix(), Temp: models.TempData{Min: 5, Max: 11}},
		},
	}

	fixtureAlert = models.Alert{Event: "Yellow wind warning", SenderName: "Met Office"}
)

func TestNowView(t *testing.T) {
	m := NewModel("London", "metric", true, TabNow, nil)
	m.now = func() time.Time { return fixtureNow }
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}})

	view := ansi.Strip(updated.View().Content)

	for _, want := range []string{
		"gust › London, GB",
		"Mon 2 Nov · 20:42 · metric",
		"1 Now",
		"Light rain",
		"feels like 6°",
		"high 11°  low 5°",
		"24 km/h ↗ SW",
		"gusts 41",
		"0.6 mm/h",
		"80% next hour",
		"dew point 6°",
		"↑ 06:57  ↓ 16:28",
		"set 4h 14m ago",
		"uv 0.0 · visibility 10 km · cloud 75% · pressure 1012 hPa",
		"NEXT 12 HOURS",
		"now",
		"› it might be raining right now - don't go out without an umbrella!",
		"updated just now",
	} {
		assert.Contains(t, view, want)
	}
	assert.NotContains(t, view, "Alerts ●")
}

func TestNowViewHidesTipWhenTipsAreOff(t *testing.T) {
	m := NewModel("London", "metric", false, TabNow, nil)
	m.now = func() time.Time { return fixtureNow }
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}})

	assert.NotContains(t, ansi.Strip(updated.View().Content), "umbrella")
}

func TestAlertsBadge(t *testing.T) {
	withAlert := *fixtureWeather
	withAlert.Alerts = []models.Alert{fixtureAlert}

	m := NewModel("London", "metric", true, TabNow, nil)
	m.now = func() time.Time { return fixtureNow }
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: &withAlert}})

	assert.Contains(t, ansi.Strip(updated.View().Content), "4 Alerts ●1")
}

func TestLoadingView(t *testing.T) {
	m := NewModel("London", "metric", true, TabNow, nil)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})

	assert.Contains(t, ansi.Strip(updated.View().Content), "fetching London…")
}

func TestFetchErrorShowsInline(t *testing.T) {
	m := NewModel("London", "metric", true, TabNow, nil)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, cmd := updated.Update(weatherMsg{err: errors.New("rate limit exceeded")})

	assert.Nil(t, cmd)
	assert.Contains(t, ansi.Strip(updated.View().Content), "rate limit exceeded")
}

func TestTabKeys(t *testing.T) {
	tests := []struct {
		name  string
		start Tab
		key   tea.KeyPressMsg
		want  Tab
	}{
		{"number picks a tab", TabNow, tea.KeyPressMsg{Code: '3', Text: "3"}, TabDaily},
		{"tab moves right", TabNow, tea.KeyPressMsg{Code: tea.KeyTab}, TabHourly},
		{"tab wraps at the end", TabAlerts, tea.KeyPressMsg{Code: tea.KeyTab}, TabNow},
		{"shift+tab wraps at the start", TabNow, tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift}, TabAlerts},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel("London", "metric", true, tt.start, nil)
			updated, _ := m.Update(tt.key)
			assert.Equal(t, tt.want, updated.(Model).tab)
		})
	}
}

func TestTabFromView(t *testing.T) {
	tests := []struct {
		view string
		want Tab
	}{
		{"hourly", TabHourly},
		{"daily", TabDaily},
		{"now", TabNow},
		{"default", TabNow},
		{"full", TabNow},
		{"", TabNow},
	}

	for _, tt := range tests {
		t.Run(tt.view, func(t *testing.T) {
			assert.Equal(t, tt.want, TabFromView(tt.view))
		})
	}
}
