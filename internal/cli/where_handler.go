package cli

import (
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/ui/output"
)

func handleWhere(cfg *config.Config) error {
	city := cfg.DefaultCity
	if city == "" {
		city = "(not set)"
	}
	output.PrintWhereConfig(city, cfg.Units, cfg.DefaultView, cfg.ApiUrl, cfg.ShowTips)
	return nil
}
