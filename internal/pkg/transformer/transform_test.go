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

package transformer

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	terraformConfig "github.com/wtschreiter/terraformsopsbackend/internal/pkg/config"
)

func TestToSops(t *testing.T) {
	if !assert.NoError(t, godotenv.Load("../../../.testenv")) {
		return
	}
	agePublicKey := os.Getenv("TRANSFORM_AGE_PUBLIC_KEY")
	if !assert.NotEmpty(t, agePublicKey, "require TRANSFORM_AGE_PUBLIC_KEY") {
		return
	}
	agePrivateKey := os.Getenv("TRANSFORM_AGE_PRIVATE_KEY")
	if !assert.NotEmpty(t, agePrivateKey, "require TRANSFORM_AGE_PRIVATE_KEY") {
		return
	}
	vaultAddr := os.Getenv("TRANSFORM_VAULT_ADDRESS")
	if !assert.NotEmpty(t, vaultAddr, "require TRANSFORM_VAULT_ADDRESS") {
		return
	}
	vaultAppRoleID := os.Getenv("TRANSFORM_VAULT_APP_ROLE_ID")
	if !assert.NotEmpty(t, vaultAppRoleID, "require TRANSFORM_VAULT_APP_ROLE_ID") {
		return
	}
	vaultAppRoleSecretID := os.Getenv("TRANSFORM_VAULT_APP_ROLE_SECRET_ID")
	if !assert.NotEmpty(t, vaultAppRoleSecretID, "require TRANSFORM_VAULT_APP_ROLE_SECRET_ID") {
		return
	}
	vaultKeyMount := os.Getenv("TRANSFORM_VAULT_TRANSIT_MOUNT")
	if !assert.NotEmpty(t, vaultKeyMount, "require TRANSFORM_VAULT_TRANSIT_MOUNT") {
		return
	}
	vaultKeyName := os.Getenv("TRANSFORM_VAULT_TRANSIT_NAME")
	if !assert.NotEmpty(t, vaultKeyName, "require TRANSFORM_VAULT_TRANSIT_NAME") {
		return
	}
	transformer := New()
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
			agePublicKey:  agePublicKey,
			agePrivateKey: agePrivateKey,
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
			agePublicKey:         agePublicKey,
			vaultAddr:            vaultAddr,
			vaultAppRoleID:       vaultAppRoleID,
			vaultAppRoleSecretID: vaultAppRoleSecretID,
			vaultKeyMount:        vaultKeyMount,
			vaultKeyName:         vaultKeyName,
		},
		{
			name:                 "error on missing vault key mount",
			agePublicKey:         agePublicKey,
			vaultAddr:            vaultAddr,
			vaultAppRoleID:       vaultAppRoleID,
			vaultAppRoleSecretID: vaultAppRoleSecretID,
			vaultKeyMount:        "not/defined/mount",
			vaultKeyName:         vaultKeyName,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var unencryptedJSON []byte
			var unencryptedTFState tfstate
			var encryptedJSON []byte
			var encryptedTFState tfstate
			var decryptedJSON []byte
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
			unencryptedJSON, err = os.ReadFile("fixtures/tfstates/unencrypted.tfstate")
			if !assert.NoError(t, err) {
				return
			}
			err = json.Unmarshal(unencryptedJSON, &unencryptedTFState)
			if !assert.NoError(t, err) {
				return
			}

			// Act
			if err = transformer.ToSops(config, unencryptedJSON, func(sopsResult []byte) { encryptedJSON = sopsResult }); (err != nil) != tt.wantErr {
				t.Errorf("TransformToSops() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Validate result
			assert.True(t, len(encryptedJSON) > 0)
			err = json.Unmarshal(encryptedJSON, &encryptedTFState)
			if !assert.NoError(t, err) {
				fmt.Println(string(encryptedJSON))
				return
			}
			assert.Equal(t, unencryptedTFState.Version, encryptedTFState.Version)
			assert.Equal(t, unencryptedTFState.TerraformVersion, encryptedTFState.TerraformVersion)
			assert.Equal(t, unencryptedTFState.Serial, encryptedTFState.Serial)
			assert.Equal(t, unencryptedTFState.Lineage, encryptedTFState.Lineage)
			assert.NotEqual(t, unencryptedTFState.Outputs.Password.Value, encryptedTFState.Outputs.Password.Value)

			// Act reverse
			err = transformer.FromSops(config, encryptedJSON, func(result []byte) error { decryptedJSON = result; return nil })
			if !assert.NoError(t, err) {
				return
			}

			// Validate reverse
			err = json.Unmarshal(decryptedJSON, &decryptedTFState)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, unencryptedTFState, decryptedTFState)
		})
	}
}

