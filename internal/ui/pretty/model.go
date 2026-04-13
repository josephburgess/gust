package pretty

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/styles"
)

const (
	TabCurrent = 0
	TabHourly  = 1
	TabDaily   = 2
	TabAlerts  = 3
)

var tabNames = []string{"Current", "Hourly", "Daily", "Alerts"}

// local tab styles — distinct from the setup wizard
var (
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.Foam).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Iris).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(styles.Subtle).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(styles.Muted).
				Padding(0, 1)

	dimTabStyle = lipgloss.NewStyle().
			Foreground(styles.Muted).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Overlay).
			Padding(0, 1)

	footerStyle = lipgloss.NewStyle().
			Foreground(styles.Subtle).
			Italic(true)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.Rose)
)

type Model struct {
	city      *models.City
	weather   *models.OneCallResponse
	cfg       *config.Config
	units     string
	activeTab int
	viewport  viewport.Model
	ready     bool
	width     int
	height    int
}

func New(city *models.City, weather *models.OneCallResponse, cfg *config.Config, units string, startTab int) Model {
	return Model{
		city:      city,
		weather:   weather,
		cfg:       cfg,
		units:     units,
		activeTab: startTab,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) setTab(tab int) {
	if tab < 0 {
		tab = len(tabNames) - 1
	}
	if tab >= len(tabNames) {
		tab = 0
	}
	m.activeTab = tab
	m.viewport.SetContent(m.tabContent())
	m.viewport.GotoTop()
}

func (m *Model) initViewport() {
	headerH := lipgloss.Height(m.tabBar())
	footerH := lipgloss.Height(m.footer())
	m.viewport = viewport.New(m.width, m.height-headerH-footerH)
	m.viewport.SetContent(m.tabContent())
	m.ready = true
}
