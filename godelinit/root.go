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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/google/go-github/github"
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/godel/framework/builtintasks/installupdate"
	"github.com/palantir/godel/framework/builtintasks/installupdate/layout"
	"github.com/palantir/godel/godelgetter"
)

func rootCmd() *cobra.Command {
	const (
		versionFlagName     = "version"
		ignoreCacheFlagName = "ignore-cache"
	)

	var versionFlag string
	var ignoreCacheFlag bool

	cmd := &cobra.Command{
		Use:   "godelinit",
		Short: "Add latest version of gödel to a project",
		Long: `godelinit adds godel to a project by adding the godelw script and godel configuration directory to it.
The default behavior adds the newest release of godel on GitHub (https://github.com/palantir/godel/releases)
to the project. If a specific version of godel is desired, it can be specified using the '--version' flag.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return errors.Wrapf(err, "failed to determine working directory")
			}
			latestVersion := versionFlag
			if latestVersion == "" {
				latestVersion, err = latestGodelVersion(ignoreCacheFlag)
				if err != nil {
					return err
				}
			}
			pkgPath, checksum, err := downloadedTGZForVersion(latestVersion)
			if err != nil {
				pkgPath = fmt.Sprintf("https://palantir.bintray.com/releases/com/palantir/godel/godel/%s/godel-%s.tgz", latestVersion, latestVersion)
				checksum = ""
			}
			if err := installupdate.NewInstall(wd, godelgetter.NewPkgSrc(pkgPath, checksum), cmd.OutOrStdout()); err != nil {
				return err
			}
			// update godel.properties with checksum
			if checksum != "" {
				if err := updateChecksum(wd, checksum); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&versionFlag, versionFlagName, "", "version to install (if unspecified, latest is used)")
	cmd.Flags().BoolVar(&ignoreCacheFlag, ignoreCacheFlagName, false, "ignore cache when determining latest version")

	return cmd
}

func updateChecksum(projectDir, checksum string) error {
	wrapperSpec, err := specdir.New(projectDir, layout.WrapperSpec(), nil, specdir.Validate)
	if err != nil {
		return errors.Wrapf(err, "unable to create wrapper spec")
	}
	cfgDirPath := wrapperSpec.Path(layout.WrapperConfigDir)
	return installupdate.SetGodelChecksum(cfgDirPath, checksum)
}

func downloadedTGZForVersion(version string) (string, string, error) {
	gödelHomeSpecDir, err := layout.GodelHomeSpecDir(specdir.Create)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to create SpecDir for gödel home")
	}
	downloadsDirPath := gödelHomeSpecDir.Path(layout.DownloadsDir)
	downloadedTGZ := path.Join(downloadsDirPath, fmt.Sprintf("%s-%s.tgz", layout.AppName, version))
	if _, err := os.Stat(downloadedTGZ); err != nil {
		return "", "", errors.Wrapf(err, "failed to stat downloaded TGZ file")
	}
	checksum, err := layout.Checksum(downloadedTGZ)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to compute checksum")
	}
	return downloadedTGZ, checksum, nil
}

func latestGodelVersion(ignoreCache bool) (string, error) {
	if !ignoreCache {
		versionCfg, err := readLatestCachedVersion()
		if err == nil && storedLatestVersionValid(versionCfg) {
			return versionCfg.LatestVersion, nil
		}
	}
	client := github.NewClient(http.DefaultClient)
	rel, _, err := client.Repositories.GetLatestRelease(context.Background(), "palantir", "godel")
	if err != nil {
		return "", errors.Wrap(err, "failed to determine latest release")
	}
	latestVersion := *rel.TagName
	if err := writeLatestCachedVersion(latestVersion); err != nil {
		return "", errors.Wrapf(err, "failed to write latest version to cache")
	}
	return latestVersion, nil
}

const latestVersionFileName = "latest-version.json"

func readLatestCachedVersion() (versionConfig, error) {
	gödelHomeSpecDir, err := layout.GodelHomeSpecDir(specdir.Create)
	if err != nil {
		return versionConfig{}, errors.Wrapf(err, "failed to create SpecDir for gödel home")
	}
	cacheDirPath := gödelHomeSpecDir.Path(layout.CacheDir)
	latestVersionFile := path.Join(cacheDirPath, latestVersionFileName)

	bytes, err := ioutil.ReadFile(latestVersionFile)
	if err != nil {
		return versionConfig{}, errors.Wrapf(err, "failed to read version file")
	}
	var versionCfg versionConfig
	if err := json.Unmarshal(bytes, &versionCfg); err != nil {
		return versionConfig{}, errors.Wrapf(err, "failed to unmarshal version file")
	}
	return versionCfg, nil
}

func writeLatestCachedVersion(version string) error {
	gödelHomeSpecDir, err := layout.GodelHomeSpecDir(specdir.Create)
	if err != nil {
		return errors.Wrapf(err, "failed to create SpecDir for gödel home")
	}
	cacheDirPath := gödelHomeSpecDir.Path(layout.CacheDir)
	latestVersionFile := path.Join(cacheDirPath, latestVersionFileName)

	bytes, err := json.Marshal(versionConfig{
		LatestVersion: version,
		Timestamp:     time.Now().Unix(),
	})
	if err != nil {
		return errors.Wrapf(err, "failed to marshal version config as JSON")
	}

	if err := ioutil.WriteFile(latestVersionFile, bytes, 0644); err != nil {
		return errors.Wrap(err, "failed to write version file")
	}
	return nil
}

func storedLatestVersionValid(cfg versionConfig) bool {
	storedTime := time.Unix(cfg.Timestamp, 0)
	return !storedTime.Before(time.Now().Add(-1 * time.Hour))
}

type versionConfig struct {
	LatestVersion string `json:"latestVersion"`
	Timestamp     int64  `json:"timestamp"`
}
