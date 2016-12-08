package beater

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	config config.Config
	client publisher.Client

	closer io.Closer
	mutex  sync.Mutex
}

func requestURL(icingabeat *Icingabeat, method, path string) *http.Response {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
	}

	url := fmt.Sprintf("https://%s:%v%s", icingabeat.config.Host, icingabeat.config.Port, path)

	request, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Add("Accept", "application/json")
	request.SetBasicAuth(icingabeat.config.User, icingabeat.config.Password)
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	return response
}

func apiStatus(icingabeat *Icingabeat) bool {

	response := requestURL(icingabeat, "GET", "/v1/status")

	if response.StatusCode == 200 {
		return true
	}

	log.Println("Request:", response.Request.URL)
	log.Fatalln("Error:", response.Status)
	return false
}

// New beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Icingabeat{
		config: config,
	}
	return bt, nil
}

// Run Icingabeat
func (bt *Icingabeat) Run(b *beat.Beat) error {
	logp.Info("icingabeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()

	if apiStatus(bt) {
		logp.Info("API is here!!!")
	} else {
		logp.Info("API not available")
	}

	response := requestURL(bt, "POST", "/v1/events?queue=icingabeat&types=CheckResult")
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

	return nil
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
}
