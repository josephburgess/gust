package output

import (
	"fmt"
	"strings"
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
		"⚠️ API Rate Limit Warning\n\n"+
			"You have %s requests remaining out of %d.\n"+
			"Your rate limit will reset at %s (%d minutes from now).",
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
		"❌ API Rate Limit Reached\n\n"+
			"Sorry - you have used all %d available requests.\n"+
			"You must really like checking the weather!!\n"+
			"Your rate limit will reset at %s (%d minutes from now).\n\n"+
			"💡 If you think the limits are too low please get in touch :)",
		limit,
		styles.TimeStyle(resetFormatted),
		minutesUntilReset,
	)))
	fmt.Println()
}

func PrintWhereConfig(defaultCity, units, view, apiURL string, showTips bool) {
	tipsStr := "off"
	if showTips {
		tipsStr = "on"
	}
	fmt.Println(styles.BoxStyle.Render(fmt.Sprintf(
		"Current configuration\n\n"+
			"Default city    %s\n"+
			"Units           %s\n"+
			"Default view    %s\n"+
			"API             %s\n"+
			"Tips            %s",
		styles.HighlightStyleF(defaultCity),
		styles.HighlightStyleF(units),
		styles.HighlightStyleF(view),
		styles.HighlightStyleF(apiURL),
		styles.HighlightStyleF(tipsStr),
	)))
}

func PrintStaleWarning(age time.Duration) {
	hours := int(age.Hours())
	minutes := int(age.Minutes()) % 60

	var ageStr string
	if hours > 0 {
		ageStr = fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		ageStr = fmt.Sprintf("%dm", minutes)
	}

	fmt.Println()
	fmt.Println(styles.BoxStyle.BorderForeground(styles.Gold).Render(fmt.Sprintf(
		"⚠️  Could not reach the weather service.\n"+
			"Showing last known data from %s ago.",
		styles.WarningStyle(ageStr),
	)))
	fmt.Println()
}

func PrintCityNotFound(city string, suggestions []string) {
	fmt.Println(styles.ErrorStyle(fmt.Sprintf("❌ City %q not found.", city)))

	if len(suggestions) == 0 {
		fmt.Println(styles.InfoStyle("No suggestions available — try a different spelling or add a country code (e.g. \"London,GB\")."))
		return
	}

	fmt.Println()
	fmt.Println(styles.InfoStyle("Did you mean one of these?"))
	for _, s := range suggestions {
		fmt.Printf("  %s %s\n", styles.WarningStyle("→"), s)
	}
	fmt.Println()
	fmt.Println(styles.HintStyle.Render("Tip: add a country code to be specific, e.g. gust \"London,GB\""))
}

func PrintQuotaStatus(limit, used int, resetAt *time.Time, unlimited bool) {
	if unlimited {
		fmt.Println(styles.BoxStyle.Render(
			"API Quota\n\n" +
				styles.SuccessStyle("Unlimited") + " — using your own OpenWeather API key",
		))
		return
	}

	if limit <= 0 {
		PrintWarning("Could not retrieve quota information.")
		return
	}

	const barWidth = 20
	filledCount := min(int(float64(used)/float64(limit)*float64(barWidth)), barWidth)
	emptyCount := barWidth - filledCount

	filled := styles.HighlightStyleF(strings.Repeat("█", filledCount))
	empty := strings.Repeat("░", emptyCount)

	percentage := float64(used) / float64(limit) * 100

	var usageText string
	switch {
	case percentage >= 90:
		usageText = styles.ErrorStyle(fmt.Sprintf("%.0f%% used", percentage))
	case percentage >= 75:
		usageText = styles.WarningStyle(fmt.Sprintf("%.0f%% used", percentage))
	default:
		usageText = styles.InfoStyle(fmt.Sprintf("%.0f%% used", percentage))
	}

	remaining := limit - used

	var resetLine string
	if resetAt != nil {
		timeUntil := time.Until(*resetAt)
		hours := int(timeUntil.Hours())
		minutes := int(timeUntil.Minutes()) % 60
		resetLine = fmt.Sprintf("\nResets in %s (%s)",
			styles.TimeStyle(fmt.Sprintf("%dh %dm", hours, minutes)),
			styles.TimeStyle(resetAt.Local().Format("15:04")),
		)
	}

	fmt.Println(styles.BoxStyle.Render(fmt.Sprintf(
		"API Quota\n\n"+
			"[%s%s] %s\n"+
			"%s used of %s daily requests (%s remaining)%s",
		filled, empty, usageText,
		styles.HighlightStyleF(fmt.Sprintf("%d", used)),
		styles.HighlightStyleF(fmt.Sprintf("%d", limit)),
		styles.HighlightStyleF(fmt.Sprintf("%d", remaining)),
		resetLine,
	)))
}
