package setup

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	StatePreferences
	StateApiKeyOption
	StateApiKeyInput
	StateComplete
)

// PrefGroup is the focused panel on the preferences step.
type PrefGroup int

const (
	GroupUnits PrefGroup = iota
	GroupView
	GroupTips
)

// option is one radio row: value is what the config stores.
type option struct{ value, label, detail string }

var (
	unitOptions = []option{{"metric", "metric", "°C · km/h"}, {"imperial", "imperial", "°F · mph"}, {"standard", "standard", "K · km/h"}}
	viewOptions = []option{{"now", "Now", ""}, {"hourly", "Hourly", ""}, {"daily", "5-day", ""}, {"compact", "compact", "print and exit"}}
	tipOptions  = []option{{"on", "on", "a tip under the weather"}, {"off", "off", ""}}
)

// current state of wizard
type Model struct {
	Config          *config.Config
	State           SetupState
	CityInput       textinput.Model
	CitySearchQuery string
	CityOptions     []models.City
	CityCursor      int
	SearchErr       string
	Client          *api.Client
	UnitCursor      int
	ViewCursor      int
	TipCursor       int
	PrefGroup       PrefGroup
	NeedsAuth       bool
	// WantsGitHub asks the runner to sign in with GitHub after the wizard closes
	WantsGitHub   bool
	Width, Height int
	Quitting      bool
	Spinner       components.SpinnerModel
	ApiKeyCursor  int
	ApiKeyInput   textinput.Model
}

func optionIndex(options []option, value string) int {
	for i, o := range options {
		if o.value == value {
			return i
		}
	}
	return 0
}

// creates a new setup model
func NewModel(cfg *config.Config, needsAuth bool, client *api.Client) Model {
	inputStyles := textinput.DefaultDarkStyles()
	for _, state := range []*textinput.StyleState{&inputStyles.Focused, &inputStyles.Blurred} {
		state.Prompt = lipgloss.NewStyle().Foreground(styles.Love)
		state.Text = lipgloss.NewStyle().Foreground(styles.Text)
		state.Placeholder = lipgloss.NewStyle().Foreground(styles.Subtle)
	}
	inputStyles.Cursor.Color = styles.Gold

	ti := textinput.New()
	ti.Prompt = "› "
	ti.Placeholder = "Wherever the wind blows..."
	ti.Focus()
	ti.CharLimit = 50
	ti.SetStyles(inputStyles)

	apiKeyInput := textinput.New()
	apiKeyInput.Prompt = "› "
	apiKeyInput.Placeholder = "paste your key"
	apiKeyInput.CharLimit = 64
	apiKeyInput.EchoMode = textinput.EchoPassword
	apiKeyInput.SetStyles(inputStyles)

	tipCursor := 0
	if cfg.DefaultCity != "" && !cfg.ShowTips {
		tipCursor = 1
	}

	return Model{
		Config:      cfg,
		State:       StateCity,
		CityInput:   ti,
		CityOptions: []models.City{},
		Client:      client,
		UnitCursor:  optionIndex(unitOptions, cfg.Units),
		ViewCursor:  optionIndex(viewOptions, cfg.DefaultView),
		TipCursor:   tipCursor,
		NeedsAuth:   needsAuth,
		Spinner:     components.NewSpinner(),
		ApiKeyInput: apiKeyInput,
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
