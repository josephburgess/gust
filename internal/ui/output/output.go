package output

import (
	"fmt"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/josephburgess/gust/internal/ui/styles"
)

func PrintError(message string) {
	lipgloss.Println(lipgloss.NewStyle().Foreground(styles.Love).Bold(true).Render(message))
}

func PrintSuccess(message string) {
	lipgloss.Println(lipgloss.NewStyle().Foreground(styles.Foam).Render("✓ " + message))
}

func PrintInfo(message string) {
	lipgloss.Println(lipgloss.NewStyle().Foreground(styles.Subtle).Render(message))
}

func PrintRateLimitWarning(remaining, limit int, resetTime time.Time) {
	timeUntilReset := time.Until(resetTime)
	minutesUntilReset := int(timeUntilReset.Minutes())
	resetFormatted := resetTime.Format("15:04")

	fmt.Println()
	lipgloss.Println(styles.BoxStyle.Render(fmt.Sprintf(
		"API rate limit warning\n\n"+
			"You have %s requests remaining out of %d.\n"+
			"Your rate limit will reset at %s (%d minutes from now).",
		lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%d", remaining)),
		limit,
		lipgloss.NewStyle().Foreground(styles.Gold).Render(resetFormatted),
		minutesUntilReset,
	)))
	fmt.Println()
}

func PrintRateLimitError(limit int, resetTime time.Time) {
	timeUntilReset := time.Until(resetTime)
	minutesUntilReset := int(timeUntilReset.Minutes())
	resetFormatted := resetTime.Format("15:04")

	fmt.Println()
	lipgloss.Println(styles.BoxStyle.BorderForeground(styles.Love).Render(fmt.Sprintf(
		"API rate limit reached\n\n"+
			"Sorry - you have used all %d available requests.\n"+
			"You must really like checking the weather!!\n"+
			"Your rate limit will reset at %s (%d minutes from now).\n\n"+
			"If you think the limits are too low please get in touch :)",
		limit,
		lipgloss.NewStyle().Foreground(styles.Gold).Render(resetFormatted),
		minutesUntilReset,
	)))
	fmt.Println()
}

// going to implement this later - will create an api key status check endpoint
/*
func PrintRateLimitStatus(remaining, limit int) {
	if limit <= 0 {
		return
	}

	const barWidth = 20
	used := limit - remaining

	filledCount := min(int(float64(used) / float64(limit) * barWidth), barWidth)
	emptyCount := barWidth - filledCount

	filled := styles.HighlightStyleF(strings.Repeat("█", filledCount))
	empty := strings.Repeat("░", emptyCount)

	percentage := float64(used) / float64(limit) * 100

	var usageText string
	if percentage >= 90 {
		usageText = styles.ErrorStyle(fmt.Sprintf("%.0f%% used", percentage))
	} else if percentage >= 75 {
		usageText = styles.WarningStyle(fmt.Sprintf("%.0f%% used", percentage))
	} else {
		usageText = styles.InfoStyle(fmt.Sprintf("%.0f%% used", percentage))
	}

	fmt.Printf("API Usage: [%s%s] %s (%d/%d)\n", filled, empty, usageText, used, limit)
}
*/
