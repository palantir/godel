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

package testfuncs

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/palantir/distgo/dister/disterfactory"
	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/config"
	"github.com/palantir/distgo/dockerbuilder/dockerbuilderfactory"
	"github.com/palantir/distgo/projectversioner/projectversionerfactory"
	"github.com/palantir/distgo/publisher/publisherfactory"
)

func NewProjectParam(t *testing.T, projectConfig config.ProjectConfig, projectDir, failMsg string) distgo.ProjectParam {
	projectVersionerFactory, err := projectversionerfactory.New(nil, nil)
	require.NoError(t, err, failMsg)
	disterFactory, err := disterfactory.New(nil, nil)
	require.NoError(t, err, failMsg)
	defaultDisterCfg, err := disterfactory.DefaultConfig()
	require.NoError(t, err, failMsg)
	dockerBuilderFactory, err := dockerbuilderfactory.New(nil, nil)
	require.NoError(t, err, failMsg)
	publisherFactory, err := publisherfactory.New(nil, nil)
	require.NoError(t, err, failMsg)

	projectParam, err := projectConfig.ToParam(projectDir, projectVersionerFactory, disterFactory, defaultDisterCfg, dockerBuilderFactory, publisherFactory)
	require.NoError(t, err, failMsg)
	return projectParam
}

func NewProjectParamReturnError(t *testing.T, projectConfig config.ProjectConfig, projectDir, failMsg string) (distgo.ProjectParam, error) {
	projectVersionerFactory, err := projectversionerfactory.New(nil, nil)
	require.NoError(t, err, failMsg)
	disterFactory, err := disterfactory.New(nil, nil)
	require.NoError(t, err, failMsg)
	defaultDisterCfg, err := disterfactory.DefaultConfig()
	require.NoError(t, err, failMsg)
	dockerBuilderFactory, err := dockerbuilderfactory.New(nil, nil)
	require.NoError(t, err, failMsg)
	publisherFactory, err := publisherfactory.New(nil, nil)
	require.NoError(t, err, failMsg)

	return projectConfig.ToParam(projectDir, projectVersionerFactory, disterFactory, defaultDisterCfg, dockerBuilderFactory, publisherFactory)
}
