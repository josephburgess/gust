package models

import (
	"fmt"
)

type City struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type OneCallResponse struct {
	Lat            float64        `json:"lat"`
	Lon            float64        `json:"lon"`
	Timezone       string         `json:"timezone"`
	TimezoneOffset int            `json:"timezone_offset"`
	Current        CurrentWeather `json:"current"`
	Minutely       []MinuteData   `json:"minutely"`
	Hourly         []HourData     `json:"hourly"`
	Daily          []DayData      `json:"daily"`
	Alerts         []Alert        `json:"alerts"`
}

type WeatherCondition struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type RainData struct {
	OneHour float64 `json:"1h"`
}

type SnowData struct {
	OneHour float64 `json:"1h"`
}

type CurrentWeather struct {
	Dt         int64              `json:"dt"`
	Sunrise    int64              `json:"sunrise"`
	Sunset     int64              `json:"sunset"`
	Temp       float64            `json:"temp"`
	FeelsLike  float64            `json:"feels_like"`
	Pressure   int                `json:"pressure"`
	Humidity   int                `json:"humidity"`
	DewPoint   float64            `json:"dew_point"`
	UVI        float64            `json:"uvi"`
	Clouds     int                `json:"clouds"`
	Visibility int                `json:"visibility"`
	WindSpeed  float64            `json:"wind_speed"`
	WindGust   float64            `json:"wind_gust"`
	WindDeg    int                `json:"wind_deg"`
	Rain       *RainData          `json:"rain,omitempty"`
	Snow       *SnowData          `json:"snow,omitempty"`
	Weather    []WeatherCondition `json:"weather"`
}

type MinuteData struct {
	Dt            int64   `json:"dt"`
	Precipitation float64 `json:"precipitation"`
}

type HourData struct {
	Dt         int64              `json:"dt"`
	Temp       float64            `json:"temp"`
	FeelsLike  float64            `json:"feels_like"`
	Pressure   int                `json:"pressure"`
	Humidity   int                `json:"humidity"`
	DewPoint   float64            `json:"dew_point"`
	UVI        float64            `json:"uvi"`
	Clouds     int                `json:"clouds"`
	Visibility int                `json:"visibility"`
	WindSpeed  float64            `json:"wind_speed"`
	WindGust   float64            `json:"wind_gust"`
	WindDeg    int                `json:"wind_deg"`
	Pop        float64            `json:"pop"`
	Rain       *RainData          `json:"rain,omitempty"`
	Snow       *SnowData          `json:"snow,omitempty"`
	Weather    []WeatherCondition `json:"weather"`
}

type TempData struct {
	Day   float64 `json:"day"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

type FeelsLikeData struct {
	Day   float64 `json:"day"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

type DayData struct {
	Dt        int64              `json:"dt"`
	Sunrise   int64              `json:"sunrise"`
	Sunset    int64              `json:"sunset"`
	Moonrise  int64              `json:"moonrise"`
	Moonset   int64              `json:"moonset"`
	MoonPhase float64            `json:"moon_phase"`
	Summary   string             `json:"summary"`
	Temp      TempData           `json:"temp"`
	FeelsLike FeelsLikeData      `json:"feels_like"`
	Pressure  int                `json:"pressure"`
	Humidity  int                `json:"humidity"`
	DewPoint  float64            `json:"dew_point"`
	WindSpeed float64            `json:"wind_speed"`
	WindGust  float64            `json:"wind_gust"`
	WindDeg   int                `json:"wind_deg"`
	Clouds    int                `json:"clouds"`
	UVI       float64            `json:"uvi"`
	Pop       float64            `json:"pop"`
	Rain      float64            `json:"rain,omitempty"`
	Snow      float64            `json:"snow,omitempty"`
	Weather   []WeatherCondition `json:"weather"`
}

type Alert struct {
	SenderName  string   `json:"sender_name"`
	Event       string   `json:"event"`
	Start       int64    `json:"start"`
	End         int64    `json:"end"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
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

func GetWindDirection(degrees int) string {
	normalizedDegrees := ((degrees % 360) + 360) % 360
	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	index := int((float64(normalizedDegrees)+22.5)/45) % 8
	return directions[index]
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
