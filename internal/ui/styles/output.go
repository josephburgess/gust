package styles

import (
	"fmt"
)

func PrintError(message string) {
	fmt.Println(ErrorStyle("❌ " + message))
}

func PrintSuccess(message string) {
	fmt.Println(SuccessStyle("✅ " + message))
}

func PrintInfo(message string) {
	fmt.Println(InfoStyle("ℹ️  " + message))
}

func PrintWarning(message string) {
	fmt.Println(WarningStyle("⚠️ " + message))
}

func PrintHeader(title string) {
	fmt.Printf("\n%s\n%s\n", HeaderStyle(title), Divider(len(title)*2))
}

func PrintBoxedMessage(message string) {
	fmt.Println(BoxStyle.Render(message))
}
