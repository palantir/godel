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

package assetapi

import (
	"encoding/json"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type AssetType string

const (
	Dister        AssetType = "dister"
	Publisher     AssetType = "publisher"
	DockerBuilder AssetType = "docker-builder"
)

const AssetTypeCommand = "asset-type"

func NewAssetTypeCmd(assetType AssetType) *cobra.Command {
	return &cobra.Command{
		Use:   AssetTypeCommand,
		Short: "Prints the JSON representation of the asset type",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, err := json.Marshal(assetType)
			if err != nil {
				return errors.Wrapf(err, "failed to marshal JSON")
			}
			cmd.Print(string(jsonOutput))
			return nil
		},
	}
}

func LoadAssets(assets []string) (map[AssetType][]string, error) {
	loadedAssets := make(map[AssetType][]string)
	for _, currAsset := range assets {
		assetType, err := getAssetType(currAsset)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get asset type for asset %s", currAsset)
		}
		loadedAssets[assetType] = append(loadedAssets[assetType], currAsset)
	}
	return loadedAssets, nil
}

func getAssetType(assetPath string) (AssetType, error) {
	cmd := exec.Command(assetPath, AssetTypeCommand)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "failed to run command %v, output: %s", cmd.Args, string(outputBytes))
	}

	var assetType AssetType
	if err := json.Unmarshal(outputBytes, &assetType); err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal JSON")
	}

	switch assetType {
	case Dister, Publisher, DockerBuilder:
		return assetType, nil
	default:
		return "", errors.Errorf("unrecognized asset type: %s", assetType)
	}
}
