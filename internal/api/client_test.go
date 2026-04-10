package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client := NewClient("https://example.com", "test-key", "metric")
	assert.Equal(t, "https://example.com", client.baseURL)
	assert.Equal(t, "test-key", client.apiKey)
	assert.NotNil(t, client.client)
}

func TestGetWeather_OK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/weather/London", r.URL.Path)
		assert.Equal(t, "test-api-key", r.URL.Query().Get("api_key"))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"city":{"name":"London","lat":51.5,"lon":-0.1},"weather":{"lat":51.5,"lon":-0.1,"current":{"temp":283.15,"weather":[{"id":800,"description":"clear sky"}]}}}`))
	}))
	defer server.Close()

	resp, err := NewClient(server.URL, "test-api-key", "metric").GetWeather("London")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "London", resp.City.Name)
	assert.Equal(t, 283.15, resp.Weather.Current.Temp)
	assert.Equal(t, "clear sky", resp.Weather.Current.Weather[0].Description)
}

func TestGetWeather_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := NewClient(server.URL, "key", "metric").GetWeather("Nowhere")
	require.Error(t, err)
	var notFound *CityNotFoundError
	assert.ErrorAs(t, err, &notFound)
	assert.Equal(t, "Nowhere", notFound.City)
}

func TestGetWeather_RateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("rate limit exceeded"))
	}))
	defer server.Close()

	_, err := NewClient(server.URL, "key", "metric").GetWeather("London")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit")
}

func TestExtractRateLimitInfo(t *testing.T) {
	resetTime := time.Now().Add(time.Hour).UTC().Truncate(time.Second)

	t.Run("all headers present", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-RateLimit-Limit", "50")
			w.Header().Set("X-RateLimit-Remaining", "42")
			w.Header().Set("X-RateLimit-Reset", resetTime.Format(time.RFC3339))
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"city":{"name":"X"},"weather":{}}`))
		}))
		defer server.Close()

		client := NewClient(server.URL, "key", "metric")
		client.GetWeather("X")

		assert.Equal(t, 50, client.RateLimitInfo.Limit)
		assert.Equal(t, 42, client.RateLimitInfo.Remaining)
		assert.Equal(t, resetTime, client.RateLimitInfo.ResetTime)
	})

	t.Run("missing headers leave values zero", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"city":{"name":"X"},"weather":{}}`))
		}))
		defer server.Close()

		client := NewClient(server.URL, "key", "metric")
		client.GetWeather("X")

		assert.Equal(t, 0, client.RateLimitInfo.Limit)
		assert.Equal(t, 0, client.RateLimitInfo.Remaining)
	})

	t.Run("malformed reset time falls back to one hour from now", func(t *testing.T) {
		before := time.Now()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-RateLimit-Reset", "not-a-time")
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"city":{"name":"X"},"weather":{}}`))
		}))
		defer server.Close()

		client := NewClient(server.URL, "key", "metric")
		client.GetWeather("X")

		assert.True(t, client.RateLimitInfo.ResetTime.After(before.Add(55*time.Minute)))
	})
}
