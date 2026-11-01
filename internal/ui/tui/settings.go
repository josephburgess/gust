package tui

import (
	"io"
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/ui/setup"
	"github.com/josephburgess/gust/internal/ui/styles"
	"github.com/josephburgess/gust/internal/ui/widgets"
)

type settingKind int

const (
	settingCity settingKind = iota
	settingUnits
	settingView
	settingTips
	settingVia
	settingAPI
	settingConfigFile
	settingSetup
)

// subView is a setup screen opened from settings.
type subView int

const (
	subNone subView = iota
	subCity
	subSignIn
	subSetup
)

type setting struct {
	kind         settingKind
	group, label string
	help, keys   string
}

var settings = []setting{
	{settingCity, "location", "default city", "The city gust opens on. c on the dashboard changes the city for this session only.", "enter change"},
	{settingUnits, "display", "units", "Units for temperature and wind. A change saves and fetches again.", "←→ change"},
	{settingView, "display", "opens on", "The tab a plain gust run opens on. compact prints two lines and exits.", "←→ change"},
	{settingTips, "display", "tips", "A one-line tip under the weather.", "←→ change"},
	{settingVia, "account", "weather via", "Sign in with GitHub through breeze, or use your own OpenWeatherMap key.", "enter change"},
	{settingAPI, "advanced", "api server", "The breeze server gust talks to. Set it with gust --api URL.", ""},
	{settingConfigFile, "advanced", "config file", "Open the config file in $EDITOR. gust reads it again when you close the editor.", "enter open"},
	{settingSetup, "advanced", "run setup again", "Run the whole setup wizard again.", "enter run"},
}

var (
	unitValues = []string{"metric", "imperial", "standard"}
	viewValues = []string{"now", "hourly", "daily", "compact"}
)

// editorDoneMsg comes back when $EDITOR closes.
type editorDoneMsg struct{ err error }

// githubDoneMsg comes back when the GitHub sign in finishes.
type githubDoneMsg struct {
	auth *config.AuthConfig
	err  error
}

// githubSignIn runs config.Authenticate outside the TUI, because it prints and opens a browser.
type githubSignIn struct {
	apiURL string
	result *config.AuthConfig
}

func (g *githubSignIn) Run() error {
	auth, err := config.Authenticate(g.apiURL)
	if err != nil {
		return err
	}
	g.result = auth
	return config.SaveAuthConfig(auth)
}

func (g *githubSignIn) SetStdin(io.Reader)  {}
func (g *githubSignIn) SetStdout(io.Writer) {}
func (g *githubSignIn) SetStderr(io.Writer) {}

func cycle(values []string, current string, step int) string {
	i := 0
	for j, v := range values {
		if v == current {
			i = j
		}
	}
	return values[(i+step+len(values))%len(values)]
}

func (m Model) settingValue(kind settingKind) string {
	switch kind {
	case settingCity:
		return m.cfg.DefaultCity
	case settingUnits:
		return m.cfg.Units
	case settingView:
		return viewLabel(m.cfg.DefaultView)
	case settingTips:
		if m.cfg.ShowTips {
			return "on"
		}
		return "off"
	case settingVia:
		if m.ownKey() {
			return "your key ••••" + m.auth.APIKey[max(len(m.auth.APIKey)-4, 0):]
		}
		if m.auth.GithubUser != "" {
			return "GitHub (" + m.auth.GithubUser + ")"
		}
		return "GitHub"
	case settingAPI:
		return m.cfg.ApiUrl
	case settingConfigFile:
		path, _ := config.GetConfigPath()
		if home, err := os.UserHomeDir(); err == nil && strings.HasPrefix(path, home) {
			return "~" + strings.TrimPrefix(path, home)
		}
		return path
	}
	return ""
}

func viewLabel(view string) string {
	switch view {
	case "hourly":
		return "Hourly"
	case "daily":
		return "5-day"
	case "compact":
		return "compact"
	default:
		return "Now"
	}
}

