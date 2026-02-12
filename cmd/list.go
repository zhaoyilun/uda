package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

var listCmd = &cli.Command{
	Name:  "list",
	Usage: "List all Python environments",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	},
}
