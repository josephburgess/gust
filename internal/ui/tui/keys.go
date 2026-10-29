package tui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
	"github.com/josephburgess/gust/internal/ui/styles"
)

var keys = struct {
	Tabs, Next, Prev, Down, Up, City, Units, Refresh, Help, Quit key.Binding
}{
	Tabs:    key.NewBinding(key.WithKeys("1", "2", "3", "4"), key.WithHelp("1-4", "tabs")),
	Next:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next tab")),
	Prev:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "previous tab")),
	Down:    key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j/↓", "down")),
	Up:      key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k/↑", "up")),
	City:    key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "city")),
	Units:   key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "units")),
	Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
	Help:    key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	Quit:    key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}

func helpView() string {
	h := help.New()
	h.Styles.FullKey = lipgloss.NewStyle().Foreground(styles.Gold)
	h.Styles.FullDesc = lipgloss.NewStyle().Foreground(styles.Subtle)
	h.Styles.FullSeparator = lipgloss.NewStyle().Foreground(styles.HighlightMed)
	return h.FullHelpView([][]key.Binding{
		{keys.Tabs, keys.Next, keys.Prev},
		{keys.Down, keys.Up},
		{keys.City, keys.Units, keys.Refresh},
		{keys.Help, keys.Quit},
	})
}

// nextUnits cycles metric, imperial and standard.
func nextUnits(units string) string {
	switch units {
	case "metric":
		return "imperial"
	case "imperial":
		return "standard"
	default:
		return "metric"
	}
}
