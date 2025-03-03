package ui

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

var (
	base    = lipgloss.Color("#191724")
	surface = lipgloss.Color("#1f1d2e")
	overlay = lipgloss.Color("#26233a")
	muted   = lipgloss.Color("#6e6a86")
	subtle  = lipgloss.Color("#908caa")
	text    = lipgloss.Color("#e0def4")
	love    = lipgloss.Color("#eb6f92")
	gold    = lipgloss.Color("#f6c177")
	rose    = lipgloss.Color("#ebbcba")
	pine    = lipgloss.Color("#31748f")
	foam    = lipgloss.Color("#9ccfd8")
	iris    = lipgloss.Color("#c4a7e7")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(rose)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(gold)

	highlightStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(text)

	cursorStyle = lipgloss.NewStyle().
			Foreground(love)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(foam)

	hintStyle = lipgloss.NewStyle().
			Foreground(subtle).
			Italic(true)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(iris).
			Padding(1, 2)

	HeaderStyle    = color.New(color.FgHiCyan, color.Bold).SprintFunc()
	TempStyle      = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	HighlightStyle = color.New(color.FgHiWhite).SprintFunc()
	InfoStyle      = color.New(color.FgHiBlue).SprintFunc()
	TimeStyle      = color.New(color.FgHiYellow).SprintFunc()
	AlertStyle     = color.New(color.FgHiRed, color.Bold).SprintFunc()
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
