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

package monitoring

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/backend"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/config"
)

// Server interface to handle a monitoring server
type Server interface {
	Start()
}

// NewMonitoringServer Server using the given server config
func NewMonitoringServer(config config.ServerConfig, backend backend.Client) Server {
	return &server{
		config:        config,
		backend:       backend,
		requestLogger: config.Logger().Named("frontend"),
	}
}

type server struct {
	config        config.ServerConfig
	backend       backend.Client
	requestLogger hclog.Logger
}

func (s server) Start() {
	monitoringMux := http.NewServeMux()
	monitoringMux.Handle("/metrics", promhttp.Handler())
	monitoringMux.HandleFunc("/liveness", s.newLivenessRequestHandler())
	monitoringMux.HandleFunc("/readiness", s.newReadinessRequestHandler())
	adminSrv := &http.Server{
		Addr:         fmt.Sprintf(":%s", "2112"),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		Handler:      monitoringMux,
	}
	s.config.Logger().Info("Start monitoring service", "port", "2112")
	log.Fatal(adminSrv.ListenAndServe())
}

func (s server) newLivenessRequestHandler() func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, incomingRequest *http.Request) {
		defer func() {
			if flusher, ok := responseWriter.(http.Flusher); ok {
				s.requestLogger.Trace("Flush response writer")
				flusher.Flush()
			}
			probeRequestCounter.WithLabelValues("liveness", fmt.Sprint(http.StatusOK)).Inc()
		}()
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte("liveness OK"))
	}
}

func (s server) newReadinessRequestHandler() func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, incomingRequest *http.Request) {
		statusCode := http.StatusInternalServerError
		defer func() {
			if flusher, ok := responseWriter.(http.Flusher); ok {
				s.requestLogger.Trace("Flush response writer")
				flusher.Flush()
			}
			probeRequestCounter.WithLabelValues("readiness", fmt.Sprint(statusCode)).Inc()
		}()
		backendRequest, err := retryablehttp.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", s.config.BackendURL(), "/-/readiness"), []byte{})
		if err != nil {
			http.Error(responseWriter, err.Error(), statusCode)
			return
		}
		backendResponse, err := s.backend.Send(backendRequest)
		if err != nil {
			http.Error(responseWriter, err.Error(), statusCode)
			return
		}
		responseBody, err := readBody(backendResponse.Body)
		if err != nil {
			http.Error(responseWriter, err.Error(), statusCode)
			return
		}
		statusCode = backendResponse.StatusCode
		responseWriter.WriteHeader(statusCode)
		responseWriter.Write(responseBody)
	}
}

func readBody(body io.ReadCloser) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, body); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
