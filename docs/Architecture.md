gödel consists of the `godelw` wrapper script, the project configuration files in `godel/config` and the `godel` Go
executable.

godelw
------
`godelw` is a bash script that is used as the entry point for gödel. The purpose of the script is to locate and invoke
the real `godel` executable and to provide it with the configuration specified in the current project.

If the executable exists at the expected location, it is invoked. The `godelw` script is built as part of the
distribution and embeds the checksum of the `godel` executable in its source and verifies that the checksum of the
executable it is invoking matches the checksum stored in its source.

If the executable does not exist, then the wrapper script gets the URL for the distribution from the
`godel/config/config.properties` file, downloads the distribution from that location (using `wget` or `curl`), verifies
the checksum of the downloaded distribution (if provided), expands the distribution into the expected location and then
invokes it.

This setup is used so that the `godel` executable itself does not need to be checked in as part of a project.

The `godelw` script invokes the `godel` executable with the `--wrapper` flag that provides the location of the invoking
wrapper script. This value is used by `godel` to determine the location from which the project configuration should be
loaded. All of the user-provided arguments to `godelw` are passed on directly to `godel`.

godel
-----
`godel` is a single Go executable. It contains several Go libraries that use [`amalgomate`](https://github.com/palantir/amalgomate)
to combine disparate Go `main` programs into a single library. Thus, `godel` contains and can run as several disparate
Go programs such as `errcheck`, `go-junit-report` and others.

Many `godel` tasks act as an orchestrator of other sub-tasks. For example, the `check` command uses `okgo` to run
several different checks like `deadcode` and `errcheck`. The nifty thing is that, because `godel` embeds the
functionality of these programs and can run as them, the programs do not need to be separately installed. Instead, when
`godel` needs the functionality of one of these programs, it re-invokes itself as a sub-process with specific arguments
that instruct it to run as a specific sub-program.

This can be observed by manually invoking `godel` with special syntax:

```
./godelw __check __errcheck --help
Usage of /Users/nmiyake/.godel/dists/godel-0.9.0/bin/darwin-amd64/godel:
  -abspath
        print absolute paths to files
  -asserts
        if true, check for ignored type assertion results
  -blank
        if true, check for errors assigned to blank identifier
  -ignore value
        comma-separated list of pairs of the form pkg:regex
            the regex is used to ignore names within pkg (default "fmt:.*")
  -ignorepkg string
        comma-separated list of package paths to ignore
  -ignoretests
        if true, checking of _test.go files is disabled
  -tags value
        space-separated list of build tags to include
  -verbose
        produce more verbose logging
```

When `godel` is invoked with `__check __errcheck`, it runs in a manner that is identical to running the stand-alone
`errcheck` program. Thus, when `godel` needs the functionality of `errcheck`, it can invoke itself as a subprocess with
these arguments to get the functionality.

Package scope definition
------------------------
gödel defines global configuration that is used to specify the scope of the files that it should deal with. Projects
often have certain classes of files that should not be operated on by automated tooling (for example, the `vendor`
directory or generated source directories). The `godel/config/exclude.yml` file is used to define rules for excluding
certain paths from consideration. The `cfgcli` package is used so that all of the tasks (including those defined by
other stand-alone applications) can receive this configuration and thus operate on the same set of files.

Task encapsulation
------------------
gödel is composed of many independent tasks. Some of the simple tasks are implemented directly in the gödel project.
However, in order to prevent the project from becoming a single monolithic program, the large pieces of core
functionality are implemented as independent sub-programs that also expose their functionality as a library. For
example, `distgo` handles building products, creating distributions and running publish operations, while `okgo` handles
running all of the code checks.

Most of the sub-programs define a configuration file from which its configuration is read. The configuration file
typically contains program-specific configuration and an `exclude` object that specifies files and directories that
should be excluded from consideration.

However, when run as part of `godel`, the `exclude` configuration is typically defined globally in `exclude.yml`.
Requiring every separate sub-program to copy the `exclude` block would be cumbersome and error-prone. In order to deal
with this, sub-programs also accept configuration as a JSON string provided by the `--json` flag and can either combine
the JSON configuration with the file-based one or override the file-based configuration using the values provided in the
JSON. The `github.com/palantir/pkg/cli/cfgcli` provides a common API to do this.

Task API
--------
Tasks that are implemented as sub-programs use the `github.com/palantir/pkg/cli/cfgcli` package as an API to manage
configuration. Tasks should use declarative file-based configuration. Most tasks have a corresponding configuration file
in `godel/config` -- for example, the `check` task is configured using `godel/config/check.yml`, while the `distgo`
tasks (`build`, `dist`, `publish`, etc.) are configured using `godel/config/dist.yml`.

When a task is invoked using `godel`, `godel` configures the global variables in the `cfgcli` package to store the
proper values for the configuration file and JSON configuration for the task. The JSON configuration is populated with
the `exclude` contents specified in `godel/config/exclude.yml` so that sub-programs can use the global excludes.
