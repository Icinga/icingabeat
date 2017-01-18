package beater

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"time"

	"github.com/icinga/icingabeat/config"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

// Statuspoller type
type Statuspoller struct {
	icingabeat *Icingabeat
	config     config.Config

	done chan struct{}
}

// NewStatuspoller ...
func NewStatuspoller(bt *Icingabeat, cfg config.Config) *Statuspoller {
	statuspoller := &Statuspoller{
		icingabeat: bt,
		done:       make(chan struct{}),
		config:     cfg,
	}
	return statuspoller
}

// Run evenstream receiver
func (sp *Statuspoller) Run() error {
	host := sp.config.Host + ":" + strconv.Itoa(sp.config.Port)
	var URL *url.URL

	URL, err := url.Parse("https://" + host)
	if err != nil {
		logp.Info("Invalid request URL")
	}

	URL.Path += "/v1/status"

	for {
		ticker := time.NewTicker(sp.config.Statuspoller.Interval)
		response, responseErr := requestURL(sp.icingabeat, "GET", URL)
		var event common.MapStr

		if responseErr == nil {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				logp.Warn("Response body invalid: %v", err)
			}

			if err := json.Unmarshal(body, &event); err != nil {
				logp.Info("Unmarshal problem %v", err)

				if err == io.ErrUnexpectedEOF || err == io.EOF {
					break
				}
				continue
			}

			for _, result := range event {
				switch statustype := result.(type) {
				case []interface{}:

					for _, status := range statustype {
						statusevent := common.MapStr{
							"@timestamp": common.Time(time.Now()),
							"type":       "icingabeat.status",
						}

						for key, value := range status.(map[string]interface{}) {
							if key != "perfdata" {
								statusevent.Put(key, value)
							}
						}

						sp.icingabeat.client.PublishEvent(statusevent)
						logp.Info("Event sent")
					}
				}
			}
		} else {
			logp.Info("Error connecting to API: %v", responseErr)
		}

		select {
		case <-sp.done:
			return nil
		case <-ticker.C:
		}
	}
	return nil
}

// Stop statuspoller
func (sp *Statuspoller) Stop() {
	close(sp.done)
}
