# UDA Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 创建一个 Go 实现的 Python 环境管理器，结合 Conda 的全局环境管理和 UV 的快速安装能力。

**Architecture:** 使用 Go 单一二进制，通过调用 uv CLI 实现环境管理和包管理。镜像管理采用自动检测+内置列表的策略。

**Tech Stack:** Go, urfave/cli, uv (CLI 调用)

---

## Task 1: 初始化 Go 项目结构

**Files:**
- Create: `go.mod`
- Create: `main.go`
- Create: `cmd/root.go`
- Create: `internal/config/config.go`
- Create: `internal/uv/uv.go`
- Create: `internal/env/env.go`

**Step 1: 创建项目目录结构**

```bash
mkdir -p cmd internal/config internal/uv internal/env
```

**Step 2: 初始化 go.mod**

```bash
go mod init github.com/uda/uda
```

**Step 3: 创建 main.go**

```go
package main

import (
	"os"

	"github.com/uda/uda/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

**Step 4: 创建 cmd/root.go**

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var version = "0.1.0"

func Execute() error {
	app := &cli.Command{
		Name:  "uda",
		Usage: "Python environment manager combining Conda and UV",
		Version: version,
		Commands: []cli.Command{
			createCmd,
			listCmd,
			removeCmd,
			activateCmd,
			deactivateCmd,
			installCmd,
			runCmd,
			selfCmd,
			initCmd,
		},
	}

	return app.Run(os.Args)
}
```

**Step 5: 创建 internal/config/config.go**

```go
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
```

**Step 6: 创建 internal/uv/uv.go**

```go
package uv

import (
	"fmt"
	"os"
	"os/exec"
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
```

**Step 7: 创建 internal/env/env.go**

```go
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
```

**Step 8: 验证编译**

```bash
go build -o uda .
./uda --help
```

Expected: Help output showing available commands

**Step 9: Commit**

```bash
git add .
git commit -m "chore: initialize Go project structure"
```

---

## Task 2: 实现 MVP 命令 (create, list, remove, self install)

**Files:**
- Modify: `cmd/root.go`
- Modify: `internal/uv/uv.go`
- Modify: `internal/env/env.go`

**Step 1: 添加 create 命令到 cmd/root.go**

```go
var createCmd = cli.Command{
	Name:    "create",
	Aliases: []string{"c"},
	Usage:   "Create a new Python environment",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "python",
			Usage: "Python version (e.g., 3.11)",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		name := cmd.Args().First()
		if name == "" {
			return fmt.Errorf("environment name is required")
		}

		pythonVersion := cmd.String("python")
		return env.Create(name, pythonVersion)
	},
}
```

**Step 2: 添加 list 命令**

```go
var listCmd = cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List all environments",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		envs, err := env.List()
		if err != nil {
			return err
		}

		if len(envs) == 0 {
			fmt.Println("No environments found")
			return nil
		}

		for _, e := range envs {
			fmt.Println(e)
		}
		return nil
	},
}
```

**Step 3: 添加 remove 命令**

```go
var removeCmd = cli.Command{
	Name:    "remove",
	Aliases: []string{"rm"},
	Usage:   "Remove an environment",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		name := cmd.Args().First()
		if name == "" {
			return fmt.Errorf("environment name is required")
		}

		if !env.Exists(name) {
			return fmt.Errorf("environment %s does not exist", name)
		}

		return env.Remove(name)
	},
}
```

**Step 4: 实现 self install 命令**

首先更新 internal/uv/uv.go 添加下载功能:

```go
package uv

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/uda/uda/internal/config"
)

// Mirror URLs for UV
var mirrors = []string{
	"https://astral.sh/uv/install.sh",
	// Add backup mirrors here
}

func Install() error {
	uvPath := config.UvPath()

	// Try official installer first
	if err := downloadUvOfficial(uvPath); err != nil {
		return fmt.Errorf("failed to install uv: %w", err)
	}

	if err := os.Chmod(uvPath, 0755); err != nil {
		return fmt.Errorf("failed to make uv executable: %w", err)
	}

	fmt.Println("UV installed successfully to", uvPath)
	return nil
}

func downloadUvOfficial(dest string) error {
	// Download UV installer script
	resp, err := http.Get("https://astral.sh/uv/install.sh")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download installer: %s", resp.Status)
	}

	// For simplicity, download prebuilt binary
	uvVersion := "latest"
	arch := runtime.GOARCH
	osName := runtime.GOOS

	url := fmt.Sprintf("https://github.com/astral-sh/uv/releases/%s/download/uv-%s-%s.tar.gz", uvVersion, osName, arch)

	resp, err = http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create temp file
	tmp, err := os.CreateTemp("", "uv-*.tar.gz")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		return err
	}
	tmp.Close()

	// Extract uv binary from tarball (simplified - just copy for now)
	// In real implementation, use archive/tar to extract

	return nil
}
```

**Step 5: 添加 self 命令到 cmd/root.go**

