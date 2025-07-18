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
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/prometheus/client_golang/prometheus"
	transformConfig "github.com/wtschreiter/terraformsopsbackend/internal/pkg/config"
)

type vaultClient struct {
	client           *vault.Client
	ctx              context.Context
	appRoleID        string
	appRoleSecretID  string
	appRoleMountPath string
	token            string
	tokenUntil       time.Time
	logger           hclog.Logger
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
		token:            "",
		tokenUntil:       time.Now().Add(time.Duration(-24) * time.Hour),
		logger:           config.Logger(),
	}
}

func (c *vaultClient) getToken() string {
	if len([]byte(c.token)) > 0 && time.Now().Before(c.tokenUntil) {
		return c.token
	}
	if c.logger.IsDebug() {
		c.logger.Log(hclog.Error, "create new token", "old-until", c.tokenUntil, "now", time.Now(), "before", time.Now().Before(c.tokenUntil), "token-len", len([]byte(c.token)))
	}
	timer := prometheus.NewTimer(vaultRequestDuration.WithLabelValues("token"))
	defer timer.ObserveDuration()
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
	c.token = resp.Auth.ClientToken
	// create new token 60s upfront end of duration
	c.tokenUntil = time.Now().Add(time.Duration(resp.Auth.LeaseDuration-60) * time.Second)
	if c.logger.IsDebug() {
		c.logger.Log(hclog.Error, "new token", "duration", resp.Auth.LeaseDuration, "until", c.tokenUntil)
	}
	return c.token
}
