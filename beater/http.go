package beater

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/elastic/beats/libbeat/logp"
)

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
