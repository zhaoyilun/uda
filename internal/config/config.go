package config

import (
	"os"
	"path/filepath"
)

var HomeDir = filepath.Join(os.Getenv("HOME"), ".uda")

func Init() error {
	dirs := []string{
		HomeDir,
		filepath.Join(HomeDir, "envs"),
		filepath.Join(HomeDir, "cache"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func UvPath() string {
	return filepath.Join(HomeDir, "uv")
}

func EnvsPath() string {
	return filepath.Join(HomeDir, "envs")
}

func EnvPath(name string) string {
	return filepath.Join(EnvsPath(), name)
}
