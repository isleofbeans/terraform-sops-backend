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

package server

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/config"
)

func Test_server_buildBackendRequest(t *testing.T) {
	tests := []struct {
		name             string
		incomingMethod   string
		hasIncomingQuery bool
		backendWithPort  bool
		isLockMethod     bool
		isUnlockMethod   bool
		expectsTransform bool
		transformErr     error
		wantErr          bool
	}{
		{
			name:           "GET",
			incomingMethod: methodGet,
		},
		{
			name:            "GET with backend port",
			incomingMethod:  methodGet,
			backendWithPort: true,
		},
		{
			name:             "POST",
			incomingMethod:   methodPost,
			hasIncomingQuery: true,
			expectsTransform: true,
		},
		{
			name:             "POST with failing transform",
			incomingMethod:   methodPost,
			hasIncomingQuery: true,
			expectsTransform: true,
			transformErr:     fmt.Errorf("Expected transform error"),
			wantErr:          true,
		},
		{
			name:           "LOCK",
			incomingMethod: methodLock,
			isLockMethod:   true,
		},
		{
			name:           "UNLOCK",
			incomingMethod: methodUnlock,
			isUnlockMethod: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := randConfig(t, tt.backendWithPort)
			var transformer *simpleTestTransformer
			if tt.expectsTransform {
				transformer = randAllowToSopsTransformer(t, tt.transformErr)
			} else {
				transformer = randAllowNothingTransformer(t)
			}
			incomingRequestBuilder := randRequestBuilder(tt.incomingMethod, tt.hasIncomingQuery)
			s := server{
				config:      config,
				transformer: transformer,
			}
			got, err := s.buildBackendRequest(incomingRequestBuilder.buildRequest())
			if (err != nil) != tt.wantErr {
				t.Errorf("server.buildBackendRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Nil(t, got)
				return
			}
			if tt.isLockMethod {
				assert.Equal(t, config.BackendLockMethod(), got.Method)
			} else if tt.isUnlockMethod {
				assert.Equal(t, config.BackendUnlockMethod(), got.Method)
			} else {
				assert.Equal(t, incomingRequestBuilder.requestMethod, got.Method)
			}
			expectedURL, _ := url.Parse(config.BackendURL())
			assert.Equal(t, expectedURL.Scheme, got.URL.Scheme)
			assert.Equal(t, expectedURL.Host, got.URL.Host)
			assert.Equal(t, expectedURL.Port(), got.URL.Port())
			assert.Equal(t, incomingRequestBuilder.requestPath, got.URL.Path)
			assert.Equal(t, incomingRequestBuilder.requestRawQuery, got.URL.RawQuery)
			assert.Equal(t, incomingRequestBuilder.buildRequest().Header, got.Header)
			gotBody, _ := got.BodyBytes()
			if tt.expectsTransform {
				assert.Equal(t, transformer.output, gotBody)
			} else {
				assert.Equal(t, []byte(incomingRequestBuilder.requestBody), gotBody)
			}
		})
	}
}

func Test_server_writeResponse(t *testing.T) {
	tests := []struct {
		name              string
		incomingMethod    string
		expectsTransform  bool
		transformErr      error
		backendStatusCode int
	}{
		{
			name:              "GET OK",
			incomingMethod:    methodGet,
			expectsTransform:  true,
			backendStatusCode: http.StatusOK,
		},
		{
			name:              "GET NOT OK",
			incomingMethod:    methodGet,
			expectsTransform:  true,
			backendStatusCode: 418,
		},
		{
			name:              "GET OK failed from SOPS",
			incomingMethod:    methodGet,
			expectsTransform:  true,
			transformErr:      fmt.Errorf("Expected transform error"),
			backendStatusCode: http.StatusOK,
		},
		{
			name:              "POST OK",
			incomingMethod:    methodPost,
			backendStatusCode: http.StatusOK,
		},
		{
			name:              "LOCK OK",
			incomingMethod:    methodLock,
			backendStatusCode: http.StatusOK,
		},
		{
			name:              "UNLOCK OK",
			incomingMethod:    methodUnlock,
			backendStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := randConfig(t, false)
			var transformer *simpleTestTransformer
			if tt.expectsTransform {
				transformer = randAllowFromSopsTransformer(t, tt.transformErr)
			} else {
				transformer = randAllowNothingTransformer(t)
			}
			responseWriter := &simpleResponseWriter{}
			backendResponse := randResponse(tt.backendStatusCode)
			s := server{
				config:        config,
				transformer:   transformer,
				requestLogger: config.Logger().Named("frontend"),
			}
			s.writeResponse(responseWriter, backendResponse.build(), tt.incomingMethod, "/test")
			assert.Equal(t, tt.backendStatusCode, responseWriter.statusCode)
			expectedResponseHeader := http.Header{
				"Content-Type":                       []string{backendResponse.responseContentType},
				backendResponse.responseRandomHeader: []string{backendResponse.responseRandomHeaderValue},
			}
			assert.Equal(t, expectedResponseHeader, responseWriter.header)
			if tt.expectsTransform && tt.transformErr == nil {
				assert.Equal(t, string(transformer.output), responseWriter.body.String())
			} else {
				assert.Equal(t, backendResponse.responseBody, responseWriter.body.String())
			}
		})
	}
}

