package tui

import "math"

func tempUnit(units string) string {
	switch units {
	case "imperial":
		return "°F"
	case "metric":
		return "°C"
	default:
		return "K"
	}
}

// windSpeed converts the API speed (m/s, or mph for imperial) to the display unit.
func windSpeed(speed float64, units string) float64 {
	if units == "imperial" {
		return speed
	}
	return speed * 3.6
}

func windUnit(units string) string {
	if units == "imperial" {
		return "mph"
	}
	return "km/h"
}

// kmh is the speed in km/h, so the wind bar keeps one scale for all units.
func kmh(speed float64, units string) float64 {
	if units == "imperial" {
		return speed * 1.609
	}
	return speed * 3.6
}

func deg(t float64) int {
	return int(math.Round(t))
}

func minMax(vals []float64) (float64, float64) {
	lo, hi := vals[0], vals[0]
	for _, v := range vals[1:] {
		lo = math.Min(lo, v)
		hi = math.Max(hi, v)
	}
	return lo, hi
}

func clamp(v, lo, hi int) int {
	return min(max(v, lo), hi)
}
