package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

var removeCmd = &cli.Command{
	Name:  "remove",
	Usage: "Remove a Python environment",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	},
}
