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
Getting package from https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0-rc12/godel-2.0.0-rc12.tgz...
 0 B / 9.32 MiB    0.00% 134.82 KiB / 9.32 MiB    1.41% 14s 365.88 KiB / 9.32 MiB    3.83% 10s 620.17 KiB / 9.32 MiB    6.50% 8s 897.71 KiB / 9.32 MiB    9.40% 7s 1.13 MiB / 9.32 MiB   12.10% 7s 1.40 MiB / 9.32 MiB   14.97% 6s 1.65 MiB / 9.32 MiB   17.71% 6s 1.93 MiB / 9.32 MiB   20.70% 6s 2.14 MiB / 9.32 MiB   22.94% 6s 2.36 MiB / 9.32 MiB   25.31% 5s 2.55 MiB / 9.32 MiB   27.38% 5s 2.75 MiB / 9.32 MiB   29.50% 5s 2.98 MiB / 9.32 MiB   31.93% 5s 3.25 MiB / 9.32 MiB   34.81% 5s 3.45 MiB / 9.32 MiB   37.06% 5s 3.64 MiB / 9.32 MiB   38.99% 5s 3.90 MiB / 9.32 MiB   41.81% 4s 4.13 MiB / 9.32 MiB   44.30% 4s 4.38 MiB / 9.32 MiB   47.00% 4s 4.60 MiB / 9.32 MiB   49.37% 4s 4.87 MiB / 9.32 MiB   52.25% 3s 5.14 MiB / 9.32 MiB   55.09% 3s 5.38 MiB / 9.32 MiB   57.75% 3s 5.64 MiB / 9.32 MiB   60.45% 3s 5.84 MiB / 9.32 MiB   62.65% 3s 6.15 MiB / 9.32 MiB   65.94% 2s 6.40 MiB / 9.32 MiB   68.68% 2s 6.61 MiB / 9.32 MiB   70.93% 2s 6.93 MiB / 9.32 MiB   74.35% 2s 7.28 MiB / 9.32 MiB   78.09% 1s 7.66 MiB / 9.32 MiB   82.13% 1s 7.97 MiB / 9.32 MiB   85.49% 1s 8.18 MiB / 9.32 MiB   87.73% 8.50 MiB / 9.32 MiB   91.22% 8.79 MiB / 9.32 MiB   94.32% 9.11 MiB / 9.32 MiB   97.75% 9.32 MiB / 9.32 MiB   99.93% 9.32 MiB / 9.32 MiB  100.00% 7s
