package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/errhandler"
	"github.com/josephburgess/gust/internal/models"
)

type WeatherResponse struct {
	City    *models.City            `json:"city"`
	Weather *models.OneCallResponse `json:"weather"`
}

var (
	headerStyle    = color.New(color.FgHiCyan, color.Bold).SprintFunc()
	tempStyle      = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	highlightStyle = color.New(color.FgHiWhite).SprintFunc()
	infoStyle      = color.New(color.FgHiBlue).SprintFunc()
	timeStyle      = color.New(color.FgHiYellow).SprintFunc()
	alertStyle     = color.New(color.FgHiRed, color.Bold).SprintFunc()
)

func main() {
	_ = godotenv.Load()

	// Command-line flags
	cityPtr := flag.String("city", "", "Name of the city")
	setDefaultCityPtr := flag.String("default", "", "Set a new default city")
	loginPtr := flag.Bool("login", false, "Authenticate with GitHub")
	logoutPtr := flag.Bool("logout", false, "Log out and remove authentication")
	apiURLPtr := flag.String("api", "", "Set custom API server URL")
	fullPtr := flag.Bool("full", false, "Show full weather report including daily and hourly forecasts")
	dailyPtr := flag.Bool("daily", false, "Show daily forecast only")
	hourlyPtr := flag.Bool("hourly", false, "Show hourly forecast only")
	alertsPtr := flag.Bool("alerts", false, "Show weather alerts only")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		errhandler.CheckFatal(err, "Failed to load configuration")
	}

	if *apiURLPtr != "" {
		cfg.APIURL = *apiURLPtr
		if err := cfg.Save(); err != nil {
			errhandler.CheckFatal(err, "Failed to save configuration")
		}
		fmt.Printf("API server URL set to: %s\n", *apiURLPtr)
		return
	}

	if cfg.APIURL == "" {
		cfg.APIURL = "https://gust.ngrok.io"
		if err := cfg.Save(); err != nil {
			errhandler.CheckFatal(err, "Failed to save configuration")
		}
	}

	if *logoutPtr {
		authConfig, _ := config.LoadAuthConfig()
		if authConfig != nil {
			authConfigPath, err := config.GetAuthConfigPath()
			if err != nil {
				log.Fatalf("Failed to get auth config path: %v", err)
			}

			if err := os.Remove(authConfigPath); err != nil {
				log.Fatalf("Failed to remove auth config file: %v", err)
			}
			fmt.Println("Logged out successfully")
		} else {
			fmt.Println("Not currently logged in")
		}
		return
	}

	if *loginPtr {
		fmt.Println("Starting GitHub authentication...")
		authConfig, err := config.Authenticate(cfg.APIURL)
		if err != nil {
			errhandler.CheckFatal(err, "Authentication failed")
		}

		if err := config.SaveAuthConfig(authConfig); err != nil {
			errhandler.CheckFatal(err, "Failed to save authentication")
		}

		fmt.Printf("Successfully authenticated as %s\n", authConfig.GithubUser)
		return
	}

	if *setDefaultCityPtr != "" {
		cfg.DefaultCity = *setDefaultCityPtr
		if err := cfg.Save(); err != nil {
			errhandler.CheckFatal(err, "Failed to save configuration")
		}
		fmt.Printf("Default city set to: %s\n", *setDefaultCityPtr)
		return
	}

	authConfig, err := config.LoadAuthConfig()
	if err != nil {
		errhandler.CheckFatal(err, "Failed to load authentication")
	}

	if authConfig == nil {
		fmt.Println("You need to authenticate with GitHub before using Gust.")
		fmt.Println("Run 'gust --login' to authenticate.")
		os.Exit(1)
	}

	var cityName string
	args := flag.Args()
	if *cityPtr != "" {
		cityName = *cityPtr
	} else if len(args) > 0 {
		cityName = strings.Join(args, " ")
	} else {
		cityName = cfg.DefaultCity
	}

	if cityName == "" {
		fmt.Println("No city specified and no default city set.")
		fmt.Println("Specify a city: gust [city name]")
		fmt.Println("Or set a default city: gust --default \"London\"")
		os.Exit(1)
	}

	weather, err := getWeather(cfg.APIURL, authConfig.APIKey, cityName)
	errhandler.CheckFatal(err, "Failed to get weather data")

	if *alertsPtr {
		displayAlerts(weather.City, weather.Weather)
	} else if *hourlyPtr {
		displayHourlyForecast(weather.City, weather.Weather)
	} else if *dailyPtr {
		displayDailyForecast(weather.City, weather.Weather)
	} else if *fullPtr {
		displayFullWeather(weather.City, weather.Weather)
	} else {
		displayCurrentWeather(weather.City, weather.Weather)
	}
}

