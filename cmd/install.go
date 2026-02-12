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

var installCmd = &cli.Command{
	Name:    "install",
	Aliases: []string{"add", "i"},
	Usage:   "Install packages into an environment",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "env",
			Usage: "Environment name",
		},
		&cli.StringFlag{
			Name:    "requirements",
			Aliases: []string{"r"},
			Usage:   "Requirements file",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		envName := cmd.String("env")
		reqFile := cmd.String("requirements")

		// Try to get env from VIRTUAL_ENV if not specified
		if envName == "" {
			virtualEnv := os.Getenv("VIRTUAL_ENV")
			if virtualEnv != "" && !filepath.IsAbs(virtualEnv) {
				// VIRTUAL_ENV is a path, extract env name
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

		// Build uv pip install command
		args := []string{"pip", "install"}

		if reqFile != "" {
			args = append(args, "-r", reqFile)
		} else {
			args = append(args, cmd.Args().Slice()...)
		}

		if len(cmd.Args().Slice()) == 0 && reqFile == "" {
			return fmt.Errorf("no packages specified")
		}

		return uv.RunUvWithPython(python, args...)
	},
}
