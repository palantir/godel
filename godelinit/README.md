godelinit
=========
`godelinit` is a CLI that can be used to add gödel to any project. Once it is installed, invoking `godelinit` will add
the latest version of gödel to the working directory.

Installation
------------
`go get github.com/palantir/godel/godelinit`

Usage
-----
Run `godelinit` in the root directory of a Go project to add the `godelw` script and `godel` configuration directory
from the latest gödel release (as listed on https://github.com/palantir/godel/releases) to the project.

Running `godelinit` in the root directory of a project that already uses gödel will update that project to use the
latest version (if it is already on the latest version, it will be a no-op).

The `--version` flag can be used to specify the version of gödel that should be installed. This flag should be used if
the latest released version is not the desired one.

If a version is not specified using the `--version` flag, `godelinit` determines the latest release by querying the
GitHub API of the palantir/godel GitHub project. The version of the latest release is stored in a cache directory in the
gödel home directory. If the latest version was determined within the last hour, the cached version will be used (and
thus an API call will not be made). The `--ignore-cache` flag can be used to ignore the cached version and force it to
be retrieved using the GitHub API even if a valid cache entry exists.
