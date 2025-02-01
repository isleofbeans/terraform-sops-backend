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

package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/backend"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/config"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/transformer"
)

// Server interface to handle a terraform SOPS backend server
type Server interface {
	Start()
}

// New Server using the given server config
func New(config config.ServerConfig, backend backend.Client, transformer transformer.SOPSTransformer) Server {
	return &server{
		config:        config,
		backend:       backend,
		transformer:   transformer,
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
	backend       backend.Client
	transformer   transformer.SOPSTransformer
	requestLogger hclog.Logger
}

func (s server) Start() {
	http.HandleFunc("/", s.newRequestHandler())
	s.config.Logger().Trace("Used configuration", "config", s.config.String())
	s.config.Logger().Info("Start service", "port", s.config.ServerPort(), "vault_addr", s.config.VaultAddr(), "has_private_age_key", len(s.config.AgePrivateKey()) > 0)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.config.ServerPort()), nil))
}

func (s server) newRequestHandler() func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, incomingRequest *http.Request) {

		if s.requestLogger.IsDebug() {
			s.requestLogger.Debug("incoming request", "method", incomingRequest.Method, "uri", buildIncomingURI(incomingRequest.URL))
		}

		if !isSupportedRequestMethod(incomingRequest.Method) {
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

		backendResponse, err := s.backend.Send(backendRequest)
		if err != nil {
			s.requestLogger.Warn("Can not perform backend request", "error", err)
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		s.writeResponse(responseWriter, backendResponse, incomingRequest.Method)
	}
}

func (s server) buildBackendRequest(incomingRequest *http.Request) (*retryablehttp.Request, error) {
	body, err := readBody(incomingRequest.Body)
	if err != nil {
		return nil, err
	}
	method := incomingRequest.Method
	if method == methodLock {
		method = s.config.BackendLockMethod()
	} else if method == methodUnlock {
		method = s.config.BackendUnlockMethod()
	} else if method == methodPost && len(body) > 0 {
		if err := s.transformer.ToSops(s.config, body, func(result []byte) { body = result }); err != nil {
			return nil, err
		}
	}
	backendRequest, err := retryablehttp.NewRequest(method, fmt.Sprintf("%s%s", s.config.BackendURL(), incomingRequest.URL.Path), body)
	if err != nil {
		return nil, err
	}
	copyHeader(incomingRequest.Header, backendRequest.Header, ignoredRequestHeaders)
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
	copyHeader(backendResponse.Header, responseWriter.Header(), ignoredResponseHeaders)
	responseBody, err := readBody(backendResponse.Body)
	if err != nil {
		s.requestLogger.Warn("Can not read backend response body", "error", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	if requestMethod == methodGet && len(responseBody) > 0 {
		s.requestLogger.Trace("Decrypt response body with", "length", len(responseBody))
		if err := s.transformer.FromSops(s.config, responseBody, func(result []byte) error { responseBody = result; return nil }); err != nil {
			s.requestLogger.Warn("Can not decrypt body. Leave body unchanged", "error", err)
		}
		s.requestLogger.Trace("Decrypted response body with", "length", len(responseBody))
	}
	if backendResponse.StatusCode/100 != 2 {
		s.requestLogger.Warn("Unexpected backendResponse", "status-code", backendResponse.StatusCode)
	}
	responseWriter.WriteHeader(backendResponse.StatusCode)
	responseWriter.Write(responseBody)
}

func readBody(body io.ReadCloser) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, body); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func copyHeader(from, to http.Header, ignoredHeaders ignoredHeaders) {
	for k, vs := range from {
		if isIgnoredHeader(k, ignoredHeaders) {
			continue
		}
		for _, v := range vs {
			to.Add(k, v)
		}
	}
}

func isIgnoredHeader(header string, ignoredHeaders ignoredHeaders) bool {
	_, ok := ignoredHeaders[header]
	return ok
}

func isSupportedRequestMethod(requestMethod string) bool {
	_, ok := supportedRequestMethods[requestMethod]
	return ok
}

func buildIncomingURI(url *url.URL) (result string) {
	result = url.Path
	if len(url.Query().Encode()) > 0 {
		result = fmt.Sprintf("%s?%s", result, url.Query().Encode())
	}
	return
}
