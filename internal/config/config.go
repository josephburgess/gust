package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	APIKey      string `json:"api_key"`
	DefaultCity string `json:"default_city"`
}

type GetConfigPathFunc func() (string, error)

// holding var for the fn, so we can mock in tests
var GetConfigPath GetConfigPathFunc = defaultGetConfigPath

// get the gust .config path
func defaultGetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "gust")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}

	return filepath.Join(configDir, "config.json"), nil
}

// load the config from file
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{}, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("could not decode config file: %w", err)
	}

	return &config, nil
}

// save to cfg file
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("could not create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("could not encode config: %w", err)
	}

	return nil
}

// interactive prompt for first startup
func PromptForConfiguration() (*Config, error) {
	reader := bufio.NewReader(os.Stdin)
	config := &Config{}

	fmt.Println("Welcome to Gust! Let's set up your configuration.")

	fmt.Println("\nYou'll need an OpenWeather API key.")
	fmt.Println("You can get one at: https://home.openweathermap.org/api_keys")
	fmt.Print("Enter your OpenWeather API key: ")
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}
	config.APIKey = strings.TrimSpace(apiKey)

	fmt.Print("\nEnter your default city: ")
	defaultCity, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}
	config.DefaultCity = strings.TrimSpace(defaultCity)

	fmt.Println("\nConfiguration complete!")
	return config, nil
}
