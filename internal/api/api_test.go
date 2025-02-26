package api

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/josephburgess/gust/internal/models"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestGetCoordinates(t *testing.T) {
	originalClient := http.DefaultClient
	defer func() { http.DefaultClient = originalClient }()

	tests := []struct {
		name           string
		city           string
		apiKey         string
		mockResponse   string
		mockStatusCode int
		wantCity       *models.City
		wantErr        bool
	}{
		{
			name:           "successful request",
			city:           "London",
			apiKey:         "test-api-key",
			mockResponse:   `[{"name":"London","lat":51.5074,"lon":-0.1278}]`,
			mockStatusCode: http.StatusOK,
			wantCity: &models.City{
				Name: "London",
				Lat:  51.5074,
				Lon:  -0.1278,
			},
			wantErr: false,
		},
		{
			name:           "city not found",
			city:           "NonExistentCity",
			apiKey:         "test-api-key",
			mockResponse:   `[]`,
			mockStatusCode: http.StatusOK,
			wantCity:       nil,
			wantErr:        true,
		},
		{
			name:           "api error",
			city:           "London",
			apiKey:         "invalid-key",
			mockResponse:   `{"cod":401, "message": "Invalid API key"}`,
			mockStatusCode: http.StatusUnauthorized,
			wantCity:       nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			http.DefaultClient = &http.Client{
				Transport: &mockTransport{
					mockResponse:   tt.mockResponse,
					mockStatusCode: tt.mockStatusCode,
				},
			}

			got, err := GetCoordinates(tt.city, tt.apiKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetCoordinates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Name != tt.wantCity.Name ||
					got.Lat != tt.wantCity.Lat ||
					got.Lon != tt.wantCity.Lon {
					t.Errorf("GetCoordinates() got = %v, want %v", got, tt.wantCity)
				}
			}
		})
	}
}

func TestGetWeather(t *testing.T) {
	originalClient := http.DefaultClient
	defer func() { http.DefaultClient = originalClient }()

	tests := []struct {
		name           string
		lat            float64
		lon            float64
		apiKey         string
		mockResponse   string
		mockStatusCode int
		wantWeather    *models.Weather
		wantErr        bool
	}{
		{
			name:           "successful request",
			lat:            51.5074,
			lon:            -0.1278,
			apiKey:         "test-api-key",
			mockResponse:   `{"main":{"temp":283.15},"weather":[{"description":"cloudy"}]}`,
			mockStatusCode: http.StatusOK,
			wantWeather: &models.Weather{
				Temp:        283.15,
				Description: "cloudy",
			},
			wantErr: false,
		},
		{
			name:           "missing weather data",
			lat:            51.5074,
			lon:            -0.1278,
			apiKey:         "test-api-key",
			mockResponse:   `{"main":{"temp":283.15},"weather":[]}`,
			mockStatusCode: http.StatusOK,
			wantWeather:    nil,
			wantErr:        true,
		},
		{
			name:           "api error",
			lat:            51.5074,
			lon:            -0.1278,
			apiKey:         "invalid-key",
			mockResponse:   `{"cod":401, "message": "Invalid API key"}`,
			mockStatusCode: http.StatusUnauthorized,
			wantWeather:    nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			http.DefaultClient = &http.Client{
				Transport: &mockTransport{
					mockResponse:   tt.mockResponse,
					mockStatusCode: tt.mockStatusCode,
				},
			}

			got, err := GetWeather(tt.lat, tt.lon, tt.apiKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetWeather() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Temp != tt.wantWeather.Temp ||
					got.Description != tt.wantWeather.Description {
					t.Errorf("GetWeather() got = %v, want %v", got, tt.wantWeather)
				}
			}
		})
	}
}

type mockTransport struct {
	mockResponse   string
	mockStatusCode int
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	response := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: m.mockStatusCode,
	}
	response.Header.Set("Content-Type", "application/json")
	response.Body = io.NopCloser(strings.NewReader(m.mockResponse))
	return response, nil
}
