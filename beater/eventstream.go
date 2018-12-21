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

	"github.com/elastic/beats/libbeat/beat"
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

// BuildEventstreamEvent ...
func BuildEventstreamEvent(e []byte) beat.Event {

	var event beat.Event
	var icingaEvent map[string]interface{}

	if err := json.Unmarshal(e, &icingaEvent); err != nil {
		logp.Warn("Error decoding event: %v", err)
	}

	event.Timestamp = time.Now()
	event.Fields = common.MapStr{}

	for key, value := range icingaEvent {
		event.Fields.Put(target_key+key, value)
	}

	logp.Debug("icingabeat", "Type: %v", icingaEvent["type"])
	switch icingaEvent["type"] {
	case "CheckResult", "StateChange", "Notification":
		checkResult := icingaEvent["check_result"].(map[string]interface{})
		event.Fields.Put(target_key+"check_result.execution_start", FloatToTimestamp(checkResult["execution_start"].(float64)))
		event.Fields.Put(target_key+"check_result.execution_end", FloatToTimestamp(checkResult["execution_end"].(float64)))
		event.Fields.Put(target_key+"check_result.schedule_start", FloatToTimestamp(checkResult["schedule_start"].(float64)))
		event.Fields.Put(target_key+"check_result.schedule_end", FloatToTimestamp(checkResult["schedule_end"].(float64)))
		event.Fields.Delete(target_key+"check_result.performance_data")

	case "AcknowledgementSet":
		event.Delete("comment")
		event.Fields.Put(target_key+"comment.text", icingaEvent["comment"])
		event.Fields.Put(target_key+"expiry", FloatToTimestamp(icingaEvent["expiry"].(float64)))

	case "CommentAdded", "CommentRemoved":
		comment := icingaEvent["comment"].(map[string]interface{})
		event.Fields.Put(target_key+"comment.entry_time", FloatToTimestamp(comment["entry_time"].(float64)))
		event.Fields.Put(target_key+"comment.expire_time", FloatToTimestamp(comment["expire_time"].(float64)))

	case "DowntimeAdded", "DowntimeRemoved", "DowntimeStarted", "DowntimeTriggered":
		downtime := icingaEvent["downtime"].(map[string]interface{})
		event.Fields.Put(target_key+"downtime.end_time", FloatToTimestamp(downtime["end_time"].(float64)))
		event.Fields.Put(target_key+"downtime.entry_time", FloatToTimestamp(downtime["entry_time"].(float64)))
		event.Fields.Put(target_key+"downtime.start_time", FloatToTimestamp(downtime["start_time"].(float64)))
		event.Fields.Put(target_key+"downtime.trigger_time", FloatToTimestamp(downtime["trigger_time"].(float64)))
	}

	event.Fields.Put("type", "icingabeat.event."+strings.ToLower(icingaEvent["type"].(string)))
	event.Fields.Put("timestamp", FloatToTimestamp(icingaEvent["timestamp"].(float64)))

	return event
}

// FloatToTimestamp ...
func FloatToTimestamp(stamp float64) time.Time {
	sec := int64(stamp)
	nsec := int64((stamp - float64(int64(stamp))) * 1e9)
	t := time.Unix(sec, nsec)

	return t
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

				es.icingabeat.client.Publish(BuildEventstreamEvent(line))
				logp.Debug("icingabeat.eventstream", "Event sent")
			}

			select {
			case <-es.done:
				return nil
			default:
			}
		} else {
			logp.Err("Error connecting to API: %v", responseErr)
		}

		select {
		case <-es.done:
			defer response.Body.Close()
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
