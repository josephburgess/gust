package setup

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/ui/styles"
)

// rep's the current step of the wizard
type SetupState int

const (
	StateCity SetupState = iota
	StateUnits
	StateView
	StateAuth
	StateComplete
)

const asciiLogo = `
                        __
      ____ ___  _______/ /_
     / ** '/ / / / **_/ __/
    / /_/ / /_/ (__  ) /_   _
    \__, /\__,_/____/\__/  (_)
   /____/                      `

// setup ui styles
var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(styles.Rose)
	boxStyle          = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(styles.Iris).Padding(1, 2)
	subtitleStyle     = lipgloss.NewStyle().Foreground(styles.Gold)
	highlightStyle    = lipgloss.NewStyle().Bold(true).Foreground(styles.Text)
	cursorStyle       = lipgloss.NewStyle().Foreground(styles.Love)
	selectedItemStyle = lipgloss.NewStyle().Foreground(styles.Foam)
	hintStyle         = lipgloss.NewStyle().Foreground(styles.Subtle).Italic(true)
)

// current state of wizard
type Model struct {
	Config        *config.Config
	State         SetupState
	CityInput     textinput.Model
	UnitOptions   []string
	UnitCursor    int
	ViewOptions   []string
	ViewCursor    int
	AuthOptions   []string
	AuthCursor    int
	NeedsAuth     bool
	Width, Height int
	Quitting      bool
}

// creates a new setup model
func NewModel(cfg *config.Config, needsAuth bool) Model {
	ti := textinput.New()
	ti.Placeholder = "Wherever the wind blows..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = len(ti.Placeholder)
	ti.PromptStyle = lipgloss.NewStyle().Foreground(styles.Love)
	ti.TextStyle = lipgloss.NewStyle().Foreground(styles.Text)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(styles.Gold)

	unitCursor := 0
	switch cfg.Units {
	case "imperial":
		unitCursor = 1
	case "standard":
		unitCursor = 2
	}

	viewCursor := 0
	switch cfg.DefaultView {
	case "compact":
		viewCursor = 1
	case "daily":
		viewCursor = 2
	case "hourly":
		viewCursor = 3
	case "full":
		viewCursor = 4
	}

	return Model{
		Config:      cfg,
		State:       StateCity,
		CityInput:   ti,
		UnitOptions: []string{"metric (°C, km/h) 🌡️", "imperial (°F, mph) 🌡️", "standard (K, m/s) 🌡️"},
		UnitCursor:  unitCursor,
		ViewOptions: []string{
			"detailed 🌤️",
			"compact 📊",
			"daily (5-day) 📆",
			"hourly (24-hour forecast) 🕒",
			"full (current + daily + alerts) 📋",
		},
		ViewCursor:  viewCursor,
		AuthOptions: []string{"Yes, authenticate with GitHub 🔑", "No, I'll do it later ⏱️"},
		AuthCursor:  0,
		NeedsAuth:   needsAuth,
		Quitting:    false,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type (
	AuthenticateMsg  struct{}
	SetupCompleteMsg struct{}
)
