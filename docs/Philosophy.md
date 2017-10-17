Philosophy
==========
The design of gödel is based on some core philosophical principles. Although they are not particularly novel, these
principles drive most of the design decisions of gödel.

Builds and deployment are part of a project
-------------------------------------------
Building and deployment should be core concerns for all projects. Code is not useful in a vacuum, and it is all too
often the case that developers focus on writing code or implementing functionality without giving deep thought to how
the software will eventually be built and deployed. Building the binaries for a product and creating its distribution is
just as critical to a product as its core code, and the code and mechanisms for performing these tasks should be held to
the same standard as the core code for a project.

gödel provides tasks for building and publishing and includes them as part of the core project.

Separate configuration and logic
--------------------------------
Configuration and logic should be distinct. The place where this is violated most often is build scripts -- typically,
scripts that are written to build and distribute a project conflate the configuration for the project with the logic for
performing actions such as building and publishing the product. This makes it hard to track changes, since changes in
these files may be either due to changing configuration or changing logic. It also makes logic harder to test and
re-use. If a script is useful, other projects will tend to copy them. However, if they mingle logic and configuration,
they tend to fork as they are copied, and it's hard (or impossible) to roll out updates later in a uniform manner.

gödel establishes clear separation between configuration and logic. All of the project-specific configuration is stored
in config files in `godel/config` and all of the generic logic for running actions is in gödel. This allows the logic to
be tested and generalized, and also makes it possible to roll out updates cleanly across multiple different projects.

