package output

import (
	"fmt"
	"time"

	"github.com/josephburgess/gust/internal/ui/styles"
)

func PrintError(message string) {
	fmt.Println(styles.ErrorStyle("❌ " + message))
}

func PrintSuccess(message string) {
	fmt.Println(styles.SuccessStyle("✅ " + message))
}

func PrintInfo(message string) {
	fmt.Println(styles.InfoStyle(message))
}

func PrintWarning(message string) {
	fmt.Println(styles.WarningStyle("⚠️ " + message))
}

func PrintHeader(title string) {
	fmt.Printf("\n%s\n%s\n", styles.HeaderStyle(title), styles.Divider(len(title)*2))
}

func PrintBoxedMessage(message string) {
	fmt.Println(styles.BoxStyle.Render(message))
}

func PrintRateLimitWarning(remaining, limit int, resetTime time.Time) {
	timeUntilReset := time.Until(resetTime)
	minutesUntilReset := int(timeUntilReset.Minutes())
	resetFormatted := resetTime.Format("15:04")

	fmt.Println()
	fmt.Println(styles.BoxStyle.Render(fmt.Sprintf(
		"%s API Rate Limit Warning\n\n"+
			"You have %s requests remaining out of %d.\n"+
			"Your rate limit will reset at %s (%d minutes from now).",
		styles.WarningStyle("⚠️"),
		styles.HighlightStyleF(fmt.Sprintf("%d", remaining)),
		limit,
		styles.TimeStyle(resetFormatted),
		minutesUntilReset,
	)))
	fmt.Println()
}

func PrintRateLimitError(limit int, resetTime time.Time) {
	timeUntilReset := time.Until(resetTime)
	minutesUntilReset := int(timeUntilReset.Minutes())
	resetFormatted := resetTime.Format("15:04")

	fmt.Println()
	fmt.Println(styles.BoxStyle.BorderForeground(styles.Love).Render(fmt.Sprintf(
		"%s API Rate Limit Reached\n\n"+
			"You have used all %d available requests.\n"+
			"Your rate limit will reset at %s (%d minutes from now).\n\n"+
			styles.ErrorStyle("❌"),
		limit,
		styles.TimeStyle(resetFormatted),
		minutesUntilReset,
		styles.InfoStyle("💡"),
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
