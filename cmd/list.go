package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/uda/uda/internal/env"
)

var listCmd = &cli.Command{
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
