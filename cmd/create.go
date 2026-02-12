package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

var createCmd = &cli.Command{
	Name:  "create",
	Usage: "Create a new Python environment",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	},
}
