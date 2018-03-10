// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package completion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirRaw(t *testing.T) {
	assert.Equal(t, "", dirRaw(""))
	assert.Equal(t, "./", dirRaw("./"))
	assert.Equal(t, "test/", dirRaw("test/"))
	assert.Equal(t, "test/", dirRaw("test/ing"))
	assert.Equal(t, "test/ing/", dirRaw("test/ing/"))
	assert.Equal(t, "test/ing/", dirRaw("test/ing/still"))
	assert.Equal(t, "/", dirRaw("/"))
	assert.Equal(t, "/", dirRaw("/test"))
	assert.Equal(t, "/test/", dirRaw("/test/"))
	assert.Equal(t, "/test/", dirRaw("/test/ing"))
	assert.Equal(t, "/test/ing/", dirRaw("/test/ing/"))
	assert.Equal(t, "/test/ing/", dirRaw("/test/ing/still"))
}
