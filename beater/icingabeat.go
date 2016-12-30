package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/icinga/icingabeat/config"
)

// Icingabeat type
type Icingabeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

// New beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Icingabeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

// Run Icingabeat
func (bt *Icingabeat) Run(b *beat.Beat) error {
	var eventstream *Eventstream

	logp.Info("icingabeat is running! Hit CTRL-C to stop it.")
	bt.client = b.Publisher.Connect()

	eventstream = NewEventstream(bt, bt.config)
	logp.Info("hostname: %v", bt.config.Host)
	go eventstream.Run()

	for {
		select {
		case <-bt.done:
			return nil
		}
	}
}

// Stop Icingabeat
func (bt *Icingabeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
