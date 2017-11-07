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

package installupdate_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/palantir/godel/framework/builtintasks/installupdate"
)

func TestUpdateUsesPathFromExecutable(t *testing.T) {
	err := installupdate.Update("invalid-script", ioutil.Discard)
	assert.EqualError(t, err, "wrapper script invalid-script is not in a valid location: godelw does not exist")
}