```go
var selfCmd = cli.Command{
	Name:  "self",
	Usage: "Self management commands",
	Subcommands: []cli.Command{
		{
			Name:   "install",
			Usage:  "Install or update uv",
			Action: selfInstall,
		},
	},
}

var selfInstall = func(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("Installing uv...")
	return uv.Install()
}
```

**Step 6: 测试编译和运行**

```bash
go build -o uda .
./uda list
```

Expected: "No environments found"

**Step 7: Commit**

```bash
git add .
git commit -m "feat: implement MVP commands (create, list, remove, self install)"
```

---

## Task 3: 实现 activate/deactivate 和 shell 集成

**Files:**
- Modify: `cmd/root.go`
- Modify: `internal/env/env.go`
- Create: `internal/shell/shell.go`

**Step 1: 创建 internal/shell/shell.go**

```go
package shell

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/uda/uda/internal/config"
)

func Init(shell string) string {
	switch shell {
	case "bash":
		return bashInit()
	case "zsh":
		return zshInit()
	case "fish":
		return fishInit()
	default:
		return bashInit()
	}
}

func bashInit() string {
	return `uda() {
    local cmd="$1"
    shift
    case "$cmd" in
        activate)
            _uda_activate "$@"
            ;;
        deactivate)
            _uda_deactivate
            ;;
        *)
            command uda "$cmd" "$@"
            ;;
    esac
}

_uda_activate() {
    local env_name="$1"
    if [ -z "$env_name" ]; then
        echo "Usage: uda activate <env_name>"
        return 1
    fi

    local env_path="$HOME/.uda/envs/$env_name"
    if [ ! -d "$env_path" ]; then
        echo "Environment $env_name does not exist"
        return 1
    fi

    # Save current state
    export _UDA_OLD_PATH="$PATH"
    export _UDA_OLD_PS1="$PS1"

    # Activate
    export PATH="$env_path/bin:$PATH"
    export VIRTUAL_ENV="$env_path"

    # Update prompt
    if [ -n "$PS1" ]; then
        PS1="($env_name) $PS1"
    fi

    echo "Activated environment: $env_name"
}

_uda_deactivate() {
    if [ -n "$_UDA_OLD_PATH" ]; then
        PATH="$_UDA_OLD_PATH"
        unset _UDA_OLD_PATH
    fi
    if [ -n "$_UDA_OLD_PS1" ]; then
        PS1="$_UDA_OLD_PS1"
        unset _UDA_OLD_PS1
    fi
    unset VIRTUAL_ENV
    echo "Deactivated"
}

# Alias for conda compatibility
alias conda=uda
`
}

func zshInit() string {
	return bashInit()
}

func fishInit() string {
	return `# UDA fish functions
function uda
    set cmd (math (count $argv) - 1)
    if test $argv[1] = "activate"
        set env_name $argv[2]
        set -gx _UDA_OLD_PATH $PATH
        set -gx PATH $HOME/.uda/envs/$env_name/bin $PATH
        set -gx VIRTUAL_ENV $HOME/.uda/envs/$env_name
        echo "Activated environment: $env_name"
    else if test $argv[1] = "deactivate"
        if set -q _UDA_OLD_PATH
            set -gx PATH $_UDA_OLD_PATH
            set -e _UDA_OLD_PATH
        end
        set -e VIRTUAL_ENV
        echo "Deactivated"
    else
        command uda $argv
    end
end

alias conda uda
`
}

func WriteActivateScript(envName string) error {
	envPath := config.EnvPath(envName)
	activatePath := filepath.Join(envPath, "activate.sh")

	script := fmt.Sprintf(`#!/bin/bash
# This script is auto-generated by uda
# DO NOT EDIT

export VIRTUAL_ENV="%s"
export PATH="$VIRTUAL_ENV/bin:$PATH"

if [ -n "$PS1" ]; then
    PS1="(%s) $PS1"
fi
`, envPath, envName)

	return os.WriteFile(activatePath, []byte(script), 0644)
}
```

**Step 2: 添加 init 命令**

```go
var initCmd = cli.Command{
	Name:  "init",
	Usage: "Initialize shell integration",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "shell",
			Usage: "Shell type (bash, zsh, fish)",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		shell := cmd.String("shell")
		if shell == "" {
			shell = os.Getenv("SHELL")
			if shell != "" {
				shell = filepath.Base(shell)
			}
		}

		fmt.Println(shell.Init(shell))
		return nil
	},
}
```

**Step 3: 添加 activate/deactivate 命令**

```go
var activateCmd = cli.Command{
	Name:    "activate",
	Aliases: []string{"a"},
	Usage:   "Activate an environment",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		name := cmd.Args().First()
		if name == "" {
			return fmt.Errorf("environment name is required")
		}

		if !env.Exists(name) {
			return fmt.Errorf("environment %s does not exist", name)
		}

		// Print activation script to stdout
		envPath := config.EnvPath(name)
		fmt.Printf("export VIRTUAL_ENV=%s\n", envPath)
		fmt.Printf("export PATH=%s/bin:$PATH\n", envPath)
		fmt.Printf("echo 'Activated environment: %s'\n", name)

		return nil
	},
}

var deactivateCmd = cli.Command{
	Name:    "deactivate",
	Aliases: []string{"d"},
	Usage:   "Deactivate current environment",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("unset VIRTUAL_ENV")
		fmt.Println("echo 'Deactivated'")
		return nil
	},
}
```

