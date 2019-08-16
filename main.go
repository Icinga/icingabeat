package main

import (
	"os"

	"github.com/icinga/icingabeat/cmd"

	_ "github.com/icinga/icingabeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
