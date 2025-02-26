package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

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

func fetchJSON(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

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

func GetCurrentWeather(lat, lon float64, apiKey string) (*models.Weather, error) {
	requestURL := fmt.Sprintf(weatherURL, lat, lon, apiKey)

	var result struct {
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
			SeaLevel  int     `json:"sea_level,omitempty"`
			GrndLevel int     `json:"grnd_level,omitempty"`
		} `json:"main"`
		Visibility int `json:"visibility"`
		Wind       struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
			Gust  float64 `json:"gust,omitempty"`
		} `json:"wind"`
		Rain struct {
			OneHour float64 `json:"1h,omitempty"`
		} `json:"rain,omitempty"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Sys struct {
			Sunrise int64  `json:"sunrise"`
			Sunset  int64  `json:"sunset"`
			Country string `json:"country"`
		} `json:"sys"`
		Name string `json:"name"`
	}

	if err := fetchJSON(requestURL, &result); err != nil {
		return nil, err
	}

	if len(result.Weather) == 0 {
		return nil, fmt.Errorf("no weather data available")
	}

	weather := &models.Weather{
		ID:          result.Weather[0].ID,
		Icon:        result.Weather[0].Icon,
		Description: result.Weather[0].Description,
		Temp:        result.Main.Temp,
		FeelsLike:   result.Main.FeelsLike,
		TempMin:     result.Main.TempMin,
		TempMax:     result.Main.TempMax,
		Humidity:    result.Main.Humidity,
		Pressure:    result.Main.Pressure,
		Visibility:  result.Visibility,
		WindSpeed:   result.Wind.Speed,
		WindDeg:     result.Wind.Deg,
		WindGust:    result.Wind.Gust,
		Rain1h:      result.Rain.OneHour,
		Clouds:      result.Clouds.All,
		Sunrise:     result.Sys.Sunrise,
		Sunset:      result.Sys.Sunset,
	}

	return weather, nil
}

// GetForecast gets weather forecast data
func GetForecast(lat, lon float64, apiKey string) ([]models.ForecastItem, error) {
	requestURL := fmt.Sprintf(forecastURL, lat, lon, apiKey)

	var result struct {
		List []struct {
			Dt   int64 `json:"dt"`
			Main struct {
				Temp      float64 `json:"temp"`
				FeelsLike float64 `json:"feels_like"`
				TempMin   float64 `json:"temp_min"`
				TempMax   float64 `json:"temp_max"`
				Pressure  int     `json:"pressure"`
				Humidity  int     `json:"humidity"`
			} `json:"main"`
			Weather []struct {
				ID          int    `json:"id"`
				Main        string `json:"main"`
				Description string `json:"description"`
				Icon        string `json:"icon"`
			} `json:"weather"`
			Wind struct {
				Speed float64 `json:"speed"`
				Deg   int     `json:"deg"`
			} `json:"wind"`
			Pop float64 `json:"pop"`
		} `json:"list"`
	}

	if err := fetchJSON(requestURL, &result); err != nil {
		return nil, err
	}

	var forecasts []models.ForecastItem

	// Get one forecast per day (using noon forecast)
	uniqueDays := make(map[string]bool)
	for _, item := range result.List {
		// Skip if there's no weather data
		if len(item.Weather) == 0 {
			continue
		}

		// Get the day string
		day := time.Unix(item.Dt, 0).Format("2006-01-02")

		// Only take the first forecast of each day
		if !uniqueDays[day] {
			uniqueDays[day] = true
			forecast := models.ForecastItem{
				DateTime:    item.Dt,
				TempMin:     item.Main.TempMin,
				TempMax:     item.Main.TempMax,
				WeatherID:   item.Weather[0].ID,
				Icon:        item.Weather[0].Icon,
				Description: item.Weather[0].Description,
				Humidity:    item.Main.Humidity,
				WindSpeed:   item.Wind.Speed,
				WindDeg:     item.Wind.Deg,
				Pop:         item.Pop,
			}
			forecasts = append(forecasts, forecast)

			// Limit to 5 days
			if len(forecasts) >= 5 {
				break
			}
		}
	}

	return forecasts, nil
}

// GetWeatherAndForecast gets both current weather and forecast data
func GetWeatherAndForecast(lat, lon float64, apiKey string) (*models.Weather, []models.ForecastItem, error) {
	weather, err := GetCurrentWeather(lat, lon, apiKey)
	if err != nil {
		return nil, nil, err
	}

	forecast, err := GetForecast(lat, lon, apiKey)
	if err != nil {
		return weather, nil, err
	}

	return weather, forecast, nil
}

// Legacy function to maintain compatibility
func GetWeather(lat, lon float64, apiKey string) (*models.Weather, error) {
	return GetCurrentWeather(lat, lon, apiKey)
}
