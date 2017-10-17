Summary
-------
The `github.com/palantir/godel/pkg/products` package can be used to write integration tests that run against the build
artifacts of a product and test tags can be used to define test sets for integration tests.

Tutorial start state
--------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist.yml` is configured to build `echgo`
* Project is tagged as 0.0.1
* `godel/config/dist.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* Go files have license headers
* `godel/config/generate.yml` is configured to generate string function
* `godel/config/exclude.yml` is configured to ignore all `.+_string.go` files

([Link](https://github.com/nmiyake/echgo/tree/1982133dbe7c811f1e2d71f4dcc25ff20f84146a))

Write tests that run using build artifacts
------------------------------------------

`echgo` currently has unit tests that test the contracts of the `echgo` package. Unit tests are a great way to test the
API contracts of packages, and in an ideal world all of the packages for a project having tests that verify the package
APIs would be sufficient to ensure the correctness of an entire program.

However, in many cases there exists behavior that can only be tested in a true end-to-end workflow. For example, `echgo`
currently has some logic in its `main.go` file that parses the command-line flags, determines what functions to call
based on flags and ultimately prints the output to the console. If we want to test things such as what happens when
invalid values are supplied as flags, how multiple command-line arguments are parsed or the exit codes of the program,
there is not a straightforward way to write that test.

The `github.com/palantir/godel/pkg/products` packages provides functionality that makes it easy to write such tests for
projects that use gödel to build their products. The `products` package provides functions that ensure that specified
products are built using the build configuration defined for the product and provides a path to the built executable
that can be used for testing.

We need to add `github.com/palantir/godel/pkg/products` as a vendored dependency for the project. Start by cloning the
gödel project:

```
➜ mkdir -p $GOPATH/github.com/palantir && cd $_
➜ git clone https://github.com/palantir/godel.git
Cloning into 'godel'...
remote: Counting objects: 3210, done.
remote: Total 3210 (delta 0), reused 0 (delta 0), pack-reused 3209
Receiving objects: 100% (3210/3210), 3.65 MiB | 1.92 MiB/s, done.
Resolving deltas: 100% (1382/1382), done.
Checking connectivity... done.
```

There are multiple different ways to vendor dependencies. For the purposes of this tutorial, we will forego formal
vendoring and vendor the dependency manually.

```
➜ cd $GOPATH/src/github.com/nmiyake/echgo
➜ mkdir -p vendor/github.com/palantir/godel/pkg/products
➜ cp $GOPATH/src/github.com/palantir/godel/pkg/products/* vendor/github.com/palantir/godel/pkg/products/
```

Run the following to define a test that tests the behavior of invoking `echgo` with an invalid echo type and run the
test (this test is still in the iteration phase, so it simply prints the result of the output rather than asserting
against it):

```
➜ mkdir -p integration_test
➜ echo '// Copyright (c) 2017 Author Name
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

package integration_test

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/palantir/godel/pkg/products"
)

func TestInvalidType(t *testing.T) {
	echgoPath, err := products.Bin("echgo")
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(echgoPath, "-type", "invalid", "foo")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("cmd %v failed with error %v. Output: %s", cmd.Args, err, string(output))
	}
	fmt.Printf("%q", string(output))
	fmt.Println()
}' > integration_test/integration_test.go
➜ go test -v ./integration_test
=== RUN   TestInvalidType
"invalid echo type: invalid\n"
--- PASS: TestInvalidType (1.43s)
PASS
ok  	github.com/nmiyake/echgo/integration_test	1.566s
```

The `products.Bin("echgo")` call uses gödel to build the `echgo` product (if needed) and returns a path to the binary
that was built. Because this is a path to a valid binary, `exec.Command` can be use to invoke it. This allows the test
to specify arguments, hook up input/output streams, check error values and assert various behavior.

In this case, the output seems reasonable -- it prints `invalid echo type: invalid\n`. However, note that the error was
`nil` -- this is a bug. If the specified echo type was invalid, then the program should return with a non-zero exit
code, which should cause `cmd.CombinedOutput` to return an error.

Fix the bug by updating `main.go` and then re-run the test:

```
➜ echo '// Copyright (c) 2017 Author Name
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
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nmiyake/echgo/echo"
)

var version = "none"

func main() {
	versionVar := flag.Bool("version", false, "print version")
	typeVar := flag.String("type", echo.Simple.String(), "type of echo")
	flag.Parse()
	if *versionVar {
		fmt.Println("echgo version:", version)
		return
	}
	typ, err := echo.TypeFrom(*typeVar)
	if err != nil {
		fmt.Println("invalid echo type:", *typeVar)
		os.Exit(1)
	}
	echoer := echo.NewEchoer(typ)
	fmt.Println(echoer.Echo(strings.Join(flag.Args(), " ")))
}' > main.go
➜ go test -v ./integration_test
=== RUN   TestInvalidType
"invalid echo type: invalid\n"
--- FAIL: TestInvalidType (1.38s)
	integration_test.go:33: cmd [/Volumes/git/go/src/github.com/nmiyake/echgo/build/0.0.2-4-g1982133.dirty/darwin-amd64/echgo -type invalid foo] failed with error exit status 1. Output: invalid echo type: invalid
FAIL
exit status 1
FAIL	github.com/nmiyake/echgo/integration_test	1.564s
```

We can see that the test now fails as expected. Since this is the expected behavior, update the test to pass when this
happens and run the test again:

```
➜ echo '// Copyright (c) 2017 Author Name
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

package integration_test

import (
	"os/exec"
	"testing"

	"github.com/palantir/godel/pkg/products"
)

func TestInvalidType(t *testing.T) {
	echgoPath, err := products.Bin("echgo")
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(echgoPath, "-type", "invalid", "foo")
	output, err := cmd.CombinedOutput()
	gotOutput := string(output)
	if err == nil {
		t.Errorf("expected command %v to fail. Output: %s", cmd.Args, gotOutput)
	}
	wantOutput := "invalid echo type: invalid\\n"
	if wantOutput != gotOutput {
		t.Errorf("invalid output: want %q, got %q", wantOutput, gotOutput)
	}
	wantErr := "exit status 1"
	gotErr := err.Error()
	if wantErr != gotErr {
		t.Errorf("invalid error output: want %q, got %q", wantErr, gotErr)
	}
}' > integration_test/integration_test.go
➜ go test -v ./integration_test
=== RUN   TestInvalidType
--- PASS: TestInvalidType (0.51s)
PASS
ok  	github.com/nmiyake/echgo/integration_test	0.707s
```

We can see that the test now passes. The test will now run when `./godelw test` is invoked.

One thing to note about this construction is that the `go build` and `go install` commands will currently not work for
`integration_test` because the directory contains only tests:

```
➜ go build ./integration_test
go build github.com/nmiyake/echgo/integration_test: no non-test Go files in /Volumes/git/go/src/github.com/nmiyake/echgo/integration_test
```

We can work around this by adding a `doc.go` file to the directory to act as a placeholder:

```
➜ echo '// Copyright (c) 2017 Author Name
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

// Package integration contains integration tests.
package integration' > integration_test/doc.go
```

Verify that building the directory no longer fails:

```
➜ go build ./integration_test
```

Run `./godelw test` to verify that this test is run:

```
➜ ./godelw test
ok  	github.com/nmiyake/echgo                 	0.189s [no tests to run]
ok  	github.com/nmiyake/echgo/echo            	0.201s
ok  	github.com/nmiyake/echgo/generator       	0.179s [no tests to run]
ok  	github.com/nmiyake/echgo/integration_test	0.614s
```

The configuration in `godel/config/test.yml` can be used to group tests into tags. Update the configuration as follows:

```
➜ echo 'tags:
  integration:
    names:
      - "^integration_test$"' > godel/config/test.yml
```

This configuration defines a tag named "integration" that matches any directories named "integration_test". Run the
following command to run only the tests that match the "integration" tag:

```
➜ ./godelw test --tags=integration
ok  	github.com/nmiyake/echgo/integration_test	0.583s
```

By default, the `./godelw test` task runs all tests (all tagged and untagged tests). Multiple tags can be specified by
separating them with a comma. Specifying `all` will run all tagged tests, while specifying `none` will run all tests
that do not match any tags.

Commit these changes by running the following:

```
➜ git add godel main.go integration_test vendor
➜ git commit -m "Add integration tests"
[master 676aad3] Add integration tests
 6 files changed, 236 insertions(+), 1 deletion(-)
 create mode 100644 integration_test/doc.go
 create mode 100644 integration_test/integration_test.go
 create mode 100644 vendor/github.com/palantir/godel/pkg/products/products.go
 create mode 100644 vendor/github.com/palantir/godel/pkg/products/products_test.go
```

Tutorial end state
------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist.yml` is configured to build `echgo`
* Project is tagged as 0.0.1
* `godel/config/dist.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* Go files have license headers
* `godel/config/generate.yml` is configured to generate string function
* `godel/config/exclude.yml` is configured to ignore all `.+_string.go` files
* `integration_test` contains integration tests
* `godel/config/test.yml` is configured to specify the "integration" tag

([Link](https://github.com/nmiyake/echgo/tree/676aad36a5c355af826397be682f49bbb4a9ed20))

Tutorial next step
------------------

[Sync documentation with GitHub wiki](https://github.com/palantir/godel/wiki/GitHub-wiki)
