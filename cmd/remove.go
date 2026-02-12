package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/uda/uda/internal/env"
)

var removeCmd = &cli.Command{
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

		fmt.Printf("Removing environment %s...\n", name)
		return env.Remove(name)
	},
}