func Test_server_newRequestHandler(t *testing.T) {
	tests := []struct {
		name                        string
		incomingMethod              string
		hasIncomingQuery            bool
		expectsFromSops             bool
		expectsToSops               bool
		transformErr                error
		backendStatusCode           int
		expoectedResponseStatusCode int
	}{
		{
			name:                        "GET OK",
			incomingMethod:              methodGet,
			expectsFromSops:             true,
			backendStatusCode:           http.StatusOK,
			expoectedResponseStatusCode: http.StatusOK,
		},
		{
			name:                        "POST OK",
			incomingMethod:              methodPost,
			expectsToSops:               true,
			backendStatusCode:           http.StatusOK,
			expoectedResponseStatusCode: http.StatusOK,
		},
		{
			name:                        "LOCK OK",
			incomingMethod:              methodLock,
			backendStatusCode:           http.StatusOK,
			expoectedResponseStatusCode: http.StatusOK,
		},
		{
			name:                        "UNLOCK OK",
			incomingMethod:              methodUnlock,
			backendStatusCode:           http.StatusOK,
			expoectedResponseStatusCode: http.StatusOK,
		},
		{
			name:                        "Unsupported Method",
			incomingMethod:              randString(5),
			expoectedResponseStatusCode: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := randConfig(t, false)
			var transformer *simpleTestTransformer
			if tt.expectsToSops && tt.expectsFromSops {
				transformer = randAllowAllTransformer(t, tt.transformErr)
			} else if tt.expectsToSops {
				transformer = randAllowToSopsTransformer(t, tt.transformErr)
			} else if tt.expectsFromSops {
				transformer = randAllowFromSopsTransformer(t, tt.transformErr)
			} else {
				transformer = randAllowNothingTransformer(t)
			}
			responseWriter := &simpleResponseWriter{}
			incomingRequestBuilder := randRequestBuilder(tt.incomingMethod, tt.hasIncomingQuery)
			backendResponseBuilder := randResponse(tt.backendStatusCode)
			backendClient := simpleTestBackendClient{responseBuilder: backendResponseBuilder}
			s := server{
				config:        config,
				transformer:   transformer,
				backend:       backendClient,
				requestLogger: config.Logger().Named("frontend"),
			}
			s.newRequestHandler()(responseWriter, incomingRequestBuilder.buildRequest())

			assert.Equal(t, tt.expoectedResponseStatusCode, responseWriter.statusCode)
		})
	}
}

var (
	testLogger   hclog.Logger = newTestLogger()
	allowedRunes []rune       = []rune("abcdefghijklmnopqrstuvwxyz")
)

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = allowedRunes[rand.Intn(len(allowedRunes))]
	}
	return string(b)
}

func randBackendURL(withPort bool) string {
	if withPort {
		return fmt.Sprintf("https://%s.test:%v", randString(8), rand.Intn(1000)+1000)
	}
	return fmt.Sprintf("https://%s.test", randString(8))
}

func randBackendLockMethod() string {
	return strings.ToUpper(randString(5))
}

func randBackendUnlockMethod() string {
	return strings.ToUpper(randString(5))
}

