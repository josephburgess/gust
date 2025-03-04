package cli

import (
	"fmt"
	"strings"
)

func determineCityName(cityFlag string, args []string, defaultCity string) string {
	if cityFlag != "" {
		return cityFlag
	}
	if len(args) > 0 {
		return strings.Join(args, " ") // might be multi word city
	}
	return defaultCity
}

func handleMissingCity() error {
	fmt.Println("No city specified and no default city set.")
	fmt.Println("Specify a city: gust [city name]")
	fmt.Println("Or set a default city: gust --default \"London\"")
	fmt.Println("Or run the setup wizard: gust --setup")
	return fmt.Errorf("no city provided")
}
