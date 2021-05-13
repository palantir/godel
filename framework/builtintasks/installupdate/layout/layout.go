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

package layout

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	s "github.com/palantir/pkg/specdir"
)

// AllPaths returns a map that contains all of the paths in the provided directory. The paths are relative to the
// directory. The boolean key is true if the path is a directory, false otherwise.
func AllPaths(dir string) (map[string]bool, error) {
	m := make(map[string]bool)
	return m, allPaths(m, nil, dir)
}

func allPaths(paths map[string]bool, pathStack []string, dir string) error {
	fis, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		p := filepath.Join(dir, fi.Name())
		key := filepath.Join(strings.Join(pathStack, "/"), fi.Name())
		// record current path
		paths[key] = fi.IsDir()
		if fi.IsDir() {
			// if path is directory, recursively record paths
			if err := allPaths(paths, append(pathStack, fi.Name()), p); err != nil {
				return err
			}
		}
	}
	return nil
}

func GodelDistLayout(version string, mode s.Mode) (s.SpecDir, error) {
	rootDir, err := GodelHomePath()
	if err != nil {
		return nil, err
	}

	values := s.TemplateValues{
		godelHomeTemplate: filepath.Base(rootDir),
		versionTemplate:   version,
	}
	for key, value := range AppSpecTemplate(version) {
		values[key] = value
	}

	return s.New(rootDir, godelHomeSpec([]s.FileNodeProvider{AppSpec()}), values, mode)
}

const (
	godelHomeTemplate = "godel-home"
	godelHomeEnvVar   = "GODEL_HOME"
	defaultGodelHome  = ".godel"

	AssetsDir    = "assets"
	CacheDir     = "cache"
	DistsDir     = "dists"
	ConfigsDir   = "configs"
	DownloadsDir = "downloads"
	PluginsDir   = "plugins"
)

// GodelHomePath returns the path to the gödel home directory. If $GODEL_HOME is set as an environment variable, that
// value is used. Otherwise, the value is "$HOME/{{defaultGodelHome}}"
func GodelHomePath() (string, error) {
	// check the environment variable
	if godelHomeDir := os.Getenv(godelHomeEnvVar); godelHomeDir != "" {
		return godelHomeDir, nil
	}
	// if not present, create from home directory
	if userHomeDir := os.Getenv("HOME"); userHomeDir != "" {
		return filepath.Join(userHomeDir, defaultGodelHome), nil
	}
	return "", fmt.Errorf("failed to get %s home directory", AppName)
}

func GodelHomeSpecDir(mode s.Mode) (s.SpecDir, error) {
	rootDir, err := GodelHomePath()
	if err != nil {
		return nil, err
	}
	values := s.TemplateValues{
		godelHomeTemplate: filepath.Base(rootDir),
	}

	return s.New(rootDir, GodelHomeSpec(), values, mode)
}

func GodelHomeSpec() s.LayoutSpec {
	return godelHomeSpec(nil)
}

func godelHomeSpec(providers []s.FileNodeProvider) s.LayoutSpec {
	return s.NewLayoutSpec(
		s.Dir(s.TemplateName(godelHomeTemplate), "",
			s.Dir(s.LiteralName("assets"), AssetsDir),
			s.Dir(s.LiteralName("cache"), CacheDir),
			s.Dir(s.LiteralName("configs"), ConfigsDir),
			s.Dir(s.LiteralName("dists"), DistsDir, providers...),
			s.Dir(s.LiteralName("downloads"), DownloadsDir),
			s.Dir(s.LiteralName("plugins"), PluginsDir),
		),
		true,
	)
}

const (
	AppName         = "godel"
	AppDir          = "gödel-app"
	AppExecutable   = "app-executable"
	osTemplate      = "os"
	archTemplate    = "arch"
	versionTemplate = "version"
)

func AppSpecDir(rootDir, version string) (s.SpecDir, error) {
	return s.New(rootDir, AppSpec(), AppSpecTemplate(version), s.Validate)
}

func AppSpec() s.LayoutSpec {
	return s.NewLayoutSpec(
		s.Dir(s.CompositeName(s.LiteralName(AppName+"-"), s.TemplateName(versionTemplate)), AppDir,
			s.Dir(s.LiteralName("bin"), "",
				s.Dir(s.CompositeName(s.TemplateName(osTemplate), s.LiteralName("-"), s.TemplateName(archTemplate)), "",
					s.File(s.LiteralName(AppName), AppExecutable),
				),
			),
			WrapperSpec(),
		),
		true,
	)
}

func AppSpecTemplate(version string) s.TemplateValues {
	return s.TemplateValues{
		versionTemplate: version,
		osTemplate:      runtime.GOOS,
		archTemplate:    runtime.GOARCH,
	}
}

const (
	WrapperDir        = "wrapper-dir"
	WrapperScriptFile = "wrapper-script"
	WrapperAppDir     = "wrapper-app"
	WrapperConfigDir  = "config"
	WrapperName       = "godelw"
)

func WrapperSpec() s.LayoutSpec {
	return s.NewLayoutSpec(
		s.Dir(s.LiteralName("wrapper"), WrapperDir,
			s.File(s.LiteralName(WrapperName), WrapperScriptFile),
			s.Dir(s.LiteralName(AppName), WrapperAppDir,
				s.Dir(s.LiteralName(WrapperConfigDir), WrapperConfigDir),
			),
		),
		false,
	)
}
