package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	app, cli := NewApp()

	assert.NotNil(t, app)
	assert.NotNil(t, cli)

	fields := []struct {
		name  string
		value any
	}{
		{"City", cli.City},
		{"Default", cli.Default},
		{"ApiUrl", cli.ApiUrl},
		{"Login", cli.Login},
		{"Setup", cli.Setup},
		{"Compact", cli.Compact},
		{"Detailed", cli.Detailed},
		{"Full", cli.Full},
		{"Daily", cli.Daily},
		{"Hourly", cli.Hourly},
		{"Alerts", cli.Alerts},
		{"Units", cli.Units},
		{"Pretty", cli.Pretty},
	}

	for _, field := range fields {
		switch v := field.value.(type) {
		case string:
			assert.Empty(t, v, "%s should be empty", field.name)
		case bool:
			assert.False(t, v, "%s should be false", field.name)
		}
	}

	assert.Empty(t, cli.Args)
}
