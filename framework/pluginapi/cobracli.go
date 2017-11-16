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
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

const usage = "Prints the gödel plugin.Info struct in its JSON-serialized form"

func CobraInfoCmd(info Info) *cobra.Command {
	return &cobra.Command{
		Use:    InfoCommandName,
		Short:  usage,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return infoAction(info, cmd.OutOrStdout())
		},
	}
}

func infoAction(info Info, w io.Writer) error {
	bytes, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal plugin info: %v", err)
	}
	fmt.Fprint(w, string(bytes))
	return nil
}
