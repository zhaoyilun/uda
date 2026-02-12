package uv

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/uda/uda/internal/config"
	"github.com/uda/uda/internal/mirror"
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

// Install downloads and installs uv binary with mirror support
func Install() error {
	return installWithMirror(false)
}

func installWithMirror(fallbackToOfficial bool) error {
	uvPath := config.UvPath()

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(uvPath), 0755); err != nil {
		return err
	}

	// Get mirror URL
	mirrorURL := ""
	if !fallbackToOfficial {
		mirrorURL = mirror.GetMirror()
		if mirrorURL == "" {
			// Try to find working mirror
			var err error
			mirrorURL, err = mirror.FindWorkingMirror()
			if err != nil {
				mirrorURL = "https://astral.sh" // Fallback to official
			}
		}
	}

	// Use mirror for download
	uvURL := getUvDownloadURL(mirrorURL)
	fmt.Printf("Downloading uv from %s...\n", uvURL)

	// Download the file
	resp, err := http.Get(uvURL)
	if err != nil {
		// Try fallback to official if mirror fails
		if mirrorURL != "" && mirrorURL != "https://astral.sh" {
			fmt.Println("Mirror failed, trying official source...")
			return installWithMirror(true) // Retry with official
		}
		return fmt.Errorf("failed to download uv: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Try fallback to official if mirror fails
		if mirrorURL != "" && mirrorURL != "https://astral.sh" {
			fmt.Println("Mirror failed, trying official source...")
			return installWithMirror(true) // Retry with official
		}
		return fmt.Errorf("failed to download uv: %s", resp.Status)
	}

	// Create destination file
	out, err := os.Create(uvPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// Make executable
	if err := os.Chmod(uvPath, 0755); err != nil {
		return err
	}

	fmt.Println("UV installed successfully to", uvPath)
	return nil
}

func getUvDownloadURL(mirrorURL string) string {
	arch := runtime.GOARCH
	os := runtime.GOOS
	ext := ""
	if os == "windows" {
		ext = ".exe"
	}

	// Construct download URL based on mirror
	if mirrorURL != "" && mirrorURL != "https://astral.sh" {
		// For mirrors, use the same path structure
		return fmt.Sprintf("%s/uv/releases/latest/download/uv-%s-%s%s", mirrorURL, os, arch, ext)
	}

	// Official URL
	return fmt.Sprintf("https://github.com/astral-sh/uv/releases/latest/download/uv-%s-%s%s", os, arch, ext)
}

// InstallPython installs a specific Python version
func InstallPython(version string) error {
	// First check if Python is already installed
	uv, err := FindUv()
	if err != nil {
		return err
	}

	// Run uv python install
	cmd := exec.Command(uv, "python", "install", version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// RunUvWithPython runs uv with a specific Python interpreter
func RunUvWithPython(pythonPath string, args ...string) error {
	uv, err := FindUv()
	if err != nil {
		return err
	}

	// Insert --python flag after the uv command
	fullArgs := []string{}
	foundPython := false
	for _, arg := range args {
		if arg == "--python" {
			foundPython = true
		}
		fullArgs = append(fullArgs, arg)
	}

	// If --python not already in args, add it
	if !foundPython && pythonPath != "" {
		fullArgs = append([]string{"--python", pythonPath}, fullArgs...)
	}

	cmd := exec.Command(uv, fullArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
