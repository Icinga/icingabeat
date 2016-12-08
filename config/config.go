// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

// Config options
type Config struct {
	Period   time.Duration `config:"period"`
	Host     string        `config:"host"`
	Port     int           `config:"port"`
	User     string        `config:"user"`
	Password string        `config:"password"`
}

// DefaultConfig values
var DefaultConfig = Config{
	Period: 1 * time.Second,
	Host:   "localhost",
	Port:   5665,
}
