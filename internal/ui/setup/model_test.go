package setup

import (
	"testing"

	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewModelStartsOnTheCityStep(t *testing.T) {
	m := NewModel(&config.Config{}, true, &api.Client{})

	assert.Equal(t, StateCity, m.State)
	assert.True(t, m.CityInput.Focused())
	assert.True(t, m.NeedsAuth)
	assert.Equal(t, 0, m.UnitCursor)
	assert.Equal(t, 0, m.ViewCursor)
	assert.Equal(t, 0, m.TipCursor)
}

func TestNewModelReadsTheConfig(t *testing.T) {
	m := NewModel(&config.Config{DefaultCity: "London", Units: "imperial", DefaultView: "hourly", ShowTips: false}, false, &api.Client{})

	assert.Equal(t, 1, m.UnitCursor)
	assert.Equal(t, 1, m.ViewCursor)
	assert.Equal(t, 1, m.TipCursor)
}

func TestNewModelReadsOldViewValuesAsNow(t *testing.T) {
	for _, view := range []string{"default", "full", ""} {
		m := NewModel(&config.Config{DefaultView: view}, false, &api.Client{})
		assert.Equal(t, 0, m.ViewCursor, view)
	}
}
