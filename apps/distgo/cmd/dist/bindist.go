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

package dist

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/binspec"
	"github.com/palantir/godel/apps/distgo/templating"
)

const binDistInitSh = `#!/bin/bash
set -euo pipefail

BIN_DIR="$(cd "$(dirname "$0")" && pwd)"

# determine OS
OS=""
case "$(uname)" in
  Darwin*)
    OS=darwin
    ;;
  Linux*)
    OS=linux
    ;;
  *)
    echo "Unsupported operating system: $(uname)"
    exit 1
    ;;
esac

# determine executable location based on OS
CMD=$BIN_DIR/$OS-amd64/{{.ProductName}}

# verify that executable exists
if [ ! -e "$CMD" ]; then
    echo "Executable $CMD does not exist"
    exit 1
fi

# invoke appropriate executable
$CMD "$@"
`

func binDist(buildSpecWithDeps params.ProductBuildSpecWithDeps, distCfg params.Dist, outputProductDir string, stdout io.Writer) (Packager, error) {
	buildSpec := buildSpecWithDeps.Spec
	binDistInfo, ok := distCfg.Info.(*params.BinDistInfo)
	if !ok {
		binDistInfo = &params.BinDistInfo{}
		distCfg.Info = binDistInfo
	}
	fmt.Fprintf(stdout, "Creating bin distribution for %v at %v\n", buildSpecWithDeps.Spec.ProductName, ArtifactPath(buildSpec, distCfg))

	binSpec := binspec.New(buildSpec.Build.OSArchs, buildSpec.ProductName)
	binDir := path.Join(outputProductDir, "bin")
	binSpecDir, err := specdir.New(binDir, binSpec, nil, specdir.Create)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create directory structure for %v", binDir)
	}
	if err := copyBuildArtifactsToBinDir(buildSpecWithDeps, binSpecDir); err != nil {
		return nil, errors.Wrapf(err, "failed to copy artifacts to bin dir")
	}

	if !binDistInfo.OmitInitSh {
		var initShTemplateBytes []byte
		if binDistInfo.InitShTemplateFile != "" {
			initShTemplateFilePath := path.Join(buildSpec.ProjectDir, binDistInfo.InitShTemplateFile)
			var err error
			initShTemplateBytes, err = ioutil.ReadFile(initShTemplateFilePath)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to read init.sh template file %v", initShTemplateFilePath)
			}
		} else {
			initShTemplateBytes = []byte(binDistInitSh)
		}
		initShBuf := bytes.Buffer{}
		t := template.Must(template.New("init.sh").Parse(string(initShTemplateBytes)))
		if err := t.Execute(&initShBuf, templating.ConvertSpec(buildSpec, distCfg)); err != nil {
			return nil, errors.Wrapf(err, "failed to execute template %v on template %v", t, buildSpec)
		}
		if err := ioutil.WriteFile(path.Join(binDir, buildSpec.ProductName+".sh"), initShBuf.Bytes(), 0755); err != nil {
			return nil, errors.Wrapf(err, "failed to write init.sh")
		}
	}

	return tgzPackager(buildSpec, distCfg, outputProductDir), nil
}
