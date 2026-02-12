package env

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/uda/uda/internal/config"
)

func List() ([]string, error) {
	entries, err := os.ReadDir(config.EnvsPath())
	if err != nil {
		return nil, err
	}

	var envs []string
	for _, entry := range entries {
		if entry.IsDir() {
			envs = append(envs, entry.Name())
		}
	}
	return envs, nil
}

func Exists(name string) bool {
	_, err := os.Stat(config.EnvPath(name))
	return err == nil
}

func Create(name string, pythonVersion string) error {
	envPath := config.EnvPath(name)
	if err := os.MkdirAll(envPath, 0755); err != nil {
		return fmt.Errorf("failed to create env directory: %w", err)
	}

	// Create virtual environment using uv
	args := []string{"venv", envPath}
	if pythonVersion != "" {
		args = append(args, "--python", pythonVersion)
	}

	// TODO: Call uv to create venv
	_ = args
	return nil
}

func Remove(name string) error {
	envPath := config.EnvPath(name)
	return os.RemoveAll(envPath)
}
