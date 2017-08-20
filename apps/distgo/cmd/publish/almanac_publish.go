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

package publish

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/palantir/godel/apps/distgo/params"
)

func almanacPublish(artifactURL string, almanacInfo AlmanacInfo, buildSpec params.ProductBuildSpec, distCfg params.Dist, stdout io.Writer) error {
	client := http.DefaultClient
	if err := almanacInfo.CheckProduct(client, buildSpec.ProductName); err != nil {
		if err := almanacInfo.CreateProduct(client, buildSpec.ProductName); err != nil {
			return fmt.Errorf("failed to create product %s: %v", buildSpec.ProductName, err)
		}
	}

	if err := almanacInfo.CheckProductBranch(client, buildSpec.ProductName, buildSpec.VersionInfo.Branch); err != nil {
		if err := almanacInfo.CreateProductBranch(client, buildSpec.ProductName, buildSpec.VersionInfo.Branch); err != nil {
			return fmt.Errorf("failed to create branch %s for product %s: %v", buildSpec.VersionInfo.Branch, buildSpec.ProductName, err)
		}
	}

	if bytes, err := almanacInfo.GetUnit(client, buildSpec.ProductName, buildSpec.VersionInfo.Branch, buildSpec.VersionInfo.Revision); err == nil {
		// unit exists -- check if URL matches artifact URL
		if almanacURLMatches(artifactURL, bytes) {
			fmt.Fprintf(stdout, "Unit for product %s branch %s revision %s with URL %s already exists; skipping publish.\n", buildSpec.ProductName, buildSpec.VersionInfo.Branch, buildSpec.VersionInfo.Revision, artifactURL)
			return nil
		}
		return fmt.Errorf("unit for product %s branch %s revision %s already exists; not overwriting it", buildSpec.ProductName, buildSpec.VersionInfo.Branch, buildSpec.VersionInfo.Revision)
	}

	if err := almanacInfo.CreateUnit(client, AlmanacUnit{
		Product:  buildSpec.ProductName,
		Branch:   buildSpec.VersionInfo.Branch,
		Revision: buildSpec.VersionInfo.Revision,
		Metadata: distCfg.Publish.Almanac.Metadata,
		Tags:     distCfg.Publish.Almanac.Tags,
		URL:      artifactURL,
	}, buildSpec.ProductVersion, stdout); err != nil {
		return fmt.Errorf("failed to publish unit: %v", err)
	}

	if almanacInfo.Release {
		if err := almanacInfo.ReleaseProduct(client, buildSpec.ProductName, buildSpec.VersionInfo.Branch, buildSpec.VersionInfo.Revision); err != nil {
			return fmt.Errorf("failed to release unit: %v", err)
		}
	}

	return nil
}

// Returns true if the provided almanacResponse is a JSON object that contains a key named "url" that has a string value
// that matches the provided url, false otherwise.
func almanacURLMatches(url string, almanacResponse []byte) bool {
	var jsonMap map[string]*json.RawMessage
	if err := json.Unmarshal(almanacResponse, &jsonMap); err == nil {
		if urlJSON, ok := jsonMap["url"]; ok && urlJSON != nil {
			var dstURL string
			if err := json.Unmarshal(*urlJSON, &dstURL); err == nil {
				return url == dstURL
			}
		}
	}
	return false
}
