package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const (
	baseURL    = "https://api.openweathermap.org/"
	geoCodeURL = baseURL + "geo/1.0/direct?q=%s&limit=1&appid=%s"
	weatherURL = baseURL + "data/2.5/weather?lat=%f&lon=%f&appid=%s"
)

type City struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type Weather struct {
	Temp        float64 `json:"temp"`
	Description string  `json:"description"`
}

func checkError(err error, message string) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func getApiKey() string {
	err := godotenv.Load()
	checkError(err, "Error loading .env")
	return os.Getenv("OPENWEATHER_API_KEY")
}

func main() {
	cityName := "London"
	apiKey := getApiKey()

	city, err := getCoordinates(cityName, apiKey)
	checkError(err, "Didn't get coords")

	weather, err := getWeather(city.Lat, city.Lon, apiKey)
	checkError(err, "Didn't get weather")

	fmt.Printf("Weather in %s:\n", city.Name)
	fmt.Printf("Temperature: %.1f°C\n", weather.Temp-273.15) // Convert from Kelvin
	fmt.Printf("Conditions: %s\n", weather.Description)
}

func getCoordinates(city, apiKey string) (*City, error) {
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf(geoCodeURL, encodedCity, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cities []City
	err = json.Unmarshal(body, &cities)
	if err != nil {
		return nil, err
	}

	if len(cities) == 0 {
		return nil, fmt.Errorf("no coordinates found for %s", city)
	}

	return &cities[0], nil
}

func getWeather(lat float64, lon float64, apiKey string) (*Weather, error) {
	url := fmt.Sprintf(weatherURL, lat, lon, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	weather := &Weather{
		Temp:        result.Main.Temp,
		Description: result.Weather[0].Description,
	}

	return weather, nil
}
