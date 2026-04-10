package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCache(t *testing.T) *Cache {
	t.Helper()
	dir := t.TempDir()
	return &Cache{dir: dir}
}

func testWeatherResponse() *api.WeatherResponse {
	return &api.WeatherResponse{
		City: &models.City{
			Name:    "London",
			Country: "GB",
			Lat:     51.5074,
			Lon:     -0.1278,
		},
		Weather: &models.OneCallResponse{
			Lat:      51.5074,
			Lon:      -0.1278,
			Timezone: "Europe/London",
			Current: models.CurrentWeather{
				Temp:      15.5,
				FeelsLike: 14.0,
				Humidity:  72,
				Weather: []models.WeatherCondition{
					{ID: 800, Main: "Clear", Description: "clear sky"},
				},
			},
		},
	}
}

func TestCache_SetAndGet(t *testing.T) {
	c := newTestCache(t)
	data := testWeatherResponse()

	err := c.Set("London", "metric", data)
	require.NoError(t, err)

	got, age, ok := c.Get("London", "metric")
	require.True(t, ok)
	assert.Less(t, age, time.Second)
	assert.Equal(t, data.City.Name, got.City.Name)
	assert.Equal(t, data.Weather.Current.Temp, got.Weather.Current.Temp)
}

func TestCache_Miss(t *testing.T) {
	c := newTestCache(t)

	_, _, ok := c.Get("Paris", "metric")
	assert.False(t, ok)
}

func TestCache_Expired(t *testing.T) {
	c := newTestCache(t)
	data := testWeatherResponse()

	// write an entry with a CachedAt well in the past
	e := entry{
		Data:     data,
		CachedAt: time.Now().Add(-(TTL + time.Minute)),
	}
	b, err := json.Marshal(e)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(c.filePath("London", "metric"), b, 0644))

	_, _, ok := c.Get("London", "metric")
	assert.False(t, ok)

	// expired file should be deleted
	_, err = os.Stat(c.filePath("London", "metric"))
	assert.True(t, os.IsNotExist(err))
}

func TestCache_CorruptFile(t *testing.T) {
	c := newTestCache(t)
	path := c.filePath("London", "metric")

	require.NoError(t, os.WriteFile(path, []byte("not valid json"), 0644))

	_, _, ok := c.Get("London", "metric")
	assert.False(t, ok)
}

func TestCache_NormalisesKey(t *testing.T) {
	c := newTestCache(t)
	data := testWeatherResponse()

	// set with mixed case and spaces
	require.NoError(t, c.Set("New York", "imperial", data))

	// file should exist with normalised name
	expected := filepath.Join(c.dir, "new_york-imperial.json")
	_, err := os.Stat(expected)
	assert.NoError(t, err)
}

func TestCache_EmptyUnits(t *testing.T) {
	c := newTestCache(t)
	data := testWeatherResponse()

	require.NoError(t, c.Set("Tokyo", "", data))

	expected := filepath.Join(c.dir, "tokyo.json")
	_, err := os.Stat(expected)
	assert.NoError(t, err)

	got, _, ok := c.Get("Tokyo", "")
	require.True(t, ok)
	assert.Equal(t, "London", got.City.Name)
}

func TestCache_DifferentUnitsAreIndependent(t *testing.T) {
	c := newTestCache(t)

	metric := testWeatherResponse()
	metric.Weather.Current.Temp = 15.0

	imperial := testWeatherResponse()
	imperial.Weather.Current.Temp = 59.0

	require.NoError(t, c.Set("London", "metric", metric))
	require.NoError(t, c.Set("London", "imperial", imperial))

	gotMetric, _, ok := c.Get("London", "metric")
	require.True(t, ok)
	assert.Equal(t, 15.0, gotMetric.Weather.Current.Temp)

	gotImperial, _, ok := c.Get("London", "imperial")
	require.True(t, ok)
	assert.Equal(t, 59.0, gotImperial.Weather.Current.Temp)
}

func TestFormatAge(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{0, "just now"},
		{30 * time.Second, "just now"},
		{59 * time.Second, "just now"},
		{time.Minute, "1m ago"},
		{90 * time.Second, "1m ago"},
		{2 * time.Minute, "2m ago"},
		{9*time.Minute + 59*time.Second, "9m ago"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, FormatAge(tt.d), "FormatAge(%v)", tt.d)
	}
}
