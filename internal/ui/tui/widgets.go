package tui

import (
	"math"
	"strings"
)

var sparkChars = []rune("▁▂▃▄▅▆▇█")

// Bar returns a width-cell gauge of v against max.
func Bar(v, max float64, width int) string {
	filled := 0
	if max > 0 {
		filled = clamp(int(math.Round(v/max*float64(width))), 0, width)
	}
	return strings.Repeat("▰", filled) + strings.Repeat("▱", width-filled)
}

// Spark returns one block char per value, scaled to the min and max of vals.
func Spark(vals []float64) string {
	if len(vals) == 0 {
		return ""
	}
	lo, hi := minMax(vals)
	var sb strings.Builder
	for _, v := range vals {
		idx := len(sparkChars) / 2
		if hi > lo {
			idx = int(math.Round((v - lo) / (hi - lo) * float64(len(sparkChars)-1)))
		}
		sb.WriteRune(sparkChars[idx])
	}
	return sb.String()
}

func HBlocks(n int) string {
	return strings.Repeat("█", max(n, 0))
}

// RangeBar places lo..hi on a track that spans weekMin..weekMax. The fill is at least 1 cell.
func RangeBar(lo, hi, weekMin, weekMax float64, width int) (pre, fill, post string) {
	start, end := 0, width
	if span := weekMax - weekMin; span > 0 {
		start = clamp(int(math.Round((lo-weekMin)/span*float64(width))), 0, width)
		end = clamp(int(math.Round((hi-weekMin)/span*float64(width))), 0, width)
	}
	if end-start < 1 {
		if start >= width {
			start = width - 1
		}
		end = start + 1
	}
	return strings.Repeat("·", start), strings.Repeat("━", end-start), strings.Repeat("·", width-end)
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
