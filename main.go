package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/icinga/icingabeat/beater"
)

func main() {
	err := beat.Run("icingabeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
