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

	"github.com/getsops/sops/v3/hcvault"
	"github.com/getsops/sops/v3/keyservice"
	"github.com/prometheus/client_golang/prometheus"
	transformConfig "github.com/wtschreiter/terraformsopsbackend/internal/pkg/config"
)

var (
	keyServiceServerCache *keyServiceServer
)

type keyServiceServer struct {
	parent      keyservice.Server
	config      transformConfig.VaultConfig
	vaultClient *vaultClient
}

func cachedKeyServiceServer(config transformConfig.VaultConfig) keyservice.KeyServiceServer {
	if keyServiceServerCache == nil {
		keyServiceServerCache = newKeyServiceServer(config, keyservice.Server{})
	}
	return keyServiceServerCache
}

func newKeyServiceServer(config transformConfig.VaultConfig, parent keyservice.Server) *keyServiceServer {
	vaultClient := newVaultClient(config)
	return &keyServiceServer{
		parent:      parent,
		config:      config,
		vaultClient: vaultClient,
	}
}

func (ks *keyServiceServer) encryptWithVault(key *keyservice.VaultKey, plaintext []byte) ([]byte, error) {
	vaultKey := hcvault.MasterKey{
		VaultAddress: key.VaultAddress,
		EnginePath:   key.EnginePath,
		KeyName:      key.KeyName,
	}
	hcvault.Token(ks.vaultClient.getToken()).ApplyToMasterKey(&vaultKey)
	timer := prometheus.NewTimer(vaultRequestDuration.WithLabelValues("encrypt"))
	err := vaultKey.Encrypt(plaintext)
	timer.ObserveDuration()
	if err != nil {
		return nil, err
	}
	return []byte(vaultKey.EncryptedKey), nil
}

func (ks *keyServiceServer) decryptWithVault(key *keyservice.VaultKey, ciphertext []byte) ([]byte, error) {
	vaultKey := hcvault.MasterKey{
		VaultAddress: key.VaultAddress,
		EnginePath:   key.EnginePath,
		KeyName:      key.KeyName,
	}
	vaultKey.EncryptedKey = string(ciphertext)
	hcvault.Token(ks.vaultClient.getToken()).ApplyToMasterKey(&vaultKey)
	timer := prometheus.NewTimer(vaultRequestDuration.WithLabelValues("decrypt"))
	plaintext, err := vaultKey.Decrypt()
	timer.ObserveDuration()
	return []byte(plaintext), err
}

func (ks *keyServiceServer) Encrypt(ctx context.Context,
	req *keyservice.EncryptRequest) (*keyservice.EncryptResponse, error) {

	switch k := req.Key.KeyType.(type) {
	case *keyservice.Key_VaultKey:
		timer := prometheus.NewTimer(keyServiceRequestDuration.WithLabelValues("encrypt", "vault"))
		defer timer.ObserveDuration()
		ciphertext, err := ks.encryptWithVault(k.VaultKey, req.Plaintext)
		if err != nil {
			return nil, err
		}
		return &keyservice.EncryptResponse{
			Ciphertext: ciphertext,
		}, nil
	default:
		timer := prometheus.NewTimer(keyServiceRequestDuration.WithLabelValues("encrypt", "age"))
		defer timer.ObserveDuration()
		return ks.parent.Encrypt(ctx, req)
	}
}

func (ks *keyServiceServer) Decrypt(ctx context.Context,
	req *keyservice.DecryptRequest) (*keyservice.DecryptResponse, error) {

	switch k := req.Key.KeyType.(type) {
	case *keyservice.Key_VaultKey:
		timer := prometheus.NewTimer(keyServiceRequestDuration.WithLabelValues("decrypt", "vault"))
		defer timer.ObserveDuration()
		plaintext, err := ks.decryptWithVault(k.VaultKey, req.Ciphertext)
		if err != nil {
			return nil, err
		}
		return &keyservice.DecryptResponse{
			Plaintext: plaintext,
		}, nil
	default:
		timer := prometheus.NewTimer(keyServiceRequestDuration.WithLabelValues("decrypt", "age"))
		defer timer.ObserveDuration()
		return ks.parent.Decrypt(ctx, req)
	}
}