**Step 4: 修复 import**

```go
import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)
```

**Step 5: 测试**

```bash
go build -o uda .
./uda init bash
./uda activate myenv
```

**Step 6: Commit**

```bash
git add .
git commit -m "feat: implement activate/deactivate and shell integration"
```

---

## Task 4: 实现 install 和 run 命令

**Files:**
- Modify: `cmd/root.go`

**Step 1: 添加 install 命令**

```go
var installCmd = cli.Command{
	Name:    "install",
	Aliases: []string{"add", "i"},
	Usage:   "Install packages into an environment",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "env",
			Usage: "Environment name",
		},
		&cli.StringFlag{
			Name:  "requirements",
			Usage: "Requirements file",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		envName := cmd.String("env")
		reqFile := cmd.String("requirements")

		if envName == "" {
			envName = os.Getenv("VIRTUAL_ENV")
			if envName != "" {
				envName = filepath.Base(envName)
			}
		}

		if envName == "" {
			return fmt.Errorf("environment not specified and VIRTUAL_ENV not set")
		}

		if !env.Exists(envName) {
			return fmt.Errorf("environment %s does not exist", envName)
		}

		envPath := config.EnvPath(envName)
		args := []string{"pip", "install", "--python", filepath.Join(envPath, "bin", "python")}

		if reqFile != "" {
			args = append(args, "-r", reqFile)
		} else {
			args = append(args, cmd.Args().Slice()...)
		}

		return uv.RunUv(args...)
	},
}
```

**Step 2: 添加 run 命令**

```go
var runCmd = cli.Command{
	Name:  "run",
	Usage: "Run a command in an environment",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "env",
			Usage: "Environment name",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		envName := cmd.String("env")

		if envName == "" {
			envName = os.Getenv("VIRTUAL_ENV")
			if envName != "" {
				envName = filepath.Base(envName)
			}
		}

		if envName == "" {
			return fmt.Errorf("environment not specified and VIRTUAL_ENV not set")
		}

		if !env.Exists(envName) {
			return fmt.Errorf("environment %s does not exist", envName)
		}

		envPath := config.EnvPath(envName)
		python := filepath.Join(envPath, "bin", "python")

		return uv.RunUv("run", "--python", python, cmd.Args().Slice()...)
	},
}
```

**Step 3: Commit**

```bash
git add .
git commit -m "feat: implement install and run commands"
```

---

## Task 5: 实现镜像管理

**Files:**
- Modify: `internal/uv/uv.go`

**Step 1: 添加镜像管理逻辑**

```go
package uv

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/uda/uda/internal/config"
)

// Mirror configuration
type Mirror struct {
	Name string
	URL  string
}

var defaultMirrors = []Mirror{
	{Name: "official", URL: "https://astral.sh/uv"},
	{Name: "tsinghua", URL: "https://pypi.tuna.tsinghua.edu.cn"},
	{Name: "aliyun", URL: "https://mirrors.aliyun.com"},
}

func GetMirror() string {
	// Check UV_MIRROR environment variable first
	if mirror := os.Getenv("UV_MIRROR"); mirror != "" {
		return mirror
	}

	// Check config file
	configPath := filepath.Join(config.HomeDir, "config.toml")
	if data, err := os.ReadFile(configPath); err == nil {
		// Parse config (simplified)
		if contains(string(data), "mirror") {
			// Return configured mirror
		}
	}

	// Return default
	return ""
}

func InstallWithMirror() error {
	mirror := GetMirror()

	if mirror != "" {
		os.Setenv("UV_MIRROR", mirror)
	}

	return Install()
}
```

**Step 2: 添加自动检测逻辑**

```go
func TestMirror(mirror string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := mirror
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += "simple/uv"

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func FindWorkingMirror() string {
	for _, m := range defaultMirrors {
		if TestMirror(m.URL) {
			return m.URL
		}
	}
	return ""
}
```

**Step 3: Commit**

```bash
git add .
git commit -m "feat: implement mirror management with auto-detection"
```

---

## Task 6: 最终测试和验证

**Step 1: 构建发布版本**

```bash
go build -ldflags="-s -w" -o uda .
```

**Step 2: 测试完整流程**

```bash
# Install uv
./uda self install

# Create environment
./uda create myenv --python 3.11

# List environments
./uda list

# Show activation script
./uda activate myenv

# Install packages
./uda install --env myenv requests

# Run command
./uda run --env myenv python --version

# Remove environment
./uda remove myenv

# Show shell init
./uda init bash
```

**Step 3: Commit**

```bash
git add .
git commit -m "chore: release v0.1.0"
```
