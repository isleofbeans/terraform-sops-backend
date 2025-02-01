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
	"context"
	"log"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	transformConfig "github.com/wtschreiter/terraformsopsbackend/internal/pkg/config"
)

type vaultClient struct {
	client           *vault.Client
	ctx              context.Context
	appRoleID        string
	appRoleSecretID  string
	appRoleMountPath string
}

func newVaultClient(config transformConfig.VaultConfig) *vaultClient {
	client, err := vault.New(
		vault.WithAddress(config.VaultAddr()),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	return &vaultClient{
		client:           client,
		ctx:              ctx,
		appRoleID:        config.VaultAppRoleID(),
		appRoleSecretID:  config.VaultAppRoleSecretID(),
		appRoleMountPath: "approle",
	}
}

func (c *vaultClient) getToken() string {
	resp, err := c.client.Auth.AppRoleLogin(
		c.ctx,
		schema.AppRoleLoginRequest{
			RoleId:   c.appRoleID,
			SecretId: c.appRoleSecretID,
		},
	)
	if err != nil {
		return ""
	}
	return resp.Auth.ClientToken
}
