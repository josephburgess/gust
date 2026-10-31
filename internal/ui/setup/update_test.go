package setup

import (
	"errors"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/stretchr/testify/assert"
)

var (
	enter   = tea.KeyPressMsg{Code: tea.KeyEnter}
	esc     = tea.KeyPressMsg{Code: tea.KeyEscape}
	down    = tea.KeyPressMsg{Code: tea.KeyDown}
	up      = tea.KeyPressMsg{Code: tea.KeyUp}
	tab     = tea.KeyPressMsg{Code: tea.KeyTab}
	londons = []models.City{{Name: "London", State: "England", Country: "GB"}, {Name: "London", State: "Ontario", Country: "CA"}}
)

func TestCityEnterSearches(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})
	m.CityInput.SetValue("london")

	updated, cmd := m.Update(enter)

	assert.Equal(t, StateCitySearch, updated.(Model).State)
	assert.Equal(t, "london", updated.(Model).CitySearchQuery)
	assert.NotNil(t, cmd)
}

func TestCityEnterWithNoTextDoesNothing(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})

	updated, cmd := m.Update(enter)

	assert.Equal(t, StateCity, updated.(Model).State)
	assert.Nil(t, cmd)
}

func TestCitySearchResult(t *testing.T) {
	tests := []struct {
		name      string
		result    CitiesSearchResult
		wantState SetupState
		wantErr   string
	}{
		{"matches", CitiesSearchResult{cities: londons}, StateCitySelect, ""},
		{"no matches", CitiesSearchResult{cities: []models.City{}}, StateCity, `no matches for "london"`},
		{"error", CitiesSearchResult{err: errors.New("boom")}, StateCity, "boom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(&config.Config{}, false, &api.Client{})
			m.State = StateCitySearch
			m.CitySearchQuery = "london"

			updated, _ := m.Update(tt.result)

			assert.Equal(t, tt.wantState, updated.(Model).State)
			assert.Equal(t, tt.wantErr, updated.(Model).SearchErr)
		})
	}
}

func TestCitySearchResultIsIgnoredAfterTheSearchEnds(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})
	m.State = StatePreferences

	updated, _ := m.Update(CitiesSearchResult{cities: londons})

	assert.Equal(t, StatePreferences, updated.(Model).State)
}

func TestCitySelectArrowsMoveTheCursor(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})
	m.State = StateCitySelect
	m.CityOptions = londons

	updated, _ := m.Update(down)
	updated, _ = updated.Update(down)
	assert.Equal(t, 1, updated.(Model).CityCursor)

	updated, _ = updated.Update(up)
	updated, _ = updated.Update(up)
	assert.Equal(t, 0, updated.(Model).CityCursor)
}

func TestCitySelectEnterPicksTheCity(t *testing.T) {
	cfg := &config.Config{}
	m := NewModel(cfg, false, &api.Client{})
	m.State = StateCitySelect
	m.CitySearchQuery = "london"
	m.CityInput.SetValue("london")
	m.CityOptions = londons
	m.CityCursor = 1

	updated, _ := m.Update(enter)

	assert.Equal(t, StatePreferences, updated.(Model).State)
	assert.Equal(t, "London", cfg.DefaultCity)
}

func TestCitySelectSearchesAgainWhenTheTextChanges(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})
	m.State = StateCitySelect
	m.CitySearchQuery = "london"
	m.CityInput.SetValue("paris")
	m.CityOptions = londons

	updated, cmd := m.Update(enter)

	assert.Equal(t, StateCitySearch, updated.(Model).State)
	assert.Equal(t, "paris", updated.(Model).CitySearchQuery)
	assert.NotNil(t, cmd)
}

func TestSearchCitiesWithoutAClient(t *testing.T) {
	m := NewModel(&config.Config{}, false, nil)
	m.CitySearchQuery = "London"

	result := m.searchCities()().(CitiesSearchResult)

	assert.ErrorContains(t, result.err, "API client not initialized")
}

func TestPreferencesTabMovesTheFocus(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})
	m.State = StatePreferences

	updated, _ := m.Update(tab)
	assert.Equal(t, GroupView, updated.(Model).PrefGroup)
	updated, _ = updated.Update(tab)
	updated, _ = updated.Update(tab)
	assert.Equal(t, GroupUnits, updated.(Model).PrefGroup)
	updated, _ = updated.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	assert.Equal(t, GroupTips, updated.(Model).PrefGroup)
}

