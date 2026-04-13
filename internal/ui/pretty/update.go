package pretty

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.initViewport()
		} else {
			headerH := lipglossHeight(m.tabBar())
			footerH := lipglossHeight(m.footer())
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - headerH - footerH
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit

		case "tab", "right", "l":
			m.setTab(m.activeTab + 1)
			return m, nil

		case "shift+tab", "left", "h":
			m.setTab(m.activeTab - 1)
			return m, nil

		case "1":
			m.setTab(TabCurrent)
			return m, nil
		case "2":
			m.setTab(TabHourly)
			return m, nil
		case "3":
			m.setTab(TabDaily)
			return m, nil
		case "4":
			m.setTab(TabAlerts)
			return m, nil
		}
	}

	// delegate remaining keys (j/k, up/down, page up/down) to viewport
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// lipglossHeight returns the rendered height of a string.
func lipglossHeight(s string) int {
	count := 0
	for _, c := range s {
		if c == '\n' {
			count++
		}
	}
	return count + 1
}
