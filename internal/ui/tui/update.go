package tui

import (
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
)

const clockInterval = 30 * time.Second

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchCmd(), clockTick())
}

func (m Model) fetchCmd() tea.Cmd {
	city, units, fetch := m.city, m.units, m.fetch
	return func() tea.Msg {
		resp, err := fetch(city, units)
		return weatherMsg{resp: resp, err: err}
	}
}

func clockTick() tea.Cmd {
	return tea.Tick(clockInterval, func(t time.Time) tea.Msg { return clockMsg(t) })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m = m.scrollToCursor()
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	case weatherMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.location = msg.resp.City
			m.weather = msg.resp.Weather
			m.fetchedAt = m.now()
		}
	case clockMsg:
		return m, clockTick()
	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch key := msg.String(); key {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "1", "2", "3", "4":
		m.tab = Tab(key[0] - '1')
	case "tab":
		m.tab = (m.tab + 1) % Tab(len(tabNames))
	case "shift+tab":
		m.tab = (m.tab + Tab(len(tabNames)) - 1) % Tab(len(tabNames))
	case "j", "down":
		return m.moveCursor(1), nil
	case "k", "up":
		return m.moveCursor(-1), nil
	}
	return m, nil
}

func (m Model) moveCursor(step int) Model {
	if m.weather == nil {
		return m
	}
	switch m.tab {
	case TabHourly:
		m.hourCursor = clamp(m.hourCursor+step, 0, max(len(m.hours())-1, 0))
		return m.scrollToCursor()
	case TabDaily:
		m.dayCursor = clamp(m.dayCursor+step, 0, max(len(m.days())-1, 0))
	}
	return m
}
