package beater

import (
	"bufio"
	"encoding/json"
	"io"
	"net/url"
	"strconv"
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
	queue := "icingabeat"
	host := es.config.Host + ":" + strconv.Itoa(es.config.Port)
	var URL *url.URL

	URL, err := url.Parse("https://" + host)
	if err != nil {
		logp.Info("Invalid request URL")
	}

	URL.Path += "/v1/events/"

	parameters := url.Values{}
	parameters.Add("queue", queue)

	if es.config.Eventstream.Filter != "" {
		parameters.Add("filter", es.config.Eventstream.Filter)
	}

	for _, eventType := range es.config.Eventstream.Types {
		parameters.Add("types", eventType)
	}

	URL.RawQuery = parameters.Encode()

	for {

		ticker := time.NewTicker(es.config.Eventstream.RetryInterval)
		response, responseErr := requestURL(es.icingabeat, "POST", URL)

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
				documentType := strings.ToLower(event["type"].(string))
				event["type"] = "icingabeat.event." + documentType
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
