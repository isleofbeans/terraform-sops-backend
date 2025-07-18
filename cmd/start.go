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

package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/backend"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/config"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/monitoring"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/server"
	"github.com/wtschreiter/terraformsopsbackend/internal/pkg/transformer"
)

const (
	traceLevel           string = "TRACE"
	debugLevel           string = "DEBUG"
	infoLevel            string = "INFO"
	warnLevel            string = "WARN"
	errorLevel           string = "ERROR"
	offLevel             string = "OFF"
	cobraKeyAgePublicKey string = "age-public-key"
	viperKeyAgePublicKey string = "transform.age.public_key"

	cobraKeyAgePrivateKey string = "age-private-key"
	viperKeyAgePrivateKey string = "transform.age.private_key"

	cobraKeyVaultAddr string = "vault-addr"
	viperKeyVaultAddr string = "transform.vault.address"

	cobraKeyVaultAppRoleID string = "vault-app-role-id"
	viperKeyVaultAppRoleID string = "transform.vault.app_role.id"

	cobraKeyVaultAppRoleSecretID string = "vault-app-role-secret-id"
	viperKeyVaultAppRoleSecretID string = "transform.vault.app_role.secret_id"

	cobraKeyVaultTransitMount string = "vault-transit-mount"
	viperKeyVaultTransitMount string = "transform.vault.transit.mount"

	cobraKeyVaultTransitName string = "vault-transit-name"
	viperKeyVaultTransitName string = "transform.vault.transit.name"

	cobraKeyServerPort string = "port"
	viperKeyServerPort string = "server.port"

	cobraKeyBackendURL string = "backend-url"
	viperKeyBackendURL string = "backend.url"

	cobraKeyBackendLockMethod string = "backend-lock-method"
	viperKeyBackendLockMethod string = "backend.lock_method"

	cobraKeyBackendUnlockMethod string = "backend-unlock-method"
	viperKeyBackendUnlockMethod string = "backend.unlock_method"

	cobraKeyBackendReadinessProbePath string = "backend-readiness-probe-path"
	viperKeyBackendReadinessProbePath string = "backend.readiness_probe.path"
)

var (
	allowedLogLevel map[string]hclog.Level = map[string]hclog.Level{
		traceLevel: hclog.Trace,
		debugLevel: hclog.Debug,
		infoLevel:  hclog.Info,
		warnLevel:  hclog.Warn,
		errorLevel: hclog.Error,
		offLevel:   hclog.Off,
	}
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starting the service",
	Long: `Starts the web service for the terraform SOPS backend.

This uses the default interface to accept incoming traffic`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cfgFile)
		config, err := newServerConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, config)
			cmd.Usage()
			os.Exit(200)
		}
		if config.Logger().IsDebug() {
			fmt.Fprintln(os.Stderr, config)
		}
		backendClient := backend.New(config.Logger().Named("backend"))
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			monitoring.NewMonitoringServer(
				config,
				backendClient,
			).Start()
		}()
		go func() {
			defer wg.Done()
			server.New(
				config,
				backendClient,
				transformer.New(),
			).Start()
		}()
		wg.Wait()
	},
}

func initStartCmd() {
	rootCmd.AddCommand(startCmd)

	registerStringParameter(startCmd, cobraKeyAgePublicKey, viperKeyAgePublicKey, "public AGE key to encrypt terraform state", true)
	registerStringParameter(startCmd, cobraKeyAgePrivateKey, viperKeyAgePrivateKey, "private AGE key to decrypt terraform state", false)
	registerStringParameter(startCmd, cobraKeyVaultAddr, viperKeyVaultAddr, "vault address to de- and encrypt terraform state", false)
	registerStringParameter(startCmd, cobraKeyVaultAppRoleID, viperKeyVaultAppRoleID, "(required if --vault-addr != \"\") AppRole ID to authenticate with vault", false)
	registerStringParameter(startCmd, cobraKeyVaultAppRoleSecretID, viperKeyVaultAppRoleSecretID, "(required if --vault-addr != \"\") AppRole secret ID to authenticate with vault", false)
	registerStringParameterWithDefault(startCmd, cobraKeyVaultTransitMount, viperKeyVaultTransitMount, "mount point of the transit engine to use", false, "sops")
	registerStringParameterWithDefault(startCmd, cobraKeyVaultTransitName, viperKeyVaultTransitName, "name of the transit engine secret to use", false, "terraform")
	registerStringParameterWithDefault(startCmd, cobraKeyServerPort, viperKeyServerPort, "port the service is listening to", false, "8080")
	registerStringParameter(startCmd, cobraKeyBackendURL, viperKeyBackendURL, "base url to connect with the backend terraform state server", true)
	registerStringParameterWithDefault(startCmd, cobraKeyBackendLockMethod, viperKeyBackendLockMethod, "lock method to use with the backend terraform state server", false, "LOCK")
	registerStringParameterWithDefault(startCmd, cobraKeyBackendUnlockMethod, viperKeyBackendUnlockMethod, "unlock method to use with the backend terraform state server", false, "UNLOCK")
	registerStringParameterWithDefault(startCmd, cobraKeyBackendReadinessProbePath, viperKeyBackendReadinessProbePath, "path to probe backend for readiness.", false, "/")

	//-------

	registerBoolParameterWithDefault(startCmd, "log-json", "log.json", "if logging has to use json format", false)
	registerStringParameterWithDefault(startCmd, "log-level", "log.level", fmt.Sprintf("active log level one of [%s]", strings.Join([]string{
		traceLevel,
		debugLevel,
		infoLevel,
		warnLevel,
		errorLevel,
		offLevel,
	}, ", ")), false, infoLevel)

}