func isCycle(kind settingKind) bool {
	return kind == settingUnits || kind == settingView || kind == settingTips
}

func (m Model) openSettings() Model {
	m.screen = screenSettings
	m.settingsErr = ""
	return m
}

// saveConfig writes the config and keeps any error for the settings screen.
func (m Model) saveConfig() Model {
	m.settingsErr = ""
	if err := m.cfg.Save(); err != nil {
		m.settingsErr = err.Error()
	}
	return m
}

func (m Model) changeSetting(step int) (Model, tea.Cmd) {
	switch settings[m.settingsCursor].kind {
	case settingUnits:
		m.cfg.Units = cycle(unitValues, m.cfg.Units, step)
		m.units = m.cfg.Units
		m = m.saveConfig()
		return m.refetch(false)
	case settingView:
		m.cfg.DefaultView = cycle(viewValues, m.cfg.DefaultView, step)
	case settingTips:
		m.cfg.ShowTips = !m.cfg.ShowTips
	default:
		return m, nil
	}
	return m.saveConfig(), nil
}

// openSub opens a setup screen on a copy of the config, so esc leaves the config alone.
func (m Model) openSub(kind subView, state setup.SetupState) (Model, tea.Cmd) {
	cfg := *m.cfg
	sub := setup.NewModel(&cfg, kind != subCity, api.NewClient(m.cfg.ApiUrl, m.auth.APIKey, m.cfg.Units))
	sub.State = state
	sub.Width, sub.Height = m.width, m.height
	m.sub, m.subKind = &sub, kind
	return m, sub.Init()
}

func (m Model) activateSetting() (Model, tea.Cmd) {
	switch settings[m.settingsCursor].kind {
	case settingCity:
		return m.openSub(subCity, setup.StateCity)
	case settingVia:
		return m.openSub(subSignIn, setup.StateApiKeyOption)
	case settingSetup:
		return m.openSub(subSetup, setup.StateCity)
	case settingConfigFile:
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
		path, err := config.GetConfigPath()
		if err != nil {
			m.settingsErr = err.Error()
			return m, nil
		}
		return m, tea.ExecProcess(exec.Command(editor, path), func(err error) tea.Msg { return editorDoneMsg{err} })
	}
	if isCycle(settings[m.settingsCursor].kind) {
		return m.changeSetting(1)
	}
	return m, nil
}

func (m Model) handleSettingsKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.screen = screenMain
	case "up", "k":
		m.settingsCursor = max(m.settingsCursor-1, 0)
	case "down", "j":
		m.settingsCursor = min(m.settingsCursor+1, len(settings)-1)
	case "left", "h":
		return m.changeSetting(-1)
	case "right", "l":
		return m.changeSetting(1)
	case "enter":
		return m.activateSetting()
	}
	return m, nil
}

// updateSub passes a message to the open setup screen, and closes it when it is done.
func (m Model) updateSub(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "esc" && m.sub.State == firstSubState(m.subKind) {
		m.sub = nil
		return m, nil
	}

	next, cmd := m.sub.Update(msg)
	sub := next.(setup.Model)
	m.sub = &sub

	switch {
	case sub.Quitting:
		m.sub = nil
		return m, nil
	case m.subKind == subCity && sub.State == setup.StatePreferences:
		m.cfg.DefaultCity = sub.Config.DefaultCity
		m.sub = nil
		return m.saveConfig(), nil
	case m.subKind != subCity && sub.State == setup.StateComplete && (m.subKind == subSignIn || isEnter(msg)):
		return m.finishSub(sub)
	}
	return m, cmd
}

func firstSubState(kind subView) setup.SetupState {
	if kind == subSignIn {
		return setup.StateApiKeyOption
	}
	return setup.StateCity
}

func isEnter(msg tea.Msg) bool {
	key, ok := msg.(tea.KeyPressMsg)
	return ok && key.String() == "enter"
}

