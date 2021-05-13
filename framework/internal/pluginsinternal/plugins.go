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

package pluginsinternal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/palantir/godel/v2/framework/artifactresolver"
	"github.com/palantir/godel/v2/framework/internal/pathsinternal"
	"github.com/palantir/godel/v2/pkg/osarch"
	"github.com/pkg/errors"
)

const (
	IndentSpaces = 4
)

func ResolveAssets(assetsDir, downloadsDir string, assetParams []artifactresolver.LocatorWithResolverParam, osArch osarch.OSArch, defaultResolvers []artifactresolver.Resolver, stdout io.Writer) ([]artifactresolver.Locator, error) {
	if len(assetParams) == 0 {
		return nil, nil
	}

	var assets []artifactresolver.Locator
	assetErrors := make(map[artifactresolver.Locator]error)
	for _, currAsset := range assetParams {
		currAssetLocator, ok := ResolveAndVerify(
			currAsset,
			assetErrors,
			assetsDir,
			downloadsDir,
			defaultResolvers,
			osArch,
			stdout,
		)
		if !ok {
			continue
		}
		assets = append(assets, currAssetLocator)
	}
	SortLocators(assets)

	if len(assetErrors) == 0 {
		return assets, nil
	}

	// encountered errors: summarize and return
	var sortedKeys []artifactresolver.Locator
	for k := range assetErrors {
		sortedKeys = append(sortedKeys, k)
	}
	SortLocators(sortedKeys)

	errStringsParts := []string{fmt.Sprintf("failed to resolve %d asset(s):", len(assetErrors))}
	for _, k := range sortedKeys {
		errStringsParts = append(errStringsParts, assetErrors[k].Error())
	}
	return nil, errors.New(strings.Join(errStringsParts, "\n"+strings.Repeat(" ", IndentSpaces)))
}

func ResolveAndVerify(
	currArtifact artifactresolver.LocatorWithResolverParam,
	artifactErrors map[artifactresolver.Locator]error,
	dstBaseDir, downloadsDir string,
	defaultResolvers []artifactresolver.Resolver,
	osArch osarch.OSArch,
	stdout io.Writer) (currLocator artifactresolver.Locator, ok bool) {

	currLocator = currArtifact.LocatorWithChecksums.Locator
	currDstPath := filepath.Join(dstBaseDir, pathsinternal.PluginFileName(currLocator))

	if _, err := os.Stat(currDstPath); os.IsNotExist(err) {
		tgzDstPath := filepath.Join(downloadsDir, pathsinternal.PluginFileName(currLocator)+".tgz")
		if err := artifactresolver.ResolveArtifactTGZ(currArtifact, defaultResolvers, osArch, tgzDstPath, stdout); err != nil {
			artifactErrors[currLocator] = err
			return currLocator, false
		}

		if err := func() (rErr error) {
			pluginFile, err := os.OpenFile(currDstPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return errors.Wrapf(err, "failed to create file %s", currDstPath)
			}
			defer func() {
				if err := pluginFile.Close(); err != nil {
					rErr = errors.Wrapf(err, "failed to close file %s", currDstPath)
				}
			}()

			tgzFile, err := os.Open(tgzDstPath)
			if err != nil {
				return errors.Wrapf(err, "failed to open %s for reading", tgzDstPath)
			}
			defer func() {
				if err := tgzFile.Close(); err != nil {
					rErr = errors.Wrapf(err, "failed to close file %s", tgzDstPath)
				}
			}()

			if err := artifactresolver.CopySingleFileTGZContent(pluginFile, tgzFile); err != nil {
				return errors.Wrapf(err, "failed to copy file out of TGZ from %s to %s", tgzDstPath, currDstPath)
			}
			return nil
		}(); err != nil {
			artifactErrors[currLocator] = errors.Wrapf(err, "failed to extract artifact from archive into destination")
			return currLocator, false
		}
	}

	if wantChecksum, ok := currArtifact.LocatorWithChecksums.Checksums[osArch]; ok {
		gotChecksum, err := artifactresolver.SHA256ChecksumFile(currDstPath)
		if err != nil {
			artifactErrors[currLocator] = errors.Wrapf(err, "failed to compute checksum for plugin")
			return currLocator, false
		}
		if gotChecksum != wantChecksum {
			artifactErrors[currLocator] = errors.Errorf("failed to verify checksum for %s: want %s, got %s", currDstPath, wantChecksum, gotChecksum)
			return currLocator, false
		}
	}
	return currLocator, true
}

func SortLocators(locs []artifactresolver.Locator) {
	sort.Slice(locs, func(i, j int) bool {
		return locs[i].String() < locs[j].String()
	})
}

func Uniquify(in []string) []string {
	if in == nil {
		return nil
	}
	var out []string
	seen := make(map[string]struct{})
	for _, curr := range in {
		if _, ok := seen[curr]; ok {
			continue
		}
		out = append(out, curr)
		seen[curr] = struct{}{}
	}
	return out
}
