// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pluginapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const usage = "Prints the gödel plugin.Info struct in its JSON-serialized form"

func CobraInfoCmd(info PluginInfo) *cobra.Command {
	return &cobra.Command{
		Use:    PluginInfoCommandName,
		Short:  usage,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return infoAction(info, cmd.OutOrStdout())
		},
	}
}

type UpgradeConfigFn func(cfg []byte) ([]byte, error)

func CobraUpgradeConfigCmd(upgradeFn UpgradeConfigFn) *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade-config [base64-config-content]",
		Short: "Print the upgraded version of the provided config",
		Long: `Prints the base64-encoded representation of the updated version of the provided configuration, which is
provided as a base64-encoded string. If the provided configuration is a valid representation of the newest version, it 
is printed unmodified. If the upgrade fails, exits with a non-0 exit code and prints the error.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.Errorf("input must be exactly one argument, was %d: %v", len(args), args)
			}
			cfgBytes, err := base64.StdEncoding.DecodeString(args[0])
			if err != nil {
				return errors.Wrapf(err, "failed to decode input as base64")
			}

			upgradedCfg, err := upgradeFn(cfgBytes)
			if err != nil {
				return err
			}
			cmd.Print(base64.StdEncoding.EncodeToString(upgradedCfg))
			return nil
		},
	}
}

func infoAction(info PluginInfo, w io.Writer) error {
	bytes, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal plugin info: %v", err)
	}
	fmt.Fprint(w, string(bytes))
	return nil
}
