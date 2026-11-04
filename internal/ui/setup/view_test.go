package setup

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestView(t *testing.T) {
	tests := []struct {
		name    string
		state   SetupState
		prepare func(m *Model)
		want    []string
	}{
		{"city", StateCity, func(m *Model) {}, []string{
			"gust setup", "● city ─── ○ preferences ─── ○ sign in", "Where's home?",
			"Your default city. Change it any time with c or s.", "Wherever the wind blows...",
			"not there? add a country code, like london,gb", "step 1 of 3",
		}},
		{"search", StateCitySearch, func(m *Model) { m.CitySearchQuery = "london" }, []string{`searching for "london"…`}},
		{"search error", StateCity, func(m *Model) { m.SearchErr = `no matches for "zzz"` }, []string{`no matches for "zzz"`}},
		{"results", StateCitySelect, func(m *Model) {
			m.CityOptions = []models.City{{Name: "London", State: "England", Country: "GB", Lat: 51.5074, Lon: -0.1278}, {Name: "London", Country: "CA"}}
		}, []string{"2 MATCHES", "▸ London  England, GB", "51.51, -0.13", "  London  CA"}},
		{"preferences", StatePreferences, func(m *Model) {}, []string{
			"✓ city ─── ● preferences ─── ○ sign in", "How do you like it?", "UNITS", "OPENS ON", "TIPS",
			"● metric  °C · km/h", "○ imperial  °F · mph", "PREVIEW", "9°C", "opens on Now",
			"› don't go out without an umbrella!", "step 2 of 3",
		}},
		{"preferences in imperial without tips", StatePreferences, func(m *Model) { m.UnitCursor, m.TipCursor = 1, 1 }, []string{"49°F", "15 mph"}},
		{"sign in", StateApiKeyOption, func(m *Model) {}, []string{
			"Last bit: how should gust get weather?", "Sign in with GitHub  recommended",
			"Use my own OpenWeatherMap key", "step 3 of 3",
		}},
		{"key input", StateApiKeyInput, func(m *Model) { m.ApiKeyCursor = 1 }, []string{"API KEY", "openweathermap.org/api/one-call-3", "saved to ~/.config/gust/auth.json"}},
		{"done", StateComplete, func(m *Model) {
			m.Config.DefaultCity, m.Config.Units, m.Config.DefaultView, m.Config.ShowTips = "London", "metric", "daily", true
			m.ApiKeyCursor = 1
			m.ApiKeyInput.SetValue("abcd1234")
		}, []string{
			"✓ city ─── ✓ preferences ─── ✓ sign in", "✓ You're set.", "city          London",
			"opens on      5-day", "tips          on", "weather via   your key ••••1234",
			"HANDY KEYS", "s settings", "enter  show me the weather", "done",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(&config.Config{}, false, &api.Client{})
			updated, _ := m.Update(tea.WindowSizeMsg{Width: 140, Height: 44})
			m = updated.(Model)
			m.State = tt.state
			tt.prepare(&m)

			view := ansi.Strip(m.View().Content)

			for _, want := range tt.want {
				assert.Contains(t, view, want)
			}
			for _, r := range view {
				assert.Less(t, r, rune(0x1F000), "emoji in view")
			}
		})
	}
}

func TestPreferencesPreviewHidesTheTipWhenTipsAreOff(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})
	m.State = StatePreferences
	m.TipCursor = 1

	assert.NotContains(t, ansi.Strip(m.View().Content), "umbrella")
}

func TestWeatherViaGitHub(t *testing.T) {
	m := NewModel(&config.Config{}, true, &api.Client{})
	m.State = StateComplete

	assert.Contains(t, ansi.Strip(m.View().Content), "weather via   GitHub")
}
