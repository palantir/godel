nobadfuncs
==========
`nobadfuncs` verifies that a set of specified functions are not referenced in the packages being checked. It can be used
to blacklist specific functions that should not typically be referenced or called. It is possible to explicitly allow
uses of black-listed functions by adding a comment to the line before the calling line.

Usage
-----
`nobadfuncs` takes the path to the packages that should be checked for function calls. It also takes configuration (as
JSON) that specifies the blacklisted functions.

The function signatures that are blacklisted are full function signatures consisting of the fully qualified package name
or receiver, name, parameter types and return types. Examples:

```
func (*net/http.Client).Do(*net/http.Request) (*net/http.Response, error)
func fmt.Println(...interface{}) (int, error)
```

`nobadfuncs` can be run with the `--all` flag to print all of the function references in the provided packages. The output
can be used as the basis for determining the signatures for blacklist functions.

Examples
========

```bash
> nobadfuncs --all .
/Volumes/.../src/github.com/palantir/checks/nobadfuncs/nobadfuncs.go:54:13: func github.com/palantir/pkg/cli.NewApp(...github.com/palantir/pkg/cli.Option) *github.com/palantir/pkg/cli.App
/Volumes/.../src/github.com/palantir/checks/nobadfuncs/nobadfuncs.go:54:24: func github.com/palantir/pkg/cli.DebugHandler(github.com/palantir/pkg/cli.ErrorStringer) github.com/palantir/pkg/cli.Option
```

```bash
> nobadfuncs --config '{"func os.Exit(int)": "do not call os.Exit directly"}' .
/Volumes/.../src/github.com/palantir/checks/nobadfuncs/nobadfuncs.go:85:5: do not call os.Exit directly
```
