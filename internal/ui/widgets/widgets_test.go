package widgets

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

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

func TestBigTemp(t *testing.T) {
	tests := []struct {
		name string
		temp int
		want []string
	}{
		{"single digit", 9, []string{
			"█▀▀▀█ █▀█",
			"█   █ ▀▀▀",
			"▀▀▀▀█    ",
			"    █    ",
			"▀▀▀▀▀    ",
		}},
		{"negative", -1, []string{
			"     ▀█   █▀█",
			"      █   ▀▀▀",
			"▀▀▀   █      ",
			"      █      ",
			"    ▀▀▀▀▀    ",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, BigTemp(tt.temp))
		})
	}
}

func TestIcon(t *testing.T) {
	tests := []struct {
		name  string
		id    int
		night bool
		want  string
	}{
		{"storm", 211, false, "/_"},
		{"drizzle", 301, false, "' ' '"},
		{"rain", 501, false, "‚'‚'‚'"},
		{"snow", 601, false, "*  *  *"},
		{"fog", 741, false, "_ - _ -"},
		{"clear day", 800, false, "― (   ) ―"},
		{"clear night", 800, true, "_.-."},
		{"cloudy", 803, false, "(___.__)__)"},
		{"unknown falls back to cloudy", 999, false, "(___.__)__)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := Icon(tt.id, tt.night)
			assert.Len(t, lines, 5)
			for _, line := range lines {
				assert.Equal(t, IconWidth, ansi.StringWidth(line))
			}
			assert.Contains(t, ansi.Strip(strings.Join(lines, "\n")), tt.want)
		})
	}
}

func TestPanel(t *testing.T) {
	tests := []struct {
		name    string
		focused bool
	}{
		{"blurred", false},
		{"focused", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := strings.Split(ansi.Strip(Panel("wind", "24 km/h", 20, tt.focused)), "\n")
			assert.Equal(t, []string{
				"╭──────────────────╮",
				"│ WIND             │",
				"│ 24 km/h          │",
				"╰──────────────────╯",
			}, lines)
		})
	}
}

func TestWordmark(t *testing.T) {
	lines := Wordmark()
	assert.Len(t, lines, 3)
	for _, line := range lines {
		assert.Equal(t, ansi.StringWidth(lines[0]), ansi.StringWidth(line))
	}
}

func TestStepRail(t *testing.T) {
	tests := []struct {
		name    string
		current int
		want    string
	}{
		{"first step", 0, "● city ─── ○ preferences ─── ○ sign in"},
		{"middle step", 1, "✓ city ─── ● preferences ─── ○ sign in"},
		{"all done", 3, "✓ city ─── ✓ preferences ─── ✓ sign in"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ansi.Strip(StepRail([]string{"city", "preferences", "sign in"}, tt.current))
			assert.Equal(t, tt.want, got)
		})
	}
}
