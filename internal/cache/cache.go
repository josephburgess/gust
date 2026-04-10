package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/josephburgess/gust/internal/api"
)

const TTL = 10 * time.Minute

type entry struct {
	Data     *api.WeatherResponse `json:"data"`
	CachedAt time.Time            `json:"cached_at"`
}

type Cache struct {
	dir string
}

func New() (*Cache, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine home directory: %w", err)
	}
	dir := filepath.Join(homeDir, ".config", "gust", "cache")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("could not create cache directory: %w", err)
	}
	return &Cache{dir: dir}, nil
}

// Get returns cached weather with its age. False if entry doesn't exist or expired.
func (c *Cache) Get(city, units string) (*api.WeatherResponse, time.Duration, bool) {
	data, err := os.ReadFile(c.filePath(city, units))
	if err != nil {
		return nil, 0, false
	}

	var e entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, 0, false
	}

	age := time.Since(e.CachedAt)
	if age > TTL {
		os.Remove(c.filePath(city, units))
		return nil, 0, false
	}

	return e.Data, age, true
}

func (c *Cache) Set(city, units string, data *api.WeatherResponse) error {
	e := entry{
		Data:     data,
		CachedAt: time.Now(),
	}
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return os.WriteFile(c.filePath(city, units), b, 0644)
}

func (c *Cache) filePath(city, units string) string {
	key := strings.ToLower(city)
	key = strings.NewReplacer(" ", "_", "/", "_", "\\", "_", "..", "_").Replace(key)
	if units != "" {
		key += "-" + units
	}
	return filepath.Join(c.dir, key+".json")
}

// FormatAge converts a cache age into a human-readable string.
func FormatAge(d time.Duration) string {
	minutes := int(d.Minutes())
	switch {
	case minutes < 1:
		return "just now"
	case minutes == 1:
		return "1m ago"
	default:
		return fmt.Sprintf("%dm ago", minutes)
	}
}
