package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

var deactivateCmd = &cli.Command{
	Name:  "deactivate",
	Usage: "Deactivate the current Python environment",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	},
}
