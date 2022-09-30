// Copyright (c) 2019 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build module
// +build module

// This file exists only to smooth the transition for modules. Having this file makes it such that other modules that
// consume this module will not have import path conflicts caused by github.com/palantir/pkg.
package main

import (
	_ "github.com/palantir/pkg"
)
