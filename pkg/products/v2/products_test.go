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

package products_test

import (
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/palantir/godel/pkg/products/v2"
)

func TestMain(t *testing.M) {
	// run once to ensure that godel is properly set up before tests run
	_, _ = products.List()
	os.Exit(t.Run())
}

func TestList(t *testing.T) {
	p, err := products.List()
	requireNoError(t, err, "failed: %v", err)
	if !reflect.DeepEqual([]string{"godel"}, p) {
		t.Fatalf("expected %v, got %v", []string{"godel"}, p)
	}
}

func TestBin(t *testing.T) {
	bin, err := products.Bin("godel")
	requireNoError(t, err, "failed: %v", err)
	cmd := exec.Command(bin, "version")
	output, err := cmd.CombinedOutput()
	requireNoError(t, err, "Command %v failed.\nOutput:\n%s\nError:\n%v", cmd.Args, string(output), err)

	if !strings.HasPrefix(string(output), "godel version") {
		t.Fatalf("output %q did not have prefix %q", string(output), "godel version")
	}
}

func TestDist(t *testing.T) {
	dist, err := products.Dist("godel")
	requireNoError(t, err, "failed: %v", err)
	if !strings.HasSuffix(dist, ".tgz") {
		t.Fatalf("output %q did not have suffix %q", dist, ".tgz")
	}
}

func requireNoError(t *testing.T, err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}
	t.Fatalf(msg, args...)
}
