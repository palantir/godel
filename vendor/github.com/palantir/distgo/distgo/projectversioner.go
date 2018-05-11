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

package distgo

type ProjectVersioner interface {
	// TypeName returns the type of this project versioner.
	TypeName() (string, error)

	// ProjectVersion returns the string that should be used as the version for the project in the given directory.
	ProjectVersion(projectDir string) (string, error)
}

type ProjectVersionerFactory interface {
	Types() []string
	NewProjectVersioner(typeName string, cfgYMLBytes []byte) (ProjectVersioner, error)
	ConfigUpgrader(typeName string) (ConfigUpgrader, error)
}
