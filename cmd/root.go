package cmd

import (
	cmd "github.com/elastic/beats/libbeat/cmd"
	"github.com/icinga/icingabeat/beater"
)

// Name of this beat
var Name = "icingabeat"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmd(Name, "", beater.New)
