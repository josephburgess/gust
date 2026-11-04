package cli

import (
	"github.com/alecthomas/kong"
)

type CLI struct {
	// config flags
	Default string `name:"default" short:"D" help:"Set a new default city"`
	ApiUrl  string `name:"api" short:"A" help:"Set custom API server URL"`
	Login   bool   `name:"login" short:"L" help:"Authenticate with GitHub"`
	Setup   bool   `name:"setup" short:"S" help:"Run the setup wizard"`
	Units   string `name:"units" short:"U" help:"Temperature units (metric, imperial, standard)"`
	ApiKey  string `name:"api-key" short:"K" help:"Set your api key (either gust or openweathermap)"`

	// tab flags: open the dashboard on a tab, or pick the tab to print
	Now    bool `name:"now" short:"n" help:"Open on the Now tab"`
	Hourly bool `name:"hourly" short:"r" help:"Open on the Hourly tab"`
	Daily  bool `name:"daily" short:"y" help:"Open on the 5-day tab"`
	Alerts bool `name:"alerts" short:"a" help:"Open on the Alerts tab"`

	// print flags
	Compact bool `name:"compact" short:"c" help:"Print a two-line summary and exit"`
	Print   bool `name:"print" short:"p" help:"Print the tab and exit, without the dashboard"`

	// args (city name)
	Args []string `arg:"" optional:"" help:"City name (can be multiple words)"`
}

func NewApp() (*kong.Kong, *CLI) {
	cli := &CLI{}
	parser := kong.Must(cli,
		kong.Name("gust"),
		kong.Description("Simple terminal weather"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
	)
	return parser, cli
}
