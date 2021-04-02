Summary
-------
Install and run `godelinit` to add gödel to a project.

Tutorial start state
--------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository and Go module

Add gödel to a project
----------------------
The godelinit program can be used to add gödel to a new project. The godelinit program determines the latest release of
gödel and installs it into the current directory, downloading it if needed.

Install the `godelinit` program using `go get` and run the program to install gödel:

```
➜ go install github.com/palantir/godel/v2/godelinit@latest
go: downloading github.com/palantir/godel/v2 v2.36.0
go: downloading github.com/palantir/godel v2.20.0+incompatible
go: downloading github.com/spf13/cobra v1.1.3
go: downloading github.com/palantir/pkg v1.0.1
go: downloading github.com/pkg/errors v0.8.1
go: downloading github.com/cheggaaa/pb/v3 v3.0.2
go: downloading github.com/mholt/archiver/v3 v3.3.0
go: downloading github.com/nmiyake/pkg/dirs v1.0.0
go: downloading github.com/palantir/pkg/cobracli v1.0.1
go: downloading github.com/palantir/pkg/specdir v1.0.1
go: downloading github.com/rogpeppe/go-internal v1.7.0
go: downloading github.com/VividCortex/ewma v1.1.1
go: downloading github.com/fatih/color v1.7.0
go: downloading github.com/mattn/go-colorable v0.1.2
go: downloading github.com/mattn/go-isatty v0.0.8
go: downloading github.com/mattn/go-runewidth v0.0.4
go: downloading github.com/spf13/pflag v1.0.5
go: downloading github.com/nmiyake/pkg/errorstringer v1.0.0
go: downloading github.com/andybalholm/brotli v0.0.0-20190621154722-5f990b63d2d6
go: downloading github.com/dsnet/compress v0.0.1
go: downloading github.com/golang/snappy v0.0.1
go: downloading github.com/klauspost/compress v1.9.2
go: downloading github.com/klauspost/pgzip v1.2.1
go: downloading github.com/nwaples/rardecode v1.0.0
go: downloading github.com/pierrec/lz4 v2.0.5+incompatible
go: downloading github.com/ulikunitz/xz v0.5.6
go: downloading github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8
go: downloading golang.org/x/sys v0.0.0-20190624142023-c5567b49c5d0
go: downloading github.com/golang/gddo v0.0.0-20190419222130-af0f2af80721
➜ godelinit
Getting package from https://github.com/palantir/godel/releases/download/v2.36.0/godel-2.36.0.tgz...
1.11 MiB / 11.68 MiB [------>__________________________________________________________] 9.53% ? p/s3.25 MiB / 11.68 MiB [----------------->______________________________________________] 27.85% ? p/s5.21 MiB / 11.68 MiB [---------------------------->___________________________________] 44.56% ? p/s5.83 MiB / 11.68 MiB [---------------------------->____________________________] 49.91% 7.85 MiB p/s6.75 MiB / 11.68 MiB [-------------------------------->________________________] 57.80% 7.85 MiB p/s8.69 MiB / 11.68 MiB [------------------------------------------>______________] 74.39% 7.85 MiB p/s10.35 MiB / 11.68 MiB [------------------------------------------------->______] 88.56% 7.83 MiB p/s11.63 MiB / 11.68 MiB [------------------------------------------------------->] 99.53% 7.83 MiB p/s11.68 MiB / 11.68 MiB [-------------------------------------------------------] 100.00% 8.30 MiB p/s
```

Run `./godelw version` to verify that gödel was installed correctly. If this is the first run, this invocation will
download all of the plugins and assets:

