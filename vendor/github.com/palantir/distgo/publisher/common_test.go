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

package publisher_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/publisher"
)

func TestSetConfigValueString(t *testing.T) {
	const valProvidedInFlag = "specified-as-flag"
	flag := distgo.PublisherFlag{
		Name: "foo",
		Type: distgo.StringFlag,
	}
	flagVals := map[distgo.PublisherFlagName]interface{}{
		flag.Name: valProvidedInFlag,
	}
	cfg := struct {
		FooVal string
	}{
		FooVal: "original-value",
	}

	err := publisher.SetConfigValue(flagVals, flag, &cfg.FooVal)
	require.NoError(t, err)

	assert.Equal(t, valProvidedInFlag, cfg.FooVal)
}

func TestSetConfigValueBool(t *testing.T) {
	const valProvidedInFlag = true
	flag := distgo.PublisherFlag{
		Name: "foo",
		Type: distgo.BoolFlag,
	}
	flagVals := map[distgo.PublisherFlagName]interface{}{
		flag.Name: valProvidedInFlag,
	}
	cfg := struct {
		FooVal bool
	}{}

	err := publisher.SetConfigValue(flagVals, flag, &cfg.FooVal)
	require.NoError(t, err)

	assert.Equal(t, valProvidedInFlag, cfg.FooVal)
}

func TestSetConfigValueDoesNotChangeIfFlagNotSet(t *testing.T) {
	const originalValue = "original-value"
	flag := distgo.PublisherFlag{
		Name: "foo",
		Type: distgo.StringFlag,
	}
	flagVals := map[distgo.PublisherFlagName]interface{}{}
	cfg := struct {
		FooVal string
	}{
		FooVal: originalValue,
	}

	err := publisher.SetConfigValue(flagVals, flag, &cfg.FooVal)
	require.NoError(t, err)

	assert.Equal(t, originalValue, cfg.FooVal)
}

func TestSetConfigValueFailsIfTypeDoesNotMatch(t *testing.T) {
	const valProvidedInFlag = "specified-as-flag"
	flag := distgo.PublisherFlag{
		Name: "foo",
		Type: distgo.BoolFlag,
	}
	flagVals := map[distgo.PublisherFlagName]interface{}{
		flag.Name: valProvidedInFlag,
	}
	cfg := struct {
		FooVal bool
	}{}

	err := publisher.SetConfigValue(flagVals, flag, &cfg.FooVal)
	assert.EqualError(t, err, "flagValType is not assignable to the provided configValPtr: !string.AssignableTo(bool)")
}

func TestSetConfigValueFailsIfProvidedValNotPointer(t *testing.T) {
	const valProvidedInFlag = "specified-as-flag"
	flag := distgo.PublisherFlag{
		Name: "foo",
		Type: distgo.StringFlag,
	}
	flagVals := map[distgo.PublisherFlagName]interface{}{
		flag.Name: valProvidedInFlag,
	}
	cfg := struct {
		FooVal string
	}{
		FooVal: "original-value",
	}

	err := publisher.SetConfigValue(flagVals, flag, cfg.FooVal)
	assert.EqualError(t, err, `configValPtr type "string" is not a pointer type`)
}
