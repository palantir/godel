// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package completion

// Completion script sees these codes and uses compgen to complete the right thing.
const (
	FilepathCode       = 99
	DirectoryCode      = 98
	CustomFilepathCode = 97
)
