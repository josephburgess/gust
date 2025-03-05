package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josephburgess/gust/internal/ui/styles"
)

type SpinnerModel struct {
	Spinner spinner.Model
}

func NewSpinner() SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(styles.Foam)
	return SpinnerModel{
		Spinner: s,
	}
}

func NewCustomSpinner(spinnerType spinner.Spinner, color lipgloss.Color) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinnerType
	s.Style = lipgloss.NewStyle().Foreground(color)
	return SpinnerModel{
		Spinner: s,
	}
}

func (s SpinnerModel) Tick() tea.Cmd {
	return s.Spinner.Tick
}

func (s SpinnerModel) Update(msg tea.Msg) (SpinnerModel, tea.Cmd) {
	var cmd tea.Cmd
	spinner, cmd := s.Spinner.Update(msg)
	s.Spinner = spinner
	return s, cmd
}

func (s SpinnerModel) View() string {
	return s.Spinner.View()
}

var WeatherEmojis = spinner.Spinner{
	Frames: []string{"☀️ ", "⛅️ ", "☁️ ", "🌧️ ", "⛈️ ", "❄️ ", "🌪️ ", "🌈 "},
	FPS:    time.Second / 5,
}

type WeatherFetchModel struct {
	Spinner  SpinnerModel
	City     string
	Message  string
	Weather  any
	CityData any
	Err      error
	Done     bool
	FetchCmd tea.Cmd
}

func NewWeatherFetchModel(city string, message string, fetchCmd tea.Cmd) WeatherFetchModel {
	spinnerModel := NewCustomSpinner(WeatherEmojis, styles.Foam)

	return WeatherFetchModel{
		Spinner:  spinnerModel,
		City:     city,
		Message:  message,
		FetchCmd: fetchCmd,
	}
}

func (m WeatherFetchModel) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick(),
		m.FetchCmd,
	)
}

func (m WeatherFetchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m WeatherFetchModel) View() string {
	if m.Done {
		return ""
	}

	if m.Message != "" {
		return fmt.Sprintf("%s %s", m.Spinner.View(), m.Message)
	}

	return fmt.Sprintf("%s Fetching weather for %s...", m.Spinner.View(), m.City)
}
