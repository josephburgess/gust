package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/josephburgess/gust/internal/models"
)

type WeatherResponse struct {
	City    *models.City            `json:"city"`
	Weather *models.OneCallResponse `json:"weather"`
}

type RateLimitInfo struct {
	Limit     int
	Remaining int
	ResetTime time.Time
}

type Client struct {
	baseURL       string
	apiKey        string
	units         string
	client        *http.Client
	RateLimitInfo *RateLimitInfo
}

func NewClient(baseURL, apiKey string, units string) *Client {
	return &Client{
		baseURL:       baseURL,
		apiKey:        apiKey,
		units:         units,
		client:        &http.Client{},
		RateLimitInfo: &RateLimitInfo{},
	}
}

func (c *Client) extractRateLimitInfo(resp *http.Response) {
	if c.RateLimitInfo == nil {
		c.RateLimitInfo = &RateLimitInfo{}
	}

	if limit := resp.Header.Get("X-RateLimit-Limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			c.RateLimitInfo.Limit = val
		}
	}

	if remaining := resp.Header.Get("X-RateLimit-Remaining"); remaining != "" {
		if val, err := strconv.Atoi(remaining); err == nil {
			c.RateLimitInfo.Remaining = val
		} else {
			c.RateLimitInfo.Remaining = 0
		}
	}

	if reset := resp.Header.Get("X-RateLimit-Reset"); reset != "" {
		resetTime, err := time.Parse(time.RFC3339, reset)
		if err == nil {
			c.RateLimitInfo.ResetTime = resetTime
		} else {
			c.RateLimitInfo.ResetTime = time.Now().Add(time.Hour)
		}
	}

	if c.RateLimitInfo.Limit > 0 {
		fmt.Printf("DEBUG: Rate limit: %d/%d, Reset: %s\n",
			c.RateLimitInfo.Remaining,
			c.RateLimitInfo.Limit,
			c.RateLimitInfo.ResetTime.Format(time.RFC3339))
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

	c.extractRateLimitInfo(resp)

	if resp.StatusCode == http.StatusTooManyRequests {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("rate limit exceeded: %s", string(body))
	}

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

	c.extractRateLimitInfo(resp)

	if resp.StatusCode == http.StatusTooManyRequests {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("rate limit exceeded: %s", string(body))
	}

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
