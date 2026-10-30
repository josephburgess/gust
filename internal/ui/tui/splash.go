package tui

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/josephburgess/gust/internal/ui/styles"
)

type screen int

const (
	screenSplash screen = iota
	screenMain
)

const (
	// minSplash stops the splash from flashing when the fetch is fast.
	minSplash      = 400 * time.Millisecond
	streakInterval = 120 * time.Millisecond
)

var (
	streakTop = "─────   ──  ────────    ───   ──────  ─    ────"
	streakLow = "   ───────   ─    ─────   ──  ──────────  ─────"
)

// splashMinMsg ends the minimum splash time. seq drops timers from an older splash.
type splashMinMsg struct{ seq int }

type streakMsg struct{}

func splashMinTick(seq int) tea.Cmd {
	return tea.Tick(minSplash, func(time.Time) tea.Msg { return splashMinMsg{seq: seq} })
}

func streakTick() tea.Cmd {
	return tea.Tick(streakInterval, func(time.Time) tea.Msg { return streakMsg{} })
}

func rotate(s string, n int) string {
	r := []rune(s)
	n %= len(r)
	return string(r[n:]) + string(r[:n])
}

// showSplash puts the splash up and starts its timers.
func (m Model) showSplash() (Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	if m.screen != screenSplash || m.splashSeq == 0 {
		cmds = append(cmds, streakTick())
	}
	m.screen = screenSplash
	m.splashSeq++
	m.splashMinDone = false
	return m, tea.Batch(append(cmds, splashMinTick(m.splashSeq))...)
}

// leaveSplash goes to the tabs once the weather is in and the minimum time is up.
func (m Model) leaveSplash() Model {
	if m.screen == screenSplash && m.splashMinDone && !m.loading && m.err == nil && m.weather != nil {
		m.screen = screenMain
	}
	return m
}

func (m Model) splashView() string {
	status := m.spinner.View() + " " + fg(styles.Subtle).Render("fetching "+m.cityLabel()+"…")
	if m.err != nil {
		status = fg(styles.Love).Render(m.err.Error()) + "\n\n" +
			fg(styles.Subtle).Render("r retry · c city · q quit")
	}

	content := lipgloss.JoinVertical(lipgloss.Center,
		fg(styles.Foam).Render(rotate(streakTop, m.streakFrame)),
		fg(styles.Rose).Bold(true).Render(strings.Join(Wordmark(), "\n")),
		fg(styles.Subtle).Render(rotate(streakLow, 2*m.streakFrame)),
		"",
		fg(styles.Subtle).Italic(true).Render("simple terminal weather"),
		"",
		status,
	)
	return lipgloss.Place(m.width, max(m.height-1, 0), lipgloss.Center, lipgloss.Center, content)
}

func (m Model) splashKeyHints() []keyHint {
	if m.err != nil {
		return []keyHint{{"r", "retry"}, {"c", "city"}, {"q", "quit"}}
	}
	return []keyHint{{"q", "cancel"}}
}
