package main

import (
	"os"

	"github.com/icinga/icingabeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
