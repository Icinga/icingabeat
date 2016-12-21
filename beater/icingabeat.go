package beater

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

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

	closer io.Closer
	mutex  sync.Mutex
}

func requestURL(bt *Icingabeat, method, path string) (*http.Response, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
	}

	url := fmt.Sprintf("https://%s:%v%s", bt.config.Host, bt.config.Port, path)

	request, err := http.NewRequest(method, url, nil)

	if err != nil {
		logp.Info("Request:", err)
	}

	request.Header.Add("Accept", "application/json")
	request.SetBasicAuth(bt.config.User, bt.config.Password)
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	return response, err
}

func connectionTest(icingabeat *Icingabeat) (bool, error) {
	response, err := requestURL(icingabeat, "GET", "/v1")

	if err != nil {
		return false, err
	}

	logp.Debug("Connection test request URL:", response.Request.URL.String())
	logp.Debug("Connection test status:", response.Status)

	return true, nil
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
	logp.Info("icingabeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	apiAvailable, connectionerr := connectionTest(bt)

	for {
		ticker := time.NewTicker(2 * time.Second)

		if apiAvailable {
			logp.Info("API seems available")

			response, _ := requestURL(bt, "POST", "/v1/events?queue=icingabeat&types=CheckResult")
			reader := bufio.NewReader(response.Body)
			bt.mutex.Lock()
			bt.closer = response.Body
			bt.mutex.Unlock()

			for {
				line, err := reader.ReadBytes('\n')
				if err != nil {
					bt.mutex.Lock()
					tst := bt.closer == nil
					bt.mutex.Unlock()

					if tst {
						break
					}
					logp.Err("Error reading line %#v", err)
				}

				var event common.MapStr

				if err := json.Unmarshal(line, &event); err != nil {
					logp.Info("Unmarshal problem %v", err)
					bt.mutex.Lock()
					tst := bt.closer == nil
					bt.mutex.Unlock()

					if tst {
						break
					}
					continue
				}
				event["@timestamp"] = common.Time(time.Now())
				event["type"] = b.Name
				bt.client.PublishEvent(event)
				logp.Info("Event sent")
			}
		} else {
			logp.Info("Cannot connect to API", connectionerr)
		}

		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
	}
	//return nil
}

// Stop Icingabeat
func (bt *Icingabeat) Stop() {
	bt.client.Close()
	bt.mutex.Lock()
	if bt.closer != nil {
		bt.closer.Close()
		bt.closer = nil
	}
	bt.mutex.Unlock()
	close(bt.done)
	logp.Info("CTRL+C hit!!!")
}
