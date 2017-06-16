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
	"io"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/binspec"
	"github.com/palantir/godel/apps/distgo/pkg/slsspec"
	"github.com/palantir/godel/apps/distgo/templating"
)

const slsDistInitSh = `#!/bin/bash
#
# Copyright 2015 Palantir Technologies
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# <http://www.apache.org/licenses/LICENSE-2.0>
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Daemonizes a service in accordance with the Linux Standard Base Core Specification 3.1
# https://refspecs.linuxbase.org/LSB_3.1.0/LSB-Core-generic/LSB-Core-generic/iniscrptact.html

set -uo pipefail

# uses SERVICE_HOME when set, else, traverse up two directories respecting symlinks
SERVICE_HOME=${SERVICE_HOME:-$(cd "$(dirname "$0")/../../" && pwd)}
cd "$SERVICE_HOME"

# either linux-amd64 or darwin-amd64
OS_ARCH="$(uname -s | awk '{print tolower($0)}')-amd64"

ACTION=$1
SERVICE="{{.ProductName}}"
SERVICE_CMD="$SERVICE_HOME/service/bin/$OS_ARCH/$SERVICE {{.Dist.ServiceArgs}}"
PIDFILE="var/run/$SERVICE.pid"

if [ -f service/bin/config.sh ]; then
    source service/bin/config.sh
fi

# Returns 0 if the service's pid is running
is_process_active() {
    local PID=$1
    if [ -z "$PID" ]; then
        return 1
    fi
    ps -o command "$PID" | grep -q "$SERVICE"
}

case $ACTION in
start)
    printf "%-50s" "Running '$SERVICE'..."

    if service/bin/init.sh status > /dev/null 2>&1; then
        printf "%s\n" "Already running ($(cat $PIDFILE))"
        exit 0
    fi

    # ensure log and pid directories exist
    mkdir -p "var/log"
    mkdir -p "var/run"

    PID=$($SERVICE_CMD > var/log/$SERVICE-startup.log 2>&1 & echo $!)
    echo $PID > $PIDFILE
    sleep 1
    if is_process_active $PID; then
        printf "%s\n" "Started ($PID)"
        exit 0
    else
        rm -f $PIDFILE
        printf "%s\n" "Failed"
        exit 1
    fi
;;
status)
    printf "%-50s" "Checking '$SERVICE'..."
    if [ -f $PIDFILE ]; then
        PID=$(cat $PIDFILE)
        if is_process_active $PID; then
            printf "%s\n" "Running ($PID)"
            exit 0
        fi
        printf "%s\n" "Process dead but pidfile exists"
        exit 1
    else
        printf "%s\n" "Service not running"
        exit 3
    fi
;;
stop)
    printf "%-50s" "Stopping '$SERVICE'..."

    if service/bin/init.sh status > /dev/null 2>&1; then
        PID=$(cat $PIDFILE)
        kill -s TERM $PID

        STOP_TIMEOUT=90
        COUNTER=0
        while is_process_active $PID && [ "$COUNTER" -lt "$STOP_TIMEOUT" ]; do
            sleep 1
            let COUNTER=COUNTER+1
            if [ $((COUNTER%5)) == 0 ]; then
                if [ "$COUNTER" -eq "5" ]; then
                    printf "\n" # first time get a new line to get off Stopping printf
                fi
                printf "%s\n" "Waiting for '$SERVICE' ($PID) to stop"
            fi
        done
        if is_process_active $PID; then
            printf "%-60s" "$SERVICE failed to stop, sending KILL signal..."
            kill -s KILL $PID
        fi
        printf "%s\n" "Stopped ($PID)"
    else
        printf "%s\n" "Service not running"
    fi
    rm -f $PIDFILE
    exit 0
;;
console)
    if service/bin/init.sh status > /dev/null 2>&1; then
        echo "Process is already running"
        exit 1
    fi
    trap "service/bin/init.sh stop > /dev/null 2>&1" TERM INT EXIT
    mkdir -p "$(dirname $PIDFILE)"

    $SERVICE_CMD &
    echo $! > $PIDFILE
    wait
;;
restart)
    service/bin/init.sh stop
    service/bin/init.sh start
;;
reload){{if .Dist.Reloadable}}
    printf "%-50s" "Reloading '$SERVICE'..."
    if service/bin/init.sh status > /dev/null 2>&1; then
        PID=$(cat $PIDFILE)
        if ! kill -s HUP $PID; then
            printf "%s\n" "Failed to send HUP ($PID)"
            exit 1
        fi
        printf "%s\n" "Reloaded ($PID)"
    else
        printf "%s\n" "Service not running"
        exit 7
    fi
{{else}}
    printf "%s\n" "'$SERVICE' does not support reload"
    exit 3
{{end}};;
*)
    echo "Usage: $0 status|start|stop|console|restart{{if .Dist.Reloadable}}|reload{{end}}"
    exit 1
esac
`

type slsDister params.SLSDistInfo

func (s *slsDister) NumArtifacts() int {
	return 1
}

