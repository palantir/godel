Summary
-------
`./godelw format` formats all of the Go files in a project by running `gofmt` and `ptimports` on them.

Tutorial start state
--------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files

([Link](https://github.com/nmiyake/echgo/tree/0c815bfd02711336f5ec0124377c95829667928a))

Format code
-----------

Update the program code to put the echo functionality into a separate package and call it from the main project:

```
➜ mkdir -p echo
➜ echo 'package echo
func Echo(in string) string {
	fmt.Println(strings.Join(os.Args[1:], " "))
}' > echo/echo.go
➜ echo 'package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/nmiyake/echgo/echo"
)

func main() {
	fmt.Println(echo.Echo(strings.Join(os.Args[1:], " ")))
}' > main.go
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

This command formats all of the files using `gofmt` and `ptimports`. Verify that this command modified the files:

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
index f7d4dbe..3e055c2 100644
--- a/echo/echo.go
+++ b/echo/echo.go
@@ -1,4 +1,5 @@
 package echo
+
 func Echo(in string) string {
 	fmt.Println(strings.Join(os.Args[1:], " "))
 }
diff --git a/main.go b/main.go
index 1d9820b..9087898 100644
--- a/main.go
+++ b/main.go
@@ -4,6 +4,7 @@ import (
 	"fmt"
 	"os"
 	"strings"
+
 	"github.com/nmiyake/echgo/echo"
 )

```

In `main.go`, note how the `github.com/nmiyake/echgo/echo` import has been separated from the other imports.
`ptimports` groups imports into distinct sections consisting of the standard library imports, external imports and
imports that are part of the same package.

Commit the formatted files:

```
➜ git add main.go echo
➜ git commit -m "Create echo package"
[master 24f63f7] Create echo package
 2 files changed, 8 insertions(+), 1 deletion(-)
 create mode 100644 echo/echo.go
➜ git status
On branch master
nothing to commit, working directory clean
```

Tutorial end state
------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go`

([Link](https://github.com/nmiyake/echgo/tree/24f63f727542c7189c82f04f7e2a4aa38c090137))

Tutorial next step
------------------

[Run static checks on code](https://github.com/palantir/godel/wiki/Check)

More
----

### Differences between `./godelw format` and `gofmt -w ./...`

`./godelw format` has the following advantages over running `gofmt -w ./...`

* Only formats files that are part of the project (does not format files in excluded directories)
* Runs `ptimports` in addition to `gofmt`

`ptimports` groups Go imports into 3 groups: the Go standard library, packages that are not part of the current project
and packages that are part of the current project.

The output of `./godelw format` is compatible with `gofmt` (that is, it is guaranteed that applying `gofmt` to a file
formatted by `./godelw format` will be a no-op).

### `--list` flag

Running `./godelw format` with the `--list` (or `-l`) flag will output the files that would be changed if
`./godelw format` were run without actually applying the changes.
