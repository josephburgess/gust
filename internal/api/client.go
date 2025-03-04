package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/josephburgess/gust/internal/models"
)

type WeatherResponse struct {
	City    *models.City            `json:"city"`
	Weather *models.OneCallResponse `json:"weather"`
}

type Client struct {
	baseURL string
	apiKey  string
	units   string
	client  *http.Client
}

func NewClient(baseURL, apiKey string, units string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		units:   units,
		client:  &http.Client{},
	}
}

func (c *Client) GetWeather(cityName string) (*WeatherResponse, error) {
	endpoint := fmt.Sprintf(
		"%s/api/weather/%s?api_key=%s",
		c.baseURL,
		url.QueryEscape(cityName),
		c.apiKey,
	)

	if c.units != "" {
		endpoint = fmt.Sprintf("%s&units=%s", endpoint, c.units)
	}

	resp, err := c.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var response WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return &response, nil
}

func (c *Client) SearchCities(query string) ([]models.City, error) {
	endpoint := fmt.Sprintf(
		"%s/api/cities/search?q=%s",
		c.baseURL,
		url.QueryEscape(query),
	)

	resp, err := c.client.Get(endpoint)
	if err != nil {
		fmt.Println("API request failed:", err)
		return nil, fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var cities []models.City
	if err := json.NewDecoder(resp.Body).Decode(&cities); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return cities, nil
}
