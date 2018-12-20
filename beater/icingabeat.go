package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/icinga/icingabeat/config"
)

// Icingabeat type
type Icingabeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

var target_key = "icinga."

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
	logp.Info("icingabeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	if len(bt.config.Eventstream.Types) > 0 {
		var eventstream *Eventstream
		eventstream = NewEventstream(bt, bt.config)
		go eventstream.Run()
	}

	if bt.config.Statuspoller.Interval > 0 {
		var statuspoller *Statuspoller
		statuspoller = NewStatuspoller(bt, bt.config)
		go statuspoller.Run()
	}

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
