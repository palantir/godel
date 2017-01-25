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
	"fmt"
	"io"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
)

type DockerPublishInfo struct{}

func (d DockerPublishInfo) Publish(buildSpec params.ProductBuildSpec, paths ProductPaths, stdout io.Writer) (string, error) {
	dockerPush := exec.Command("docker", "push", paths.artifactPath)
	if output, err := dockerPush.CombinedOutput(); err != nil {
		fmt.Println(string(output))
		return "", errors.Wrapf(err, "docker push failed for image %s", paths.artifactPath)
	}
	return paths.artifactPath, nil
}
