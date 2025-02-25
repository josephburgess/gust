package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const (
	baseUrl = "http://api.openweathermap.org/"
	geoCode = baseUrl + "geo/1.0/direct?q=%s&limit=1&appid=%s"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("error loading .env")
	}

	city := flag.String("city", "New York", "City name for weather lookup")
	flag.Parse()
	apiKey := os.Getenv("OPENWEATHER_API_KEY")

	encodedCity := url.QueryEscape(*city)
	url := fmt.Sprintf(geoCode, encodedCity, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("err making request", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("error reading response", err)
	}

	fmt.Println(url)
	fmt.Println("Fetching weather for:", *city)
	fmt.Println("response:", string(body))
}