func TestFromSops(t *testing.T) {
	if !assert.NoError(t, godotenv.Load("../../../.testenv")) {
		return
	}
	agePublicKey := os.Getenv("TRANSFORM_AGE_PUBLIC_KEY")
	if !assert.NotEmpty(t, agePublicKey, "require TRANSFORM_AGE_PUBLIC_KEY") {
		return
	}
	agePrivateKey := os.Getenv("TRANSFORM_AGE_PRIVATE_KEY")
	if !assert.NotEmpty(t, agePrivateKey, "require TRANSFORM_AGE_PRIVATE_KEY") {
		return
	}
	vaultAddr := os.Getenv("TRANSFORM_VAULT_ADDRESS")
	if !assert.NotEmpty(t, vaultAddr, "require TRANSFORM_VAULT_ADDRESS") {
		return
	}
	vaultAppRoleID := os.Getenv("TRANSFORM_VAULT_APP_ROLE_ID")
	if !assert.NotEmpty(t, vaultAppRoleID, "require TRANSFORM_VAULT_APP_ROLE_ID") {
		return
	}
	vaultAppRoleSecretID := os.Getenv("TRANSFORM_VAULT_APP_ROLE_SECRET_ID")
	if !assert.NotEmpty(t, vaultAppRoleSecretID, "require TRANSFORM_VAULT_APP_ROLE_SECRET_ID") {
		return
	}
	vaultKeyMount := os.Getenv("TRANSFORM_VAULT_TRANSIT_MOUNT")
	if !assert.NotEmpty(t, vaultKeyMount, "require TRANSFORM_VAULT_TRANSIT_MOUNT") {
		return
	}
	vaultKeyName := os.Getenv("TRANSFORM_VAULT_TRANSIT_NAME")
	if !assert.NotEmpty(t, vaultKeyName, "require TRANSFORM_VAULT_TRANSIT_NAME") {
		return
	}
	transformer := New()
	withVaultToSops := newConfig(
		agePublicKey,
		"",
		vaultAddr,
		vaultAppRoleID,
		vaultAppRoleSecretID,
		vaultKeyMount,
		vaultKeyName,
	)
	withoutVaultToSops := newConfig(
		agePublicKey,
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
			vaultAddr:            vaultAddr,
			vaultAppRoleID:       vaultAppRoleID,
			vaultAppRoleSecretID: vaultAppRoleSecretID,
		},
		{
			name:          "plain positive without vault",
			toSopsConfig:  withoutVaultToSops,
			agePrivateKey: agePrivateKey,
		},
		{
			name:          "positive with vault fall back to age",
			toSopsConfig:  withVaultToSops,
			vaultAddr:     vaultAddr,
			agePrivateKey: agePrivateKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var unencryptedJSON []byte
			var unencryptedTFState tfstate
			var encryptedJSON []byte
			var decryptedJSON []byte
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
			unencryptedJSON, err = os.ReadFile("fixtures/tfstates/unencrypted.tfstate")
			if !assert.NoError(t, err) {
				return
			}
			err = json.Unmarshal(unencryptedJSON, &unencryptedTFState)
			if !assert.NoError(t, err) {
				return
			}
			err = transformer.ToSops(tt.toSopsConfig, unencryptedJSON, func(sopsResult []byte) { encryptedJSON = sopsResult })
			if !assert.NoError(t, err) {
				return
			}

			// Act
			if err = transformer.FromSops(config, encryptedJSON, func(sopsResult []byte) error { decryptedJSON = sopsResult; return nil }); (err != nil) != tt.wantErr {
				t.Errorf("TransformFromSops() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Validate result
			err = json.Unmarshal(decryptedJSON, &decryptedTFState)
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
