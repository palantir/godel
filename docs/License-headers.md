Summary
-------
`./godelw license` updates the Go files in the project to have a specific license header based on configuration.

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

([Link](https://github.com/nmiyake/echgo/tree/55182ff79dd28048782fb240920d6f2d90b453da))

License
-------

Many open-source projects require specific license headers to be part of every source file. This can be enforced using
the `license` task and configuration.

First, add the license as a license file:

```
➜ curl http://www.apache.org/licenses/LICENSE-2.0.txt | sed '/./,$!d' > LICENSE
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 11358  100 11358    0     0   150k      0 --:--:-- --:--:-- --:--:--  151k
```

Run the following to configure a license header:

```
➜ echo 'header: |
  // Copyright (c) 2017 Author Name
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
  // limitations under the License.' > godel/config/license.yml
```

Run `./godelw license` to apply this license to all of the Go files in the project:

```
➜ ./godelw license
```

Verify that this updated the Go files:

```
➜ git status
On branch master
Your branch is up-to-date with 'origin/master'.
Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git checkout -- <file>..." to discard changes in working directory)

	modified:   echo/echo.go
	modified:   echo/echo_test.go
	modified:   echo/echoer.go
	modified:   godel/config/license.yml
	modified:   main.go

Untracked files:
  (use "git add <file>..." to include in what will be committed)

	LICENSE

no changes added to commit (use "git add" and/or "git commit -a")
➜ cat echo/echo.go
// Copyright (c) 2017 Author Name
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

package echo

func NewEchoer() Echoer {
	return &simpleEchoer{}
}

type simpleEchoer struct{}

func (e *simpleEchoer) Echo(in string) string {
	return in
}
```

Commit the changes to the repository:

```
➜ git add LICENSE echo godel main.go
➜ git commit -m "Add LICENSE and license headers"
[master 0239b28] Add LICENSE and license headers
 6 files changed, 271 insertions(+)
 create mode 100644 LICENSE
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

([Link](https://github.com/nmiyake/echgo/tree/0239b282904d05bb9eef6c3c3edfe1c28f888ad3))

Tutorial next step
------------------

[Go generate tasks](https://github.com/palantir/godel/wiki/Generate)

More
----

### Remove license headers

In some instances, it may be desirable to remove the license headers from all of the files. For example, if you are
changing the type of license for the repository, you will want to remove all of the license headers that are already
present before adding new headers.

Run the following command:

```
➜ ./godelw license --remove
```

Verify that this removed the headers:

```
➜ git status
On branch master
Your branch is ahead of 'origin/master' by 1 commit.
  (use "git push" to publish your local commits)
Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git checkout -- <file>..." to discard changes in working directory)

	modified:   echo/echo.go
	modified:   echo/echo_test.go
	modified:   echo/echoer.go
	modified:   main.go

no changes added to commit (use "git add" and/or "git commit -a")
➜ cat echo/echo.go
package echo

func NewEchoer() Echoer {
	return &simpleEchoer{}
}

type simpleEchoer struct{}

func (e *simpleEchoer) Echo(in string) string {
	return in
}
```

Revert these changes by running the following:

```
➜ git checkout -- echo main.go
```

### Specify custom license headers for specific paths

In some instances, a project may contain certain files or directories that should have a different license header from
other files -- for example, if a file or directory is based on a file from another project, it may be necessary to have
a custom header to provide attribution for the original authors.

Run the following command to remove the existing headers:

```
➜ ./godelw license --remove
```

Once that is done, update the license configuration as follows:

```
➜ echo 'header: |
  // Copyright (c) 2017 Author Name
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
custom-headers:
  - name: echo
    header: |
      // Copyright 2017 Author Name. All rights reserved.
      // Licensed under the MIT License. See LICENSE in the project root
      // for license information.
    paths:
      - echo' > godel/config/license.yml
```

This configuration specifies that the paths that match `echo` (which includes all paths within the `echo` directory)
should use the custom header named "echo", while all of the other files should use the standard header. Run the
following to apply the license and verify that it behaved as expected:

```
➜ ./godelw license
➜ cat echo/echo.go
// Copyright 2017 Author Name. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package echo

func NewEchoer() Echoer {
	return &simpleEchoer{}
}

type simpleEchoer struct{}

func (e *simpleEchoer) Echo(in string) string {
	return in
}
➜ cat main.go
// Copyright (c) 2017 Author Name
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
	"strings"

	"github.com/nmiyake/echgo/echo"
)

var version = "none"

func main() {
	versionVar := flag.Bool("version", false, "print version")
	flag.Parse()
	if *versionVar {
		fmt.Println("echgo version:", version)
		return
	}
	echoer := echo.NewEchoer()
	fmt.Println(echoer.Echo(strings.Join(flag.Args(), " ")))
}
```

Revert these changes by running the following:

```
➜ git checkout -- echo godel
```
