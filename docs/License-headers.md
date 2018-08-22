Summary
-------
`./godelw license` updates the Go files in the project to have a specific license header based on configuration.

Tutorial start state
--------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist-plugin.yml` is configured to build `echgo2`
* Project is tagged as 0.0.1
* `godel/config/dist-plugin.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* `dockerctx` directory exists and `godel/config/dist-plugin.yml` is configured to build Docker images for the product

License
-------
Many open-source projects require specific license headers to be part of every source file. This can be enforced using
the `license` task and configuration.

First, add the license as a license file:

```
➜ curl http://www.apache.org/licenses/LICENSE-2.0.txt | sed '/./,$!d' > LICENSE
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0100 11358  100 11358    0     0  35792      0 --:--:-- --:--:-- --:--:-- 35716
```

Run the following to configure a license header:

```
➜ echo 'header: |
  // Copyright (c) {{YEAR}} Author Name. All rights reserved.
  // Use of this source code is governed by the Apache License, Version 2.0
  // that can be found in the LICENSE file.' > godel/config/license-plugin.yml
```

Run `./godelw license` to apply this license to all of the Go files in the project:

```
➜ ./godelw license
```

Verify that this updated the Go files:

```
➜ git status
On branch master
Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git checkout -- <file>..." to discard changes in working directory)

	modified:   echo/echo.go
	modified:   echo/echo_test.go
	modified:   echo/echoer.go
	modified:   godel/config/license-plugin.yml
	modified:   main.go

Untracked files:
  (use "git add <file>..." to include in what will be committed)

	LICENSE

no changes added to commit (use "git add" and/or "git commit -a")
➜ cat echo/echo.go
// Copyright (c) 2018 Author Name. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package echo

func NewEchoer() Echoer {
	return &simpleEchoer{}
}

type simpleEchoer struct{}

func (e *simpleEchoer) Echo(in string) string {
	return in
}
```

Note that the "{{YEAR}}" in the license header was automatically replaced with the year at the time that the operation
is run (in this case, 2018). This template is rendered once when adding the license and is not otherwise modified
(and thus the license year will generally match the creation year for the file).

Commit the changes to the repository:

```
➜ git add LICENSE echo godel main.go
➜ git commit -m "Add LICENSE and license headers"
[master 2187646] Add LICENSE and license headers
 6 files changed, 221 insertions(+)
 create mode 100644 LICENSE
```

Tutorial end state
------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist-plugin.yml` is configured to build `echgo2`
* Project is tagged as 0.0.1
* `godel/config/dist-plugin.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* `dockerctx` directory exists and `godel/config/dist-plugin.yml` is configured to build Docker images for the product
* Go files have license headers

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

If a license contains the "{{YEAR}}" placeholder, any 4-digit year will match.

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
  // Copyright (c) {{YEAR}} Author Name. All rights reserved.
  // Use of this source code is governed by the Apache License, Version 2.0
  // that can be found in the LICENSE file.
custom-headers:
  - name: echo
    header: |
      // Copyright {{YEAR}} Author Name. All rights reserved.
      // Licensed under the MIT License. See LICENSE in the project root
      // for license information.
    paths:
      - echo' > godel/config/license-plugin.yml
```

This configuration specifies that the paths that match `echo` (which includes all paths within the `echo` directory)
should use the custom header named "echo", while all of the other files should use the standard header. Run the
following to apply the license and verify that it behaved as expected:

```
➜ ./godelw license
➜ cat echo/echo.go
// Copyright 2018 Author Name. All rights reserved.
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
// Copyright (c) 2018 Author Name. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/nmiyake/echgo2/echo"
)

var version = "none"

func main() {
	versionVar := flag.Bool("version", false, "print version")
	flag.Parse()
	if *versionVar {
		fmt.Println("echgo2 version:", version)
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
