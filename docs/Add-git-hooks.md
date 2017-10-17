Summary
-------
`./godelw git-hooks` installs a [Git commit hook](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks) that ensures
that all of the files in a project are formatted using `./godelw format` before they are committed (requires the project
to be a Git repository).

Tutorial start state
--------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`

([Link](https://github.com/nmiyake/echgo/tree/6a73370d5b9c8c32ce1a5218938c922f1218db30))

Create Git commit hook
----------------------

Install the Git hooks for gödel in the current project by running the following:

```
➜ ./godelw git-hooks
```

With the repository initialized and hooks installed, we start writing code. Run the following to generate the initial
version of a `main.go` file echoes the arguments provided by the user:

```
➜ echo 'package main
import "fmt"
import "strings"
import "os"
func main() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}' > main.go
```

This is valid Go that compiles and runs properly:

```
➜ go run main.go foo
foo
```

However, if we attempt to add and commit this file to the repository, it will fail:

```
➜ git add main.go
➜ git commit -m "Add main.go"
Unformatted files exist -- run ./godelw format to format these files:
  /Volumes/git/go/src/github.com/nmiyake/echgo/main.go
```

This is because the commit hook has determined that `main.go` is not formatted properly. We can run `./godelw format`
(this is covered in more detail in the [Format](https://github.com/palantir/godel/wiki/Format) section of the tutorial)
to format the file and then verify that adding and committing the file works:

```
➜ ./godelw format
➜ git add main.go
➜ git commit -m "Add main.go"
[master 9e77ec4] Add main.go
 1 file changed, 11 insertions(+)
 create mode 100644 main.go
➜ git status
On branch master
nothing to commit, working directory clean
```

We now have a repository that contains the first version of our `echgo` program and have a commit hook that ensures that
all of the code we check in for our program will be properly formatted.

Tutorial end state
------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`

([Link](https://github.com/nmiyake/echgo/tree/9e77ec4885591f5a4fd95b550da729a004e7a04a))

Tutorial next step
------------------

[Generate IDE project for Gogland](https://github.com/palantir/godel/wiki/Generate-IDE-project)

More
----

### Hook installation

Running `./godelw git-hooks` will overwrite the `.git/hooks/pre-commit` file (including any previous customizations).

### Uninstalling the hook

The commit hook can be uninstalled by removing the generated commit hook file at `.git/hooks/pre-commit`.
