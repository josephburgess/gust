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
		Daily: func() []models.DayData {
			pops := []float64{0.8, 0.6, 0.1, 0.4, 0.3}
			nights := []float64{6, 4, 1, 5, 7}
			days := make([]models.DayData, 5)
			for i := range days {
				day := fixtureNow.AddDate(0, 0, i)
				days[i] = models.DayData{
					Dt:        day.Unix(),
					Sunrise:   day.Add(-13 * time.Hour).Unix(),
					Sunset:    day.Add(-4 * time.Hour).Unix(),
					Summary:   "Rain clearing by the afternoon",
					Temp:      models.TempData{Min: 5 - float64(i), Max: 11 + float64(i), Morn: 6, Day: 10, Eve: 8, Night: nights[i]},
					Pop:       pops[i],
					Rain:      2.1,
					WindSpeed: 5.6,
					WindGust:  11.4,
					WindDeg:   225,
					UVI:       1.2,
					Humidity:  80,
					Weather:   []models.WeatherCondition{{ID: 500, Description: "light rain"}},
				}
			}
			return days
		}(),
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

func TestHourlyView(t *testing.T) {
	m := NewModel("London", "metric", true, TabHourly, nil)
	m.now = func() time.Time { return fixtureNow }
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}})

	view := ansi.Strip(updated.View().Content)

	for _, want := range []string{
		" ▸ 20:42  9°",
		"21:42",
		"── Tue 3 Nov ──",
		"light rain",
		"80% ▮▮▮▮▮▮▮▮▯▯",
		"–",
		"j/k move",
		"20:42 · 1 of 24",
	} {
		assert.Contains(t, view, want)
	}
}

func TestHourlyCursorScrollsOnShortTerminals(t *testing.T) {
	m := NewModel("London", "metric", true, TabHourly, nil)
	m.now = func() time.Time { return fixtureNow }
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 12})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}})
	for range 20 {
		updated, _ = updated.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	}

	view := ansi.Strip(updated.View().Content)

	assert.Equal(t, 20, updated.(Model).hourCursor)
	assert.Contains(t, view, " ▸ 16:42")
	assert.NotContains(t, view, "20:42  9°")
}

func TestHourlyCursorStopsAtTheEnds(t *testing.T) {
	m := NewModel("London", "metric", true, TabHourly, nil)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}})
	updated, _ = updated.Update(tea.KeyPressMsg{Code: tea.KeyUp})

	assert.Equal(t, 0, updated.(Model).hourCursor)
}

func TestDailyView(t *testing.T) {
	m := NewModel("London", "metric", true, TabDaily, nil)
	m.now = func() time.Time { return fixtureNow }
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}})

	view := ansi.Strip(updated.View().Content)

	for _, want := range []string{
		" ▸ Today",
		"Tue 3",
		"5° ",
		"━",
		" 11°",
		"☂ 80%",
		"20 km/h ↗ SW",
		"Monday 2 November",
		"Rain clearing by the afternoon",
		"morning",
		"night",
		"rain 2.1 mm · wind 20 km/h · gusts 41 · uv 1.2 · humidity 80%",
		"› driest day is wed 4, 10% chance of rain. coldest night is wed 4 at 1°",
		"week 1°–15°",
	} {
		assert.Contains(t, view, want)
	}
}

func TestDailyCursorMovesTheDetailBox(t *testing.T) {
	m := NewModel("London", "metric", true, TabDaily, nil)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}})
	updated, _ = updated.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})

	view := ansi.Strip(updated.View().Content)

	assert.Contains(t, view, " ▸ Tue 3")
	assert.Contains(t, view, "Tuesday 3 November")
}

func TestAlertsView(t *testing.T) {
	withAlert := *fixtureWeather
	withAlert.Alerts = []models.Alert{{
		Event:       "Yellow wind warning",
		SenderName:  "Met Office",
		Start:       fixtureNow.Unix(),
		End:         fixtureNow.Add(12 * time.Hour).Unix(),
		Description: "Strong winds may cause some disruption.",
	}}

	m := NewModel("London", "metric", true, TabAlerts, nil)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: &withAlert}})

	view := ansi.Strip(updated.View().Content)

	assert.Contains(t, view, "Yellow wind warning")
	assert.Contains(t, view, "Met Office · Mon 2 Nov 20:42 – Tue 3 Nov 08:42")
	assert.Contains(t, view, "Strong winds may cause some disruption.")
}

