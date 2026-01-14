// Copyright 2026 The Terraform SOPS backend Authors
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
	"io"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/config"
)

func TestNew(t *testing.T) {
	type args struct {
		config config.ServerConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no mTLS",
			args: args{
				config: &testConfig{},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				err    error
				client Client
			)
			if client, err = New(tt.args.config); tt.wantErr(t, err) && err != nil {
				assert.NotNil(t, client)
			}
		})
	}
}

type testConfig struct {
	test     *testing.T
	logger   hclog.Logger
	mTLSCert []byte
	mTLSKey  []byte
}

func (t *testConfig) AgePublicKey() string {
	assert.FailNow(t.test, "unexpected AgePublicKey called")
	return ""
}

func (t *testConfig) AgePrivateKey() string {
	assert.FailNow(t.test, "unexpected AgePrivateKey called")
	return ""
}

func (t *testConfig) VaultAddr() string {
	assert.FailNow(t.test, "unexpected VaultAddr called")
	return ""
}

func (t *testConfig) VaultKeyMount() string {
	assert.FailNow(t.test, "unexpected VaultKeyMount called")
	return ""
}

func (t *testConfig) VaultKeyName() string {
	assert.FailNow(t.test, "unexpected VaultKeyName called")
	return ""
}

func (t *testConfig) VaultAppRoleID() string {
	assert.FailNow(t.test, "unexpected VaultAppRoleID called")
	return ""
}

func (t *testConfig) VaultAppRoleSecretID() string {
	assert.FailNow(t.test, "unexpected VaultAppRoleSecretID called")
	return ""
}

func (t *testConfig) Logger() hclog.Logger {
	if t.logger == nil {
		t.logger = newTestHCLogger()
	}
	return t.logger
}

func (t *testConfig) ServerPort() string {
	assert.FailNow(t.test, "unexpected ServerPort called")
	return ""
}

func (t *testConfig) BackendURL() string {
	assert.FailNow(t.test, "unexpected BackendURL called")
	return ""
}

func (t *testConfig) BackendMTLSCert() []byte {
	return t.mTLSCert
}

func (t *testConfig) BackendMTLSKey() []byte {
	return t.mTLSKey
}

func (t *testConfig) BackendLockMethod() string {
	assert.FailNow(t.test, "unexpected BackendLockMethod called")
	return ""
}

func (t *testConfig) BackendUnlockMethod() string {
	assert.FailNow(t.test, "unexpected BackendUnlockMethod called")
	return ""
}

func (t *testConfig) BackendReadinessProbePath() string {
	assert.FailNow(t.test, "unexpected BackendReadinessProbePath called")
	return ""
}

func (t *testConfig) String() string {
	return "testConfig"
}

func newTestHCLogger() hclog.Logger {
	logOutput := io.Writer(os.Stderr)

	return hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:              "unit-test",
		Level:             hclog.Trace,
		Output:            logOutput,
		IndependentLevels: true,
		JSONFormat:        false,
	})
}