// finishSub saves what the sign in or setup screen chose, then signs in with GitHub if asked.
func (m Model) finishSub(sub setup.Model) (Model, tea.Cmd) {
	m.sub = nil
	if m.subKind == subSetup {
		*m.cfg = *sub.Config
		m.units = m.cfg.Units
		m = m.saveConfig()
	}
	if sub.WantsGitHub {
		signIn := &githubSignIn{apiURL: m.cfg.ApiUrl}
		return m, tea.Exec(signIn, func(err error) tea.Msg { return githubDoneMsg{auth: signIn.result, err: err} })
	}
	if auth, err := config.LoadAuthConfig(); err == nil && auth != nil {
		*m.auth = *auth
	}
	return m.refetch(false)
}

func (m Model) settingsView() string {
	if m.sub != nil {
		return m.sub.View().Content
	}

	left := fg(styles.Rose).Bold(true).Render("gust") + fg(styles.Subtle).Render(" › ") + lipgloss.NewStyle().Bold(true).Render("settings")
	right := fg(styles.Gold).Render("esc") + fg(styles.Subtle).Render(" back to "+m.cityLabel())
	header := lipgloss.NewStyle().Padding(0, pad).Render(spread(left, right, m.width-2*pad))
	rule := fg(styles.Overlay).Render(strings.Repeat("─", m.width))

	listWidth := min(75, m.width-2*pad-30)
	var rows []string
	group := ""
	for i, s := range settings {
		if s.group != group {
			if group != "" {
				rows = append(rows, "")
			}
			rows = append(rows, fg(styles.Iris).Render(strings.ToUpper(s.group)))
			group = s.group
		}
		rows = append(rows, m.settingRow(i, listWidth))
	}
	if m.settingsErr != "" {
		rows = append(rows, "", fg(styles.Love).Width(listWidth).Render(m.settingsErr))
	}
	list := lipgloss.NewStyle().Width(listWidth).Render(strings.Join(rows, "\n"))

	helpWidth := m.width - 2*pad - listWidth - 3
	s := settings[m.settingsCursor]
	helpBody := lipgloss.NewStyle().Bold(true).Render(s.label) + "\n" +
		fg(styles.Subtle).Width(helpWidth-4).Render(s.help)
	if s.keys != "" {
		k, label, _ := strings.Cut(s.keys, " ")
		helpBody += "\n\n" + fg(styles.Gold).Render(k) + " " + fg(styles.Subtle).Render(label)
	}
	help := widgets.Panel("help", helpBody, helpWidth, false)

	body := lipgloss.NewStyle().Padding(1, pad).Width(m.width).Height(m.bodyHeight()).MaxHeight(m.bodyHeight()).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, list, "   ", help))

	return lipgloss.JoinVertical(lipgloss.Left, header, rule, body, m.keyBar())
}

func (m Model) settingRow(i, width int) string {
	s := settings[i]
	st := rowStyle(i == m.settingsCursor)
	valueColor := styles.Text
	if i == m.settingsCursor {
		valueColor = styles.Foam
	}

	// cursor (3) + label (18) + cycle arrows (4)
	text := ansi.Truncate(m.settingValue(s.kind), width-25, "…")
	value := st(fg(valueColor)).Render(text)
	if isCycle(s.kind) {
		value = st(fg(styles.Gold)).Render("◀ ") + value + st(fg(styles.Gold)).Render(" ▶")
	}
	row := cursorCell(i == m.settingsCursor) +
		st(fg(styles.Subtle).Width(18)).Render(s.label) + value
	return row + st(lipgloss.NewStyle()).Render(strings.Repeat(" ", max(width-lipgloss.Width(row), 0)))
}

func (m Model) settingsKeyHints() []keyHint {
	return []keyHint{{"↑↓", "move"}, {"←→", "change"}, {"enter", "open"}, {"esc", "back"}}
}
