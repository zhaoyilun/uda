package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

var selfCmd = &cli.Command{
	Name:  "self",
	Usage: "Manage uda itself",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	},
}
