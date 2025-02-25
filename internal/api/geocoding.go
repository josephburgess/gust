package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/josephburgess/gust/internal/models"
)

const (
	baseURL    = "https://api.openweathermap.org/"
	geoCodeURL = baseURL + "geo/1.0/direct?q=%s&limit=1&appid=%s"
)

func GetCoordinates(city, apiKey string) (*models.City, error) {
	encodedCity := url.QueryEscape(city)
	requestURL := fmt.Sprintf(geoCodeURL, encodedCity, apiKey)

	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cities []models.City
	err = json.Unmarshal(body, &cities)
	if err != nil {
		return nil, err
	}

	if len(cities) == 0 {
		return nil, fmt.Errorf("no coordinates found for %s", city)
	}

	return &cities[0], nil
}