func getWeather(apiURL, apiKey, cityName string) (*WeatherResponse, error) {
	url := fmt.Sprintf("%s/api/weather/%s?api_key=%s",
		apiURL, url.QueryEscape(cityName), apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var response WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return &response, nil
}

func displayCurrentWeather(city *models.City, weather *models.OneCallResponse) {
	current := weather.Current

	fmt.Printf("\n%s\n", headerStyle(fmt.Sprintf("WEATHER FOR %s", strings.ToUpper(city.Name))))
	fmt.Println(strings.Repeat("─", 50))

	if len(current.Weather) > 0 {
		weatherCond := current.Weather[0]
		fmt.Printf("Current Conditions: %s %s\n\n",
			highlightStyle(weatherCond.Description),
			models.GetWeatherEmoji(weatherCond.ID))

		fmt.Printf("Temperature: %s %s (Feels like: %.1f°C)\n",
			tempStyle(fmt.Sprintf("%.1f°C", models.KelvinToCelsius(current.Temp))),
			"🌡️",
			models.KelvinToCelsius(current.FeelsLike))

		fmt.Printf("Humidity: %d%% %s\n", current.Humidity, "💧")

		fmt.Printf("UV Index: %.1f ☀️\n", current.UVI)

		if current.WindGust > 0 {
			fmt.Printf("Wind: %.1f km/h %s %s (Gusts: %.1f km/h)\n",
				current.WindSpeed*3.6, // m/s to km/h
				models.GetWindDirection(current.WindDeg),
				"💨",
				current.WindGust*3.6)
		} else {
			fmt.Printf("Wind: %.1f km/h %s %s\n",
				current.WindSpeed*3.6, // m/s to km/h
				models.GetWindDirection(current.WindDeg),
				"💨")
		}

		if current.Clouds > 0 {
			fmt.Printf("Cloud coverage: %d%% ☁️\n", current.Clouds)
		}

		if current.Rain != nil && current.Rain.OneHour > 0 {
			fmt.Printf("Rain: %.1f mm (last hour) 🌧️\n", current.Rain.OneHour)
		}

		if current.Snow != nil && current.Snow.OneHour > 0 {
			fmt.Printf("Snow: %.1f mm (last hour) ❄️\n", current.Snow.OneHour)
		}

		fmt.Printf("Visibility: %s\n", models.VisibilityToString(current.Visibility))

		fmt.Printf("Sunrise: %s %s  Sunset: %s %s\n\n",
			time.Unix(current.Sunrise, 0).Format("15:04"),
			"🌅",
			time.Unix(current.Sunset, 0).Format("15:04"),
			"🌇")
	}

	if len(weather.Alerts) > 0 {
		alertStyle := color.New(color.FgHiRed, color.Bold).SprintFunc()
		fmt.Printf("%s Use 'gust --alerts %s' to view them.\n",
			alertStyle(fmt.Sprintf("⚠️  There are %d weather alerts for this area.", len(weather.Alerts))),
			city.Name)
	}
}

func displayDailyForecast(city *models.City, weather *models.OneCallResponse) {
	fmt.Printf("\n%s\n", headerStyle(fmt.Sprintf("7-DAY FORECAST FOR %s", strings.ToUpper(city.Name))))
	fmt.Println(strings.Repeat("─", 50))

	if len(weather.Daily) > 0 {
		for i, day := range weather.Daily {
			if i >= 5 {
				break
			}

			date := time.Unix(day.Dt, 0).Format("Mon Jan 2")

			if i > 0 {
				fmt.Println()
			}

			fmt.Printf("%s: %s\n",
				highlightStyle(date),
				day.Summary)

			fmt.Printf("  High/Low: %s/%s %s\n",
				tempStyle(fmt.Sprintf("%.1f°C", models.KelvinToCelsius(day.Temp.Max))),
				tempStyle(fmt.Sprintf("%.1f°C", models.KelvinToCelsius(day.Temp.Min))),
				"🌡️")

			fmt.Printf("  Morning: %.1f°C  Day: %.1f°C  Evening: %.1f°C  Night: %.1f°C\n",
				models.KelvinToCelsius(day.Temp.Morn),
				models.KelvinToCelsius(day.Temp.Day),
				models.KelvinToCelsius(day.Temp.Eve),
				models.KelvinToCelsius(day.Temp.Night))

			if len(day.Weather) > 0 {
				weather := day.Weather[0]
				condition := fmt.Sprintf("%s %s", weather.Description, models.GetWeatherEmoji(weather.ID))
				fmt.Printf("  Conditions: %s\n", infoStyle(condition))
			}

			if day.Pop > 0 {
				fmt.Printf("  Precipitation: %d%% chance\n", int(day.Pop*100))
			}
			if day.Rain > 0 {
				fmt.Printf("  Rain: %.1f mm 🌧️\n", day.Rain)
			}

			if day.Snow > 0 {
				fmt.Printf("  Snow: %.1f mm ❄️\n", day.Snow)
			}

			fmt.Printf("  Wind: %.1f km/h %s\n",
				day.WindSpeed*3.6,
				models.GetWindDirection(day.WindDeg))

			fmt.Printf("  UV Index: %.1f\n", day.UVI)
		}
		fmt.Println()
	}
}

func displayHourlyForecast(city *models.City, weather *models.OneCallResponse) {
	fmt.Printf("\n%s\n", headerStyle(fmt.Sprintf("HOURLY FORECAST FOR %s", strings.ToUpper(city.Name))))
	fmt.Println(strings.Repeat("─", 50))

	if len(weather.Hourly) > 0 {
		hourLimit := 24
		hourLimit = int(math.Min(float64(hourLimit), float64(len(weather.Hourly))))

		currentDay := ""

		for i := 0; i < hourLimit; i++ {
			hour := weather.Hourly[i]

			t := time.Unix(hour.Dt, 0)
			day := t.Format("Mon Jan 2")
			hourStr := t.Format("15:04")

			if day != currentDay {
				if currentDay != "" {
					fmt.Println()
				}
				fmt.Printf("%s:\n", highlightStyle(day))
				currentDay = day
			}

			if len(hour.Weather) == 0 {
				continue
			}

			weatherCond := hour.Weather[0]

			temp := tempStyle(fmt.Sprintf("%.1f°C", models.KelvinToCelsius(hour.Temp)))

			popStr := ""
			if hour.Pop > 0 {
				popStr = fmt.Sprintf(" (%.0f%% chance of precipitation)", hour.Pop*100)
			}

			fmt.Printf("  %s: %s %s  %s%s\n",
				hourStr,
				temp,
				models.GetWeatherEmoji(weatherCond.ID),
				infoStyle(weatherCond.Description),
				popStr)

			if hour.Rain != nil && hour.Rain.OneHour > 0 {
				fmt.Printf("       Rain: %.1f mm/h\n", hour.Rain.OneHour)
			}

			if hour.Snow != nil && hour.Snow.OneHour > 0 {
				fmt.Printf("       Snow: %.1f mm/h\n", hour.Snow.OneHour)
			}
		}
		fmt.Println()
	}
}

func displayAlerts(city *models.City, weather *models.OneCallResponse) {
	fmt.Printf("\n%s\n", headerStyle(fmt.Sprintf("WEATHER ALERTS FOR %s", strings.ToUpper(city.Name))))
	fmt.Println(strings.Repeat("─", 50))

	if len(weather.Alerts) == 0 {
		fmt.Println("No weather alerts for this area.")
		return
	}

	for i, alert := range weather.Alerts {
		if i > 0 {
			fmt.Println(strings.Repeat("─", 50))
		}

		fmt.Printf("%s\n", alertStyle(fmt.Sprintf("⚠️  %s", alert.Event)))
		fmt.Printf("Issued by: %s\n", alert.SenderName)
		fmt.Printf("Valid: %s to %s\n\n",
			timeStyle(time.Unix(alert.Start, 0).Format("Mon Jan 2 15:04")),
			timeStyle(time.Unix(alert.End, 0).Format("Mon Jan 2 15:04")))

		fmt.Println(alert.Description)
		fmt.Println()
	}
}

func displayFullWeather(city *models.City, weather *models.OneCallResponse) {
	displayCurrentWeather(city, weather)
	fmt.Println()
	if len(weather.Alerts) > 0 {
		displayAlerts(city, weather)
		fmt.Println()
	}
	displayDailyForecast(city, weather)
	fmt.Println()
}
