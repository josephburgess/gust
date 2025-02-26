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

func (w Weather) Emoji() string {
	switch {
	case w.ID >= 200 && w.ID <= 232:
		return "⛈" // storm
	case w.ID >= 300 && w.ID <= 321:
		return "🌦" // drizzle
	case w.ID >= 500 && w.ID <= 531:
		return "🌧" // rain
	case w.ID >= 600 && w.ID <= 622:
		return "❄️" // snow
	case w.ID >= 700 && w.ID <= 781:
		return "🌫" // fog
	case w.ID == 800:
		return "☀️" // clear
	case w.ID >= 801 && w.ID <= 804:
		return "☁️" // cloudy
	default:
		return "🌡" // the rest
	}
}

func (fi ForecastItem) Emoji() string {
	switch {
	case fi.WeatherID >= 200 && fi.WeatherID <= 232:
		return "⛈" // storm
	case fi.WeatherID >= 300 && fi.WeatherID <= 321:
		return "🌦" // drizzle
	case fi.WeatherID >= 500 && fi.WeatherID <= 531:
		return "🌧" // rain
	case fi.WeatherID >= 600 && fi.WeatherID <= 622:
		return "❄️" // snow
	case fi.WeatherID >= 700 && fi.WeatherID <= 781:
		return "🌫" // fog
	case fi.WeatherID == 800:
		return "☀️" // clear
	case fi.WeatherID >= 801 && fi.WeatherID <= 804:
		return "☁️" // cloudy
	default:
		return "🌡" // the rest
	}
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
