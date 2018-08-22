Summary
-------
Install and run `godelinit` to add gödel to a project.

Tutorial start state
--------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository

Add gödel to a project
----------------------
The godelinit program can be used to add gödel to a new project. The godelinit program determines the latest release of
gödel and installs it into the current directory, downloading it if needed.

Install the `godelinit` program using `go get` and run the program to install gödel:

```
➜ go get github.com/palantir/godel/godelinit
➜ godelinit
Getting package from https://palantir.bintray.com/releases/com/palantir/godel/godel/2.9.1/godel-2.9.1.tgz...
 0 B / 9.33 MiB    0.00% 639.10 KiB / 9.33 MiB    6.69% 2s 1.58 MiB / 9.33 MiB   16.89% 1s 1.58 MiB / 9.33 MiB   16.89% 1s 2.03 MiB / 9.33 MiB   21.72% 2s 2.03 MiB / 9.33 MiB   21.72% 2s 2.62 MiB / 9.33 MiB   28.13% 3s 3.77 MiB / 9.33 MiB   40.37% 2s 4.99 MiB / 9.33 MiB   53.48% 1s 6.14 MiB / 9.33 MiB   65.84% 7.25 MiB / 9.33 MiB   77.67% 8.50 MiB / 9.33 MiB   91.07% 9.33 MiB / 9.33 MiB  100.00% 2s
```

Run `./godelw version` to verify that gödel was installed correctly. If this is the first run, this invocation will
download all of the plugins and assets:

