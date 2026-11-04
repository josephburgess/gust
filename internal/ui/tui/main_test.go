package tui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/josephburgess/gust/internal/config"
)

// TestMain keeps the settings saves out of the real ~/.config/gust.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "gust-tui-test")
	if err != nil {
		panic(err)
	}
	config.GetConfigPath = func() (string, error) { return filepath.Join(dir, "config.json"), nil }
	config.GetAuthConfigPath = func() (string, error) { return filepath.Join(dir, "auth.json"), nil }

	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}
