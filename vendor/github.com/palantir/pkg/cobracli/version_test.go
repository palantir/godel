// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cobracli_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/pkg/cobracli"
)

func TestVersionCmd(t *testing.T) {
	outBuf := &bytes.Buffer{}
	cmd := cobracli.VersionCmd("foo-app", "1.0.0")
	cmd.SetOutput(outBuf)
	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "foo-app version 1.0.0\n", outBuf.String())
}
