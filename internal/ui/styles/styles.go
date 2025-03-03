package styles

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

var (
	Base    = lipgloss.Color("#191724")
	Surface = lipgloss.Color("#1f1d2e")
	Overlay = lipgloss.Color("#26233a")
	Muted   = lipgloss.Color("#6e6a86")
	Subtle  = lipgloss.Color("#908caa")
	Text    = lipgloss.Color("#e0def4")
	Love    = lipgloss.Color("#eb6f92")
	Gold    = lipgloss.Color("#f6c177")
	Rose    = lipgloss.Color("#ebbcba")
	Pine    = lipgloss.Color("#31748f")
	Foam    = lipgloss.Color("#9ccfd8")
	Iris    = lipgloss.Color("#c4a7e7")
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Rose)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Gold)

	HighlightStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Text)

	CursorStyle = lipgloss.NewStyle().
			Foreground(Love)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Foam)

	HintStyle = lipgloss.NewStyle().
			Foreground(Subtle).
			Italic(true)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Iris).
			Padding(1, 2)
)

var (
	HeaderStyle     = color.New(color.FgHiCyan, color.Bold).SprintFunc()
	TempStyle       = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	HighlightStyleF = color.New(color.FgHiWhite).SprintFunc()
	InfoStyle       = color.New(color.FgHiBlue).SprintFunc()
	TimeStyle       = color.New(color.FgHiYellow).SprintFunc()
	AlertStyle      = color.New(color.FgHiRed, color.Bold).SprintFunc()
)

func Divider(len int) string {
	return strings.Repeat("─", len)
}

func ExitWithError(message string, err error) {
	log.Printf("%s: %v", message, err)
	os.Exit(1)
}

func FormatHeader(title string) string {
	return fmt.Sprintf("\n%s\n%s\n", HeaderStyle(title), Divider(len(title)*2))
}
