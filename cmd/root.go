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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "terraform-sops-backend",
		Short: "SOPS encrypting terraform http backend",
		Long: `terraform-sops-backend serves as an intermediate terraform HTTP backend.
It is encrypting and decrypting the state file before it is passing it on
to the backend terraform HTTP backend`,
		DisableAutoGenTag: true,
	}
	cfgFile       string
	cmdViper      *viper.Viper
	viperReplacer replacer
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	viperReplacer = replacer{}
	cmdViper = viper.NewWithOptions(viper.EnvKeyReplacer(viperReplacer))
	cobra.OnInitialize(initConfig)

	initRootCmd()
	initStartCmd()

}

func initRootCmd() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/etc/terraform-sops-backend/conf.yaml", "config file")
}

func initConfig() {
	cmdViper.SetConfigFile(cfgFile)

	cmdViper.AutomaticEnv()

	cmdViper.ReadInConfig()
}

type replacer struct{}

func (replacer) Replace(r string) string {
	return strings.Replace(r, ".", "_", -1)
}

func registerStringParameter(cmd *cobra.Command, cobraKey, viperKey, helpText string, required bool) {
	registerStringParameterWithDefault(cmd, cobraKey, viperKey, helpText, required, "")
}

func registerStringParameterWithDefault(cmd *cobra.Command, cobraKey, viperKey, helpText string, required bool, defaultValue string) {
	requiredText := "optional"
	if required {
		requiredText = "required"
	}
	cmd.Flags().String(cobraKey, defaultValue, fmt.Sprintf("%s (%s) %s", envVarName(viperKey), requiredText, helpText))
	cmdViper.BindPFlag(viperKey, cmd.Flags().Lookup(cobraKey))
}

func registerBoolParameterWithDefault(cmd *cobra.Command, cobraKey, viperKey, helpText string, defaultValue bool) {
	cmd.Flags().Bool(cobraKey, defaultValue, fmt.Sprintf("%s (optional) %s", envVarName(viperKey), helpText))
	cmdViper.BindPFlag(viperKey, cmd.Flags().Lookup(cobraKey))
}

func envVarName(viperKey string) string {
	return strings.ToUpper(viperReplacer.Replace(viperKey))
}
