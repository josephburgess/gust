package ui

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	HeaderStyle    = color.New(color.FgHiCyan, color.Bold).SprintFunc()
	TempStyle      = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	HighlightStyle = color.New(color.FgHiWhite).SprintFunc()
	InfoStyle      = color.New(color.FgHiBlue).SprintFunc()
	TimeStyle      = color.New(color.FgHiYellow).SprintFunc()
	AlertStyle     = color.New(color.FgHiRed, color.Bold).SprintFunc()
)

func Divider() string {
	return strings.Repeat("─", 50)
}

func ExitWithError(message string, err error) {
	log.Printf("%s: %v", message, err)
	os.Exit(1)
}

func FormatHeader(title string) string {
	return fmt.Sprintf("\n%s\n%s\n", HeaderStyle(title), Divider())
}
