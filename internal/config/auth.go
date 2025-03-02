package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/josephburgess/gust/internal/templates"
)

type AuthConfig struct {
	APIKey     string    `json:"api_key"`
	ServerURL  string    `json:"server_url"`
	LastAuth   time.Time `json:"last_auth"`
	GithubUser string    `json:"github_user"`
}

type GetAuthConfigPathFunc func() (string, error)

var GetAuthConfigPath GetAuthConfigPathFunc = defaultGetAuthConfigPath

func defaultGetAuthConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}

	return filepath.Join(homeDir, ".config", "gust", "auth.json"), nil
}

func Authenticate(serverURL string) (*AuthConfig, error) {
	if serverURL == "" {
		serverURL = "https://gust.ngrok.io"
	}

	authDone := make(chan *AuthConfig)
	errorChan := make(chan error)

	port := 9876
	server := &http.Server{Addr: fmt.Sprintf(":%d", port)}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errorChan <- fmt.Errorf("no auth code received")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Authentication failed: No code provided"))
			return
		}

		apiConfig, err := exchangeCodeForAPIKey(serverURL, code, port)
		if err != nil {
			errorChan <- err
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("Authentication failed: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = templates.RenderSuccessTemplate(w, apiConfig.GithubUser, apiConfig.APIKey, serverURL)
		if err != nil {
			errorChan <- fmt.Errorf("failed to render success template: %w", err)
			return
		}

		authDone <- apiConfig

		go func() {
			time.Sleep(100 * time.Millisecond)
			server.Shutdown(context.TODO())
		}()
	})

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errorChan <- err
		}
	}()

	authURL, err := getAuthURL(serverURL, port)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth URL: %w", err)
	}

	fmt.Printf("Opening browser for GitHub authentication...\n")
	if err := openBrowser(authURL); err != nil {
		fmt.Printf("Could not open browser automatically. Please open this URL manually:\n%s\n", authURL)
	}

	select {
	case config := <-authDone:
		return config, nil
	case err := <-errorChan:
		return nil, err
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("authentication timed out")
	}
}

func getAuthURL(serverURL string, port int) (string, error) {
	url := fmt.Sprintf("%s/api/auth/request?callback_port=%d", serverURL, port)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to contact auth server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	var response struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return response.URL, nil
}

func exchangeCodeForAPIKey(serverURL, code string, port int) (*AuthConfig, error) {
	url := fmt.Sprintf("%s/api/auth/exchange", serverURL)

	reqBody, err := json.Marshal(map[string]any{
		"code":          code,
		"callback_port": port,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		APIKey     string `json:"api_key"`
		GithubUser string `json:"github_user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &AuthConfig{
		APIKey:     response.APIKey,
		ServerURL:  serverURL,
		LastAuth:   time.Now(),
		GithubUser: response.GithubUser,
	}, nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

func SaveAuthConfig(config *AuthConfig) error {
	configPath, err := GetAuthConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create auth config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode auth config: %w", err)
	}

	return nil
}

func LoadAuthConfig() (*AuthConfig, error) {
	configPath, err := GetAuthConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open auth config file: %w", err)
	}
	defer file.Close()

	var config AuthConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode auth config: %w", err)
	}

	return &config, nil
}
