package cli

import (
	"testing"

	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/ui/tui"
	"github.com/stretchr/testify/assert"
)

func TestChooseMode(t *testing.T) {
	tests := []struct {
		name     string
		cli      *CLI
		view     string
		tty      bool
		wantMode runMode
		wantTab  tui.Tab
	}{
		{"terminal opens the dashboard on opens-on", &CLI{}, "hourly", true, modeDashboard, tui.TabHourly},
		{"old default view opens now", &CLI{}, "default", true, modeDashboard, tui.TabNow},
		{"tab flag wins over opens-on", &CLI{Daily: true}, "hourly", true, modeDashboard, tui.TabDaily},
		{"alerts flag wins over the others", &CLI{Alerts: true, Hourly: true}, "now", true, modeDashboard, tui.TabAlerts},
		{"now flag", &CLI{Now: true}, "daily", true, modeDashboard, tui.TabNow},
		{"opens-on compact prints compact", &CLI{}, "compact", true, modeCompact, tui.TabNow},
		{"tab flag gets out of opens-on compact", &CLI{Hourly: true}, "compact", true, modeDashboard, tui.TabHourly},
		{"compact flag", &CLI{Compact: true}, "now", true, modeCompact, tui.TabNow},
		{"print flag", &CLI{Print: true}, "daily", true, modePrint, tui.TabDaily},
		{"print flag with a tab", &CLI{Print: true, Hourly: true}, "now", true, modePrint, tui.TabHourly},
		{"pipe prints", &CLI{}, "now", false, modePrint, tui.TabNow},
		{"pipe with a tab flag prints that tab", &CLI{Alerts: true}, "now", false, modePrint, tui.TabAlerts},
		{"pipe with opens-on compact prints compact", &CLI{}, "compact", false, modeCompact, tui.TabNow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, tab := chooseMode(tt.cli, &config.Config{DefaultView: tt.view}, tt.tty)
			assert.Equal(t, tt.wantMode, mode)
			assert.Equal(t, tt.wantTab, tab)
		})
	}
}
