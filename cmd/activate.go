package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

var activateCmd = &cli.Command{
	Name:  "activate",
	Usage: "Activate a Python environment",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	},
}
