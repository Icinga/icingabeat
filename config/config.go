// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

// Config options
type Config struct {
	Host          string        `config:"host"`
	Port          int           `config:"port"`
	User          string        `config:"user"`
	Password      string        `config:"password"`
	RetryInterval time.Duration `config:"retry_interval"`
	EventTypes    []string      `config:"event_types"`
}

// DefaultConfig values
var DefaultConfig = Config{
	RetryInterval: 1 * time.Second,
	Host:          "localhost",
	Port:          5665,
}
