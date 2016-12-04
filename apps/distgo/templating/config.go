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

package templating

import (
	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/git"
)

type Config struct {
	// {{.ProductName}}
	ProductName string

	// {{.ProductVersion}}
	ProductVersion string

	// {{.VersionInfo.Version}} is a string
	// {{.VersionInfo.Branch}} is a string
	// {{.VersionInfo.Revision}} is a string
	VersionInfo git.ProjectInfo

	// {{.Dist.Type}} is "sls" or "bin" or "rpm"
	// {{.Dist}} is a fully customizable map
	Dist params.DistInfo

	// {{.Publish.GroupID}} is a string
	// {{.Publish.Metadata}} is a map of string to string
	// {{.Publish.Tags}} is a slice of strings
	Publish params.Publish
}