```

Run `./godelw version` to verify that gödel was installed correctly. If this is the first run, this invocation will
download all of the plugins and assets:

```
➜ ./godelw version
Getting package from https://palantir.bintray.com/releases/com/palantir/distgo/dist-plugin/1.0.0-rc15/dist-plugin-1.0.0-rc15-linux-amd64.tgz...
 0 B / 4.73 MiB    0.00% 156.68 KiB / 4.73 MiB    3.23% 6s 350.82 KiB / 4.73 MiB    7.24% 5s 640.66 KiB / 4.73 MiB   13.21% 3s 952.38 KiB / 4.73 MiB   19.64% 3s 1.15 MiB / 4.73 MiB   24.30% 3s 1.38 MiB / 4.73 MiB   29.12% 2s 1.64 MiB / 4.73 MiB   34.65% 2s 1.90 MiB / 4.73 MiB   40.06% 2s 2.16 MiB / 4.73 MiB   45.70% 2s 2.47 MiB / 4.73 MiB   52.16% 1s 2.76 MiB / 4.73 MiB   58.39% 1s 3.05 MiB / 4.73 MiB   64.37% 1s 3.30 MiB / 4.73 MiB   69.75% 1s 3.51 MiB / 4.73 MiB   74.15% 3.73 MiB / 4.73 MiB   78.81% 4.02 MiB / 4.73 MiB   84.84% 4.26 MiB / 4.73 MiB   89.92% 4.51 MiB / 4.73 MiB   95.36% 4.73 MiB / 4.73 MiB  100.00% 3s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-format-plugin/format-plugin/1.0.0-rc7/format-plugin-1.0.0-rc7-linux-amd64.tgz...
 0 B / 3.32 MiB    0.00% 208.42 KiB / 3.32 MiB    6.13% 3s 554.29 KiB / 3.32 MiB   16.30% 2s 987.83 KiB / 3.32 MiB   29.06% 1s 1.32 MiB / 3.32 MiB   39.70% 1s 1.83 MiB / 3.32 MiB   55.26% 2.23 MiB / 3.32 MiB   67.07% 2.58 MiB / 3.32 MiB   77.72% 2.93 MiB / 3.32 MiB   88.36% 3.32 MiB / 3.32 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-format-asset-ptimports/ptimports-asset/1.0.0-rc6/ptimports-asset-1.0.0-rc6-linux-amd64.tgz...
 0 B / 3.60 MiB    0.00% 360.69 KiB / 3.60 MiB    9.78% 1s 750.40 KiB / 3.60 MiB   20.34% 1s 1.13 MiB / 3.60 MiB   31.34% 1s 1.54 MiB / 3.60 MiB   42.66% 1s 1.96 MiB / 3.60 MiB   54.30% 2.39 MiB / 3.60 MiB   66.37% 2.64 MiB / 3.60 MiB   73.16% 2.91 MiB / 3.60 MiB   80.71% 3.17 MiB / 3.60 MiB   87.93% 3.42 MiB / 3.60 MiB   95.04% 3.60 MiB / 3.60 MiB  100.00% 2s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-goland-plugin/goland-plugin/1.0.0-rc2/goland-plugin-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.09 MiB    0.00% 333.10 KiB / 3.09 MiB   10.53% 1s 599.62 KiB / 3.09 MiB   18.96% 1s 945.49 KiB / 3.09 MiB   29.90% 1s 1.37 MiB / 3.09 MiB   44.49% 1s 1.71 MiB / 3.09 MiB   55.42% 2.15 MiB / 3.09 MiB   69.51% 2.54 MiB / 3.09 MiB   82.34% 2.85 MiB / 3.09 MiB   92.39% 3.09 MiB / 3.09 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/okgo/check-plugin/1.0.0-rc6/check-plugin-1.0.0-rc6-linux-amd64.tgz...
 0 B / 3.53 MiB    0.00% 319.93 KiB / 3.53 MiB    8.86% 2s 570.46 KiB / 3.53 MiB   15.80% 2s 876.65 KiB / 3.53 MiB   24.28% 1s 1.22 MiB / 3.53 MiB   34.63% 1s 1.59 MiB / 3.53 MiB   45.10% 1s 1.87 MiB / 3.53 MiB   53.14% 1s 2.09 MiB / 3.53 MiB   59.30% 2.38 MiB / 3.53 MiB   67.46% 2.69 MiB / 3.53 MiB   76.26% 3.14 MiB / 3.53 MiB   89.04% 3.34 MiB / 3.53 MiB   94.77% 3.53 MiB / 3.53 MiB  100.00% 2s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-compiles/compiles-asset/1.0.0-rc3/compiles-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.71 MiB    0.00% 264.08 KiB / 3.71 MiB    6.95% 2s 637.78 KiB / 3.71 MiB   16.80% 1s 916.14 KiB / 3.71 MiB   24.13% 1s 1.22 MiB / 3.71 MiB   32.92% 1s 1.62 MiB / 3.71 MiB   43.61% 1s 2.00 MiB / 3.71 MiB   53.87% 1s 2.36 MiB / 3.71 MiB   63.71% 2.77 MiB / 3.71 MiB   74.71% 3.10 MiB / 3.71 MiB   83.51% 3.42 MiB / 3.71 MiB   92.30% 3.71 MiB / 3.71 MiB  100.00% 2s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-deadcode/deadcode-asset/1.0.0-rc2/deadcode-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.73 MiB    0.00% 276.82 KiB / 3.73 MiB    7.25% 2s 515.51 KiB / 3.73 MiB   13.50% 2s 793.87 KiB / 3.73 MiB   20.78% 2s 1.16 MiB / 3.73 MiB   30.98% 1s 1.55 MiB / 3.73 MiB   41.50% 1s 1.93 MiB / 3.73 MiB   51.70% 1s 2.36 MiB / 3.73 MiB   63.36% 2.77 MiB / 3.73 MiB   74.29% 3.09 MiB / 3.73 MiB   82.72% 3.42 MiB / 3.73 MiB   91.78% 3.73 MiB / 3.73 MiB  100.00% 2s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-errcheck/errcheck-asset/1.0.0-rc3/errcheck-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.81 MiB    0.00% 332.65 KiB / 3.81 MiB    8.52% 2s 638.84 KiB / 3.81 MiB   16.36% 2s 1016.71 KiB / 3.81 MiB   26.03% 1s 1.36 MiB / 3.81 MiB   35.60% 1s 1.77 MiB / 3.81 MiB   46.29% 1s 2.16 MiB / 3.81 MiB   56.68% 2.41 MiB / 3.81 MiB   63.09% 2.71 MiB / 3.81 MiB   70.93% 3.09 MiB / 3.81 MiB   80.91% 3.44 MiB / 3.81 MiB   90.18% 3.78 MiB / 3.81 MiB   99.03% 3.81 MiB / 3.81 MiB  100.00% 2s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-extimport/extimport-asset/1.0.0-rc2/extimport-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.36 MiB    0.00% 275.90 KiB / 3.36 MiB    8.01% 2s 582.09 KiB / 3.36 MiB   16.90% 1s 999.63 KiB / 3.36 MiB   29.02% 1s 1.40 MiB / 3.36 MiB   41.60% 1s 1.81 MiB / 3.36 MiB   53.72% 2.32 MiB / 3.36 MiB   69.07% 2.80 MiB / 3.36 MiB   83.15% 3.26 MiB / 3.36 MiB   96.88% 3.36 MiB / 3.36 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-golint/golint-asset/1.0.0-rc3/golint-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.88 MiB    0.00% 470.54 KiB / 3.88 MiB   11.85% 1s 1.03 MiB / 3.88 MiB   26.56% 1s 1.59 MiB / 3.88 MiB   40.98% 2.16 MiB / 3.88 MiB   55.70% 2.89 MiB / 3.88 MiB   74.62% 3.59 MiB / 3.88 MiB   92.44% 3.88 MiB / 3.88 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-govet/govet-asset/1.0.0-rc3/govet-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.17 MiB    0.00% 638.38 KiB / 3.17 MiB   19.70% 1.18 MiB / 3.17 MiB   37.37% 1.94 MiB / 3.17 MiB   61.42% 2.46 MiB / 3.17 MiB   77.73% 3.17 MiB / 3.17 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-importalias/importalias-asset/1.0.0-rc2/importalias-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.38 MiB    0.00% 444.31 KiB / 3.38 MiB   12.85% 1s 1017.03 KiB / 3.38 MiB   29.41% 1.70 MiB / 3.38 MiB   50.33% 2.33 MiB / 3.38 MiB   68.84% 2.92 MiB / 3.38 MiB   86.55% 3.38 MiB / 3.38 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-ineffassign/ineffassign-asset/1.0.0-rc2/ineffassign-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.35 MiB    0.00% 403.54 KiB / 3.35 MiB   11.76% 1s 749.40 KiB / 3.35 MiB   21.84% 1s 1.19 MiB / 3.35 MiB   35.63% 1s 1.60 MiB / 3.35 MiB   47.80% 1.98 MiB / 3.35 MiB   59.15% 2.57 MiB / 3.35 MiB   76.65% 2.89 MiB / 3.35 MiB   86.39% 3.35 MiB / 3.35 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-novendor/novendor-asset/1.0.0-rc4/novendor-asset-1.0.0-rc4-linux-amd64.tgz...
 0 B / 3.42 MiB    0.00% 375.45 KiB / 3.42 MiB   10.74% 1s 860.50 KiB / 3.42 MiB   24.61% 1s 1.32 MiB / 3.42 MiB   38.60% 1.71 MiB / 3.42 MiB   50.08% 2.11 MiB / 3.42 MiB   61.68% 2.69 MiB / 3.42 MiB   78.74% 3.11 MiB / 3.42 MiB   91.13% 3.42 MiB / 3.42 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-outparamcheck/outparamcheck-asset/1.0.0-rc3/outparamcheck-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.84 MiB    0.00% 276.96 KiB / 3.84 MiB    7.05% 2s 527.48 KiB / 3.84 MiB   13.42% 2s 861.51 KiB / 3.84 MiB   21.92% 2s 1.22 MiB / 3.84 MiB   31.84% 1s 1.67 MiB / 3.84 MiB   43.58% 1s 2.09 MiB / 3.84 MiB   54.51% 1s 2.54 MiB / 3.84 MiB   66.25% 2.93 MiB / 3.84 MiB   76.47% 3.32 MiB / 3.84 MiB   86.39% 3.84 MiB / 3.84 MiB  100.00% 2s
 3.84 MiB / 3.84 MiB  100.00% 2sGetting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-unconvert/unconvert-asset/1.0.0-rc3/unconvert-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.91 MiB    0.00% 555.02 KiB / 3.91 MiB   13.87% 1s 1.11 MiB / 3.91 MiB   28.47% 1s 1.59 MiB / 3.91 MiB   40.69% 2.12 MiB / 3.91 MiB   54.20% 2.85 MiB / 3.91 MiB   72.98% 3.56 MiB / 3.91 MiB   91.06% 3.64 MiB / 3.91 MiB   93.15% 3.67 MiB / 3.91 MiB   93.84% 3.72 MiB / 3.91 MiB   95.23% 3.75 MiB / 3.91 MiB   95.93% 3.80 MiB / 3.91 MiB   97.32% 3.91 MiB / 3.91 MiB  100.00% 2s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-varcheck/varcheck-asset/1.0.0-rc2/varcheck-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.75 MiB    0.00% 598.56 KiB / 3.75 MiB   15.60% 1s 1.03 MiB / 3.75 MiB   27.51% 1s 1.52 MiB / 3.75 MiB   40.57% 2.01 MiB / 3.75 MiB   53.63% 2.66 MiB / 3.75 MiB   71.04% 3.18 MiB / 3.75 MiB   84.82% 3.75 MiB / 3.75 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-license-plugin/license-plugin/1.0.0-rc1/license-plugin-1.0.0-rc1-linux-amd64.tgz...
 0 B / 3.30 MiB    0.00% 294.76 KiB / 3.30 MiB    8.71% 2s 517.61 KiB / 3.30 MiB   15.30% 2s 754.14 KiB / 3.30 MiB   22.29% 2s 966.05 KiB / 3.30 MiB   28.55% 2s 1.20 MiB / 3.30 MiB   36.39% 1s 1.40 MiB / 3.30 MiB   42.41% 1s 1.64 MiB / 3.30 MiB   49.56% 1s 1.84 MiB / 3.30 MiB   55.83% 1s 2.08 MiB / 3.30 MiB   63.10% 1s 2.34 MiB / 3.30 MiB   70.74% 2.58 MiB / 3.30 MiB   78.01% 2.80 MiB / 3.30 MiB   84.72% 3.06 MiB / 3.30 MiB   92.64% 3.30 MiB / 3.30 MiB  100.00% 2s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-test-plugin/test-plugin/1.0.0-rc5/test-plugin-1.0.0-rc5-linux-amd64.tgz...
 0 B / 3.60 MiB    0.00% 53.18 KiB / 3.60 MiB    1.44% 13s 97.01 KiB / 3.60 MiB    2.63% 14s 108.85 KiB / 3.60 MiB    2.95% 19s 180.52 KiB / 3.60 MiB    4.90% 15s 387.21 KiB / 3.60 MiB   10.51% 8s 876.42 KiB / 3.60 MiB   23.78% 3s 1.28 MiB / 3.60 MiB   35.43% 2s 1.75 MiB / 3.60 MiB   48.70% 1s 2.35 MiB / 3.60 MiB   65.32% 2.84 MiB / 3.60 MiB   78.91% 3.26 MiB / 3.60 MiB   90.56% 3.60 MiB / 3.60 MiB  100.00% 2s
