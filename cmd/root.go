package cmd

import (
	"github.com/icinga/icingabeat/beater"

	cmd "github.com/elastic/beats/libbeat/cmd"

	// Register the includes.
	_ "github.com/icinga/icingabeat/include"
)

// Name of this beat
var Name = "icingabeat"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmd(Name, "", beater.New)
