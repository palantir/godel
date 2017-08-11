Summary
-------
gödel tasks can be configured to run in a CI environment to verify, build and publish products.

Tutorial start state
--------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist.yml` is configured to build `echgo`
* Project is tagged as 0.0.1
* `godel/config/dist.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* Go files have license headers
* `godel/config/generate.yml` is configured to generate string function
* `godel/config/exclude.yml` is configured to ignore all `.+_string.go` files
* `integration_test` contains integration tests
* `godel/config/test.yml` is configured to specify the "integration" tag
* `docs` contains documentation

([Link](https://github.com/nmiyake/echgo/tree/17c7406291096306e92c6f82da2df09388766693))

CI setup
--------

Now that we have set up a project and a repository, we will configure CI (continuous integration) to verify that all of
the PRs for our project properly pass verification and so that artifacts are published for releases.

We will use CircleCI to set up CI for this project. Run the following to create a CircleCI configuration file:

```
➜ mkdir -p .circleci
➜ echo 'defaults: &defaults
  working_directory: /go/src/github.com/nmiyake/echgo
  docker:
    - image: golang:1.9.1

version: 2
jobs:
  build:
    <<: *defaults
    steps:
      - type: checkout
      - type: cache-restore
        key: godel-{{ checksum "godelw" }}
      - type: run
        name: "Verify godel version"
        command: ./godelw version
      - type: cache-save
        key: godel-{{ checksum "godelw" }}
        paths:
          - /root/.godel
      - type: run
        name: "Verify Go version"
        command: go version
      - type: run
        name: "Install project packages"
        command: go install $(./godelw packages)
      - type: run
        name: "Create test output directory"
        command: mkdir -p /tmp/test-results/"${CIRCLE_PROJECT_REPONAME}"
      - type: run
        name: "Run godel verification"
        command: ./godelw verify --apply=false --junit-output="/tmp/test-results/${CIRCLE_PROJECT_REPONAME}-tests.xml"
      - type: test-results-store
        path: /tmp/test-results
      - type: artifacts-store
        path: /tmp/test-results
        destination: test-results
      - type: run
        name: "Create distribution"
        command: ./godelw dist
      - type: artifacts-store
        path: /go/src/github.com/nmiyake/echgo/dist
  wiki:
    <<: *defaults
    steps:
      - type: checkout
      - type: cache-restore
        key: godel-{{ checksum "godelw" }}
      - type: run
        name: "Verify godel version"
        command: ./godelw version
      - type: run
        name: "Update GitHub Wiki on master branch"
        command: ./godelw github-wiki --docs-dir docs --repository=git@github.com:nmiyake/echgo.wiki.git
  publish:
    <<: *defaults
    steps:
      - type: checkout
      - type: cache-restore
        key: godel-{{ checksum "godelw" }}
      - type: run
        name: "Verify godel version"
        command: ./godelw version
      - type: run
        name: "Publish"
        command: ./godelw publish github --url https://api.github.com --user nmiyake --password $GITHUB_TOKEN --owner nmiyake --repository echgo

workflows:
  version: 2
  build-deploy:
    jobs:
      - build:
          filters:
            tags:
              only: /.*/
      - wiki:
          requires:
            - build
          filters:
            branches:
              only: master
      - publish:
          requires:
            - build
          filters:
            tags:
              only: /^[0-9]+(\.[0-9]+)+(-rc[0-9]+)?$/
            branches:
              ignore: /.*/' > .circleci/config.yml
```

The primary tasks performed by this CI are the following:

* Runs `./godelw version` to ensure that the gödel distribution is downloaded and configured in the CI environment
  * This configuration caches the distribution so that it is only downloaded when the version changes
* Runs `./godelw verify` with the `--apply=false` and `--junit-output=<path>` flags
  * This ensures that the code passes all of the required checks, runs all tests and saves the test output as a JUnit XML file
* Runs `./godelw dist` to create the distribution
* Runs `./godelw github-wiki` on the "master" branch to update documentation
  * Runs only on the "master" branch to ensure that only one branch is publishing documentation
* Runs `./godelw publish` on release tags

Commit the changes to the repository by running the following:

```
➜ git add .circleci
➜ git commit -m "Add CircleCI configuration"
[master 25d27eb] Add CircleCI configuration
 1 file changed, 89 insertions(+)
 create mode 100644 .circleci/config.yml
```

You can now configure the GitHub project to be run using CircleCI and it will run the CI process.

Verify that everything is working as expected by tagging a 1.0.0, pushing the tag and verifying that the tag kicks off a
build and publishes the artifacts:

```
➜ git push origin master
Counting objects: 4, done.
Delta compression using up to 8 threads.
Compressing objects: 100% (3/3), done.
Writing objects: 100% (4/4), 1.02 KiB | 0 bytes/s, done.
Total 4 (delta 1), reused 0 (delta 0)
remote: Resolving deltas: 100% (1/1), completed with 1 local object.
To git@github.com:nmiyake/echgo.git
   17c7406..25d27eb  master -> master
➜ git tag 1.0.0
➜ git push origin --tags
Total 0 (delta 0), reused 0 (delta 0)
To git@github.com:nmiyake/echgo.git
 * [new tag]         1.0.0 -> 1.0.0
```

Although this example was for CircleCI 2.0, the general principles/steps should be applicable in any CI system.

Tutorial end state
------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist.yml` is configured to build `echgo`
* Project is tagged as 0.0.1
* `godel/config/dist.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* Go files have license headers
* `godel/config/generate.yml` is configured to generate string function
* `godel/config/exclude.yml` is configured to ignore all `.+_string.go` files
* `integration_test` contains integration tests
* `godel/config/test.yml` is configured to specify the "integration" tag
* `docs` contains documentation
* `.circleci/config.yml` exists
* Project is tagged as 1.0.0

([Link](https://github.com/nmiyake/echgo/tree/25d27eb1763e55f228282594691798ca0c2bbe28))

Tutorial next step
------------------
[Update gödel](https://github.com/palantir/godel/wiki/Update-g%C3%B6del)

More
----

### CircleCI 2.0 without workflows

```yaml
jobs:
  build:
    working_directory: /go/src/github.com/nmiyake/echgo
    docker:
      - image: golang:1.9.1
    steps:
      - type: checkout
      - type: cache-restore
        key: godel-{{ checksum "godelw" }}
      - type: shell
        name: "Verify godel version"
        command: ./godelw version
      - type: cache-save
        key: godel-{{ checksum "godelw" }}
        paths:
          - /root/.godel
      - type: shell
        name: "Verify Go version"
        command: go version
      - type: shell
        name: "Install project packages"
        command: go install $(./godelw packages)
      - type: shell
        name: "Create test output directory"
        command: mkdir -p /tmp/test-results/"${CIRCLE_PROJECT_REPONAME}"
      - type: shell
        name: "Run godel verification"
        command: ./godelw verify --apply=false --junit-output="/tmp/test-results/${CIRCLE_PROJECT_REPONAME}-tests.xml"
      - type: test-results-store
        path: /tmp/test-results
      - type: artifacts-store
        path: /tmp/test-results
        destination: test-results
      - type: shell
        name: "Create distribution"
        command: ./godelw dist
      - type: artifacts-store
        path: /go/src/github.com/nmiyake/echgo/dist
      - type: deploy
        name: "Update GitHub Wiki on master branch"
        command: |
          set -eu
          if [ "${CIRCLE_BRANCH}" == "master" ]; then
            ./godelw github-wiki --docs-dir docs --repository=git@github.com:nmiyake/echgo.wiki.git
          else
            echo "Not master branch: skipping wiki publish"
          fi
      - type: deploy
        name: "Publish on release tags"
        command: |
          set -eu
          TAG=$(./godelw project-version)
          if [[ $TAG =~ ^[0-9]+(\.[0-9]+)+(-rc[0-9]+)?$ ]]; then
            ./godelw publish github --url https://api.github.com --user nmiyake --password $GITHUB_TOKEN --owner nmiyake --repository echgo
          else
            echo "Not a release tag: skipping publish"
          fi
```

### CircleCI 2.0 with workflows

```yaml
defaults: &defaults
  working_directory: /go/src/github.com/nmiyake/echgo
  docker:
    - image: golang:1.9.1

version: 2
jobs:
  build:
    <<: *defaults
    steps:
      - type: checkout
      - type: cache-restore
        key: godel-{{ checksum "godelw" }}
      - type: run
        name: "Verify godel version"
        command: ./godelw version
      - type: cache-save
        key: godel-{{ checksum "godelw" }}
        paths:
          - /root/.godel
      - type: run
        name: "Verify Go version"
        command: go version
      - type: run
        name: "Install project packages"
        command: go install $(./godelw packages)
      - type: run
        name: "Create test output directory"
        command: mkdir -p /tmp/test-results/"${CIRCLE_PROJECT_REPONAME}"
      - type: run
        name: "Run godel verification"
        command: ./godelw verify --apply=false --junit-output="/tmp/test-results/${CIRCLE_PROJECT_REPONAME}-tests.xml"
      - type: test-results-store
        path: /tmp/test-results
      - type: artifacts-store
        path: /tmp/test-results
        destination: test-results
      - type: run
        name: "Create distribution"
        command: ./godelw dist
      - type: artifacts-store
        path: /go/src/github.com/nmiyake/echgo/dist
  wiki:
    <<: *defaults
    steps:
      - type: checkout
      - type: cache-restore
        key: godel-{{ checksum "godelw" }}
      - type: run
        name: "Verify godel version"
        command: ./godelw version
      - type: run
        name: "Update GitHub Wiki on master branch"
        command: ./godelw github-wiki --docs-dir docs --repository=git@github.com:nmiyake/echgo.wiki.git
  publish:
    <<: *defaults
    steps:
      - type: checkout
      - type: cache-restore
        key: godel-{{ checksum "godelw" }}
      - type: run
        name: "Verify godel version"
        command: ./godelw version
      - type: run
        name: "Publish"
        command: ./godelw publish github --url https://api.github.com --user nmiyake --password $GITHUB_TOKEN --owner nmiyake --repository echgo

workflows:
  version: 2
  build-deploy:
    jobs:
      - build:
          filters:
            tags:
              only: /.*/
      - wiki:
          requires:
            - build
          filters:
            branches:
              only: master
      - publish:
          requires:
            - build
          filters:
            tags:
              only: /^[0-9]+(\.[0-9]+)+(-rc[0-9]+)?$/
            branches:
              ignore: /.*/
```

### CircleCI 1.0

```yaml
machine:
  environment:
    GODIST: "go1.9.1.linux-amd64.tar.gz"
    GOPATH: "$HOME/.go_workspace"
    IMPORT_PATH: "github.com/$CIRCLE_PROJECT_USERNAME/$CIRCLE_PROJECT_REPONAME"
    GO_PROJECT_SRC_PATH: "$GOPATH/src/$IMPORT_PATH"
  post:
    - mkdir -p download
    - test -e download/$GODIST || wget -O "download/$GODIST" "https://storage.googleapis.com/golang/$GODIST"
    # create custom Go distribution with packages built for darwin-amd64 if it is not present
    - |
      if [ ! -e download/$GODIST-custom.tgz ]; then
          sudo rm -rf /usr/local/go && \
          sudo tar -C /usr/local -xzf download/$GODIST && \
          sudo env GOOS=darwin GOARCH=amd64 /usr/local/go/bin/go install std && \
          tar -C /usr/local -czf download/$GODIST-custom.tgz go
      fi
    - sudo rm -rf /usr/local/go
    - sudo tar -C /usr/local -xzf download/$GODIST-custom.tgz

checkout:
  post:
    # ensure all tags are fetched and up-to-date
    - git tag -l | xargs git tag -d && git fetch -t

dependencies:
  override:
    - mkdir -p "$GOPATH/src/$IMPORT_PATH"
    - rsync -az --delete ./ "$GOPATH/src/$IMPORT_PATH/"
    - cd "$GO_PROJECT_SRC_PATH" && ./godelw version
  cache_directories:
    - ~/.godel
    - ~/download

test:
  pre:
    - go version
    - go get golang.org/x/tools/cmd/stringer
  override:
    - cd "$GO_PROJECT_SRC_PATH" && go install $(./godelw packages)
    - cd "$GO_PROJECT_SRC_PATH" && mkdir -p "$CIRCLE_TEST_REPORTS/$CIRCLE_PROJECT_REPONAME"
    - cd "$GO_PROJECT_SRC_PATH" && ./godelw verify --apply=false --junit-output="$CIRCLE_TEST_REPORTS/$CIRCLE_PROJECT_REPONAME/$CIRCLE_PROJECT_REPONAME-tests.xml"
    - cd "$GO_PROJECT_SRC_PATH" && ./godelw dist

deployment:
  master:
    branch: master
    commands:
      - cd "$GO_PROJECT_SRC_PATH" && ./godelw github-wiki --docs-dir docs --repository=git@github.com:nmiyake/echgo.wiki.git
  release:
    tag: /[0-9]+(\.[0-9]+)+(-rc[0-9]+)?/
    commands:
      - cd "$GO_PROJECT_SRC_PATH" && ./godelw publish github --url https://api.github.com --user nmiyake --password $GITHUB_TOKEN --owner nmiyake --repository echgo
```
