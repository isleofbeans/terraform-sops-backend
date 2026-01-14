// Copyright 2025-2026 The Terraform SOPS backend Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backend

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/config"
)

// Client is sending requests
type Client interface {
	Send(req *retryablehttp.Request) (*http.Response, error)
}

// New creates a Client
func New(config config.ServerConfig) (client Client, err error) {

	var (
		logger = config.Logger().Named("backend")
		cert   tls.Certificate
	)

	httpClient := retryablehttp.NewClient()

	if len(config.BackendMTLSCert()) > 0 || len(config.BackendMTLSKey()) > 0 {
		cert, err = tls.X509KeyPair(config.BackendMTLSCert(), config.BackendMTLSKey())
		if err != nil {
			logger.Error("error loading mTLS certificate and key from data", "err", err, "cert-data-len", len(config.BackendMTLSCert()), "key-data-len", len(config.BackendMTLSKey()))
			return
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		transport := cleanhttp.DefaultTransport()
		transport.TLSClientConfig = tlsConfig
		httpClient.HTTPClient.Transport = transport
	}

	httpClient.RetryMax = 0
	httpClient.RetryWaitMin = time.Duration(20) * time.Second
	httpClient.RetryWaitMax = time.Duration(60) * time.Second
	httpClient.Logger = logger
	return retryableHTTPClient{
		client: httpClient,
	}, nil
}

type retryableHTTPClient struct {
	client *retryablehttp.Client
}

func (c retryableHTTPClient) Send(req *retryablehttp.Request) (*http.Response, error) {
	timer := prometheus.NewTimer(requestDuration.WithLabelValues(req.Method, req.URL.Path))
	defer timer.ObserveDuration()
	resp, err := c.client.Do(req)
	responseStatusCounter.WithLabelValues(fmt.Sprintf("%vxx", resp.StatusCode/100), req.URL.Path).Inc()
	return resp, err
}
