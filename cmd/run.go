package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v3"
	"github.com/uda/uda/internal/config"
	"github.com/uda/uda/internal/env"
	"github.com/uda/uda/internal/uv"
)

var runCmd = &cli.Command{
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

		// Try to get env from VIRTUAL_ENV if not specified
		if envName == "" {
			virtualEnv := os.Getenv("VIRTUAL_ENV")
			if virtualEnv != "" && !filepath.IsAbs(virtualEnv) {
				envName = filepath.Base(virtualEnv)
			}
		}

		if envName == "" {
			return fmt.Errorf("environment not specified. Use --env or set VIRTUAL_ENV")
		}

		if !env.Exists(envName) {
			return fmt.Errorf("environment %s does not exist", envName)
		}

		envPath := config.EnvPath(envName)
		var python string
		if runtime.GOOS == "windows" {
			python = filepath.Join(envPath, "Scripts", "python.exe")
		} else {
			python = filepath.Join(envPath, "bin", "python")
		}

		// Use uv run with the specific python
		args := []string{}
		if len(cmd.Args().Slice()) > 0 {
			args = append(args, cmd.Args().Slice()...)
		} else {
			// Default to python REPL
			args = append(args, python)
		}

		return uv.RunUvWithPython(python, append([]string{"run"}, args...)...)
	},
}
