package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/uda/uda/internal/env"
	"github.com/uda/uda/internal/uv"
)

var createCmd = &cli.Command{
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

		// Check if env already exists
		if env.Exists(name) {
			return fmt.Errorf("environment %s already exists", name)
		}

		pythonVersion := cmd.String("python")

		// Install Python if specified
		if pythonVersion != "" {
			fmt.Printf("Installing Python %s...\n", pythonVersion)
			if err := uv.InstallPython(pythonVersion); err != nil {
				return fmt.Errorf("failed to install Python: %w", err)
			}
		}

		// Create environment
		fmt.Printf("Creating environment %s...\n", name)
		return env.Create(name, pythonVersion)
	},
}
