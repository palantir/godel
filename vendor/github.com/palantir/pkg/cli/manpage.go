// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

type Manpage struct {
	Source     string // e.g. Linux
	Manual     string // e.g. Linux Programmer's Manual
	BugTracker string
	SeeAlso    []ManpageRef
}

type ManpageRef struct {
	Name    string
	Section uint8
}
