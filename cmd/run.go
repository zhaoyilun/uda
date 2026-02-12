package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "Run a command in a Python environment",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	},
}
