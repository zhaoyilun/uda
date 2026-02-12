package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/uda/uda/internal/env"
	"github.com/uda/uda/internal/shell"
)

var activateCmd = &cli.Command{
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

		script, err := shell.GenerateActivateScript(name)
		if err != nil {
			return err
		}

		fmt.Print(script)
		return nil
	},
}
