package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/uda/uda/internal/uv"
)

var selfCmd = &cli.Command{
	Name:  "self",
	Usage: "Self management commands",
	Commands: []*cli.Command{
		{
			Name:   "install",
			Usage:  "Install or update uv",
			Action: selfInstall,
		},
	},
}

var selfInstall = func(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("Installing uv...")
	return uv.Install()
}
