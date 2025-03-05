package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josephburgess/gust/internal/ui/styles"
)

// message types
type successMsg[T any] struct {
	result T
}
type errorMsg struct {
	err error
}

// represents a spinner that runs a function and returns a result
type SpinnerRunnerModel[T any] struct {
	spinner  SpinnerModel
	message  string
	function func() (T, error)
	result   T
	err      error
	done     bool
}

// creates a new runner - displays a spinner while executing a func
func NewSpinnerRunner[T any](
	message string,
	spinnerType spinner.Spinner,
	color lipgloss.Color,
	fn func() (T, error),
) SpinnerRunnerModel[T] {
	return SpinnerRunnerModel[T]{
		spinner:  NewCustomSpinner(spinnerType, color),
		message:  message,
		function: fn,
	}
}

// initialises the spinner runner and starts func
func (m SpinnerRunnerModel[T]) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick(),
		func() tea.Msg {
			result, err := m.function()
			if err != nil {
				return errorMsg{err: err}
			}
			return successMsg[T]{result: result}
		},
	)
}

// handles messages and updates state
func (m SpinnerRunnerModel[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case errorMsg:
		m.err = msg.err
		m.done = true
		return m, tea.Quit
	case successMsg[T]:
		m.result = msg.result
		m.done = true
		return m, tea.Quit
	}
	return m, nil
}

// render
func (m SpinnerRunnerModel[T]) View() string {
	if m.done {
		return ""
	}
	return fmt.Sprintf("%s %s", m.spinner.View(), styles.ProgressMessageStyle.Render(m.message))
}

// custom func that creates and runs a SpinnerRunnerModel
func RunWithSpinner[T any](message string, spinnerType spinner.Spinner, color lipgloss.Color, fn func() (T, error)) (T, error) {
	model := NewSpinnerRunner(message, spinnerType, color, fn)
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		var zero T
		return zero, fmt.Errorf("error running spinner: %w", err)
	}

	if m, ok := finalModel.(SpinnerRunnerModel[T]); ok {
		if m.err != nil {
			var zero T
			return zero, m.err
		}
		return m.result, nil
	}

	var zero T
	return zero, fmt.Errorf("unexpected error in spinner")
}
