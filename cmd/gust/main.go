package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/josephburgess/gust/internal/cli"
	"github.com/josephburgess/gust/internal/ui/styles"
)

// Version is set at build time via -ldflags "-X main.Version=x.y.z"
var Version = "dev"

func main() {
	_ = godotenv.Load()
	app, cliInstance := cli.NewApp(Version)
	ctx, err := app.Parse(os.Args[1:])
	if err != nil {
		styles.ExitWithError("Failed to parse command line arguments", err)
	}

	if err := cli.Run(ctx, cliInstance); err != nil {
		styles.ExitWithError(fmt.Sprintf("Command failed: %s", ctx.Command()), err)
	}
}
