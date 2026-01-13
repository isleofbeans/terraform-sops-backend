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

package transformer

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes"
	"github.com/getsops/sops/v3/age"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/cmd/sops/formats"
	sopsConfig "github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/hcvault"
	"github.com/getsops/sops/v3/keyservice"
	"github.com/getsops/sops/v3/stores/json"
	"github.com/getsops/sops/v3/version"
	"github.com/prometheus/client_golang/prometheus"
	transformConfig "github.com/wtschreiter/terraformsopsbackend/internal/pkg/config"
)

// SOPSTransformer encrypts to SOPS and decrypts from SOPS
type SOPSTransformer interface {
	ToSops(config transformConfig.TransformConfig, input []byte, handler func(result []byte)) error
	FromSops(config transformConfig.TransformConfig, input []byte, handler func(result []byte) error) error
}

// New creates a new SOPSTransformer
func New() SOPSTransformer {
	return transform{}
}

type transform struct{}

// ToSops transforms the input JSON data into a SOPS encrypted JSON data and hands it tho the handler
func (transform) ToSops(config transformConfig.TransformConfig, input []byte, handler func(result []byte)) error {
	timer := prometheus.NewTimer(transformerRequestDuration.WithLabelValues("encrypt"))
	defer timer.ObserveDuration()

	cipher := aes.NewCipher()
	inputStore := inputStore()
	outputStore := outputStore()

	branches, err := inputStore.LoadPlainFile(input)
	if err != nil {
		return err
	}
	if len(branches) < 1 {
		return fmt.Errorf("input cannot be completely empty, it must contain at least one document")
	}
	if outputStore.HasSopsTopLevelKey(branches[0]) {
		return fmt.Errorf("input is already encrypted")
	}

	var group sops.KeyGroup

	ageMasterKey, err := ageMasterKey(config)
	if err != nil {
		return err
	}
	group = append(group, ageMasterKey)

	hcvaultMasterKey, err := hcvaultMasterKey(config)
	if err != nil {
		return err
	}
	if hcvaultMasterKey != nil {
		group = append(group, hcvaultMasterKey)
	}

	tree := sops.Tree{
		Branches: branches,
		Metadata: encryptMetadata(group),
	}
	dataKey, errs := tree.GenerateDataKeyWithKeyServices([]keyservice.KeyServiceClient{keyservice.NewCustomLocalClient(cachedKeyServiceServer(config))})
	if len(errs) > 0 {
		err = fmt.Errorf("could not generate data key: %s", errs)
		return err
	}

	err = common.EncryptTree(common.EncryptTreeOpts{
		DataKey: dataKey,
		Tree:    &tree,
		Cipher:  cipher,
	})
	if err != nil {
		return err
	}

	result, err := outputStore.EmitEncryptedFile(tree)
	if err != nil {
		return fmt.Errorf("could not marshal tree: %s", err)
	}
	handler(result)
	return nil
}

// FromSops transforms the SOPS encrypted input JSON data into a decrypted JSON data and hands it tho the handler
func (transform) FromSops(config transformConfig.TransformConfig, input []byte, handler func(result []byte) error) error {
	timer := prometheus.NewTimer(transformerRequestDuration.WithLabelValues("decrypt"))
	defer timer.ObserveDuration()

	os.Setenv(age.SopsAgeKeyEnv, config.AgePrivateKey())

	store := common.StoreForFormat(formats.Json, sopsConfig.NewStoresConfig())

	// Load SOPS file and access the data key
	tree, err := store.LoadEncryptedFile(input)
	if err != nil {
		return err
	}
	key, err := tree.Metadata.GetDataKeyWithKeyServices(
		[]keyservice.KeyServiceClient{
			keyservice.NewCustomLocalClient(
				cachedKeyServiceServer(config),
			),
		},
		[]string{
			"age",
			"hc_vault",
		},
	)
	if err != nil {
		return err
	}

	// Decrypt the tree
	cipher := aes.NewCipher()
	mac, err := tree.Decrypt(key, cipher)
	if err != nil {
		return err
	}

	// Compute the hash of the cleartext tree and compare it with
	// the one that was stored in the document. If they match,
	// integrity was preserved
	originalMac, err := cipher.Decrypt(
		tree.Metadata.MessageAuthenticationCode,
		key,
		tree.Metadata.LastModified.Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("failed to decrypt original mac: %w", err)
	}
	if originalMac != mac {
		return fmt.Errorf("failed to verify data integrity. expected mac %q, got %q", originalMac, mac)
	}

	result, err := store.EmitPlainFile(tree.Branches)
	if err != nil {
		return err
	}
	return handler(result)
}

func inputStore() sops.Store {
	storesConf := sopsConfig.NewStoresConfig()
	return json.NewStore(&storesConf.JSON)
}

func outputStore() sops.Store {
	storesConf := sopsConfig.NewStoresConfig()
	return json.NewStore(&storesConf.JSON)
}

func agePublicKey(config transformConfig.AgeConfig) (string, error) {
	if config.AgePublicKey() == "" {
		return "", fmt.Errorf("configuration failure, missing public AGE key")
	}
	return config.AgePublicKey(), nil
}

func encryptMetadata(keyGroup sops.KeyGroup) sops.Metadata {
	return sops.Metadata{
		KeyGroups:         []sops.KeyGroup{keyGroup},
		UnencryptedSuffix: "",
		EncryptedSuffix:   "",
		UnencryptedRegex: fmt.Sprintf(
			"^(%s)$",
			strings.Join([]string{
				"version",
				"terraform_version",
				"serial",
				"lineage",
			}, "|")),
		EncryptedRegex:          "",
		UnencryptedCommentRegex: "",
		EncryptedCommentRegex:   "",
		MACOnlyEncrypted:        false,
		Version:                 version.Version,
		ShamirThreshold:         0,
	}
}

func ageMasterKey(config transformConfig.AgeConfig) (*age.MasterKey, error) {
	agePublicKey, err := agePublicKey(config)
	if err != nil {
		return nil, err
	}
	ageKeys, err := age.MasterKeysFromRecipients(agePublicKey)
	if err != nil {
		return nil, err
	}

	if len(ageKeys) != 1 {
		return nil, fmt.Errorf("configuration failure, expected number of age keys = 1 found %v", len(ageKeys))
	}

	return ageKeys[0], nil
}

func hcvaultMasterKey(config transformConfig.VaultConfig) (*hcvault.MasterKey, error) {
	vaultAddr := config.VaultAddr()
	if vaultAddr == "" {
		return nil, nil
	}
	hcVaultKeys, err := hcvault.NewMasterKeysFromURIs(fmt.Sprintf("%s/v1/%s/keys/%s", vaultAddr, config.VaultKeyMount(), config.VaultKeyName()))
	if err != nil {
		return nil, err
	}

	if len(hcVaultKeys) != 1 {
		return nil, fmt.Errorf("configuration failure, expected number of hcvault keys = 1 found %v", len(hcVaultKeys))
	}

	return hcVaultKeys[0], nil
}