func newServerConfig() (config.ServerConfig, error) {
	c := serverConfig{
		logger: newHCLogger("service"),
	}
	return c, config.ValidateServerConfig(c)
}

type serverConfig struct {
	logger hclog.Logger
}

func (c serverConfig) AgePublicKey() string   { return cmdViper.GetString(viperKeyAgePublicKey) }
func (c serverConfig) AgePrivateKey() string  { return cmdViper.GetString(viperKeyAgePrivateKey) }
func (c serverConfig) VaultAddr() string      { return cmdViper.GetString(viperKeyVaultAddr) }
func (c serverConfig) VaultAppRoleID() string { return cmdViper.GetString(viperKeyVaultAppRoleID) }
func (c serverConfig) VaultAppRoleSecretID() string {
	return cmdViper.GetString(viperKeyVaultAppRoleSecretID)
}
func (c serverConfig) VaultKeyMount() string { return cmdViper.GetString(viperKeyVaultTransitMount) }
func (c serverConfig) VaultKeyName() string  { return cmdViper.GetString(viperKeyVaultTransitName) }
func (c serverConfig) ServerPort() string    { return cmdViper.GetString(viperKeyServerPort) }
func (c serverConfig) BackendURL() string    { return cmdViper.GetString(viperKeyBackendURL) }
func (c serverConfig) BackendLockMethod() string {
	return cmdViper.GetString(viperKeyBackendLockMethod)
}
func (c serverConfig) BackendUnlockMethod() string {
	return cmdViper.GetString(viperKeyBackendUnlockMethod)
}
func (c serverConfig) BackendReadinessProbePath() string {
	return cmdViper.GetString(viperKeyBackendReadinessProbePath)
}
func (c serverConfig) Logger() hclog.Logger { return c.logger }
func (c serverConfig) String() string {
	return fmt.Sprintf(
		`---
server:
  port: %s
backend:
  url: %s
  lock_method: %s
  unlock_method: %s
  readiness_probe:
    path: %s
transform:
  age:
    public_key: %s
    private_key: %s
  vault:
    address: %s
    app_role:
      id: %s
      secret_id: %s
    transit:
      mount: %s
      name: %s`,
		c.presentedToStringValue(c.ServerPort()),
		c.presentedToStringValue(c.BackendURL()),
		c.presentedToStringValue(c.BackendLockMethod()),
		c.presentedToStringValue(c.BackendUnlockMethod()),
		c.presentedToStringValue(c.BackendReadinessProbePath()),
		c.presentedToStringValue(c.AgePublicKey()),
		c.hiddenToStringValue(c.AgePrivateKey()),
		c.presentedToStringValue(c.VaultAddr()),
		c.hiddenToStringValue(c.VaultAppRoleID()),
		c.hiddenToStringValue(c.VaultAppRoleSecretID()),
		c.presentedToStringValue(c.VaultKeyMount()),
		c.presentedToStringValue(c.VaultKeyName()),
	)
}
func (c serverConfig) presentedToStringValue(value string) string {
	if len(value) == 0 {
		return "\"\""
	}
	return fmt.Sprintf("\"%s\"", value)
}
func (c serverConfig) hiddenToStringValue(value string) string {
	if len(value) == 0 {
		return "\"\""
	}
	return "\"*****\""
}
func newHCLogger(name string) hclog.Logger {
	logOutput := io.Writer(os.Stderr)

	return hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:              name,
		Level:             logLevel(cmdViper.GetString("log.level")),
		Output:            logOutput,
		IndependentLevels: true,
		JSONFormat:        cmdViper.GetBool("log.json"),
	})
}

func logLevel(value string) hclog.Level {
	if level, ok := allowedLogLevel[value]; ok {
		return level
	}
	return hclog.NoLevel
}