func randConfig(t *testing.T, backendWithPort bool) config.ServerConfig {
	return &simpleTestServerConfig{
		currentTest:         t,
		backendURL:          randBackendURL(backendWithPort),
		backendLockMethod:   randBackendLockMethod(),
		backendUnlockMethod: randBackendUnlockMethod(),
	}
}

func randResponse(responseCode int) simpleResponseBuilder {

	responseRandomHeader := []rune(randString(10))
	responseRandomHeader[0] = []rune(strings.ToUpper(string([]rune{responseRandomHeader[0]})))[0]
	return simpleResponseBuilder{
		responseBody:              randString(50 + rand.Intn(50)),
		responseCode:              responseCode,
		responseContentType:       randString(20),
		responseRandomHeader:      string(responseRandomHeader),
		responseRandomHeaderValue: randString(20),
	}
}

func randAllowToSopsTransformer(t *testing.T, transformErr error) *simpleTestTransformer {
	return &simpleTestTransformer{
		currentTest: t,
		allowToSops: true,
		output:      []byte(randString(50 + rand.Intn(50))),
		err:         transformErr,
	}
}

func randAllowFromSopsTransformer(t *testing.T, transformErr error) *simpleTestTransformer {
	return &simpleTestTransformer{
		currentTest:   t,
		allowFromSops: true,
		output:        []byte(randString(50 + rand.Intn(50))),
		err:           transformErr,
	}
}

func randAllowAllTransformer(t *testing.T, transformErr error) *simpleTestTransformer {
	return &simpleTestTransformer{
		currentTest:   t,
		allowFromSops: true,
		allowToSops:   true,
		output:        []byte(randString(50 + rand.Intn(50))),
		err:           transformErr,
	}
}

func randAllowNothingTransformer(t *testing.T) *simpleTestTransformer {
	return &simpleTestTransformer{
		currentTest: t,
	}
}

func randRequestBuilder(method string, hasIncomingQuery bool) incomingRequestBuilder {
	var incomingQuery string = ""
	if hasIncomingQuery {
		incomingQuery = fmt.Sprintf("%s=%s", strings.ToUpper(randString(3)), randString(20))
	}
	return incomingRequestBuilder{
		requestMethod:            method,
		requestBody:              randString(50 + rand.Intn(50)),
		requestContentType:       randString(20),
		requestRandomHeader:      randString(10),
		requestRandomHeaderValue: randString(20),
		requestPath:              fmt.Sprintf("/%s/%s", randString(10), randString(10)),
		requestRawQuery:          incomingQuery,
	}
}

func newTestLogger() hclog.Logger {
	return hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:              "test",
		Level:             hclog.Trace,
		Output:            io.Writer(os.Stderr),
		IndependentLevels: true,
		JSONFormat:        false,
	})
}

type simpleTestServerConfig struct {
	currentTest         *testing.T
	backendURL          string
	backendLockMethod   string
	backendUnlockMethod string
}

func (c *simpleTestServerConfig) AgePublicKey() string {
	c.currentTest.Fatal("Unexpected config read AgePublicKey() ")
	return ""
}
func (c *simpleTestServerConfig) AgePrivateKey() string {
	c.currentTest.Fatal("Unexpected config read AgePublicKey() ")
	return ""
}
func (c *simpleTestServerConfig) VaultAddr() string {
	c.currentTest.Fatal("Unexpected config read AgePublicKey() ")
	return ""
}
func (c *simpleTestServerConfig) VaultAppRoleID() string {
	c.currentTest.Fatal("Unexpected config read AgePublicKey() ")
	return ""
}
func (c *simpleTestServerConfig) VaultAppRoleSecretID() string {
	c.currentTest.Fatal("Unexpected config read AgePublicKey() ")
	return ""
}
func (c *simpleTestServerConfig) VaultKeyMount() string {
	c.currentTest.Fatal("Unexpected config read AgePublicKey() ")
	return ""
}
func (c *simpleTestServerConfig) VaultKeyName() string {
	c.currentTest.Fatal("Unexpected config read AgePublicKey() ")
	return ""
}
func (c *simpleTestServerConfig) ServerPort() string {
	c.currentTest.Fatal("Unexpected config read AgePublicKey() ")
	return ""
}
func (c *simpleTestServerConfig) BackendURL() string                { return c.backendURL }
func (c *simpleTestServerConfig) BackendLockMethod() string         { return c.backendLockMethod }
func (c *simpleTestServerConfig) BackendUnlockMethod() string       { return c.backendUnlockMethod }
func (c *simpleTestServerConfig) BackendReadinessProbePath() string { return "/" }
func (c *simpleTestServerConfig) Logger() hclog.Logger              { return testLogger }
func (c *simpleTestServerConfig) String() string                    { return "test-server-config" }

