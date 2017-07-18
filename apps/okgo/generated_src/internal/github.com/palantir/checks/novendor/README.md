novendor
========

`novendor` is a program that verifies that there are no unused vendored imports in a Go project. If unused vendored
packages are found, the program prints them and exits with a non-0 exit code. An exit code of 0 indicates that there are
no unused vendored packages in the project given the parameters that were provided to the program.

A vendored package is considered unused if the package is not imported by any of the non-vendored code in the project
(including test code) either directly or transitively. It should be possible to remove any package that is reported by
this tool from vendoring and still have the project build correctly. Standard Go build rules are used to determine if a
package is vendored (basically, any package that is within a "vendor" directory is vendored).

Project Packages
================
`novendor` has a notion of "project packages". A "project package" is considered to be a top-level package for a single
project that may contain many subpackages. Although this is not an official Go concept, it captures much of how code is
organized in practice. A "project package" is considered to be the 3rd level of a package import (if the import path has
fewer than 3 parts, then all of it is considered). For example, the "project package" of
`github.com/palantir/stacktrace/cleanpath` is `github.palantir.com/palantir/stacktrace`, while the "project package" of
`gopkg.in/yaml.v2` is `gopkg.in/yaml.v2`. This scheme is not perfect -- for example, if a package named
`gopkg.in/yaml.v2/subpackage` exists, its "project package" would be `gopkg.in/yaml.v2/subpackage`, even though
conceptually `gopkg.in/yaml.v2` would probably be more appropriate. However, in practice the 3-level heuristic works for
most packages

This concept is used because in many cases projects want to vendor "project packages" as a unit. For example, consider
the packages `github.com/org/project`, `github.com/org/project/api` and `github.com/org/project/impl`. If the primary
project code only imports `github.com/org/project/api`, then technically the other 2 packages are "unused". However, a
project may still want to vendor `github.com/org/project` and all of its subdirectories because they may want to make
use of the other code later or ensure that different subdirectories of a single project are not inadvertently vendored
at different versions.

The default behavior of `novendor` works at the "project package" granularity. So, if `github.com/org/project/api` and
`github.com/org/project/impl` are both vendored but only `github.com/org/project/api` is used, `github.com/org/project`
is considered as "used" and is not reported as an unused project package. If both are unused, the default behavior will
report that `github.com/org/project` is unused.

Use the `--project-package=false` flag to turn off the "project package" behavior (this will cause all used/unused
determination and output to be done purely at the Go package level).

Usage
=====
`novendor` takes the path to the packages that are part of the project for which vendoring status should be determined.
For example, if a project consists of a top-level package and a `cmd` package, `novendor` would be run at the project
root as:

```bash
novendor . ./cmd
```

Note that, as shown above, all packages that are part of a project should be provided as arguments (if no arguments are
provided, all non-vendored packages that are found in the working directory are used).

Usage of `novendor`:

```
  -f    Include full path of unused packages (default omits path to vendor directory)
  --project-package
        Use the 'project' paradigm to interpret packages and only output projects that are unused (default true)
```

Examples
========

```bash
> novendor .
github.com/docker/go-connections
gopkg.in/mcuadros/go-syslog.v2
```

```bash
> novendor -f .
github.palantir.build/deployability/novendor/vendor/github.com/docker/go-connections
github.palantir.build/deployability/novendor/vendor/gopkg.in/mcuadros/go-syslog.v2
```

```bash
> novendor --project-package=false .
github.com/docker/docker/pkg/longpath
github.com/docker/go-connections/nat
github.com/docker/go-connections/sockets
github.com/docker/go-connections/tlsconfig
gopkg.in/mcuadros/go-syslog.v2
gopkg.in/mcuadros/go-syslog.v2/example
gopkg.in/mcuadros/go-syslog.v2/format
```

```bash
> novendor -f --project-package=false .
github.palantir.build/deployability/novendor/vendor/github.com/docker/docker/pkg/longpath
github.palantir.build/deployability/novendor/vendor/github.com/docker/go-connections/nat
github.palantir.build/deployability/novendor/vendor/github.com/docker/go-connections/sockets
github.palantir.build/deployability/novendor/vendor/github.com/docker/go-connections/tlsconfig
github.palantir.build/deployability/novendor/vendor/gopkg.in/mcuadros/go-syslog.v2
github.palantir.build/deployability/novendor/vendor/gopkg.in/mcuadros/go-syslog.v2/example
github.palantir.build/deployability/novendor/vendor/gopkg.in/mcuadros/go-syslog.v2/format
```
