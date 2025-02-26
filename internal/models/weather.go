package models

import (
	"fmt"
	"math"
	"time"
)

type City struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type Weather struct {
	ID          int     `json:"id"`
	Icon        string  `json:"icon"`
	Temp        float64 `json:"temp"`
	FeelsLike   float64 `json:"feels_like"`
	TempMin     float64 `json:"temp_min"`
	TempMax     float64 `json:"temp_max"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"speed"`
	WindDeg     int     `json:"deg"`
	WindGust    float64 `json:"gust"`
	Pressure    int     `json:"pressure"`
	Visibility  int     `json:"visibility"`
	Sunrise     int64   `json:"sunrise"`
	Sunset      int64   `json:"sunset"`
	Rain1h      float64 `json:"1h"`
	Clouds      int     `json:"all"`
}

type ForecastItem struct {
	DateTime    int64   `json:"dt"`
	TempMax     float64 `json:"temp_max"`
	TempMin     float64 `json:"temp_min"`
	Description string
	WeatherID   int
	Icon        string
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"speed"`
	WindDeg     int     `json:"deg"`
	Pop         float64 `json:"pop"`
}

func KelvinToCelsius(k float64) float64 {
	return k - 273.15
}

func GetWeatherEmoji(id int) string {
	switch {
	case id >= 200 && id <= 232:
		return "⚡" // storm
	case id >= 300 && id <= 321:
		return "🌦" // drizzle
	case id >= 500 && id <= 531:
		return "☔" // rain
	case id >= 600 && id <= 622:
		return "⛄" // snow
	case id >= 700 && id <= 781:
		return "🌫" // fog
	case id == 800:
		return "🔆" // clear
	case id >= 801 && id <= 804:
		return "🌥️" // cloudy
	default:
		return "🌡" // the rest
	}
}

func (w Weather) Emoji() string {
	return GetWeatherEmoji(w.ID)
}

func (fi ForecastItem) Emoji() string {
	return GetWeatherEmoji(fi.WeatherID)
}

func FormatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("15:04")
}

func GetWindDirection(degrees int) string {
	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	index := int((math.Mod(float64(degrees)+22.5, 360)) / 45)
	return directions[index]
}

func FormatDay(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("Mon Jan 2")
}

func VisibilityToString(meters int) string {
	if meters >= 10000 {
		return "Excellent (10+ km)"
	} else if meters >= 5000 {
		return fmt.Sprintf("Good (%.1f km)", float64(meters)/1000)
	} else if meters >= 2000 {
		return fmt.Sprintf("Moderate (%.1f km)", float64(meters)/1000)
	} else {
		return fmt.Sprintf("Poor (%.1f km)", float64(meters)/1000)
	}
}
