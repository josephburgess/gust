package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestSettingsOpenWithS(t *testing.T) {
	cfg := &config.Config{DefaultCity: "London", Units: "metric", DefaultView: "now", ShowTips: true, ApiUrl: "http://localhost:8080"}
	m := NewModel("London", cfg, TabNow, nil)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	updated, _ = updated.Update(splashMinMsg{})
	updated, _ = updated.Update(weatherMsg{resp: &api.WeatherResponse{City: fixtureCity, Weather: fixtureWeather}})
	updated, _ = updated.Update(tea.KeyPressMsg{Code: 's', Text: "s"})

	view := ansi.Strip(updated.View().Content)

	assert.Equal(t, screenSettings, updated.(Model).screen)
	for _, want := range []string{
		"gust › settings", "esc back to London, GB",
		"LOCATION", " ▸ default city      London",
		"DISPLAY", "units             ◀ metric ▶", "opens on          ◀ Now ▶", "tips              ◀ on ▶",
		"ACCOUNT", "weather via       GitHub",
		"ADVANCED", "api server        http://localhost:8080", "config file", "run setup again",
		"HELP", "The city gust opens on.", "enter change",
		"←→ change",
	} {
		assert.Contains(t, view, want)
	}
}

func TestSettingsEscGoesBack(t *testing.T) {
	m := NewModel("London", &config.Config{}, TabNow, nil)
	m.screen = screenSettings

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})

	assert.Equal(t, screenMain, updated.(Model).screen)
}

func TestSettingsCursorStopsAtTheEnds(t *testing.T) {
	m := NewModel("London", &config.Config{}, TabNow, nil)
	m.screen = screenSettings

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	assert.Equal(t, 0, updated.(Model).settingsCursor)

	for range len(settings) + 2 {
		updated, _ = updated.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	}
	assert.Equal(t, len(settings)-1, updated.(Model).settingsCursor)
}

func TestSettingsUnitsChangeSavesAndRefetches(t *testing.T) {
	cfg := &config.Config{DefaultCity: "London", Units: "metric"}
	m := NewModel("London", cfg, TabNow, nil)
	m.screen = screenSettings
	m.settingsCursor = 1
	m.loading = false

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	saved, err := config.Load()

	assert.NotNil(t, cmd)
	assert.True(t, updated.(Model).loading)
	assert.Equal(t, "imperial", cfg.Units)
	assert.Equal(t, "imperial", updated.(Model).units)
	assert.NoError(t, err)
	assert.Equal(t, "imperial", saved.Units)
}

func TestSettingsCycleRows(t *testing.T) {
	tests := []struct {
		name   string
		cursor int
		key    tea.KeyPressMsg
		check  func(*config.Config) any
		want   any
	}{
		{"opens on wraps to the end", 2, tea.KeyPressMsg{Code: tea.KeyLeft}, func(c *config.Config) any { return c.DefaultView }, "compact"},
		{"opens on moves right", 2, tea.KeyPressMsg{Code: 'l', Text: "l"}, func(c *config.Config) any { return c.DefaultView }, "hourly"},
		{"tips toggle", 3, tea.KeyPressMsg{Code: tea.KeyRight}, func(c *config.Config) any { return c.ShowTips }, false},
		{"enter cycles too", 3, tea.KeyPressMsg{Code: tea.KeyEnter}, func(c *config.Config) any { return c.ShowTips }, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{DefaultCity: "London", Units: "metric", DefaultView: "now", ShowTips: true}
			m := NewModel("London", cfg, TabNow, nil)
			m.screen = screenSettings
			m.settingsCursor = tt.cursor

			m.Update(tt.key)

			assert.Equal(t, tt.want, tt.check(cfg))
		})
	}
}

func TestSettingsCityOpensTheSetupCityScreen(t *testing.T) {
	cfg := &config.Config{DefaultCity: "London"}
	m := NewModel("London", cfg, TabNow, nil)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	m = updated.(Model)
	m.screen = screenSettings

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	assert.Contains(t, ansi.Strip(updated.View().Content), "Where's home?")

	updated, _ = updated.Update(tea.KeyPressMsg{Code: tea.KeyEscape})

	assert.Nil(t, updated.(Model).sub)
	assert.Equal(t, "London", cfg.DefaultCity)
	assert.Contains(t, ansi.Strip(updated.View().Content), "gust › settings")
}

func TestSettingsWeatherViaOpensTheSignInScreen(t *testing.T) {
	m := NewModel("London", &config.Config{}, TabNow, nil)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	m = updated.(Model)
	m.screen = screenSettings
	m.settingsCursor = 4

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	assert.Contains(t, ansi.Strip(updated.View().Content), "Sign in with GitHub")
}

func TestSettingsConfigFileOpensTheEditor(t *testing.T) {
	t.Setenv("EDITOR", "true")
	m := NewModel("London", &config.Config{}, TabNow, nil)
	m.screen = screenSettings
	m.settingsCursor = 6

	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	assert.NotNil(t, cmd)
}

func TestEditorDoneReloadsTheConfig(t *testing.T) {
	assert.NoError(t, (&config.Config{DefaultCity: "Paris", Units: "imperial"}).Save())
	cfg := &config.Config{DefaultCity: "London", Units: "metric"}
	m := NewModel("London", cfg, TabNow, nil)
	m.screen = screenSettings

	updated, _ := m.Update(editorDoneMsg{})

	assert.Equal(t, "Paris", cfg.DefaultCity)
	assert.Equal(t, "imperial", updated.(Model).units)
}

func TestSettingValueWeatherVia(t *testing.T) {
	tests := []struct {
		name string
		auth *config.AuthConfig
		want string
	}{
		{"own key", &config.AuthConfig{APIKey: "abcd1234"}, "your key ••••1234"},
		{"github", &config.AuthConfig{APIKey: "gust_abc", GithubUser: "octocat"}, "GitHub (octocat)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel("London", &config.Config{}, TabNow, nil)
			m.auth = tt.auth
			assert.Equal(t, tt.want, m.settingValue(settingVia))
		})
	}
}
