The following tutorial demonstrates how to install, configure and use gödel on an example project. We will be creating a
project called `echgo` that is a simple program that echoes user input in a variety of ways.

Following every step of the tutorial from beginning to end will show the entire end-to-end process of creating a new
project and using a variety of gödel features to configure it. It is also possible to jump to any section of interest
directly. Each step of the tutorial provides a general summary of the step, the expected preconditions before the step,
the actions to take during the step, and the conditions that should exist after the step.

The repository at https://github.com/nmiyake/echgo contains the result of walking through the tutorial.

The tutorial consists of the following steps:

* [Add gödel to a project](https://github.com/palantir/godel/wiki/Add-g%C3%B6del)
* [Add Git hooks to enforce formatting](https://github.com/palantir/godel/wiki/Add-git-hooks)
* [Generate IDE project for Gogland](https://github.com/palantir/godel/wiki/Generate-IDE-project)
* [Format Go files](https://github.com/palantir/godel/wiki/Format)
* [Run static checks on code](https://github.com/palantir/godel/wiki/Check)
* [Run tests](https://github.com/palantir/godel/wiki/Test)
* [Build](https://github.com/palantir/godel/wiki/Build)
* [Run](https://github.com/palantir/godel/wiki/Run)
* [Dist](https://github.com/palantir/godel/wiki/Dist)
* [Publish](https://github.com/palantir/godel/wiki/Publish)
* [Generate license headers](https://github.com/palantir/godel/wiki/License-headers)
* [Go generate tasks](https://github.com/palantir/godel/wiki/Generate)
* [Define excludes](https://github.com/palantir/godel/wiki/Exclude)
* [Write integration tests](https://github.com/palantir/godel/wiki/Integration-tests)
* [Sync a documentation directory with GitHub wiki](https://github.com/palantir/godel/wiki/GitHub-wiki)
* [Verify project](https://github.com/palantir/godel/wiki/Verify)
* [Set up CI to run tasks](https://github.com/palantir/godel/wiki/CI-setup)
* [Update gödel](https://github.com/palantir/godel/wiki/Update-g%C3%B6del)
* [Other commands](https://github.com/palantir/godel/wiki/Other-commands)
* [Conclusion](https://github.com/palantir/godel/wiki/Tutorial-conclusion)

This tutorial uses `github.com/nmiyake/echgo` as the project path. Some parts of the tutorial require the ability to
create and push to a repository on GitHub. Although it is possible to push a project with this path to any GitHub
repository, if you want to follow the tutorial in the most realistic manner, create your `echgo` project in a path that
is under a GitHub organization or user that you control: for example, `github.com/<user>/echgo` or
`github-enterprise.domain.com/<org>/echgo`.

Start the tutorial by creating the directory for your project and setting it to be the working directory:

```
➜ mkdir -p $GOPATH/src/github.com/nmiyake/echgo && cd $_
➜ pwd
/Volumes/git/go/src/github.com/nmiyake/echgo
```

Initialize a git repository, add a README and commit it:

```
➜ git init
➜ echo '`echgo` is a program that echoes input provided by the user.' > README.md
➜ git add README.md
➜ git commit -m "Initial commit"
[master (root-commit) 54a23e6] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md
```

Tutorial end state
------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory

([Link](https://github.com/nmiyake/echgo/tree/54a23e62a4f9983d60939fdc8ed8dd59f81ddf7c))

Tutorial next step
------------------

[Add gödel](https://github.com/palantir/godel/wiki/Add-g%C3%B6del)