godel version 2.0.0-rc12
```

Technically, this is sufficient and we have a working gödel installation. However, because the installation was
performed by downloading a distribution, the gödel installation itself does not have a checksum set (no verification was
performed on the package itself):

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0-rc12/godel-2.0.0-rc12.tgz
distributionSHA256=
```

If we are concerned about integrity, we can specify the expected checksum for the package when using godelinit. The
expected checksums for gödel releases can be found on its [Bintray](https://bintray.com/palantir/releases/godel) page.
The checksum for 2.0.0-rc5 is 08d9ed3e33e69006a9c58ec65cef0ad9bd17af4c73b5c1d1aa116e813a954314. Install the distribution
using this checksum:

```
➜ godelinit --checksum ${GODEL_CHECKSUM}
```

If the installation succeeds with a checksum specified, it is set in the properties:

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0-rc12/godel-2.0.0-rc12.tgz
distributionSHA256=6645d041f2243b146f471f79f80fa59adabc5a4a9f0ae9eb1023c3ec9dcac116
```

Commit the changes to the repository:

```
➜ git add godel godelw
➜ git commit -m "Add godel to project"
[master 0a4dd00] Add godel to project
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
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0 9546k    0 91857    0     0  86862      0  0:01:52  0:00:01  0:01:51 86821 15 9546k   15 1441k    0     0   706k      0  0:00:13  0:00:02  0:00:11  706k 30 9546k   30 2917k    0     0   959k      0  0:00:09  0:00:03  0:00:06  959k 49 9546k   49 4760k    0     0  1178k      0  0:00:08  0:00:04  0:00:04 1177k 66 9546k   66 6388k    0     0  1267k      0  0:00:07  0:00:05  0:00:02 1281k 82 9546k   82 7875k    0     0  1303k      0  0:00:07  0:00:06  0:00:01 1562k 97 9546k   97 9314k    0     0  1322k      0  0:00:07  0:00:07 --:--:-- 1574k100 9546k  100 9546k    0     0  1331k      0  0:00:07  0:00:07 --:--:-- 1605k
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
godel version 2.0.0-rc12
```

Technically, this is sufficient and we have a working gödel install. However, distributions that are downloaded manually
do not have a checksum set in `godel/config/godel.properties`:

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0-rc12/godel-2.0.0-rc12.tgz
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
godel version 2.0.0-rc12
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
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0 9546k    0   857    0     0    803      0  3:22:53  0:00:01  3:22:52   804  2 9546k    2  272k    0     0   198k      0  0:00:48  0:00:01  0:00:47  198k 17 9546k   17 1698k    0     0   717k      0  0:00:13  0:00:02  0:00:11  717k 28 9546k   28 2685k    0     0   796k      0  0:00:11  0:00:03  0:00:08  796k 39 9546k   39 3723k    0     0   850k      0  0:00:11  0:00:04  0:00:07  850k 48 9546k   48 4645k    0     0   860k      0  0:00:11  0:00:05  0:00:06 1072k 60 9546k   60 5745k    0     0   902k      0  0:00:10  0:00:06  0:00:04 1096k 75 9546k   75 7166k    0     0   972k      0  0:00:09  0:00:07  0:00:02 1093k 88 9546k   88 8491k    0     0  1015k      0  0:00:09  0:00:08  0:00:01 1163k100 9546k  100 9546k    0     0  1046k      0  0:00:09  0:00:09 --:--:-- 1228k
```

The checksum can be computed using `openssl` or `shasum` as follows:

```
➜ openssl dgst -sha256 download/godel-"${GODEL_VERSION}".tgz
SHA256(download/godel-2.0.0-rc12.tgz)= 6645d041f2243b146f471f79f80fa59adabc5a4a9f0ae9eb1023c3ec9dcac116
➜ shasum -a 256 download/godel-"${GODEL_VERSION}".tgz
6645d041f2243b146f471f79f80fa59adabc5a4a9f0ae9eb1023c3ec9dcac116  download/godel-2.0.0-rc12.tgz
```

Run the following to clean up our state:

```
➜ rm -rf download
```
