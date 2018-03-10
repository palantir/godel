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

package distgo

import (
	"fmt"
	"io"
)

const dryRunPrefix = "[DRY RUN]"

func DryRunPrintln(w io.Writer, msg string) {
	DryRunPrint(w, msg+"\n")
}

func DryRunPrint(w io.Writer, msg string) {
	fmt.Fprint(w, dryRunPrefix+" ", msg)
}

func PrintlnOrDryRunPrintln(w io.Writer, msg string, dryRun bool) {
	PrintOrDryRunPrint(w, msg+"\n", dryRun)
}

func PrintOrDryRunPrint(w io.Writer, msg string, dryRun bool) {
	if dryRun {
		DryRunPrint(w, msg)
	} else {
		fmt.Fprint(w, msg)
	}
}
