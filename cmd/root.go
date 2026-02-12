package cmd

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

var version = "0.1.0"

func Execute() error {
	app := &cli.Command{
		Name:    "uda",
		Usage:   "Python environment manager combining Conda and UV",
		Version: version,
		Commands: []*cli.Command{
			createCmd,
			listCmd,
			removeCmd,
			activateCmd,
			deactivateCmd,
			installCmd,
			runCmd,
			selfCmd,
			initCmd,
		},
	}

	return app.Run(context.Background(), os.Args)
}