func TestPreferencesArrowsChangeTheFocusedGroup(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})
	m.State = StatePreferences
	m.PrefGroup = GroupView

	updated, _ := m.Update(down)
	updated, _ = updated.Update(down)
	updated, _ = updated.Update(down)
	updated, _ = updated.Update(down)

	assert.Equal(t, 3, updated.(Model).ViewCursor)
	assert.Equal(t, 0, updated.(Model).UnitCursor)
}

func TestPreferencesEnterSavesTheChoices(t *testing.T) {
	cfg := &config.Config{}
	m := NewModel(cfg, false, &api.Client{})
	m.State = StatePreferences
	m.UnitCursor, m.ViewCursor, m.TipCursor = 1, 2, 1

	updated, _ := m.Update(enter)

	assert.Equal(t, StateApiKeyOption, updated.(Model).State)
	assert.Equal(t, "imperial", cfg.Units)
	assert.Equal(t, "daily", cfg.DefaultView)
	assert.False(t, cfg.ShowTips)
}

func TestPreferencesEscGoesBack(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})
	m.State = StatePreferences

	updated, _ := m.Update(esc)

	assert.Equal(t, StateCity, updated.(Model).State)
}

func TestSignInWithGitHub(t *testing.T) {
	tests := []struct {
		name       string
		needsAuth  bool
		wantGitHub bool
	}{
		{"signs in when there is no key", true, true},
		{"keeps the key when there is one", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(&config.Config{}, tt.needsAuth, &api.Client{})
			m.State = StateApiKeyOption

			updated, _ := m.Update(enter)

			assert.Equal(t, StateComplete, updated.(Model).State)
			assert.Equal(t, tt.wantGitHub, updated.(Model).WantsGitHub)
		})
	}
}

func TestSignInWithOwnKeyOpensTheKeyInput(t *testing.T) {
	m := NewModel(&config.Config{}, true, &api.Client{})
	m.State = StateApiKeyOption

	updated, _ := m.Update(down)
	updated, _ = updated.Update(enter)

	assert.Equal(t, StateApiKeyInput, updated.(Model).State)
	assert.True(t, updated.(Model).ApiKeyInput.Focused())
}

func TestApiKeyInputSavesTheKey(t *testing.T) {
	m := NewModel(&config.Config{ApiUrl: "http://localhost:8080"}, true, &api.Client{})
	m.State = StateApiKeyInput
	m.ApiKeyCursor = 1
	m.ApiKeyInput.SetValue("abcd1234")

	updated, _ := m.Update(enter)
	saved, err := config.LoadAuthConfig()

	assert.Equal(t, StateComplete, updated.(Model).State)
	assert.False(t, updated.(Model).NeedsAuth)
	assert.NoError(t, err)
	assert.Equal(t, "abcd1234", saved.APIKey)
}

func TestApiKeyInputIgnoresAnEmptyKey(t *testing.T) {
	m := NewModel(&config.Config{}, true, &api.Client{})
	m.State = StateApiKeyInput

	updated, _ := m.Update(enter)

	assert.Equal(t, StateApiKeyInput, updated.(Model).State)
}

func TestCompleteEnterQuits(t *testing.T) {
	m := NewModel(&config.Config{}, false, &api.Client{})
	m.State = StateComplete

	_, cmd := m.Update(enter)

	assert.Equal(t, tea.QuitMsg{}, cmd())
}

func TestCompleteEscGoesBack(t *testing.T) {
	m := NewModel(&config.Config{}, true, &api.Client{})
	m.State = StateComplete
	m.WantsGitHub = true

	updated, _ := m.Update(esc)

	assert.Equal(t, StateApiKeyOption, updated.(Model).State)
	assert.False(t, updated.(Model).WantsGitHub)
}

func TestCtrlCQuitsFromAnyStep(t *testing.T) {
	for _, state := range []SetupState{StateCity, StatePreferences, StateApiKeyInput, StateComplete} {
		m := NewModel(&config.Config{}, false, &api.Client{})
		m.State = state

		updated, cmd := m.Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})

		assert.True(t, updated.(Model).Quitting)
		assert.Equal(t, tea.QuitMsg{}, cmd())
	}
}