```
➜ ./godelw version
Getting package from https://github.com/palantir/distgo/releases/download/v1.28.0/dist-plugin-1.28.0-linux-amd64.tgz...
115.30 KiB / 7.15 MiB [->______________________________________________________________] 1.57% ? p/s1006.09 KiB / 7.15 MiB [-------->_____________________________________________________] 13.74% ? p/s2.61 MiB / 7.15 MiB [----------------------->_________________________________________] 36.47% ? p/s4.69 MiB / 7.15 MiB [-------------------------------------->___________________] 65.54% 7.63 MiB p/s6.09 MiB / 7.15 MiB [------------------------------------------------->________] 85.21% 7.63 MiB p/s7.15 MiB / 7.15 MiB [---------------------------------------------------------] 100.00% 7.87 MiB p/s
Getting package from https://github.com/palantir/godel-format-plugin/releases/download/v1.7.0/format-plugin-1.7.0-linux-amd64.tgz...
1.82 MiB / 3.72 MiB [------------------------------->_________________________________] 48.99% ? p/s3.23 MiB / 3.72 MiB [-------------------------------------------------------->________] 86.67% ? p/s3.72 MiB / 3.72 MiB [---------------------------------------------------------] 100.00% 9.51 MiB p/s
Getting package from https://github.com/palantir/godel-format-asset-ptimports/releases/download/v1.6.0/ptimports-asset-1.6.0-linux-amd64.tgz...
1.51 MiB / 5.08 MiB [------------------->_____________________________________________] 29.76% ? p/s3.39 MiB / 5.08 MiB [------------------------------------------->_____________________] 66.85% ? p/s5.08 MiB / 5.08 MiB [--------------------------------------------------------] 100.00% 13.40 MiB p/s
Getting package from https://github.com/palantir/godel-goland-plugin/releases/download/v1.3.0/goland-plugin-1.3.0-linux-amd64.tgz...
859.26 KiB / 3.36 MiB [--------------->_______________________________________________] 24.96% ? p/s1.52 MiB / 3.36 MiB [----------------------------->___________________________________] 45.32% ? p/s3.36 MiB / 3.36 MiB [---------------------------------------------------------] 100.00% 9.08 MiB p/s
Getting package from https://github.com/palantir/okgo/releases/download/v1.10.0/check-plugin-1.10.0-linux-amd64.tgz...
1.81 MiB / 3.94 MiB [----------------------------->___________________________________] 45.87% ? p/s3.54 MiB / 3.94 MiB [---------------------------------------------------------->______] 89.68% ? p/s3.94 MiB / 3.94 MiB [--------------------------------------------------------] 100.00% 18.49 MiB p/s
Getting package from https://github.com/palantir/godel-okgo-asset-compiles/releases/download/v1.7.0/compiles-asset-1.7.0-linux-amd64.tgz...
993.27 KiB / 4.51 MiB [------------->_________________________________________________] 21.53% ? p/s1.76 MiB / 4.51 MiB [------------------------->_______________________________________] 39.08% ? p/s3.51 MiB / 4.51 MiB [-------------------------------------------------->______________] 77.88% ? p/s4.51 MiB / 4.51 MiB [---------------------------------------------------------] 100.00% 8.39 MiB p/s
Getting package from https://github.com/palantir/godel-okgo-asset-deadcode/releases/download/v1.6.0/deadcode-asset-1.6.0-linux-amd64.tgz...
1.85 MiB / 4.51 MiB [-------------------------->______________________________________] 41.08% ? p/s3.70 MiB / 4.51 MiB [----------------------------------------------------->___________] 82.13% ? p/s4.51 MiB / 4.51 MiB [--------------------------------------------------------] 100.00% 15.71 MiB p/s
Getting package from https://github.com/palantir/godel-okgo-asset-errcheck/releases/download/v1.8.0/errcheck-asset-1.8.0-linux-amd64.tgz...
1.90 MiB / 4.61 MiB [-------------------------->______________________________________] 41.22% ? p/s3.86 MiB / 4.61 MiB [------------------------------------------------------>__________] 83.78% ? p/s4.61 MiB / 4.61 MiB [--------------------------------------------------------] 100.00% 15.58 MiB p/s
Getting package from https://github.com/palantir/godel-okgo-asset-golint/releases/download/v1.4.0/golint-asset-1.4.0-linux-amd64.tgz...
1.61 MiB / 4.67 MiB [---------------------->__________________________________________] 34.50% ? p/s3.20 MiB / 4.67 MiB [-------------------------------------------->____________________] 68.51% ? p/s4.67 MiB / 4.67 MiB [--------------------------------------------------------] 100.00% 12.97 MiB p/s
Getting package from https://github.com/palantir/godel-okgo-asset-govet/releases/download/v1.4.0/govet-asset-1.4.0-linux-amd64.tgz...
1.99 MiB / 3.47 MiB [------------------------------------->___________________________] 57.39% ? p/s3.47 MiB / 3.47 MiB [--------------------------------------------------------] 100.00% 22.67 MiB p/s
Getting package from https://github.com/palantir/godel-okgo-asset-importalias/releases/download/v1.4.0/importalias-asset-1.4.0-linux-amd64.tgz...
1.89 MiB / 3.82 MiB [-------------------------------->________________________________] 49.47% ? p/s3.82 MiB / 3.82 MiB [--------------------------------------------------------] 100.00% 19.49 MiB p/s
Getting package from https://github.com/palantir/godel-okgo-asset-ineffassign/releases/download/v1.4.0/ineffassign-asset-1.4.0-linux-amd64.tgz...
1.57 MiB / 3.71 MiB [--------------------------->_____________________________________] 42.42% ? p/s2.59 MiB / 3.71 MiB [--------------------------------------------->___________________] 69.83% ? p/s3.71 MiB / 3.71 MiB [--------------------------------------------------------] 100.00% 11.52 MiB p/s
Getting package from https://github.com/palantir/godel-okgo-asset-outparamcheck/releases/download/v1.8.0/outparamcheck-asset-1.8.0-linux-amd64.tgz...
1.28 MiB / 4.60 MiB [------------------>______________________________________________] 27.85% ? p/s2.65 MiB / 4.60 MiB [------------------------------------->___________________________] 57.50% ? p/s4.30 MiB / 4.60 MiB [------------------------------------------------------------>____] 93.38% ? p/s4.60 MiB / 4.60 MiB [--------------------------------------------------------] 100.00% 11.21 MiB p/s
Getting package from https://github.com/palantir/godel-okgo-asset-unconvert/releases/download/v1.7.0/unconvert-asset-1.7.0-linux-amd64.tgz...
1.29 MiB / 4.73 MiB [----------------->_______________________________________________] 27.28% ? p/s1.91 MiB / 4.73 MiB [-------------------------->______________________________________] 40.39% ? p/s2.96 MiB / 4.73 MiB [---------------------------------------->________________________] 62.48% ? p/s3.45 MiB / 4.73 MiB [------------------------------------------>_______________] 72.96% 3.60 MiB p/s3.95 MiB / 4.73 MiB [------------------------------------------------>_________] 83.45% 3.60 MiB p/s4.49 MiB / 4.73 MiB [------------------------------------------------------->__] 94.88% 3.60 MiB p/s4.73 MiB / 4.73 MiB [---------------------------------------------------------] 100.00% 4.45 MiB p/s
Getting package from https://github.com/palantir/godel-okgo-asset-varcheck/releases/download/v1.6.0/varcheck-asset-1.6.0-linux-amd64.tgz...
538.27 KiB / 4.49 MiB [------->_______________________________________________________] 11.70% ? p/s947.27 KiB / 4.49 MiB [------------>__________________________________________________] 20.58% ? p/s1.36 MiB / 4.49 MiB [------------------->_____________________________________________] 30.19% ? p/s1.80 MiB / 4.49 MiB [----------------------->__________________________________] 40.16% 2.14 MiB p/s2.25 MiB / 4.49 MiB [----------------------------->____________________________] 50.14% 2.14 MiB p/s2.59 MiB / 4.49 MiB [--------------------------------->________________________] 57.70% 2.14 MiB p/s2.94 MiB / 4.49 MiB [------------------------------------->____________________] 65.48% 2.12 MiB p/s3.46 MiB / 4.49 MiB [-------------------------------------------->_____________] 77.01% 2.12 MiB p/s4.49 MiB / 4.49 MiB [---------------------------------------------------------] 100.00% 2.93 MiB p/s
Getting package from https://github.com/palantir/godel-license-plugin/releases/download/v1.5.0/license-plugin-1.5.0-linux-amd64.tgz...
1.48 MiB / 3.66 MiB [-------------------------->______________________________________] 40.33% ? p/s2.94 MiB / 3.66 MiB [---------------------------------------------------->____________] 80.46% ? p/s3.66 MiB / 3.66 MiB [--------------------------------------------------------] 100.00% 12.19 MiB p/s
Getting package from https://github.com/palantir/godel-test-plugin/releases/download/v1.6.0/test-plugin-1.6.0-linux-amd64.tgz...
497.97 KiB / 4.06 MiB [------->_______________________________________________________] 11.97% ? p/s1.69 MiB / 4.06 MiB [--------------------------->_____________________________________] 41.69% ? p/s3.19 MiB / 4.06 MiB [--------------------------------------------------->_____________] 78.47% ? p/s4.06 MiB / 4.06 MiB [---------------------------------------------------------] 100.00% 8.14 MiB p/s
godel version 2.36.0
```

