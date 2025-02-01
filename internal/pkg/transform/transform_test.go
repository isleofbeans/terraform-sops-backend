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

package transform

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	terraformConfig "github.com/wtschreiter/terraformsopsbackend/internal/pkg/config"
)

func TestTransformToSops(t *testing.T) {
	tests := []struct {
		name                 string
		agePublicKey         string
		agePrivateKey        string
		vaultAddr            string
		vaultAppRoleID       string
		vaultAppRoleSecretID string
		vaultKeyMount        string
		vaultKeyName         string
		wantErr              bool
	}{
		{
			name:          "plain positive age",
			agePublicKey:  "age17gnuhjensr0f902238xt4jkdu9qh9anhjklfn7tr8m3ex5ltxfxqt3yx08",
			agePrivateKey: "AGE-SECRET-KEY-1Z22A6EL3ECQC96ZDMPD5KRPUX32SCAMU2DJGV3Q48PXN3ZW535VQFQDEF9",
		},
		{
			name:    "error on missing public age key",
			wantErr: true,
		},
		{
			name:         "error on corrupted public age key",
			agePublicKey: "this-is-not-a-public-age-key",
			wantErr:      true,
		},
		{
			name:                 "plain positive vault",
			agePublicKey:         "age17gnuhjensr0f902238xt4jkdu9qh9anhjklfn7tr8m3ex5ltxfxqt3yx08",
			vaultAddr:            "http://127.0.0.1:8200",
			vaultAppRoleID:       "2e8a30c1-66e3-343f-6c14-f5c228eb4171",
			vaultAppRoleSecretID: "37cf928c-0c3b-6e2d-d099-d24c94fe6a75",
			vaultKeyMount:        "sops",
			vaultKeyName:         "terraform",
		},
		{
			name:                 "error on missing vault key mount",
			agePublicKey:         "age17gnuhjensr0f902238xt4jkdu9qh9anhjklfn7tr8m3ex5ltxfxqt3yx08",
			vaultAddr:            "http://127.0.0.1:8200",
			vaultAppRoleID:       "2e8a30c1-66e3-343f-6c14-f5c228eb4171",
			vaultAppRoleSecretID: "37cf928c-0c3b-6e2d-d099-d24c94fe6a75",
			vaultKeyMount:        "not/defined/mount",
			vaultKeyName:         "terraform",
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var unencryptedJson []byte
			var unencryptedTFState tfstate
			var encryptedJson []byte
			var encryptedTFState tfstate
			var decryptedJson []byte
			var decryptedTFState tfstate
			var config terraformConfig.TransformConfig = newConfig(
				tt.agePublicKey,
				tt.agePrivateKey,
				tt.vaultAddr,
				tt.vaultAppRoleID,
				tt.vaultAppRoleSecretID,
				tt.vaultKeyMount,
				tt.vaultKeyName,
			)

			// Prepare
			unencryptedJson, err = os.ReadFile("fixtures/tfstates/unencrypted.tfstate")
			if !assert.NoError(t, err) {
				return
			}
			err = json.Unmarshal(unencryptedJson, &unencryptedTFState)
			if !assert.NoError(t, err) {
				return
			}

			// Act
			if err = TransformToSops(config, unencryptedJson, func(sopsResult []byte) { encryptedJson = sopsResult }); (err != nil) != tt.wantErr {
				t.Errorf("TransformToSops() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Validate result
			assert.True(t, len(encryptedJson) > 0)
			err = json.Unmarshal(encryptedJson, &encryptedTFState)
			if !assert.NoError(t, err) {
				fmt.Println(string(encryptedJson))
				return
			}
			assert.Equal(t, unencryptedTFState.Version, encryptedTFState.Version)
			assert.Equal(t, unencryptedTFState.TerraformVersion, encryptedTFState.TerraformVersion)
			assert.Equal(t, unencryptedTFState.Serial, encryptedTFState.Serial)
			assert.Equal(t, unencryptedTFState.Lineage, encryptedTFState.Lineage)
			assert.NotEqual(t, unencryptedTFState.Outputs.Password.Value, encryptedTFState.Outputs.Password.Value)

			// Act reverse
			err = TransformFromSops(config, encryptedJson, func(result []byte) error { decryptedJson = result; return nil })
			if !assert.NoError(t, err) {
				return
			}

			// Validate reverse
			err = json.Unmarshal(decryptedJson, &decryptedTFState)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, unencryptedTFState, decryptedTFState)
		})
	}
}

