package renderer

import (
	"time"
)

type BaseRenderer struct {
	Units string
}

func (r *BaseRenderer) GetTemperatureUnit() string {
	switch r.Units {
	case "imperial":
		return "°F"
	case "metric":
		return "°C"
	default:
		return "K"
	}
}

func (r *BaseRenderer) FormatWindSpeed(speed float64) float64 {
	switch r.Units {
	case "imperial":
		return speed
	default:
		return speed * 3.6
	}
}

func (r *BaseRenderer) GetWindSpeedUnit() string {
	switch r.Units {
	case "imperial":
		return "mph"
	default:
		return "km/h"
	}
}

func FormatDateTime(timestamp int64, format string) string {
	return time.Unix(timestamp, 0).Format(format)
}
