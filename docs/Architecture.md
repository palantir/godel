gödel consists of the `godelw` wrapper script, the project configuration files in `godel/config` and the `godel` Go
executable.

Overview
--------
At a high level, the `godel` executable is a task dispatcher: it determines the task that should be run based on the
flags and arguments provided by the user and runs the task. The tasks that can be run are a combination of the built-in
`godel` tasks (which are executed in-process) and plugin tasks, which are tasks that are supplied by plugins defined in
the project configuration.

The `godelw` script ensures that the `godel` executable is available locally (and downloads it if it is not) and
provides the wrapper script path as an argument to the invocation. The `godel` executable ensures that the plugin
executables and assets defined in project configuration are available locally (and resolves them if they are not).

godelw
------
`godelw` is a bash script that is used as the entry point for gödel. The purpose of the script is to locate and invoke
the `godel` executable and to provide it with the configuration specified in the current project. A typical invocation
takes the form `./godelw [args]`. The `godelw` wrapper script ensures that the `godel` executable of the proper version
is available locally, downloading it if it does not exist. The script then invokes the `godel` executable with the
`--wrapper` flag and provides the absolute path to the `godelw` script being invoked as the value of the flag. The
script also passes along all of the user-specified arguments. This setup is used so that the `godel` executable itself
does not need to be checked in as part of a project.

godel
-----
`godel` is a Go executable. The `godel` executable runs the task specified by the flags/arguments provided by the user.
If the `--wrapper` flag is provided, then the directory of the value is used as the project directory for the invocation
and the configuration is read from the `godel/config` directory of the project. If the `godel.yml` configuration file
specifies plugins and/or assets for the project, the `godel` executable ensures that they are available locally,
resolving them if they are not. The tasks provided by the plugins are then added to the task list for the executable.

After this point, the `godel` executable determines the task that should be run based on the flags and arguments
provided by the user and invokes the task. Tasks are either built-in tasks or plugin tasks. If a task is a built-in
task, it is defined as part of the `godel` executable and run in-process. If the task is a plugin task, `godel` invokes
the plugin executable with the appropriate arguments (as defined by the plugin API) and the plugin executable handles
the rest of the execution.

Plugins
-------
gödel provides a mechanism to define plugins for a project. Plugins are executables that provide tasks that can be run
for gödel. Because Go does not have first-class cross-platform support for native plugins, plugins are fully independent
executables (and can be written in any language). An executable becomes a gödel plugin by satisfying a particular API,
which consists of returning known structured output in response to specific arguments and by supporting a specific set
of flags and values defined by gödel's plugin API. When a plugin task is executed, the plugin executable is invoked as
a separate process.

See [Plugins](https://github.com/palantir/godel/wiki/Plugins) for more information.
