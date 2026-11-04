package styles

import (
	"log"
	"os"

	"charm.land/lipgloss/v2"
)

var (
	Base    = lipgloss.Color("#191724")
	Surface = lipgloss.Color("#1f1d2e")
	Overlay = lipgloss.Color("#26233a")
	// HighlightMed is for panel borders and empty range bar dots
	HighlightMed = lipgloss.Color("#403d52")
	Muted        = lipgloss.Color("#6e6a86")
	Subtle       = lipgloss.Color("#908caa")
	Text         = lipgloss.Color("#e0def4")
	Love         = lipgloss.Color("#eb6f92")
	Gold         = lipgloss.Color("#f6c177")
	Rose         = lipgloss.Color("#ebbcba")
	Pine         = lipgloss.Color("#31748f")
	Foam         = lipgloss.Color("#9ccfd8")
	Iris         = lipgloss.Color("#c4a7e7")
)

var (
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Iris).
			Padding(1, 2)

	ProgressMessageStyle = lipgloss.NewStyle().
				Foreground(Foam).
				Italic(true)
)

func ExitWithError(message string, err error) {
	log.Printf("%s: %v", message, err)
	os.Exit(1)
}