type simpleTestBackendClient struct {
	err             error
	responseBuilder simpleResponseBuilder
}

func (b simpleTestBackendClient) Send(r *retryablehttp.Request) (*http.Response, error) {
	if b.err != nil {
		return nil, b.err
	}
	return b.responseBuilder.build(), nil
}

type simpleTestTransformer struct {
	currentTest   *testing.T
	allowToSops   bool
	allowFromSops bool
	input         []byte
	output        []byte
	err           error
}

func (t *simpleTestTransformer) ToSops(config config.TransformConfig, input []byte, handler func(result []byte)) error {
	if !t.allowToSops {
		t.currentTest.Fatal("Unexpected ToSops cal")
		return fmt.Errorf("Unexpected method call")
	}
	t.input = input
	if t.err != nil {
		return t.err
	}
	handler(t.output)
	return nil
}

func (t *simpleTestTransformer) FromSops(config config.TransformConfig, input []byte, handler func(result []byte) error) error {
	if !t.allowFromSops {
		t.currentTest.Fatal("Unexpected FromSops cal")
		return fmt.Errorf("Unexpected method call")
	}
	t.input = input
	if t.err != nil {
		return t.err
	}
	handler(t.output)
	return nil
}

func (t *simpleTestTransformer) GetInput() []byte { return t.input }

type incomingRequestBuilder struct {
	requestMethod            string
	requestBody              string
	requestContentType       string
	requestRandomHeader      string
	requestRandomHeaderValue string
	requestPath              string
	requestRawQuery          string
}

func (r incomingRequestBuilder) buildRequest() *http.Request {
	requestRandomHeader := []rune(r.requestRandomHeader)
	requestRandomHeader[0] = []rune(strings.ToUpper(string([]rune{requestRandomHeader[0]})))[0]
	return &http.Request{
		Method: r.requestMethod,
		Body:   io.NopCloser(strings.NewReader(r.requestBody)),
		Header: http.Header{
			"Content-Type":              []string{r.requestContentType},
			"Content-Length":            []string{fmt.Sprint(len([]byte(r.requestBody)))},
			string(requestRandomHeader): []string{r.requestRandomHeaderValue},
		},
		URL: &url.URL{
			Path:     r.requestPath,
			User:     &url.Userinfo{},
			RawQuery: r.requestRawQuery,
		},
	}
}

type simpleResponseWriter struct {
	header     http.Header
	body       *bytes.Buffer
	statusCode int
	flushed    bool
}

func (w *simpleResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = http.Header{}
	}
	return w.header
}
func (w *simpleResponseWriter) Write(bodyData []byte) (int, error) {
	if w.body == nil {
		w.body = bytes.NewBuffer(nil)
	}
	return w.body.Write(bodyData)
}
func (w *simpleResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *simpleResponseWriter) Flush() {
	w.flushed = true
}

type simpleResponseBuilder struct {
	responseBody              string
	responseCode              int
	responseContentType       string
	responseRandomHeader      string
	responseRandomHeaderValue string
}

func (b simpleResponseBuilder) build() *http.Response {
	var bodyReader io.ReadCloser = io.NopCloser(strings.NewReader(b.responseBody))
	response := http.Response{
		Status:     fmt.Sprintf("%v testStatus", b.responseCode),
		StatusCode: b.responseCode,
		Header: http.Header{
			"Content-Type":         []string{b.responseContentType},
			"Content-Length":       []string{fmt.Sprint(len([]byte(b.responseBody)))},
			b.responseRandomHeader: []string{b.responseRandomHeaderValue},
		},
		Body: bodyReader,
	}
	return &response
}