```
➜ ./godelw version
Getting package from https://palantir.bintray.com/releases/com/palantir/distgo/dist-plugin/1.13.1/dist-plugin-1.13.1-linux-amd64.tgz...
 0 B / 4.76 MiB    0.00% 750.25 KiB / 4.76 MiB   15.38% 1s 1000.77 KiB / 4.76 MiB   20.52% 1s 1.30 MiB / 4.76 MiB   27.36% 1s 1.33 MiB / 4.76 MiB   27.93% 2s 1.81 MiB / 4.76 MiB   37.96% 1s 2.98 MiB / 4.76 MiB   62.50% 4.08 MiB / 4.76 MiB   85.57% 4.76 MiB / 4.76 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-format-plugin/format-plugin/1.1.1/format-plugin-1.1.1-linux-amd64.tgz...
 0 B / 3.32 MiB    0.00% 388.54 KiB / 3.32 MiB   11.43% 1s 778.25 KiB / 3.32 MiB   22.90% 1s 889.59 KiB / 3.32 MiB   26.17% 1s 1.03 MiB / 3.32 MiB   31.09% 1s 1.11 MiB / 3.32 MiB   33.54% 2s 1.62 MiB / 3.32 MiB   48.76% 1s 2.47 MiB / 3.32 MiB   74.49% 3.25 MiB / 3.32 MiB   97.90% 3.32 MiB / 3.32 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-format-asset-ptimports/ptimports-asset/1.1.0/ptimports-asset-1.1.0-linux-amd64.tgz...
 0 B / 3.60 MiB    0.00% 125.20 KiB / 3.60 MiB    3.39% 5s 304.05 KiB / 3.60 MiB    8.24% 4s 471.07 KiB / 3.60 MiB   12.77% 4s 665.92 KiB / 3.60 MiB   18.06% 3s 821.10 KiB / 3.60 MiB   22.26% 3s 821.10 KiB / 3.60 MiB   22.26% 3s 1015.95 KiB / 3.60 MiB   27.55% 3s 1.09 MiB / 3.60 MiB   30.13% 3s 1.22 MiB / 3.60 MiB   33.91% 3s 1.28 MiB / 3.60 MiB   35.42% 3s 1.37 MiB / 3.60 MiB   38.11% 3s 1.47 MiB / 3.60 MiB   40.70% 3s 1.57 MiB / 3.60 MiB   43.72% 3s 1.73 MiB / 3.60 MiB   47.92% 3s 1.74 MiB / 3.60 MiB   48.25% 3s 1.87 MiB / 3.60 MiB   52.02% 2s 2.22 MiB / 3.60 MiB   61.51% 2s 2.68 MiB / 3.60 MiB   74.34% 1s 3.17 MiB / 3.60 MiB   87.93% 3.60 MiB / 3.60 MiB  100.00% 3s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-goland-plugin/goland-plugin/1.0.0/goland-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.09 MiB    0.00% 459.22 KiB / 3.09 MiB   14.53% 1s 1.00 MiB / 3.09 MiB   32.51% 1.57 MiB / 3.09 MiB   51.00% 2.28 MiB / 3.09 MiB   73.90% 2.93 MiB / 3.09 MiB   95.03% 3.09 MiB / 3.09 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/okgo/check-plugin/1.1.1/check-plugin-1.1.1-linux-amd64.tgz...
 0 B / 3.53 MiB    0.00% 667.15 KiB / 3.53 MiB   18.48% 1.36 MiB / 3.53 MiB   38.53% 1.70 MiB / 3.53 MiB   48.23% 2.36 MiB / 3.53 MiB   67.06% 3.26 MiB / 3.53 MiB   92.51% 3.53 MiB / 3.53 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-compiles/compiles-asset/1.0.0/compiles-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.71 MiB    0.00% 737.98 KiB / 3.71 MiB   19.44% 1.44 MiB / 3.71 MiB   38.82% 2.19 MiB / 3.71 MiB   59.04% 3.10 MiB / 3.71 MiB   83.55% 3.71 MiB / 3.71 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-deadcode/deadcode-asset/1.0.0/deadcode-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.73 MiB    0.00% 1015.66 KiB / 3.73 MiB   26.59% 1.87 MiB / 3.73 MiB   50.23% 2.98 MiB / 3.73 MiB   79.80% 3.73 MiB / 3.73 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-errcheck/errcheck-asset/1.1.1/errcheck-asset-1.1.1-linux-amd64.tgz...
 0 B / 3.82 MiB    0.00% 861.31 KiB / 3.82 MiB   22.03% 1.81 MiB / 3.82 MiB   47.36% 2.64 MiB / 3.82 MiB   69.03% 3.70 MiB / 3.82 MiB   96.80% 3.82 MiB / 3.82 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-extimport/extimport-asset/1.1.0/extimport-asset-1.1.0-linux-amd64.tgz...
 0 B / 3.36 MiB    0.00% 850.30 KiB / 3.36 MiB   24.68% 1.79 MiB / 3.36 MiB   53.31% 2.73 MiB / 3.36 MiB   81.25% 3.36 MiB / 3.36 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-golint/golint-asset/1.0.0/golint-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.88 MiB    0.00% 737.55 KiB / 3.88 MiB   18.58% 1.70 MiB / 3.88 MiB   43.82% 2.63 MiB / 3.88 MiB   67.95% 3.74 MiB / 3.88 MiB   96.40% 3.88 MiB / 3.88 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-govet/govet-asset/1.0.0/govet-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.16 MiB    0.00% 945.32 KiB / 3.16 MiB   29.18% 1.96 MiB / 3.16 MiB   61.83% 2.88 MiB / 3.16 MiB   91.04% 3.16 MiB / 3.16 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-importalias/importalias-asset/1.0.0/importalias-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.38 MiB    0.00% 710.42 KiB / 3.38 MiB   20.54% 1.66 MiB / 3.38 MiB   49.06% 2.64 MiB / 3.38 MiB   78.04% 3.38 MiB / 3.38 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-ineffassign/ineffassign-asset/1.0.0/ineffassign-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.35 MiB    0.00% 973.45 KiB / 3.35 MiB   28.37% 1.87 MiB / 3.35 MiB   55.96% 2.90 MiB / 3.35 MiB   86.44% 3.35 MiB / 3.35 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-novendor/novendor-asset/1.0.0/novendor-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.41 MiB    0.00% 889.89 KiB / 3.41 MiB   25.45% 1.60 MiB / 3.41 MiB   46.95% 2.66 MiB / 3.41 MiB   77.99% 3.41 MiB / 3.41 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-outparamcheck/outparamcheck-asset/1.1.1/outparamcheck-asset-1.1.1-linux-amd64.tgz...
 0 B / 3.84 MiB    0.00% 944.91 KiB / 3.84 MiB   24.05% 1.93 MiB / 3.84 MiB   50.26% 2.91 MiB / 3.84 MiB   75.76% 3.84 MiB / 3.84 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-unconvert/unconvert-asset/1.0.0/unconvert-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.91 MiB    0.00% 917.13 KiB / 3.91 MiB   22.92% 1.86 MiB / 3.91 MiB   47.67% 2.80 MiB / 3.91 MiB   71.62% 3.67 MiB / 3.91 MiB   93.89% 3.91 MiB / 3.91 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-varcheck/varcheck-asset/1.0.0/varcheck-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.75 MiB    0.00% 943.68 KiB / 3.75 MiB   24.60% 1.93 MiB / 3.75 MiB   51.44% 2.88 MiB / 3.75 MiB   76.84% 3.75 MiB / 3.75 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-license-plugin/license-plugin/1.0.0/license-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.30 MiB    0.00% 971.85 KiB / 3.30 MiB   28.72% 1.90 MiB / 3.30 MiB   57.52% 2.85 MiB / 3.30 MiB   86.31% 3.30 MiB / 3.30 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-test-plugin/test-plugin/1.0.0/test-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.60 MiB    0.00% 777.71 KiB / 3.60 MiB   21.11% 1.79 MiB / 3.60 MiB   49.82% 2.79 MiB / 3.60 MiB   77.45% 3.60 MiB / 3.60 MiB  100.00% 0s
godel version 2.9.1
```

