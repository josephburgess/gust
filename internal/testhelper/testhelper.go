package testhelper

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MockWeatherServer(t *testing.T) (*httptest.Server, string) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		path := r.URL.Path

		switch path {
		case "/geo/1.0/direct":
			city := query.Get("q")
			apiKey := query.Get("appid")

			if apiKey != "test-api-key" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]any{
					"cod":     401,
					"message": "Invalid API key",
				})
				return
			}

			switch city {
			case "London":
				json.NewEncoder(w).Encode([]map[string]any{
					{
						"name": "London",
						"lat":  51.5074,
						"lon":  -0.1278,
					},
				})
			case "NonExistentCity":
				json.NewEncoder(w).Encode([]map[string]any{})
			default:
				json.NewEncoder(w).Encode([]map[string]any{
					{
						"name": city,
						"lat":  40.7128,
						"lon":  -74.0060,
					},
				})
			}

		case "/data/2.5/weather":
			lat := query.Get("lat")
			lon := query.Get("lon")
			apiKey := query.Get("appid")

			if apiKey != "test-api-key" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]any{
					"cod":     401,
					"message": "Invalid API key",
				})
				return
			}

			switch {
			case lat == "51.5074" && lon == "-0.1278":
				json.NewEncoder(w).Encode(map[string]any{
					"main": map[string]any{
						"temp": 283.15,
					},
					"weather": []map[string]any{
						{
							"description": "cloudy",
						},
					},
				})
			case lat == "40.7128" && lon == "-74.0060":
				json.NewEncoder(w).Encode(map[string]any{
					"main": map[string]any{
						"temp": 288.15,
					},
					"weather": []map[string]any{
						{
							"description": "partly cloudy",
						},
					},
				})
			case lat == "0" && lon == "0":
				json.NewEncoder(w).Encode(map[string]any{
					"main": map[string]any{
						"temp": 0,
					},
					"weather": []map[string]any{},
				})
			default:
				json.NewEncoder(w).Encode(map[string]any{
					"main": map[string]any{
						"temp": 293.15,
					},
					"weather": []map[string]any{
						{
							"description": "clear sky",
						},
					},
				})
			}

		default:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]any{
				"cod":     404,
				"message": "API endpoint not found",
			})
		}
	}))

	return server, server.URL
}

func CaptureOutput(f func()) (string, error) {
	return "Output captured", nil
}