Technically, this is sufficient and we have a working gödel installation. However, because the installation was
performed by downloading a distribution, the gödel installation itself does not have a checksum set (no verification was
performed on the package itself):

```
➜ cat godel/config/godel.properties
distributionURL=https://github.com/palantir/godel/releases/download/v2.36.0/godel-2.36.0.tgz
distributionSHA256=
```

If we are concerned about integrity, we can specify the expected checksum for the package when using godelinit. The
expected checksum should be computed based on a trusted distribution such as the
[GitHub release](https://github.com/palantir/godel/releases). Install the distribution using the checksum:

```
➜ godelinit --checksum ${GODEL_CHECKSUM}
```

If the installation succeeds with a checksum specified, it is set in the properties:

```
➜ cat godel/config/godel.properties
distributionURL=https://github.com/palantir/godel/releases/download/v2.36.0/godel-2.36.0.tgz
distributionSHA256=91137f4fb9e1b4491d6dd821edf6ed39eb66f21410bf645a062f687049c45492
```

Commit the changes to the repository:

```
➜ git add godel godelw
➜ git commit -m "Add godel to project"
[master fcc55fa] Add godel to project
 8 files changed, 260 insertions(+)
 create mode 100644 godel/config/check-plugin.yml
 create mode 100644 godel/config/dist-plugin.yml
 create mode 100644 godel/config/format-plugin.yml
 create mode 100644 godel/config/godel.properties
 create mode 100644 godel/config/godel.yml
 create mode 100644 godel/config/license-plugin.yml
 create mode 100644 godel/config/test-plugin.yml
 create mode 100755 godelw
```

gödel has now been added to the project and is ready to use.

Tutorial end state
------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository and Go module
* Project contains `godel` and `godelw`

Tutorial next step
------------------
[Add Git hooks to enforce formatting](https://github.com/palantir/godel/wiki/Add-git-hooks)

More
----
### Download and install gödel manually
It is possible to add gödel to a project without using godelinit by downloading and installing it manually.

Download the distribution into a temporary directory and expand it:

```
➜ mkdir -p download
➜ curl -L "https://github.com/palantir/godel/releases/download/v${GODEL_VERSION}/godel-${GODEL_VERSION}.tgz" -o download/godel-"${GODEL_VERSION}".tgz
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0100   622  100   622    0     0   2286      0 --:--:-- --:--:-- --:--:--  2286
  0 11.6M    0 16518    0     0  35986      0  0:05:40 --:--:--  0:05:40 35986 59 11.6M   59 7104k    0     0  4843k      0  0:00:02  0:00:01  0:00:01 7039k100 11.6M  100 11.6M    0     0  5893k      0  0:00:02  0:00:02 --:--:-- 7610k
➜ tar -xf download/godel-"${GODEL_VERSION}".tgz -C download
```

Copy the contents of the wrapper directory (which contains godelw and godel) to the project directory:

```
➜ cp -r download/godel-"${GODEL_VERSION}"/wrapper/* .
```

Run `./godelw version` to verify that gödel was installed correctly. This command will download the distribution if the
distribution has not previously been installed locally. However, if you have followed the tutorial steps, this version
has already been downloaded so it will not re-download:

```
➜ ./godelw version
godel version 2.36.0
```

Technically, this is sufficient and we have a working gödel install. However, distributions that are downloaded manually
do not have a checksum set in `godel/config/godel.properties`:

```
➜ cat godel/config/godel.properties
distributionURL=https://github.com/palantir/godel/releases/download/v2.36.0/godel-2.36.0.tgz
distributionSHA256=
```

For completeness, set the checksum. The checksum should be computed based on the distribution available at
https://github.com/palantir/godel/releases.

```
➜ DIST_URL='distributionURL=https://github.com/palantir/godel/releases/download/vGODEL_VERSION/godel-GODEL_VERSION.tgz
distributionSHA256=GODEL_CHECKSUM'; DIST_URL="${DIST_URL//GODEL_VERSION/$GODEL_VERSION}"; DIST_URL="${DIST_URL//GODEL_CHECKSUM/$GODEL_CHECKSUM}"
➜ echo "${DIST_URL}" > godel/config/godel.properties
➜ unset DIST_URL
```

Now that gödel has been added to the project, remove the temporary directory and unset the version variable:

```
➜ rm -rf download
```

### Copying gödel from an existing project
If you have local projects that already use gödel, you can add gödel to a another project by copying the `godelw` and
`godel/config/godel.properties` files from the project that already has gödel.

For example, assume that `{$GOPATH}/src/${PROJECT_PATH}` exists in the current state of the tutorial and we want
to create another project at `${GOPATH}/src/github.com/nmiyake/sample` that also uses gödel. This can be done as follows:

```
➜ mkdir -p ${GOPATH}/src/github.com/nmiyake/sample && cd $_
➜ cp ${GOPATH}/src/${PROJECT_PATH}/godelw .
➜ mkdir -p godel/config
➜ cp ${GOPATH}/src/${PROJECT_PATH}/godel/config/godel.properties godel/config/
➜ ./godelw update --sync
```

Verify that invoking `./godelw` works and that the `godel/config` directory has been populated with the default
configuration files:

```
➜ ./godelw version
godel version 2.36.0
➜ ls godel/config
check-plugin.yml
dist-plugin.yml
format-plugin.yml
godel.properties
godel.yml
license-plugin.yml
test-plugin.yml
```

Restore the workspace to the original state by setting the working directory back to the `echgo2` project and removing
the sample project:

```
➜ cd ${GOPATH}/src/${PROJECT_PATH}
➜ rm -rf ${GOPATH}/src/github.com/nmiyake/sample
```

The steps above take the approach of copying the `godelw` file and `godel/config/godel.properties` file and running
`update` to ensure that all of the configuration is in its default state -- if the entire `godel/config` directory was
copied from the other project, then the current project would copy the other project's configuration as well, which is
most likely not the correct thing to do.

The `update` task updates gödel to the latest release version. Flags can be used to configure the behavior to perform
different operations (for example, to update gödel to a specific version or to the version specified in
`godel/config/godel.properties`). If the distribution is already downloaded locally in the local global gödel directory
(`~/.godel` by default) and matches the checksum specified in `godel.properties`, the required files are copied from the
local cache. If the distribution does not exist in the local cache or a checksum is not specified in `godel.properties`,
then the distribution is downloaded from the distribution URL.

### Specify the checksum in `godel/config/godel.properties`
If the `godelw` script cannot locate the `godel` binary, it downloads it from the URL specified in
`godel/godel.properties`. If a checksum is specified in `godel/config/godel.properties`, it is used to verify the
integrity and validity of the distribution downloaded from the URL. The checksum is also used by the `update` task -- if
the `update` task for a given version can find a download of the distribution in the downloads directory with a matching
checksum, it will be used to update the version of gödel (otherwise, the `update` task will always download the
distribution from the specified distribution URL).

The checksum is the SHA-256 checksum of the `.tgz` archive of the product. The checksum is the SHA-256 digest of the TGZ
distribution. It can be computed on a trusted distribution such as the distribution from the
[GitHub Release](https://github.com/palantir/godel/releases).

If a user has obtained a trusted distribution locally, it is also possible to manually compute the checksum.

Download the distribution by running the following:

```
➜ mkdir -p download
➜ curl -L "https://github.com/palantir/godel/releases/download/v${GODEL_VERSION}/godel-${GODEL_VERSION}.tgz" -o download/godel-"${GODEL_VERSION}".tgz
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0100   622  100   622    0     0   3222      0 --:--:-- --:--:-- --:--:--  3222
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0 66 11.6M   66 7911k    0     0  6463k      0  0:00:01  0:00:01 --:--:-- 7895k100 11.6M  100 11.6M    0     0  7435k      0  0:00:01  0:00:01 --:--:-- 8632k
```

The checksum can be computed using `openssl` or `shasum` as follows:

```
➜ openssl dgst -sha256 download/godel-"${GODEL_VERSION}".tgz
SHA256(download/godel-2.36.0.tgz)= 91137f4fb9e1b4491d6dd821edf6ed39eb66f21410bf645a062f687049c45492
➜ shasum -a 256 download/godel-"${GODEL_VERSION}".tgz
91137f4fb9e1b4491d6dd821edf6ed39eb66f21410bf645a062f687049c45492  download/godel-2.36.0.tgz
```

Run the following to clean up our state:

```
➜ rm -rf download
```
