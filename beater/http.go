package beater

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"

	"github.com/elastic/beats/libbeat/logp"
)

func requestURL(bt *Icingabeat, method string, URL *url.URL) (*http.Response, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
	}

	logp.Info("Requested URL: %v", URL.String())

	request, err := http.NewRequest(method, URL.String(), nil)

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
