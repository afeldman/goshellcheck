package main

import (
	"os"

	"github.com/afeldman/goshellcheck/internal/cli"
)

func main() {
	cfg, err := cli.Parse(os.Args[1:])
	if err != nil {
		// Error already printed by flag package
		os.Exit(2)
	}

	exitCode := cli.Run(cfg)
	os.Exit(exitCode)
}
