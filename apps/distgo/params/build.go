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

package params

import (
	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

type Build struct {
	// Skip specifies whether the build step should be skipped entirely. Its primary use is for products that handle
	// their own build logic in the "dist" step ("dist-only" products).
	Skip bool

	// Script is the content of a script that is written to file a file and run before this product is built. The
	// contents of this value are written to a file with a header `#!/bin/bash` and executed. The script process
	// inherits the environment variables of the Go process and also has the following environment variables
	// defined:
	//
	//   PROJECT_DIR: the root directory of project
	//   PRODUCT: product name,
	//   VERSION: product version
	//   IS_SNAPSHOT: 1 if the version contains a git hash as part of the string, 0 otherwise
	Script string

	// MainPkg is the location of the main package for the product relative to the root directory. For example,
	// "./distgo/main".
	MainPkg string

	// OutputDir is the directory to which the executable is written.
	OutputDir string

	// BuildArgsScript is the content of a script that is written to a file and run before this product is built
	// to provide supplemental build arguments for the product. The contents of this value are written to a file
	// with a header `#!/bin/bash` and executed. The script process inherits the environment variables of the Go
	// process. Each line of output of the script is provided to the "build" command as a separate argument. For
	// example, the following script would add the arguments "-ldflags" "-X" "main.year=$YEAR" to the build command:
	//
	//   build-args-script: |
	//                      YEAR=$(date +%Y)
	//                      echo "-ldflags"
	//                      echo "-X"
	//                      echo "main.year=$YEAR"
	BuildArgsScript string

	// VersionVar is the path to a variable that is set with the version information for the build. For example,
	// "github.com/palantir/godel/cmd/godel.Version". If specified, it is provided to the "build" command as an
	// ldflag.
	VersionVar string

	// Environment specifies values for the environment variables that should be set for the build. For example,
	// the following sets CGO to false:
	//
	//   environment:
	//     CGO_ENABLED: "0"
	Environment map[string]string

	// OSArchs specifies the GOOS and GOARCH pairs for which the product is built. If blank, defaults to the GOOS
	// and GOARCH of the host system at runtime.
	OSArchs []osarch.OSArch
}
