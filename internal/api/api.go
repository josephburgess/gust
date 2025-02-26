package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/josephburgess/gust/internal/models"
)

var (
	baseURL    = "https://api.openweathermap.org/"
	geoCodeURL = baseURL + "geo/1.0/direct?q=%s&limit=1&appid=%s"
	weatherURL = baseURL + "data/2.5/weather?lat=%f&lon=%f&appid=%s"
)

// GetBaseURL returns the current base URL (for testing)
func GetBaseURL() string {
	return baseURL
}

// SetBaseURL sets the base URL (for testing)
func SetBaseURL(url string) {
	baseURL = url
	geoCodeURL = baseURL + "geo/1.0/direct?q=%s&limit=1&appid=%s"
	weatherURL = baseURL + "data/2.5/weather?lat=%f&lon=%f&appid=%s"
}

func fetchJSON(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("unmarshaling JSON: %w", err)
	}

	return nil
}

func GetCoordinates(city, apiKey string) (*models.City, error) {
	encodedCity := url.QueryEscape(city)
	requestURL := fmt.Sprintf(geoCodeURL, encodedCity, apiKey)

	var cities []models.City
	if err := fetchJSON(requestURL, &cities); err != nil {
		return nil, err
	}

	if len(cities) == 0 {
		return nil, fmt.Errorf("no coordinates found for %s", city)
	}

	return &cities[0], nil
}

func GetWeather(lat, lon float64, apiKey string) (*models.Weather, error) {
	requestURL := fmt.Sprintf(weatherURL, lat, lon, apiKey)

	var result struct {
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		} `json:"wind"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Visibility int `json:"visibility"`
		Sys        struct {
			Sunrise int64 `json:"sunrise"`
			Sunset  int64 `json:"sunset"`
		} `json:"sys"`
	}

	if err := fetchJSON(requestURL, &result); err != nil {
		return nil, err
	}

	if len(result.Weather) == 0 {
		return nil, fmt.Errorf("no weather data available")
	}

	weather := &models.Weather{
		Temp:        result.Main.Temp,
		Description: result.Weather[0].Description,
		FeelsLike:   result.Main.FeelsLike,
		TempMin:     result.Main.TempMin,
		TempMax:     result.Main.TempMax,
		Humidity:    result.Main.Humidity,
		WindSpeed:   result.Wind.Speed,
		WindDeg:     result.Wind.Deg,
		CloudsAll:   result.Clouds.All,
		Visibility:  result.Visibility,
		Pressure:    result.Main.Pressure,
		Sunrise:     result.Sys.Sunrise,
		Sunset:      result.Sys.Sunset,
	}

	return weather, nil
}
