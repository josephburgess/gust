package renderer

import (
	"fmt"
	"strings"

	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
)

func (r *TerminalRenderer) RenderAlerts(city *models.City, weather *models.OneCallResponse, cfg *config.Config) {
	fmt.Print(styles.FormatHeader(fmt.Sprintf("WEATHER ALERTS FOR %s", strings.ToUpper(city.Name))))

	if len(weather.Alerts) == 0 {
		fmt.Println("No weather alerts for this area.")
		return
	}

	for i, alert := range weather.Alerts {
		if i > 0 {
			fmt.Println(styles.Divider(30))
		}

		fmt.Printf("%s\n", styles.AlertStyle(fmt.Sprintf("⚠️  %s", alert.Event)))
		fmt.Printf("Issued by: %s\n", alert.SenderName)
		fmt.Printf("Valid: %s to %s\n\n",
			styles.TimeStyle(weather.LocalTime(alert.Start).Format("Mon Jan 2 15:04")),
			styles.TimeStyle(weather.LocalTime(alert.End).Format("Mon Jan 2 15:04")))

		fmt.Println(alert.Description)
		fmt.Println()
	}
}
