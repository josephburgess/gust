package cli

import (
	"github.com/alecthomas/kong"
)

type CLI struct {
	City    string `name:"city" short:"C" help:"City name"`
	Default string `name:"default" short:"D" help:"Set a new default city"`
	Login   bool   `name:"login" short:"l" help:"Authenticate with GitHub"`
	APIURL  string `name:"api" short:"a" help:"Set custom API server URL"`
	Setup   bool   `name:"setup" short:"s" help:"Run the setup wizard"`

	Compact  bool   `name:"compact" short:"c" help:"Show today's compact weather view"`
	Detailed bool   `name:"detailed" short:"d" help:"Show today's detailed weather view"`
	Full     bool   `name:"full" short:"f" help:"Show today, 5 day and weather alert forecasts"`
	Daily    bool   `name:"daily" short:"y" help:"Show 5 day forecast"`
	Hourly   bool   `name:"hourly" short:"h" help:"Show 24 hour (hourly) forecast"`
	Alerts   bool   `name:"alerts" short:"a" help:"Show weather alerts"`
	Units    string `name:"units" short:"u" help:"Temperature units (metric, imperial, standard). Metric is default"`
	Pretty   bool   `name:"pretty" short:"p" hidden:"" help:"Use the pretty UI - tbc"`

	Args []string `arg:"" optional:"" help:"City name (can be multiple words)"`
}

func NewApp() (*kong.Kong, *CLI) {
	cli := &CLI{}

	parser := kong.Must(cli,
		kong.Name("gust"),
		kong.Description("Simple terminal weather 🌤️"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
	)

	return parser, cli
}
