package main

import (
	"flag"
	"fmt"
)

func main() {
	city := flag.String("city", "New York", "City name for weather lookup")
	flag.Parse()

	fmt.Println("Fetching weather for:", *city)
}
