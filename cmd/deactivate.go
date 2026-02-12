package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/uda/uda/internal/shell"
)

var deactivateCmd = &cli.Command{
	Name:    "deactivate",
	Aliases: []string{"d"},
	Usage:   "Deactivate current environment",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		script := shell.GenerateDeactivateScript()
		fmt.Print(script)
		return nil
	},
}
