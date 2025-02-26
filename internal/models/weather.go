package models

// City represents location information from the geocoding API
type City struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country,omitempty"`
	State   string  `json:"state,omitempty"`
}

// Weather represents weather information
type Weather struct {
	Temp        float64 `json:"temp"`
	Description string  `json:"description"`
	FeelsLike   float64 `json:"feels_like,omitempty"`
	TempMin     float64 `json:"temp_min,omitempty"`
	TempMax     float64 `json:"temp_max,omitempty"`
	Humidity    int     `json:"humidity,omitempty"`
	WindSpeed   float64 `json:"wind_speed,omitempty"`
	WindDeg     int     `json:"wind_deg,omitempty"`
	CloudsAll   int     `json:"clouds_all,omitempty"`
	Visibility  int     `json:"visibility,omitempty"`
	Pressure    int     `json:"pressure,omitempty"`
	Sunrise     int64   `json:"sunrise,omitempty"`
	Sunset      int64   `json:"sunset,omitempty"`
}
