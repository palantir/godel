The following tutorial demonstrates how to install, configure and use gödel on an example project. We will be creating a
project called echgo2 that is a simple program that echoes user input in a variety of ways.

Following every step of the tutorial from beginning to end will show the entire end-to-end process of creating a new
project and using a variety of gödel features to configure it. It is also possible to jump to any section of interest
directly. Each step of the tutorial provides a general summary of the step, the expected preconditions before the step,
the actions to take during the step, and the conditions that should exist after the step.

The repository at https://github.com/nmiyake/echgo2 contains the result of walking through the tutorial.

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

It is recommended that the tutorial be run in a Docker image to provide maximal environment isolation. A base Docker
image that can be used for the tutorial is available in the `docs/templates/baseimage` directory.

The Dockerfile in the directory has some values specified using `ENV` that can be customized if desired. If you want to
run the tutorial directly on a host rather than in a Docker image, you should set/export the `ENV` values.

This tutorial uses `github.com/nmiyake/echgo2` as the project path. This is fine for the main path of the tutorial.
However, if you want to walk through publishing the repository to GitHub, you should modify this project path to be a
path to a repository that you can create/push to GitHub (although it is possible to push a project with this path to any
GitHub repository, using a project path that mirrors the GitHub location is more realistic).

Start the tutorial by building the Docker image:

```
➜ cd ${GOPATH}/src/github.com/palantir/godel/docs/templates/baseimage
➜ ./build.sh
[+] Building 7.0s (8/8) FINISHED
 => [internal] load build definition from Dockerfile                                                                  0.0s
 => => transferring dockerfile: 40B                                                                                   0.0s
 => [internal] load .dockerignore                                                                                     0.0s
 => => transferring context: 2B                                                                                       0.0s
 => [internal] load metadata for docker.io/library/golang:1.16.2                                                      1.4s
 => CACHED [1/4] FROM docker.io/library/golang:1.16.2@sha256:31447e84d4af01c218cf158072028ada82d49248fd067d1b7228857  0.0s
 => [2/4] RUN apt-get update && apt-get install -y tree                                                               5.0s
 => [3/4] RUN git config --global user.name "Tutorial User" &&     git config --global user.email "tutorial@tutorial  0.3s
 => [4/4] WORKDIR /go/src/github.com/nmiyake/echgo2                                                                   0.0s
 => exporting to image                                                                                                0.1s
 => => exporting layers                                                                                               0.1s
 => => writing image sha256:d3bc3e6f4335f5f9c94c1d259f3b28a9814fa30078e916a09c99a794b90c7942                          0.0s
 => => naming to docker.io/library/godeltutorial:setup                                                                0.0s
```

Now, run the Docker image interactively:

```
➜ docker run -it godeltutorial:setup
root@010fdbb9adec:/go/src/github.com/nmiyake/echgo2#
```

You can now follow the rest of the tutorial. The following steps are also valid on a host as long as the proper
environment variables are set and the required programs are available.

Start the tutorial by creating the directory for your project and setting it to be the working directory:

```
➜ mkdir -p ${GOPATH}/src/${PROJECT_PATH} && cd $_
➜ pwd
/go/src/github.com/nmiyake/echgo2
```

Initialize a git repository, add a README and commit it:

```
➜ git init
Initialized empty Git repository in /go/src/github.com/nmiyake/echgo2/.git/
➜ echo 'echgo2 is a program that echoes input provided by the user.' > README.md
➜ git add README.md
➜ git commit -m "Initial commit"
[master (root-commit) 711897e] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md
```

Finally, define a `go.mod` file so that a module is defined for this project:

```
➜ echo 'module github.com/nmiyake/echgo2' > go.mod
➜ git add go.mod
➜ git commit -m "Add module definition"
[master 9b2b6e9] Add module definition
 1 file changed, 1 insertion(+)
 create mode 100644 go.mod
```

Tutorial end state
------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository and Go module

Tutorial next step
------------------
[Add gödel](https://github.com/palantir/godel/wiki/Add-godel)
