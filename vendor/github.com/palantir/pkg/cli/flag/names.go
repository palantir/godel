// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag

import (
	"strings"
)

func WithPrefix(name string) string {
	if len(name) == 1 {
		return "-" + name
	}
	return "--" + name
}

func placeholderOrDefault(maybePlaceholder, name string) string {
	if maybePlaceholder != "" {
		return maybePlaceholder
	}
	return defaultPlaceholder(name)
}

func defaultPlaceholder(name string) string {
	name = strings.ToUpper(name)
	name = strings.Replace(name, "-", "_", -1)
	return name
}
