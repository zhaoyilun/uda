package uv

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/uda/uda/internal/config"
	"github.com/uda/uda/internal/mirror"
)

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

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

	// Check if this is a tar.gz file
	isTarGz := strings.HasSuffix(uvURL, ".tar.gz")

	if isTarGz {
		// Download tar.gz and extract
		tmpDir, err := os.MkdirTemp("", "uv-install")
		if err != nil {
			return fmt.Errorf("failed to create temp dir: %w", err)
		}
		defer os.RemoveAll(tmpDir)

		tarPath := filepath.Join(tmpDir, "uv.tar.gz")
		out, err := os.Create(tarPath)
		if err != nil {
			return fmt.Errorf("failed to create tar file: %w", err)
		}

		resp, err := http.Get(uvURL)
		if err != nil {
			out.Close()
			if mirrorURL != "" && mirrorURL != "https://astral.sh" {
				fmt.Println("Mirror failed, trying official source...")
				return installWithMirror(true)
			}
			return fmt.Errorf("failed to download uv: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			out.Close()
			if mirrorURL != "" && mirrorURL != "https://astral.sh" {
				fmt.Println("Mirror failed, trying official source...")
				return installWithMirror(true)
			}
			return fmt.Errorf("failed to download uv: %s", resp.Status)
		}

		_, err = io.Copy(out, resp.Body)
		out.Close()
		if err != nil {
			return fmt.Errorf("failed to save tar file: %w", err)
		}

		// Extract tar.gz
		cmd := exec.Command("tar", "-xzf", tarPath, "-C", tmpDir)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to extract tar file: %w", err)
		}

		// Find and move the uv binary
		entries, err := os.ReadDir(tmpDir)
		if err != nil {
			return fmt.Errorf("failed to read temp dir: %w", err)
		}

		found := false
		for _, entry := range entries {
			if entry.IsDir() {
				subDir := filepath.Join(tmpDir, entry.Name())
				uvBinary := filepath.Join(subDir, "uv")
				if _, err := os.Stat(uvBinary); err == nil {
					if err := CopyFile(uvBinary, uvPath); err != nil {
						return fmt.Errorf("failed to copy uv binary: %w", err)
					}
					found = true
					break
				}
			}
		}

		if !found {
			return fmt.Errorf("uv binary not found in tarball")
		}
	} else {
		// Download single binary
		resp, err := http.Get(uvURL)
		if err != nil {
			if mirrorURL != "" && mirrorURL != "https://astral.sh" {
				fmt.Println("Mirror failed, trying official source...")
				return installWithMirror(true)
			}
			return fmt.Errorf("failed to download uv: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			if mirrorURL != "" && mirrorURL != "https://astral.sh" {
				fmt.Println("Mirror failed, trying official source...")
				return installWithMirror(true)
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
	}

	// Make executable
	if err := os.Chmod(uvPath, 0755); err != nil {
		return err
	}

	fmt.Println("UV installed successfully to", uvPath)
	return nil
}

// mapGoArchToUVArch maps Go architecture names to UV download architecture names
func mapGoArchToUVArch(goArch string) string {
	switch goArch {
	case "amd64":
		return "x86_64"
	case "arm64":
		return "aarch64"
	default:
		return goArch
	}
}

func getUvDownloadURL(mirrorURL string) string {
	arch := mapGoArchToUVArch(runtime.GOARCH)
	os := runtime.GOOS

	// UV now uses tar.gz format for Linux
	if os == "linux" {
		return fmt.Sprintf("https://github.com/astral-sh/uv/releases/latest/download/uv-%s-unknown-linux-gnu.tar.gz", arch)
	}

	// Windows uses .exe
	ext := ""
	if os == "windows" {
		ext = ".exe"
	}

	// Construct download URL based on mirror
	if mirrorURL != "" && mirrorURL != "https://astral.sh" {
		// For mirrors, use the same path structure
		return fmt.Sprintf("%s/uv/releases/latest/download/uv-%s-%s%s", mirrorURL, os, arch, ext)
	}

	// Official URL for other platforms
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
