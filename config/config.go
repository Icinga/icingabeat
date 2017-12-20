// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"time"
)

// Config options
type Config struct {
	Host         string             `config:"host"`
	Port         int                `config:"port"`
	User         string             `config:"user"`
	Password     string             `config:"password"`
	SSL          SSL                `config:"ssl"`
	Eventstream  EventstreamConfig  `config:"eventstream"`
	Statuspoller StatuspollerConfig `config:"statuspoller"`
}

// SSL options
type SSL struct {
	Verify                 bool     `config:"verify"`
	CertificateAuthorities []string `config:"certificate_authorities"`
}

// EventstreamConfig optoins
type EventstreamConfig struct {
	Types         []string      `config:"types"`
	Filter        string        `config:"filter"`
	RetryInterval time.Duration `config:"retry_interval"`
}

// StatuspollerConfig options
type StatuspollerConfig struct {
	Interval time.Duration `config:"interval"`
}

// DefaultConfig values
var DefaultConfig = Config{
	Host: "localhost",
	Port: 5665,
}
