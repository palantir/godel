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

package main

import (
	"fmt"
	"os"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/amalgomate/amalgomated"

	"github.com/palantir/godel/apps/okgo"
)

func main() {
	if err := dirs.SetGoEnvVariables(); err != nil {
		fmt.Println("Failed to set Go environment variables:", err)
		os.Exit(1)
	}
	os.Exit(okgo.RunApp(os.Args, amalgomated.SelfProxyCmderSupplier()))
}
