gödel supports plugins, which are executables that provide additional tasks for a gödel project. gödel provides a
mechanism for resolving and invoking plugins on a per-project basis.

Overview
--------
In order for an executable to be a gödel plugin, it must satisfy the following properties:

* When invoked with the `_godelPluginInfo` argument, it must print the JSON representation of a valid `pluginapi.Info` struct
* The executable must support the flags/arguments it defines in the `pluginapi.Info` struct returned by `_godelPluginInfo`

The plugin executable must be packaged in a `tgz` archive, and the executable must be the only file in the archive.

These are the only requirements for an executable to satisfy the requirements of being a gödel plugin. The plugin
architecture is essentially just a mechanism for advertising the tasks provided by an executable and defining the
flags/arguments for the executable. This means that, as long as these requirements are followed, there are no particular
requirements for the plugins themselves -- they may be written in any language or use any framework as long as the
executable can be run on the host system.

Plugin Identifier
-----------------
All plugins must define an identifier that is globally unique. A plugin identifier is a string of the form
`[group]:[product]:[version]` -- for example, "com.palantir:test-plugin:1.0.0". Although there are no specific
requirements for the content of `[group]`, `[product]` or `[version]` (beyond `:` not being allowed), it is recommended
[Maven naming conventions](https://maven.apache.org/guides/mini/guide-naming-conventions.html) be used for consistency.

Plugin Concepts
---------------

Tasks
=====
A plugin provides one or more tasks. When a plugin is added to a gödel project, the tasks provided by the plugin are
added to the available task set. For example, if a project was configured to use a plugin that provided a task named
`grpc`, running `./godelw grpc` for that project would invoke the `grpc` task provided by the plugin. Tasks names must
be globally unique within a project. A plugin is considered invalid if it provides multiple tasks with the same name,
and a project's configuration is considered invalid if it contains any plugins that would cause a task to be defined
multiple times.

Assets
======
A project can specify the assets for a plugin in its configuration. Assets allow for the behavior of a plugin to be
customized -- in a way, they can be thought of as "plugins for plugins", except that they have far less structure. If
a project specifies assets in its configuration, `godel` ensures that the assets are resolved and locally available when
the plugin is run, and provides the paths to all of the assets as flag arguments to the plugin. The plugin is then
responsible for using the assets in whatever manner they see fit.

Configuration
=============
Many plugins require user-specified configuration. Plugins may specify the name of a configuration file that it requires
for configuration. For example, if a plugin defined `grpc.yml` as its configuration file, then any of the tasks invoked
for the plugin would be provided with `godel/config/grpc.yml` as the configuration file (this would be specified via a
flag to the executable). Configuration files are per-plugin rather than per-task. If a plugin does not require a
configuration file, it does not need to define one. A plugin may only declare one configuration file. No verification or
validation is done by gödel for configuration files -- if multiple plugins specify the same file, the same path will be
provided to tasks of both plugins.

Shared Configuration
====================
There are some pieces of configuration that can be considered global to all gödel tasks. For example, gödel allows users
to configure files or paths that should be considered "excluded" from all gödel tasks. Such configuration is stored in
the `godel/config/godel.yml` file, which is a file that contains the YAML representation of a
`godellauncher.GodelConfig` struct. A plugin may declare that it wants access to the global config, in which case the
`godel/config/godel.yml` path is provided as the value of a flag to the executable.

Debug Mode
==========
Plugins may want to change certain behavior when run in `debug` mode. `godel` runs in debug mode if the `--debug` flag
is provided to the executable. The built-in tasks print full stack traces (rather than just the error message) on errors
when in debug mode, and some built-in tasks offer more verbose output when run in this mode. If a plugin declares that
it supports debug mode, then a debug flag is set on plugin execution if `godel` is run in debug mode.

Project Directory
=================
There are many instances in which a task may want to know the project directory. For example, if the `format` task is
run without arguments in a subdirectory of the project (`../godelw format`), the expectation is still that all project
files should be formatted. A plugin may specify that it wants to know the project directory, in which case the path to
the project directory is provided as a value of a flag.

Verify Tasks
============
The `verify` task is a built-in gödel task that is meant to invoke all of the tasks necessary to ensure that the project
is in the correct state. The `verify` task also supports an `--apply=false` flag that, when specified, states that the
verification task should determine whether or not the project is valid, but should make a best effort to not apply any
changes.

Plugins can specify whether or not the tasks that it provides should be run as part of the `verify` task on a per-task
basis. Plugins can also specify custom flags or options that should be added depending on if `verify` is run with the
`apply` flag being `true` or `false`.
