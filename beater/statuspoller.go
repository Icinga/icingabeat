package beater

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
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

// BuildStatusEvents ...
func BuildStatusEvents(body []byte) []common.MapStr {
	var statusEvents []common.MapStr
	var icingaStatus map[string]interface{}

	if err := json.Unmarshal(body, &icingaStatus); err != nil {
		logp.Warn("Unmarshal problem %v", err)
	}

	for _, result := range icingaStatus {
		for _, status := range result.([]interface{}) {

			event := common.MapStr{}
			for key, value := range status.(map[string]interface{}) {

				switch key {
				case "status":
					for _, statusvalue := range value.(map[string]interface{}) {
						switch statusvalue.(type) {
						case map[string]interface{}:
							if len(statusvalue.(map[string]interface{})) > 0 {

								event.Put(key, value)
							}

						default:
							event.Put(key, value)
						}

					}

				case "perfdata":
					for _, perfdata := range value.([]interface{}) {
						switch perfdata.(type) {
						case string:
							logp.Debug("Perfdata is a string, skipping. (%v)", perfdata.(string))

						case interface{}:
							key = "perfdata." + perfdata.(map[string]interface{})["label"].(string)
							value = perfdata
							event.Put(key, value)

						}
					}

				case "name":
					key = "type"
					value = "icingabeat.status." + strings.ToLower(value.(string))
					event.Put(key, value)

				default:
					event.Put(key, value)
				}
			}

			if statusAvailable, _ := event.HasKey("status"); statusAvailable == true {
				event.Put("@timestamp", common.Time(time.Now()))
				statusEvents = append(statusEvents, event)
			}
		}
	}

	return statusEvents
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

		if responseErr == nil {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				logp.Warn("Response body invalid: %v", err)
			}

			processedStatusEvents := BuildStatusEvents(body)
			sp.icingabeat.client.PublishEvents(processedStatusEvents)
			logp.Debug("icingabeat.statuspoller", "Events sent: %v", len(processedStatusEvents))

		} else {
			logp.Err("Error connecting to API: %v", responseErr)
		}

		select {
		case <-sp.done:
			return nil
		case <-ticker.C:
		}
	}
}

// Stop statuspoller
func (sp *Statuspoller) Stop() {
	close(sp.done)
}
