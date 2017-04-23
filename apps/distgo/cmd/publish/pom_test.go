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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/templating"
)

func TestGeneratePOM(t *testing.T) {
	for i, currCase := range []struct {
		name       string
		cfg        templating.Config
		distType   string
		classifier string
		want       string
	}{
		{
			"Configuration",
			templating.Config{
				ProductName:    "foo",
				ProductVersion: "1.0.0",
				Publish: params.Publish{
					GroupID: "com.org.group",
				},
			},
			"tgz",
			"",
			`<?xml version="1.0" encoding="UTF-8"?>
<project xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd" xmlns="http://maven.apache.org/POM/4.0.0"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<modelVersion>4.0.0</modelVersion>
<groupId>com.org.group</groupId>
<artifactId>foo</artifactId>
<version>1.0.0</version>
<packaging>tgz</packaging>
</project>
`,
		},
	} {
		bytes, err := generatePOM(templating.Config{
			ProductName:    "foo",
			ProductVersion: "1.0.0",
			Publish: params.Publish{
				GroupID: "com.org.group",
			},
		}, "tgz")
		require.NoError(t, err)
		assert.Equal(t, currCase.want, string(bytes), "Case %d: %s", i, currCase.name)
	}
}
