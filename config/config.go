// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

// Config options
type Config struct {
	Host          string             `config:"host"`
	Port          int                `config:"port"`
	User          string             `config:"user"`
	Password      string             `config:"password"`
	RetryInterval time.Duration      `config:"retry_interval"`
	SkipSSLVerify bool               `config:"skip_ssl_verify"`
	Eventstream   EventstreamConfig  `config:"eventstream"`
	Statuspoller  StatuspollerConfig `config:"statuspoller"`
}

// EventstreamConfig optoins
type EventstreamConfig struct {
	Types  []string `config:"types"`
	Filter string   `config:"filter"`
}

// StatuspollerConfig options
type StatuspollerConfig struct {
	Interval time.Duration `config:"interval"`
}

// DefaultConfig values
var DefaultConfig = Config{
	RetryInterval: 1 * time.Second,
	Host:          "localhost",
	Port:          5665,
}
