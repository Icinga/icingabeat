package beater

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/elastic/beats/v7/libbeat/logp"
)

func requestURL(bt *Icingabeat, method string, URL *url.URL) (*http.Response, error) {

	var skipSslVerify bool
	certPool := x509.NewCertPool()

	if bt.config.SSL.Verify {
		skipSslVerify = false

		for _, ca := range bt.config.SSL.CertificateAuthorities {
			cert, err := ioutil.ReadFile(ca)
			if err != nil {
				logp.Warn("Could not load certificate: %v", err)
			}
			certPool.AppendCertsFromPEM(cert)
		}
	} else {
		skipSslVerify = true
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: skipSslVerify,
		RootCAs:            certPool,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
	}

	logp.Debug("icingabeat", "Requested URL: %v", URL.String())

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
		defer response.Body.Close()
	case 404:
		err = errors.New("404 Not Found. Missing permissions may be a reason for this.")
		defer response.Body.Close()
	}

	return response, err
}
