package components

import (
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
