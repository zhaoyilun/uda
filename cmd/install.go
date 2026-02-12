package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

var installCmd = &cli.Command{
	Name:  "install",
	Usage: "Install packages into a Python environment",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	},
}
