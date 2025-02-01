// Copyright 2025 The Terraform SOPS backend Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/config"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/transform"
)

type Server interface {
	Start()
}

func New(config config.ServerConfig) Server {
	backend := retryablehttp.NewClient()
	backend.RetryMax = 0
	backend.RetryWaitMin = time.Duration(20) * time.Second
	backend.RetryWaitMax = time.Duration(60) * time.Second
	backend.Logger = config.Logger().Named("backend")
	return &server{
		config:        config,
		backend:       backend,
		requestLogger: config.Logger().Named("frontend"),
	}
}

const (
	methodGet    = "GET"
	methodPost   = "POST"
	methodLock   = "LOCK"
	methodUnlock = "UNLOCK"
)

var (
	ignoredRequestHeaders  = ignoredHeaders{}
	ignoredResponseHeaders = ignoredHeaders{
		"Content-Length": 0,
	}
	supportedRequestMethods = supportedMethods{
		methodGet:    0,
		methodPost:   0,
		methodLock:   0,
		methodUnlock: 0,
	}
)

type ignoredHeaders map[string]int
type supportedMethods map[string]int

type server struct {
	config        config.ServerConfig
	backend       *retryablehttp.Client
	requestLogger hclog.Logger
}

func (s server) Start() {
	http.HandleFunc("/", s.newRequestHandler())
	s.config.Logger().Info("Start service", "port", s.config.ServerPort(), "vault_addr", s.config.VaultAddr(), "has_private_age_key", len(s.config.AgePrivateKey()) > 0)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.config.ServerPort()), nil))
}

func (s server) newRequestHandler() func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, incomingRequest *http.Request) {

		s.requestLogger.Debug("incoming request", "method", incomingRequest.Method, "uri", incomingRequest.URL.Path)

		if !s.isSupportedRequestMethod(incomingRequest.Method) {
			s.requestLogger.Warn("Method Not Allowed", "method", incomingRequest.Method)
			http.Error(responseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		backendRequest, err := s.buildBackendRequest(incomingRequest)
		if err != nil {
			s.requestLogger.Warn("Can not build backend request", "error", err)
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		backendResponse, err := s.backend.Do(backendRequest)
		if err != nil {
			s.requestLogger.Warn("Can not perform backend request", "error", err)
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		s.writeResponse(responseWriter, backendResponse, incomingRequest.Method)
	}
}

func (s server) buildBackendRequest(incomingRequest *http.Request) (*retryablehttp.Request, error) {
	body, err := s.readBody(incomingRequest.Body)
	if err != nil {
		return nil, err
	}
	method := incomingRequest.Method
	if method == methodLock {
		method = s.config.BackendLockMethod()
	} else if method == methodUnlock {
		method = s.config.BackendUnlockMethod()
	} else if method == methodPost && len(body) > 0 {
		if err := transform.TransformToSops(s.config, body, func(result []byte) { body = result }); err != nil {
			return nil, err
		}
	}
	backendRequest, err := retryablehttp.NewRequest(method, fmt.Sprintf("%s%s", s.config.BackendURL(), incomingRequest.URL.Path), body)
	if err != nil {
		return nil, err
	}
	s.copyHeader(incomingRequest.Header, backendRequest.Header, ignoredRequestHeaders)
	backendRequest.URL.RawQuery = incomingRequest.URL.Query().Encode()
	return backendRequest, nil
}

func (s server) writeResponse(responseWriter http.ResponseWriter, backendResponse *http.Response, requestMethod string) {
	defer func() {
		if flusher, ok := responseWriter.(http.Flusher); ok {
			s.requestLogger.Trace("Flush response writer")
			flusher.Flush()
		}
	}()
	s.copyHeader(backendResponse.Header, responseWriter.Header(), ignoredResponseHeaders)
	responseBody, err := s.readBody(backendResponse.Body)
	if err != nil {
		s.requestLogger.Warn("Can not read backend response body", "error", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	if requestMethod == methodGet && len(responseBody) > 0 {
		s.requestLogger.Trace("Decrypt response body with", "length", len(responseBody))
		if err := transform.TransformFromSops(s.config, responseBody, func(result []byte) error { responseBody = result; return nil }); err != nil {
			s.requestLogger.Warn("Can not decrypt body. Leave body unchanged", "error", err)
		}
		s.requestLogger.Trace("Decrypted response body with", "length", len(responseBody))
	}
	if backendResponse.StatusCode/100 != 2 {
		s.requestLogger.Warn("Unexpected backendResponse", "status-code", backendResponse.StatusCode)
	}
	backendResponse.Header.Set("Content-Length", fmt.Sprint(len(responseBody)))
	responseWriter.WriteHeader(backendResponse.StatusCode)
	responseWriter.Write(responseBody)
}

func (s server) readBody(body io.ReadCloser) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, body); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s server) copyHeader(from, to http.Header, ignoredHeaders ignoredHeaders) {
	for k, vs := range from {
		if s.isIgnoredHeader(k, ignoredHeaders) {
			continue
		}
		for _, v := range vs {
			to.Add(k, v)
		}
	}
}

func (s server) isIgnoredHeader(header string, ignoredHeaders ignoredHeaders) bool {
	_, ok := ignoredHeaders[header]
	return ok
}

func (s server) isSupportedRequestMethod(requestMethod string) bool {
	_, ok := supportedRequestMethods[requestMethod]
	return ok
}
