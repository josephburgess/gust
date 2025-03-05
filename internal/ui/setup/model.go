package setup

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/components"
	"github.com/josephburgess/gust/internal/ui/styles"
)

// rep's the current step of the wizard
type SetupState int

const (
	StateCity SetupState = iota
	StateCitySearch
	StateCitySelect
	StateUnits
	StateView
	StateAuth
	StateComplete
)

const asciiLogo = `
                         __
       ____ ___  _______/ /_
      / ** '/ / / / **_/ __/
     / /_/ / /_/ (__  ) /_
     \__, /\__,_/____/\__/   💨🍃
    /____/                      `

// setup ui styles
var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(styles.Rose)
	boxStyle          = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(styles.Iris).Padding(0, 1, 0, 1)
	logoBoxStyle      = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(styles.Subtle).Padding(0, 2, 2, 1).Foreground(styles.Foam)
	subtitleStyle     = lipgloss.NewStyle().Foreground(styles.Gold)
	highlightStyle    = lipgloss.NewStyle().Bold(true).Foreground(styles.Text)
	cursorStyle       = lipgloss.NewStyle().Foreground(styles.Love)
	selectedItemStyle = lipgloss.NewStyle().Foreground(styles.Foam)
	hintStyle         = lipgloss.NewStyle().Foreground(styles.Subtle).Italic(true)
)

// current state of wizard
type Model struct {
	Config          *config.Config
	State           SetupState
	CityInput       textinput.Model
	CitySearchQuery string
	CityOptions     []models.City
	CityCursor      int
	Client          *api.Client
	UnitOptions     []string
	UnitCursor      int
	ViewOptions     []string
	ViewCursor      int
	AuthOptions     []string
	AuthCursor      int
	NeedsAuth       bool
	Width, Height   int
	Quitting        bool
	Spinner         components.SpinnerModel
}

// creates a new setup model
func NewModel(cfg *config.Config, needsAuth bool, client *api.Client) Model {
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
		Config:          cfg,
		State:           StateCity,
		CityInput:       ti,
		CitySearchQuery: "",
		CityOptions:     []models.City{},
		CityCursor:      0,
		Client:          client,
		UnitOptions:     []string{"metric (°C, km/h) 🌡️", "imperial (°F, mph) 🌡️", "standard (K, m/s) 🌡️"},
		UnitCursor:      unitCursor,
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
		Spinner:     components.NewSpinner(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.Spinner.Tick(),
	)
}

type (
	AuthenticateMsg  struct{}
	SetupCompleteMsg struct{}
)
