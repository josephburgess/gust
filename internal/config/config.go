package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultCity string `json:"default_city"`
	ApiUrl      string `json:"api_url"`
	Units       string `json:"units"`
	DefaultView string `json:"default_view"`
	ShowTips    bool   `json:"show_tips"`
}

type GetConfigPathFunc func() (string, error)

// for tests
var GetConfigPath GetConfigPathFunc = defaultGetConfigPath

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

func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{Units: "metric", DefaultView: "default"}, nil
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

	if config.Units == "" {
		config.Units = "metric"
	}

	if config.DefaultView == "" {
		config.DefaultView = "default"
	}

	return &config, nil
}

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
