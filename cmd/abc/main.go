package main

import (
	"fmt"
	"os"

	"github.com/abc-metrics/abc/cmd/abc/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