Technically, this is sufficient and we have a working gödel installation. However, because the installation was
performed by downloading a distribution, the gödel installation itself does not have a checksum set (no verification was
performed on the package itself):

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.9.1/godel-2.9.1.tgz
distributionSHA256=
```

If we are concerned about integrity, we can specify the expected checksum for the package when using godelinit. The
expected checksums for gödel releases can be found on its [Bintray](https://bintray.com/palantir/releases/godel) page.
Install the distribution using the checksum:

```
➜ godelinit --checksum ${GODEL_CHECKSUM}
```

If the installation succeeds with a checksum specified, it is set in the properties:

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.9.1/godel-2.9.1.tgz
distributionSHA256=240052b05e96e95b3f9bae3f89a091ba5dc5ec808d6a8d9cf086be92f7cdd31c
```

Commit the changes to the repository:

```
➜ git add godel godelw
➜ git commit -m "Add godel to project"
[master 6e60206] Add godel to project
 8 files changed, 243 insertions(+)
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
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
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
➜ curl -L "https://palantir.bintray.com/releases/com/palantir/godel/godel/${GODEL_VERSION}/godel-${GODEL_VERSION}.tgz" -o download/godel-"${GODEL_VERSION}".tgz
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0 32 9552k   32 3060k    0     0  1855k      0  0:00:05  0:00:01  0:00:04 2787k 75 9552k   75 7252k    0     0  2850k      0  0:00:03  0:00:02  0:00:01 3637k100 9552k  100 9552k    0     0  3285k      0  0:00:02  0:00:02 --:--:-- 4054k
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
godel version 2.9.1
```

Technically, this is sufficient and we have a working gödel install. However, distributions that are downloaded manually
do not have a checksum set in `godel/config/godel.properties`:

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.9.1/godel-2.9.1.tgz
distributionSHA256=
```

For completeness, set the checksum. The expected checksums for releases are listed on the download page at
https://bintray.com/palantir/releases/godel.

```
➜ DIST_URL='distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/GODEL_VERSION/godel-GODEL_VERSION.tgz
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
godel version 2.9.1
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

The checksum is the SHA-256 checksum of the `.tgz` archive of the product. The checksum can be obtained from the
Bintray page for the distribution:

![SHA checksum](images/tutorial/sha_checksum.png)

If a user has obtained a trusted distribution locally, it is also possible to manually compute the checksum.

Download the distribution by running the following:

```
➜ mkdir -p download
➜ curl -L "https://palantir.bintray.com/releases/com/palantir/godel/godel/${GODEL_VERSION}/godel-${GODEL_VERSION}.tgz" -o download/godel-"${GODEL_VERSION}".tgz
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
 28 9552k   28 2726k    0     0  2273k      0  0:00:04  0:00:01  0:00:03 2273k 80 9552k   80 7709k    0     0  3503k      0  0:00:02  0:00:02 --:--:-- 4977k100 9552k  100 9552k    0     0  3657k      0  0:00:02  0:00:02 --:--:-- 4834k
```

The checksum can be computed using `openssl` or `shasum` as follows:

```
➜ openssl dgst -sha256 download/godel-"${GODEL_VERSION}".tgz
SHA256(download/godel-2.9.1.tgz)= 240052b05e96e95b3f9bae3f89a091ba5dc5ec808d6a8d9cf086be92f7cdd31c
➜ shasum -a 256 download/godel-"${GODEL_VERSION}".tgz
240052b05e96e95b3f9bae3f89a091ba5dc5ec808d6a8d9cf086be92f7cdd31c  download/godel-2.9.1.tgz
```

Run the following to clean up our state:

```
➜ rm -rf download
```
