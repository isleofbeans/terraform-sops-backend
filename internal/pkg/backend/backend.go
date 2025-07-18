// Copyright 2025 The Terraform SOPS backend Authors
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
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/prometheus/client_golang/prometheus"
)

// Client is sending requests
type Client interface {
	Send(req *retryablehttp.Request) (*http.Response, error)
}

// New creates a Client
func New(logger hclog.Logger) Client {
	client := retryablehttp.NewClient()
	client.RetryMax = 0
	client.RetryWaitMin = time.Duration(20) * time.Second
	client.RetryWaitMax = time.Duration(60) * time.Second
	client.Logger = logger
	return retryableHTTPClient{
		client: client,
	}
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
