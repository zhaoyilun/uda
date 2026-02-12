package main

import (
	"os"

	"github.com/uda/uda/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
