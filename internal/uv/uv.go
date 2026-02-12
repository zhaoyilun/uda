package uv

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/uda/uda/internal/config"
)

func FindUv() (string, error) {
	// Check local uv first
	localUv := config.UvPath()
	if _, err := os.Stat(localUv); err == nil {
		return localUv, nil
	}

	// Check system uv
	systemUv := "uv"
	if _, err := exec.LookPath(systemUv); err == nil {
		return systemUv, nil
	}

	return "", fmt.Errorf("uv not found, run 'uda self install' to install")
}

func RunUv(args ...string) error {
	uv, err := FindUv()
	if err != nil {
		return err
	}

	cmd := exec.Command(uv, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func GetPythonPath(envName string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(config.EnvPath(envName), "Scripts", "python.exe")
	}
	return filepath.Join(config.EnvPath(envName), "bin", "python")
}
