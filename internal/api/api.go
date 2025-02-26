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
	baseURL     = "https://api.openweathermap.org/"
	geoCodeURL  = baseURL + "geo/1.0/direct?q=%s&limit=1&appid=%s"
	weatherURL  = baseURL + "data/2.5/weather?lat=%f&lon=%f&appid=%s"
	forecastURL = baseURL + "data/2.5/forecast?lat=%f&lon=%f&appid=%s"
)

func SetBaseURL(u string) {
	baseURL = u
	geoCodeURL = baseURL + "geo/1.0/direct?q=%s&limit=1&appid=%s"
	weatherURL = baseURL + "data/2.5/weather?lat=%f&lon=%f&appid=%s"
	forecastURL = baseURL + "data/2.5/forecast?lat=%f&lon=%f&appid=%s"
}

func GetForecast(lat, lon float64, apiKey string) ([]models.ForecastItem, error) {
	requestURL := fmt.Sprintf(forecastURL, lat, lon, apiKey)

	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("forecast request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading forecast response body: %w", err)
	}

	// { "list": [ { "dt_txt": "...", "main": { "temp": ... }, "weather": [{"description": "..."}] }, ... ] }
	var result struct {
		List []struct {
			DtTxt string `json:"dt_txt"`
			Main  struct {
				Temp float64 `json:"temp"`
			} `json:"main"`
			Weather []struct {
				Description string `json:"description"`
			} `json:"weather"`
		} `json:"list"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling forecast JSON: %w", err)
	}

	var forecast []models.ForecastItem
	for _, item := range result.List {
		if len(item.Weather) == 0 {
			continue
		}
		forecast = append(forecast, models.ForecastItem{
			DateTime:    item.DtTxt,
			Temp:        item.Main.Temp,
			Description: item.Weather[0].Description,
		})
	}

	return forecast, nil
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
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	}

	if err := fetchJSON(requestURL, &result); err != nil {
		return nil, err
	}

	if len(result.Weather) == 0 {
		return nil, fmt.Errorf("no weather data available")
	}

	first := result.Weather[0]
	w := &models.Weather{
		ID:          first.ID,
		Icon:        first.Icon,
		Temp:        result.Main.Temp,
		Description: first.Description,
	}
	return w, nil
}
