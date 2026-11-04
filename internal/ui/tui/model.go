package tui

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/setup"
	"github.com/josephburgess/gust/internal/ui/styles"
)

type Tab int

const (
	TabNow Tab = iota
	TabHourly
	TabDaily
	TabAlerts
)

var tabNames = []string{"Now", "Hourly", "5-day", "Alerts"}

// TabFromView maps a stored view value to its tab. Old values ("default", "full") open Now.
func TabFromView(view string) Tab {
	switch view {
	case "hourly":
		return TabHourly
	case "daily":
		return TabDaily
	default:
		return TabNow
	}
}

// FetchFunc gets the weather for a city in the given units.
type FetchFunc func(city, units string) (*api.WeatherResponse, error)

type weatherMsg struct {
	resp *api.WeatherResponse
	err  error
}

// clockMsg redraws the header clock and the "updated" age.
type clockMsg time.Time

type Model struct {
	screen        screen
	splashSeq     int
	splashMinDone bool
	streakFrame   int

	tab        Tab
	hourCursor int
	dayCursor  int
	hourScroll int
	width      int
	height     int
	city       string
	// units is the session units: u changes it without a save
	units string
	cfg   *config.Config
	// auth is shared with the runner's fetch, so a new key takes effect on the next fetch
	auth  *config.AuthConfig
	fetch FetchFunc

	location  *models.City
	weather   *models.OneCallResponse
	fetchedAt time.Time
	loading   bool
	err       error

	spinner   spinner.Model
	cityInput textinput.Model
	prompting bool
	showHelp  bool
	now       func() time.Time

	settingsCursor int
	settingsErr    string
	sub            *setup.Model
	subKind        subView
}

func NewModel(city string, cfg *config.Config, tab Tab, fetch FetchFunc) Model {
	s := spinner.New(spinner.WithSpinner(spinner.Dot))
	s.Style = lipgloss.NewStyle().Foreground(styles.Gold)

	input := textinput.New()
	input.Prompt = "› "
	input.Placeholder = "city, or city,cc"
	inputStyles := textinput.DefaultDarkStyles()
	inputStyles.Focused.Prompt = lipgloss.NewStyle().Foreground(styles.Love).Background(styles.Surface)
	inputStyles.Focused.Text = lipgloss.NewStyle().Foreground(styles.Text).Background(styles.Surface)
	inputStyles.Focused.Placeholder = lipgloss.NewStyle().Foreground(styles.Subtle).Background(styles.Surface)
	inputStyles.Cursor.Color = styles.Gold
	input.SetStyles(inputStyles)

	return Model{
		screen:    screenSplash,
		tab:       tab,
		city:      city,
		units:     cfg.Units,
		cfg:       cfg,
		auth:      &config.AuthConfig{},
		fetch:     fetch,
		loading:   true,
		spinner:   s,
		cityInput: input,
		now:       time.Now,
	}
}

// cityLabel is the resolved "Name, CC", or the requested name before the first fetch.
func (m Model) cityLabel() string {
	if m.location == nil {
		return m.city
	}
	if m.location.Country == "" {
		return m.location.Name
	}
	return m.location.Name + ", " + m.location.Country
}

// ownKey is true when the auth holds the user's own OpenWeatherMap key. Breeze keys start with gust_.
func (m Model) ownKey() bool {
	return m.auth.APIKey != "" && !strings.HasPrefix(m.auth.APIKey, "gust_")
}

func (m Model) via() string {
	if m.ownKey() {
		return "your key"
	}
	return "breeze"
}