func (s *slsDister) ArtifactPathsInOutputDir(buildSpec params.ProductBuildSpec) []string {
	values := slsspec.TemplateValues(buildSpec.ProductName, buildSpec.ProductVersion)
	return []string{slsspec.New().RootDirName(values) + ".sls.tgz"}
}

func (s *slsDister) Dist(buildSpecWithDeps params.ProductBuildSpecWithDeps, distCfg params.Dist, outputProductDir string, spec specdir.LayoutSpec, values specdir.TemplateValues, stdout io.Writer) (Packager, error) {
	buildSpec := buildSpecWithDeps.Spec
	outputSLSDir := path.Join(buildSpec.ProjectDir, distCfg.OutputDir, spec.RootDirName(values))

	// create init.sh and manifest.yml
	specDir, err := specdir.New(outputSLSDir, spec, values, specdir.Create)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create spec for %v", outputSLSDir)
	}

	if err := s.writeSLSManifest(buildSpec, distCfg, specDir); err != nil {
		return nil, err
	}

	if err := s.writeSLSInitSh(buildSpec, distCfg, specDir); err != nil {
		return nil, errors.Wrapf(err, "failed to write init.sh")
	}

	serviceBinDir := specDir.Path(slsspec.ServiceBin)
	binSpec := binspec.New(buildSpec.Build.OSArchs, buildSpec.ProductName)
	binSpecDir, err := specdir.New(serviceBinDir, binSpec, nil, specdir.Create)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create spec for directory %v", serviceBinDir)
	}

	if err := copyBuildArtifactsToBinDir(buildSpecWithDeps, binSpecDir); err != nil {
		return nil, errors.Wrapf(err, "failed to copy artifacts to service/bin dir")
	}

	return packager(func() error {
		if err := slsspec.Validate(outputProductDir, values, s.YMLValidationExclude); err != nil {
			return errors.Wrapf(err, "distribution directory failed SLS validation")
		}

		dstArtifactPath := FullArtifactsPaths(s, buildSpec, distCfg)[0]
		if err := singlePathTGZPackager(dstArtifactPath, outputProductDir).Package(); err != nil {
			return err
		}
		return nil
	}), nil
}

func (s *slsDister) DistPackageType() string {
	return "sls.tgz"
}

func (s *slsDister) writeSLSManifest(buildSpec params.ProductBuildSpec, distCfg params.Dist, specDir specdir.SpecDir) error {
	var manifestTemplateString string
	if s.ManifestTemplateFile != "" {
		manifestTemplateFilePath := path.Join(buildSpec.ProjectDir, s.ManifestTemplateFile)
		manifestBytes, err := ioutil.ReadFile(manifestTemplateFilePath)
		if err != nil {
			return errors.Wrapf(err, "failed to read manifest template file %s", manifestTemplateFilePath)
		}
		t := template.Must(template.New("manifest").Parse(string(manifestBytes)))
		manifestBuf := bytes.Buffer{}
		if err := t.Execute(&manifestBuf, templating.ConvertSpec(buildSpec, distCfg)); err != nil {
			return errors.Wrapf(err, "failed to execute template %v on spec %v", t, buildSpec)
		}
		manifestTemplateString = manifestBuf.String()
	} else {
		var err error
		manifestTemplateString, err = params.GetManifest(distCfg.Publish.GroupID, buildSpec.ProductName, buildSpec.ProductVersion, s.ProductType, s.ManifestExtensions)
		if err != nil {
			return errors.Wrapf(err, "failed to create manifest for SLS distribution")
		}
	}
	if err := ioutil.WriteFile(specDir.Path(slsspec.Manifest), []byte(manifestTemplateString), 0644); err != nil {
		return errors.Wrapf(err, "failed to write manifest")
	}
	return nil
}

func (s *slsDister) writeSLSInitSh(buildSpec params.ProductBuildSpec, distCfg params.Dist, specDir specdir.SpecDir) error {
	var initShTemplateBytes []byte
	if s.InitShTemplateFile != "" {
		initShTemplateFilePath := path.Join(buildSpec.ProjectDir, s.InitShTemplateFile)
		var err error
		initShTemplateBytes, err = ioutil.ReadFile(initShTemplateFilePath)
		if err != nil {
			return errors.Wrapf(err, "failed to read init.sh template file %v", initShTemplateFilePath)
		}
	} else {
		initShTemplateBytes = []byte(slsDistInitSh)
	}

	initShBuf := bytes.Buffer{}
	t := template.Must(template.New("init.sh").Parse(string(initShTemplateBytes)))
	if err := t.Execute(&initShBuf, templating.ConvertSpec(buildSpec, distCfg)); err != nil {
		return errors.Wrapf(err, "failed to execute template %v on template %v", t, buildSpec)
	}
	if err := ioutil.WriteFile(specDir.Path(slsspec.InitSh), initShBuf.Bytes(), 0755); err != nil {
		return errors.Wrapf(err, "failed to write init.sh")
	}
	return nil
}
