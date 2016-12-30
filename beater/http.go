package beater

import (
	"crypto/tls"
	"errors"
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
		logp.Info("Request: %v", err)
	}

	request.Header.Add("Accept", "application/json")
	request.SetBasicAuth(bt.config.User, bt.config.Password)
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 401:
		err = errors.New("Authentication failed for user " + bt.config.User)
	}

	return response, err
}