func TestTransformFromSops(t *testing.T) {
	withVaultToSops := newConfig(
		"age17gnuhjensr0f902238xt4jkdu9qh9anhjklfn7tr8m3ex5ltxfxqt3yx08",
		"",
		"http://127.0.0.1:8200",
		"2e8a30c1-66e3-343f-6c14-f5c228eb4171",
		"37cf928c-0c3b-6e2d-d099-d24c94fe6a75",
		"sops",
		"terraform",
	)
	withoutVaultToSops := newConfig(
		"age17gnuhjensr0f902238xt4jkdu9qh9anhjklfn7tr8m3ex5ltxfxqt3yx08",
		"",
		"",
		"",
		"",
		"",
		"",
	)
	tests := []struct {
		name                 string
		toSopsConfig         terraformConfig.TransformConfig
		vaultAddr            string
		agePrivateKey        string
		vaultAppRoleID       string
		vaultAppRoleSecretID string
		wantErr              bool
	}{
		{
			name:                 "plain positive with vault",
			toSopsConfig:         withVaultToSops,
			vaultAddr:            "http://127.0.0.1:8200",
			vaultAppRoleID:       "2e8a30c1-66e3-343f-6c14-f5c228eb4171",
			vaultAppRoleSecretID: "37cf928c-0c3b-6e2d-d099-d24c94fe6a75",
		},
		{
			name:          "plain positive without vault",
			toSopsConfig:  withoutVaultToSops,
			agePrivateKey: "AGE-SECRET-KEY-1Z22A6EL3ECQC96ZDMPD5KRPUX32SCAMU2DJGV3Q48PXN3ZW535VQFQDEF9",
		},
		{
			name:          "positive with vault fall back to age",
			toSopsConfig:  withVaultToSops,
			vaultAddr:     "http://127.0.0.1:8200",
			agePrivateKey: "AGE-SECRET-KEY-1Z22A6EL3ECQC96ZDMPD5KRPUX32SCAMU2DJGV3Q48PXN3ZW535VQFQDEF9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var unencryptedJson []byte
			var unencryptedTFState tfstate
			var encryptedJson []byte
			var decryptedJson []byte
			var decryptedTFState tfstate
			var config terraformConfig.TransformConfig = newConfig(
				"",
				tt.agePrivateKey,
				tt.vaultAddr,
				tt.vaultAppRoleID,
				tt.vaultAppRoleSecretID,
				"",
				"",
			)

			// Prepare
			unencryptedJson, err = os.ReadFile("fixtures/tfstates/unencrypted.tfstate")
			if !assert.NoError(t, err) {
				return
			}
			err = json.Unmarshal(unencryptedJson, &unencryptedTFState)
			if !assert.NoError(t, err) {
				return
			}
			err = TransformToSops(tt.toSopsConfig, unencryptedJson, func(sopsResult []byte) { encryptedJson = sopsResult })
			if !assert.NoError(t, err) {
				return
			}

			// Act
			if err = TransformFromSops(config, encryptedJson, func(sopsResult []byte) error { decryptedJson = sopsResult; return nil }); (err != nil) != tt.wantErr {
				t.Errorf("TransformFromSops() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Validate result
			err = json.Unmarshal(decryptedJson, &decryptedTFState)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, unencryptedTFState, decryptedTFState)
		})
	}
}

type tfstate struct {
	Version          int    `json:"version"`
	TerraformVersion string `json:"terraform_version"`
	Serial           int    `json:"serial"`
	Lineage          string `json:"lineage"`
	Outputs          struct {
		Password struct {
			Value string `json:"value"`
		} `json:"password"`
	} `json:"outputs"`
	Resources    []map[string]interface{} `json:"resources"`
	CheckResults []map[string]interface{} `json:"check_results"`
}

type testConfig struct {
	agePrivateKey        string
	agePublicKey         string
	vaultAddr            string
	vaultAppRoleID       string
	vaultAppRoleSecretID string
	vaultKeyMount        string
	vaultKeyName         string
}

func (c testConfig) AgePrivateKey() string        { return c.agePrivateKey }
func (c testConfig) AgePublicKey() string         { return c.agePublicKey }
func (c testConfig) VaultAddr() string            { return c.vaultAddr }
func (c testConfig) VaultAppRoleID() string       { return c.vaultAppRoleID }
func (c testConfig) VaultAppRoleSecretID() string { return c.vaultAppRoleSecretID }
func (c testConfig) VaultKeyMount() string        { return c.vaultKeyMount }
func (c testConfig) VaultKeyName() string         { return c.vaultKeyName }

func newConfig(
	agePublicKey,
	agePrivateKey,
	vaultAddr,
	vaultAppRoleID,
	vaultAppRoleSecretID,
	vaultKeyMount,
	vaultKeyName string,
) terraformConfig.TransformConfig {
	return testConfig{
		agePrivateKey:        agePrivateKey,
		agePublicKey:         agePublicKey,
		vaultAddr:            vaultAddr,
		vaultAppRoleID:       vaultAppRoleID,
		vaultAppRoleSecretID: vaultAppRoleSecretID,
		vaultKeyMount:        vaultKeyMount,
		vaultKeyName:         vaultKeyName,
	}
}
