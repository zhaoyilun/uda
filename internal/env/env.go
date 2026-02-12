package env

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/uda/uda/internal/config"
	"github.com/uda/uda/internal/uv"
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

	// Get uv binary
	uvPath, err := uv.FindUv()
	if err != nil {
		return err
	}

	// Create virtual environment using uv
	args := []string{"venv", envPath}
	if pythonVersion != "" {
		args = append(args, "--python", pythonVersion)
	}

	cmd := exec.Command(uvPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create venv: %w", err)
	}

	fmt.Printf("Environment %s created successfully!\n", name)
	return nil
}

func Remove(name string) error {
	envPath := config.EnvPath(name)
	return os.RemoveAll(envPath)
}
