package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

var initCmd = &cli.Command{
	Name:  "init",
	Usage: "Initialize a new Python project",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	},
}
