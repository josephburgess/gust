package output

import (
	"fmt"

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
