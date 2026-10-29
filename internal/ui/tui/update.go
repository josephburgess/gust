package tui

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
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

// refetch starts a new fetch. clear drops the old weather, for a new city or new units.
func (m Model) refetch(clear bool) (Model, tea.Cmd) {
	if clear {
		m.weather, m.location = nil, nil
		m.hourCursor, m.dayCursor, m.hourScroll = 0, 0, 0
	}
	m.err = nil
	if m.loading {
		return m, m.fetchCmd()
	}
	m.loading = true
	return m, tea.Batch(m.spinner.Tick, m.fetchCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m = m.scrollToCursor()
	case tea.KeyPressMsg:
		if m.prompting {
			return m.handlePrompt(msg)
		}
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
	default:
		if m.prompting {
			var cmd tea.Cmd
			m.cityInput, cmd = m.cityInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m Model) handlePrompt(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.prompting = false
		m.cityInput.Blur()
		return m, nil
	case "enter":
		m.prompting = false
		m.cityInput.Blur()
		city := strings.TrimSpace(m.cityInput.Value())
		if city == "" {
			return m, nil
		}
		m.city = city
		return m.refetch(true)
	}
	var cmd tea.Cmd
	m.cityInput, cmd = m.cityInput.Update(msg)
	return m, cmd
}

func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.showHelp && msg.String() == "esc" {
		m.showHelp = false
		return m, nil
	}

	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, keys.Tabs):
		m.tab = Tab(msg.String()[0] - '1')
	case key.Matches(msg, keys.Next):
		m.tab = (m.tab + 1) % Tab(len(tabNames))
	case key.Matches(msg, keys.Prev):
		m.tab = (m.tab + Tab(len(tabNames)) - 1) % Tab(len(tabNames))
	case key.Matches(msg, keys.Down):
		return m.moveCursor(1), nil
	case key.Matches(msg, keys.Up):
		return m.moveCursor(-1), nil
	case key.Matches(msg, keys.City):
		m.prompting = true
		m.cityInput.Reset()
		return m, m.cityInput.Focus()
	case key.Matches(msg, keys.Units):
		m.units = nextUnits(m.units)
		return m.refetch(true)
	case key.Matches(msg, keys.Refresh):
		return m.refetch(false)
	case key.Matches(msg, keys.Help):
		m.showHelp = !m.showHelp
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
