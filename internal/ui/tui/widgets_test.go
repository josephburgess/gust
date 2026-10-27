package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	flatTemps   = []float64{9, 9, 9}
	risingTemps = []float64{0, 7, 14}
)

func TestBar(t *testing.T) {
	tests := []struct {
		name  string
		v     float64
		max   float64
		width int
		want  string
	}{
		{"empty", 0, 100, 5, "▱▱▱▱▱"},
		{"half", 50, 100, 4, "▰▰▱▱"},
		{"full", 100, 100, 3, "▰▰▰"},
		{"over max clamps", 150, 100, 3, "▰▰▰"},
		{"zero max is empty", 5, 0, 3, "▱▱▱"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Bar(tt.v, tt.max, tt.width))
		})
	}
}

func TestSpark(t *testing.T) {
	tests := []struct {
		name string
		vals []float64
		want string
	}{
		{"empty", nil, ""},
		{"flat uses the middle char", flatTemps, "▅▅▅"},
		{"rising spans the full range", risingTemps, "▁▅█"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Spark(tt.vals))
		})
	}
}

func TestHBlocks(t *testing.T) {
	assert.Equal(t, "███", HBlocks(3))
	assert.Equal(t, "", HBlocks(-2))
}

func TestRangeBar(t *testing.T) {
	tests := []struct {
		name               string
		lo, hi, wMin, wMax float64
		width              int
		pre, fill, post    string
	}{
		{"full week", 0, 10, 0, 10, 10, "", "━━━━━━━━━━", ""},
		{"middle", 3, 7, 0, 10, 10, "···", "━━━━", "···"},
		{"same lo and hi keeps one fill", 5, 5, 0, 10, 10, "·····", "━", "····"},
		{"at the top keeps one fill", 10, 10, 0, 10, 10, "·········", "━", ""},
		{"flat week fills the track", 4, 4, 4, 4, 5, "", "━━━━━", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre, fill, post := RangeBar(tt.lo, tt.hi, tt.wMin, tt.wMax, tt.width)
			assert.Equal(t, tt.pre, pre)
			assert.Equal(t, tt.fill, fill)
			assert.Equal(t, tt.post, post)
		})
	}
}
