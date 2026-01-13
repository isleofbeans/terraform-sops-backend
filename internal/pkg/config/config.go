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

package config

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
)

// AgeConfig provides keys to handle AGE de-/encryption
type AgeConfig interface {
	AgePublicKey() string
	AgePrivateKey() string
}

// VaultConfig provides access data to a Vault server
type VaultConfig interface {
	VaultAddr() string
	VaultKeyMount() string
	VaultKeyName() string
	VaultAppRoleID() string
	VaultAppRoleSecretID() string
	Logger() hclog.Logger
}

// TransformConfig provides transform configuration data
type TransformConfig interface {
	AgeConfig
	VaultConfig
}

// ServerConfig provides configuration to a terraform SOPS backend server
type ServerConfig interface {
	TransformConfig
	ServerPort() string
	BackendURL() string
	BackendLockMethod() string
	BackendUnlockMethod() string
	BackendReadinessProbePath() string
	Logger() hclog.Logger
	String() string
}

// ValidateServerConfig returns with error if the config is not valid
func ValidateServerConfig(config ServerConfig) error {
	if config.AgePublicKey() == "" {
		return fmt.Errorf("AGE public key required")
	}
	if config.BackendURL() == "" {
		return fmt.Errorf("backend URL required")
	}
	if config.VaultAddr() == "" && config.AgePrivateKey() == "" {
		return fmt.Errorf("vault address or AGE private key required")
	}
	if config.VaultAddr() != "" && config.VaultAppRoleID() == "" {
		return fmt.Errorf("vault AppRole ID required")
	}
	if config.VaultAddr() != "" && config.VaultAppRoleSecretID() == "" {
		return fmt.Errorf("vault AppRole secret ID required")
	}
	return nil
}
