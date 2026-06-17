package main

import (
	"os"

	"github.com/fullstorydev/subtext-verify/cli/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
