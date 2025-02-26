package models

type City struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type Weather struct {
	ID          int     `json:"id"`
	Icon        string  `json:"icon"`
	Temp        float64 `json:"temp"`
	Description string  `json:"description"`
}
type ForecastItem struct {
	DateTime    string  `json:"dateTime"`
	Temp        float64 `json:"temp"`
	Description string  `json:"description"`
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
