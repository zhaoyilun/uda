package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
	"github.com/uda/uda/internal/shell"
)

var initCmd = &cli.Command{
	Name:  "init",
	Usage: "Initialize shell integration",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "shell",
			Usage: "Shell type (bash, zsh, fish)",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		shellType := cmd.String("shell")
		// Also accept positional argument
		if shellType == "" {
			shellType = cmd.Args().First()
		}
		if shellType == "" {
			shellType = os.Getenv("SHELL")
			if shellType != "" {
				shellType = filepath.Base(shellType)
			}
		}

		if shellType == "" {
			shellType = "bash"
		}

		fmt.Print(shell.Init(shellType))
		return nil
	},
}