If a convention isn't enforced with automation, it will decay
-------------------------------------------------------------
Whether it be indentation level, coding style, license headers on files or tasks that create generated code, if a
convention is not enforced in an automated manner, it will inevitably decay over time (this is especially true for
larger code bases with multiple contributors). This is similar to the philosophy behind [gofmt](https://blog.golang.org/go-fmt-your-code).
Establishing a convention and asking people to follow it (whether it be by running a program themselves, creating a Git
hook for it, or trying to catch it in code reviews) is great. However, if the convention is not checked and enforced in
an automated manner, it is doomed to fail (a.k.a. [second law of thermodynamics](https://en.wikipedia.org/wiki/Second_law_of_thermodynamics)).

gödel tasks are designed to be able to run as part of CI to verify and enforce best practices for formatting, linting
checks, license headers and more.

Embed expertise and knowledge into tools so that it scales
----------------------------------------------------------
As people spend time working in a codebase or ecosystem, they start to develop expertise. If someone has been working in
Go for a while, knowing that a project must be in a `$GOPATH` is second nature. If `go build ./...` fails due to failing
to build something in the vendor directory, they can easily reason/Google their way to running
`go build $(go list ./... | grep -v /vendor/)` instead. One may also know tricks like the fact that running `go install std`
to build the standard library for different OS/architectures before cross-compiling saves a ton of time for repeated
cross-platform compilations. Pushing out tips and tricks like this to all of the developers that might work on a project
is hard to do, and even if an effective mechanism for it exists inertia is a powerful force and many people will not
alter their practices. However, adding this kind of information/logic to tooling provides it for free to anyone who uses
the tools.

gödel bakes in many of these kinds of optimizations and expertise into the tooling so that even people who are new to Go
can immediately create projects and start writing, building and publishing code in Go without having to learn these
kinds of optimizations themselves.

Tasks that apply should be able to verify
-----------------------------------------
Some tasks like applying correct formatting or ensuring that source files have a specific license header can be applied
automatically without user feedback. Ideally, any task that can apply a change should also have a mode that verifies
whether the state of the world matches the expected state without making any modifications. If the state of the world
differs, the task should fail and provide feedback about what needs to be changed to make it pass. Having this mode of
operation makes it much easier to run the task as part of continuous integration.

gödel provides a verification mode for all of its tasks that make modifications.

Optimize the default use case for humans
----------------------------------------
When thinking about the default behavior for commands or flags, optimize on the choice that is easier for humans that
use the tool in an interactive manner. This came up when trying to determine the default behavior for `./godelw verify`.
The command has an `--apply` flag that, when true, applies the changes detected by verify, and there was a question as
to what the default value of the flag should be. Conceptually, there's a case to be made for the default value being
false -- `verify` should only verify by default and apply changes only when specifically instructed to do so. However,
most developers invoking the command locally will always want to apply the changes that are flagged. Requiring them to
run `./godelw verify --apply=true` every time is much more onerous than just making `./godelw verify` apply changes by
default. Based on this, we made the default value of the flag `true`. This means that CI environments typically have to
invoke `./godelw verify --apply=false`. However, in CI, this value is defined once and then forgotten, so this trade-off
made sense.

gödel is designed to have sensible defaults that do what a developer would expect when invoking the command locally.

Checks must be repeatable
-------------------------
The result of a given set of checks on a specific input should stay constant over time. This property is not true for
builds that use `go get` to retrieve checks or build tools at the beginning of each build because the result of `go get`
can change over time (either because the version of the tool changes or because it becomes unavailable).

gödel contains the code for all of its checks and requires projects to declare the version of gödel that is uses as part
of its configuration, so a given version of a project is tied to a specific version of gödel and the result of the
checks stay constant over time.

Although the gödel executable itself is not included in projects, as long as the distribution for the version of gödel
used is available it will build with the same result. gödel distributions can also be downloaded/saved locally or
on-prem to avoid dependencies on external services.

Checks should be fast and idempotent
------------------------------------
Developers iterate quickly and have little patience. Tasks should complete on the order of seconds when possible. It
should also be possible to kill checks at any point without any adverse effects.

Most gödel tasks complete in under 1 second on small projects, and even on larger projects most tasks complete on the
order of seconds (checks and tests can take longer). All tasks are designed to be idempotent and fail gracefully and
clean up when terminated.

Failures should be obvious and provide contextual information
-------------------------------------------------------------
If a task or operation fails, it should provide a user-readable explanation of what failed along with any relevant
context for diagnosing or reproducing the issue.

gödel strives to provide readable error messages with context rather than just stack traces on failures. Common failure
paths have been identified and the error messages strive to include common causes and work-arounds. Failures that occur
in sub-processes include information on the command that was invoked and the environment variables provided to it so
that users can attempt to diagnose the issue manually.

Provide building blocks that can be composed
--------------------------------------------
Most tasks are composed of multiple different parts or actions, and by default tasks should provide an all-inclusive
experience that works out-of-the-box. However, tasks should also provide the flexibility to be used in other ways so
other tools can use them to compose their own tasks. When possible, internal tasks should also be structured in a way
that composes these distinct tasks so that it is possible to manually perform partial tasks or debug failures of
compound tasks.

gödel provides tasks such as `./godelw packages` and `./godelw products` that echo out state defined by gödel in a
manner that can easily be consumed by other tasks or tools. It also provides the `__` invocation mechanism that can be
used to run individual pieces of subprograms in isolation (for example, running `./godelw __check __errcheck` invokes
the exact piece of sub-functionality used by `./godelw check` when it runs `errcheck`).

Use a single source of truth with good abstractions
---------------------------------------------------
If there is state about the world that is true, define it in a single place and re-use it across the project using the
correct abstractions. For example, most build tools want to ignore the "vendor" directory and any directories that
contain generated code. These directories should be excluded from tests, checks, formatting, etc. Rather than adding
logic to each of these different tasks to ignore these directories, recognize that the abstract issue is that there are
a set of directories that are not considered part of the core source of a project -- define them in a single place and
then share that information across all tasks and logic.

gödel does this for things like exclude directories and projects. It establishes abstractions for things like products,
projects and packages and uses them consistently in its own code and exposes them as tasks so that other tools can use
them as well.

Orchestrate, don't obfuscate
----------------------------
Go has standard tooling that people understand well. Whenever possible, build tools should delegate to the built-in
tooling to do things in a standard manner.

gödel is designed to do as little work as possible -- its main concerns are reading declarative configuration and then
orchestrating standard tools to perform the actual work. When errors or failures occur, gödel strives to expose the
failure in a way that can be repeated/verified using standard tooling to verify that the failure was not introduced by
gödel.

If it's not tested, it can't be trusted
---------------------------------------
If a piece of functionality isn't tested, it can't be trusted to work. Even if it works now, there are no guarantees
that the behavior won't regress in the future.

gödel has extensive unit and integration tests that are run in CI against multiple different environments. Almost all
fixes for issues that are identified are accompanied by tests that ensure that the issue stays fixed. The code and build
is designed in a manner that almost every aspect is testable.

Anticipate user needs and enable them
-------------------------------------
If a set of checks are run and one fails, a user will probably want the ability to re-run just the failing check. If
the "build" task builds all products for all platforms by default, a user will probably want to be able to build just a
specific product for a specific platform. If a task caches results by default, the user will probably want a way to run
the task in a manner that ignores the cache.

gödel tries to anticipate all of the needs or requests that a user will have for a task enable them. The goal is to
optimize for the common case, but to provide configuration and options for customizing behavior as necessary.

Trust is hard to earn and easy to lose
--------------------------------------
Having someone opt to use a piece of software is one of the highest compliments that can be paid to it. People are
inherently skeptical of new products, and nothing is more frustrating than a build tool that becomes a source of build
issues. Build tools are often the bearer of bad news so people tend to have a negative reaction towards interacting with
them in the first place, so people have very little patience with build tools.

gödel strives to be a tool that enhances productivity and can be trusted as a core part of a development and continuous
integration setup. It is dog-fooded extensively and maintained with care.