func TestAlertsViewWithNoAlerts(t *testing.T) {
	m := NewModel("London", "metric", true, TabAlerts, nil)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}})

	assert.Contains(t, ansi.Strip(updated.View().Content), "no alerts for London, GB")
}

func TestCityPromptRefetchesTheNewCity(t *testing.T) {
	var fetched []string
	fetch := func(city, units string) (*api.WeatherResponse, error) {
		fetched = append(fetched, city)
		return &api.WeatherResponse{City: &models.City{Name: "Paris", Country: "FR"}, Weather: fixtureWeather}, nil
	}
	m := NewModel("London", "metric", true, TabNow, fetch)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}})

	updated, _ = updated.Update(tea.KeyPressMsg{Code: 'c', Text: "c"})
	assert.True(t, updated.(Model).prompting)
	for _, r := range "paris" {
		updated, _ = updated.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	updated, cmd := updated.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	assert.False(t, updated.(Model).prompting)
	assert.Contains(t, ansi.Strip(updated.View().Content), "fetching paris…")
	for _, c := range cmd().(tea.BatchMsg) {
		if msg, ok := c().(weatherMsg); ok {
			updated, _ = updated.Update(msg)
		}
	}
	assert.Equal(t, []string{"paris"}, fetched)
	assert.Contains(t, ansi.Strip(updated.View().Content), "gust › Paris, FR")
}

func TestCityPromptEscCancels(t *testing.T) {
	m := NewModel("London", "metric", true, TabNow, nil)
	updated, _ := m.Update(tea.KeyPressMsg{Code: 'c', Text: "c"})
	updated, cmd := updated.Update(tea.KeyPressMsg{Code: tea.KeyEscape})

	assert.False(t, updated.(Model).prompting)
	assert.Equal(t, "London", updated.(Model).city)
	assert.Nil(t, cmd)
}

func TestUnitsKeyCyclesAndRefetches(t *testing.T) {
	var fetchedUnits []string
	fetch := func(city, units string) (*api.WeatherResponse, error) {
		fetchedUnits = append(fetchedUnits, units)
		return &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}, nil
	}
	m := NewModel("London", "metric", true, TabNow, fetch)
	updated, _ := m.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}})

	updated, cmd := updated.Update(tea.KeyPressMsg{Code: 'u', Text: "u"})
	for _, c := range cmd().(tea.BatchMsg) {
		c()
	}

	assert.Equal(t, "imperial", updated.(Model).units)
	assert.Equal(t, []string{"imperial"}, fetchedUnits)
}

func TestNextUnits(t *testing.T) {
	assert.Equal(t, "imperial", nextUnits("metric"))
	assert.Equal(t, "standard", nextUnits("imperial"))
	assert.Equal(t, "metric", nextUnits("standard"))
}

func TestRefreshKeepsTheWeatherAndSpinsInTheKeyBar(t *testing.T) {
	m := NewModel("London", "metric", true, TabNow, nil)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}})
	updated, cmd := updated.Update(tea.KeyPressMsg{Code: 'r', Text: "r"})

	view := ansi.Strip(updated.View().Content)

	assert.NotNil(t, cmd)
	assert.Contains(t, view, "gusts 41")
	assert.Contains(t, view, "refreshing")
}

func TestRefreshErrorShowsInTheKeyBar(t *testing.T) {
	m := NewModel("London", "metric", true, TabNow, nil)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}})
	updated, _ = updated.Update(tea.KeyPressMsg{Code: 'r', Text: "r"})
	updated, cmd := updated.Update(weatherMsg{err: errors.New("rate limit exceeded")})

	view := ansi.Strip(updated.View().Content)

	assert.Nil(t, cmd)
	assert.Contains(t, view, "gusts 41")
	assert.Contains(t, view, "rate limit exceeded")
}

func TestHelpKeyTogglesTheHelp(t *testing.T) {
	m := NewModel("London", "metric", true, TabNow, nil)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}})
	updated, _ = updated.Update(tea.KeyPressMsg{Code: '?', Text: "?"})

	assert.Contains(t, ansi.Strip(updated.View().Content), "previous tab")

	updated, _ = updated.Update(tea.KeyPressMsg{Code: tea.KeyEscape})

	assert.NotContains(t, ansi.Strip(updated.View().Content), "previous tab")
}
