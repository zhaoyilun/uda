package mirror

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/uda/uda/internal/config"
)

// Mirror represents a mirror source
type Mirror struct {
	Name     string
	URL      string
	Priority int
}

// Default mirrors - can be extended
var defaultMirrors = []Mirror{
	{Name: "official", URL: "https://astral.sh", Priority: 0},
	{Name: "tsinghua", URL: "https://pypi.tuna.tsinghua.edu.cn", Priority: 1},
	{Name: "aliyun", URL: "https://mirrors.aliyun.com", Priority: 2},
}

// GetMirror returns the configured mirror URL
func GetMirror() string {
	// 1. Check UV_MIRROR environment variable
	if mirror := os.Getenv("UV_MIRROR"); mirror != "" {
		return mirror
	}

	// 2. Check config file
	cfg, err := loadConfig()
	if err == nil && cfg.Mirror != nil && cfg.Mirror.URL != "" {
		return cfg.Mirror.URL
	}

	// 3. Return default (empty = official)
	return ""
}

// FindWorkingMirror finds the first working mirror
func FindWorkingMirror() (string, error) {
	// Try mirrors in order
	for _, m := range defaultMirrors {
		if testMirror(m.URL) {
			return m.URL, nil
		}
	}
	return "", fmt.Errorf("no working mirror found")
}

// TestMirror tests if a mirror is accessible
func testMirror(url string) bool {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += "simple/"

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

// loadConfig loads configuration from file
func loadConfig() (*config.Config, error) {
	var cfg config.Config

	_, err := toml.DecodeFile(config.ConfigPath(), &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveMirror saves mirror configuration
func SaveMirror(url string) error {
	cfg := &config.Config{
		Mirror: &config.MirrorConfig{
			URL:      url,
			Priority: 0,
		},
	}

	file, err := os.Create(config.ConfigPath())
	if err != nil {
		return err
	}
	defer file.Close()

	return toml.NewEncoder(file).Encode(cfg)
}
