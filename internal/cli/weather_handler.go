package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/josephburgess/gust/internal/api"
	"github.com/josephburgess/gust/internal/config"
	"github.com/josephburgess/gust/internal/models"
	"github.com/josephburgess/gust/internal/ui/components"
	"github.com/josephburgess/gust/internal/ui/renderer"
)

type weatherResponseMsg struct {
	weather *models.OneCallResponse
	city    *models.City
}

type weatherErrorMsg struct {
	err error
}

func fetchWeatherCmd(city string, cfg *config.Config, auth *config.AuthConfig) tea.Cmd {
	return func() tea.Msg {
		client := api.NewClient(cfg.ApiUrl, auth.APIKey, cfg.Units)
		weather, err := client.GetWeather(city)
		if err != nil {
			return weatherErrorMsg{err: fmt.Errorf("failed to get weather data: %w", err)}
		}
		return weatherResponseMsg{
			weather: weather.Weather,
			city:    weather.City,
		}
	}
}

type fetchWeatherModel struct {
	components.WeatherFetchModel
	cfg      *config.Config
	auth     *config.AuthConfig
	cli      *CLI
	weather  *models.OneCallResponse
	cityData *models.City
}

func (m fetchWeatherModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case weatherResponseMsg:
		m.weather = msg.weather
		m.cityData = msg.city
		m.Done = true
		return m, tea.Quit

	case weatherErrorMsg:
		m.Err = msg.err
		m.Done = true
		return m, tea.Quit
	}

	model, cmd := m.WeatherFetchModel.Update(msg)
	if weatherModel, ok := model.(components.WeatherFetchModel); ok {
		m.WeatherFetchModel = weatherModel
		return m, cmd
	}
	return m, cmd
}

func fetchAndRenderWeather(city string, cfg *config.Config, authConfig *config.AuthConfig, cli *CLI) error {
	fetchCmd := fetchWeatherCmd(city, cfg, authConfig)
	baseModel := components.NewWeatherFetchModel(city, "", fetchCmd)
	model := fetchWeatherModel{
		WeatherFetchModel: baseModel,
		cfg:               cfg,
		auth:              authConfig,
		cli:               cli,
	}

	p := tea.NewProgram(model)
	finalModel, _ := p.Run()
	m, _ := finalModel.(fetchWeatherModel)

	if m.weather != nil && m.cityData != nil {
		weatherRenderer := renderer.NewWeatherRenderer("terminal", cfg.Units)
		renderWeatherView(cli, weatherRenderer, m.cityData, m.weather)
	}

	return nil
}

func renderWeatherView(cli *CLI, weatherRenderer renderer.WeatherRenderer, city *models.City, weather *models.OneCallResponse) {
	if cli.Alerts {
		weatherRenderer.RenderAlerts(city, weather)
	} else if cli.Hourly {
		weatherRenderer.RenderHourlyForecast(city, weather)
	} else if cli.Daily {
		weatherRenderer.RenderDailyForecast(city, weather)
	} else if cli.Full {
		weatherRenderer.RenderFullWeather(city, weather)
	} else if cli.Compact {
		weatherRenderer.RenderCompactWeather(city, weather)
	} else if cli.Detailed {
		weatherRenderer.RenderCurrentWeather(city, weather)
	} else {
		weatherRenderer.RenderCurrentWeather(city, weather)
	}
}
