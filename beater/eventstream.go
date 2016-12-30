package beater

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/icinga/icingabeat/config"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

// Eventstream type
type Eventstream struct {
	icingabeat *Icingabeat
	config     config.Config

	done   chan struct{}
	closer io.Closer
	mutex  sync.Mutex
}

// NewEventstream ...
func NewEventstream(bt *Icingabeat, cfg config.Config) *Eventstream {
	eventstream := &Eventstream{
		icingabeat: bt,
		done:       make(chan struct{}),
		config:     cfg,
	}
	return eventstream
}

// Run evenstream receiver
func (es *Eventstream) Run() error {
	types := strings.Join(es.config.Eventstream.Types, "&types=")
	logp.Info(types)
	for {

		ticker := time.NewTicker(es.config.RetryInterval)
		response, responseErr := requestURL(
			es.icingabeat,
			"POST",
			"/v1/events?queue=icingabeat&types="+types)

		if responseErr == nil {
			reader := bufio.NewReader(response.Body)
			es.mutex.Lock()
			es.closer = response.Body
			es.mutex.Unlock()

			for {
				line, err := reader.ReadBytes('\n')

				if err != nil {
					es.mutex.Lock()
					tst := es.closer == nil
					es.mutex.Unlock()

					if tst || err == io.ErrUnexpectedEOF || err == io.EOF {
						break
					}
					logp.Err("Error reading line %#v", err)
				}

				var event common.MapStr

				if err := json.Unmarshal(line, &event); err != nil {
					logp.Info("Unmarshal problem %v", err)
					es.mutex.Lock()
					tst := es.closer == nil
					es.mutex.Unlock()

					if tst || err == io.ErrUnexpectedEOF || err == io.EOF {
						break
					}
					continue
				}

				event["@timestamp"] = common.Time(time.Now())
				event["type"] = "icingabeat"
				es.icingabeat.client.PublishEvent(event)
				logp.Info("Event sent")
			}

			select {
			case <-es.done:
				return nil
			default:
			}
		} else {
			logp.Info("Error connecting to API: %v", responseErr)
		}

		select {
		case <-es.done:
			return nil
		case <-ticker.C:
		}
	}
}

// Stop eventstream receiver
func (es *Eventstream) Stop() {
	es.mutex.Lock()
	if es.closer != nil {
		es.closer.Close()
		es.closer = nil
	}
	es.mutex.Unlock()
	close(es.done)
}
