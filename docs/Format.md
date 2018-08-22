Summary
-------
`./godelw format` formats all of the Go files in a project by running `ptimports` on them.

Tutorial start state
--------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files

Format code
-----------
Update the program code to put the echo functionality into a separate package and call it from the main project:

```
➜ mkdir -p echo
➜ echo 'package echo
func Echo(in string) string {
	return in
}' > echo/echo.go
➜ SRC='package main

import "fmt"
import "os"
import "strings"
import "PROJECT_PATH"

func main() {
	fmt.Println(echo.Echo(strings.Join(os.Args[1:], " ")))
}' && SRC=${SRC//PROJECT_PATH/$PROJECT_PATH} && echo "$SRC" > main.go
```

Stage these files as a git commit:

```
➜ git add echo main.go
➜ git status
On branch master
Changes to be committed:
  (use "git reset HEAD <file>..." to unstage)

	new file:   echo/echo.go
	modified:   main.go

```

Run `./godelw format` to format all of the files in the project:

```
➜ ./godelw format
```

This command formats all of the files using `ptimports`. Verify that this command modified the files:

```
➜ git status
On branch master
Changes to be committed:
  (use "git reset HEAD <file>..." to unstage)

	new file:   echo/echo.go
	modified:   main.go

Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git checkout -- <file>..." to discard changes in working directory)

	modified:   echo/echo.go
	modified:   main.go

➜ git diff | cat
diff --git a/echo/echo.go b/echo/echo.go
index c3bf761..e7e2c75 100644
--- a/echo/echo.go
+++ b/echo/echo.go
@@ -1,4 +1,5 @@
 package echo
+
 func Echo(in string) string {
 	return in
 }
diff --git a/main.go b/main.go
index f69640d..d58038c 100644
--- a/main.go
+++ b/main.go
@@ -1,9 +1,12 @@
 package main
 
-import "fmt"
-import "os"
-import "strings"
-import "github.com/nmiyake/echgo2"
+import (
+	"fmt"
+	"os"
+	"strings"
+
+	"github.com/nmiyake/echgo2/echo"
+)
 
 func main() {
 	fmt.Println(echo.Echo(strings.Join(os.Args[1:], " ")))
```

In `main.go`, note how the imports that were on individual lines were grouped into an import block and how the standard
library imports are separated from the non-standard library imports. This is due to `ptimports`.

Commit the formatted files:

```
➜ git add main.go echo
➜ git commit -m "Create echo package"
[master a69f262] Create echo package
 2 files changed, 8 insertions(+), 1 deletion(-)
 create mode 100644 echo/echo.go
```

Tutorial end state
------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files
* Project contains `echo/echo.go`

Tutorial next step
------------------
[Run static checks on code](https://github.com/palantir/godel/wiki/Check)

More
----
### Differences between `./godelw format` and `gofmt` or `goimports`
`./godelw format` has the following advantages over `gofmt` or `goimports`:

* Only formats files that are part of the project (does not format files in excluded directories)
* Single line imports are converted into the import block format (which is generally preferable/better style)
* Performs code simplification equivalent to `gofmt -s`
* Performs import grouping provided by `goimports`

`./godelw format` is essentially equivalent to running `goimports` and `gofmt -s` on the source (with the additional
behavior of converting single-line imports to block imports).

The `ptimports` asset can be configured to match the behavior of exactly `gofmt` (without simplification), exactly the
behavior of `goimports`, or exactly the behavior of running `gofmt -s` and `goimports` using the `format-plugin.yml`
configuration.

The default configuration of `ptimports` is equivalent to the following:

```yaml
formatters:
  ptimports:
    config:
      skip-refactor: false
      skip-simplify: false
```

Setting `skip-refactor` to true disables the behavior that converts single-line imports to grouped imports, while
setting `skip-simplify` to true disables the `gofmt -s` behavior.

### `--verify` flag
Running `./godelw format` with the `--verify` flag outputs the files that would be changed if `./godelw format` were
run without actually applying the changes.

### Provide filenames as arguments
`./godelw format [flags] [files]` runs the `format` operation on the specified Go files. If `[files]` is blank, the
operation is run on all of the non-excluded project Go files.
