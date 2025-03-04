package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/josephburgess/gust/internal/cli"
	"github.com/josephburgess/gust/internal/ui/styles"
)

func main() {
	_ = godotenv.Load()
	app, cliInstance := cli.NewApp()
	ctx, err := app.Parse(os.Args[1:])
	if err != nil {
		styles.ExitWithError("Failed to parse command line arguments", err)
	}

	if err := cli.Run(ctx, cliInstance); err != nil {
		styles.ExitWithError(fmt.Sprintf("Command failed: %s", ctx.Command()), err)
	}
}
