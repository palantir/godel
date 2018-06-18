Summary
-------
Projects that use a pre-2.0 version of gödel can be updated using `godelinit`.

Tutorial start state
--------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist-plugin.yml` is configured to build `echgo2`
* Project is tagged as 0.0.1
* `godel/config/dist-plugin.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* `dockerctx` directory exists and `godel/config/dist-plugin.yml` is configured to build Docker images for the product
* Go files have license headers
* `godel/config/godel.yml` is configured to add the go-generate plugin
* `godel/config/generate-plugin.yml` is configured to generate string function
* `godel/config/godel.yml` is configured to ignore all `.+_string.go` files
* `integration_test` contains integration tests
* `godel/config/test-plugin.yml` is configured to specify the "integration" tag
* `docs` contains documentation
* `.circleci/config.yml` exists
* Project is tagged as 1.0.0
* `godelw` is updated to the latest version

Update projects using pre-2.0 version of gödel
----------------------------------------------
Projects that use a pre-2.0 version of gödel can be updated to use the latest version of gödel using the `godelinit`
program. As an exmaple, we will update the version of gödel used by `github.com/nmiyake/echgo` (which was the tutorial
project for the 1.x release of gödel) to use the latest version.

Start by getting the repository and switching to it:

```
➜ go get github.com/nmiyake/echgo
➜ cd ${GOPATH}/src/github.com/nmiyake/echgo
```

Run `./godelw version` to verify that the version of gödel used by the project is pre-1.0:

```
➜ ./godelw version
Downloading https://palantir.bintray.com/releases/com/palantir/godel/godel/0.27.0/godel-0.27.0.tgz to /root/.godel/downloads/godel-0.27.0.tgz...

     0K .......... .......... .......... .......... ..........  0% 3.91M 3s
    50K .......... .......... .......... .......... ..........  0% 3.51M 3s
   100K .......... .......... .......... .......... ..........  1% 26.5M 2s
   150K .......... .......... .......... .......... ..........  1% 6.75M 2s
   200K .......... .......... .......... .......... ..........  2% 28.0M 2s
   250K .......... .......... .......... .......... ..........  2%  150M 1s
   300K .......... .......... .......... .......... ..........  3% 33.7M 1s
   350K .......... .......... .......... .......... ..........  3% 16.5M 1s
   400K .......... .......... .......... .......... ..........  4% 7.12M 1s
   450K .......... .......... .......... .......... ..........  4% 13.1M 1s
   500K .......... .......... .......... .......... ..........  4% 12.6M 1s
   550K .......... .......... .......... .......... ..........  5% 12.7M 1s
   600K .......... .......... .......... .......... ..........  5% 21.1M 1s
   650K .......... .......... .......... .......... ..........  6% 12.4M 1s
   700K .......... .......... .......... .......... ..........  6% 14.4M 1s
   750K .......... .......... .......... .......... ..........  7% 24.5M 1s
   800K .......... .......... .......... .......... ..........  7% 15.4M 1s
   850K .......... .......... .......... .......... ..........  8% 15.7M 1s
   900K .......... .......... .......... .......... ..........  8% 20.1M 1s
   950K .......... .......... .......... .......... ..........  9% 12.7M 1s
  1000K .......... .......... .......... .......... ..........  9% 17.7M 1s
  1050K .......... .......... .......... .......... ..........  9% 19.3M 1s
  1100K .......... .......... .......... .......... .......... 10% 27.9M 1s
  1150K .......... .......... .......... .......... .......... 10% 12.9M 1s
  1200K .......... .......... .......... .......... .......... 11% 22.9M 1s
  1250K .......... .......... .......... .......... .......... 11% 21.0M 1s
  1300K .......... .......... .......... .......... .......... 12% 12.5M 1s
  1350K .......... .......... .......... .......... .......... 12% 16.7M 1s
  1400K .......... .......... .......... .......... .......... 13% 17.9M 1s
  1450K .......... .......... .......... .......... .......... 13% 28.9M 1s
  1500K .......... .......... .......... .......... .......... 14% 13.2M 1s
  1550K .......... .......... .......... .......... .......... 14% 13.5M 1s
  1600K .......... .......... .......... .......... .......... 14% 16.6M 1s
  1650K .......... .......... .......... .......... .......... 15% 20.6M 1s
  1700K .......... .......... .......... .......... .......... 15% 29.7M 1s
  1750K .......... .......... .......... .......... .......... 16% 13.7M 1s
  1800K .......... .......... .......... .......... .......... 16% 25.5M 1s
  1850K .......... .......... .......... .......... .......... 17% 5.70M 1s
  1900K .......... .......... .......... .......... .......... 17% 9.79M 1s
  1950K .......... .......... .......... .......... .......... 18% 7.90M 1s
  2000K .......... .......... .......... .......... .......... 18% 2.61M 1s
  2050K .......... .......... .......... .......... .......... 19% 5.13M 1s
  2100K .......... .......... .......... .......... .......... 19% 4.14M 1s
  2150K .......... .......... .......... .......... .......... 19% 10.4M 1s
  2200K .......... .......... .......... .......... .......... 20% 8.89M 1s
  2250K .......... .......... .......... .......... .......... 20% 11.4M 1s
  2300K .......... .......... .......... .......... .......... 21% 21.1M 1s
  2350K .......... .......... .......... .......... .......... 21% 13.6M 1s
  2400K .......... .......... .......... .......... .......... 22% 11.2M 1s
  2450K .......... .......... .......... .......... .......... 22% 18.2M 1s
  2500K .......... .......... .......... .......... .......... 23% 15.3M 1s
  2550K .......... .......... .......... .......... .......... 23% 14.9M 1s
  2600K .......... .......... .......... .......... .......... 24% 10.3M 1s
  2650K .......... .......... .......... .......... .......... 24% 17.3M 1s
  2700K .......... .......... .......... .......... .......... 24% 22.5M 1s
  2750K .......... .......... .......... .......... .......... 25% 15.5M 1s
  2800K .......... .......... .......... .......... .......... 25% 4.18M 1s
  2850K .......... .......... .......... .......... .......... 26% 30.5M 1s
  2900K .......... .......... .......... .......... .......... 26% 51.3M 1s
  2950K .......... .......... .......... .......... .......... 27% 28.1M 1s
  3000K .......... .......... .......... .......... .......... 27% 8.29M 1s
  3050K .......... .......... .......... .......... .......... 28% 9.26M 1s
  3100K .......... .......... .......... .......... .......... 28% 9.98M 1s
  3150K .......... .......... .......... .......... .......... 29% 12.8M 1s
  3200K .......... .......... .......... .......... .......... 29% 13.2M 1s
  3250K .......... .......... .......... .......... .......... 29% 10.0M 1s
  3300K .......... .......... .......... .......... .......... 30% 11.6M 1s
  3350K .......... .......... .......... .......... .......... 30% 9.34M 1s
  3400K .......... .......... .......... .......... .......... 31% 10.6M 1s
  3450K .......... .......... .......... .......... .......... 31% 13.1M 1s
  3500K .......... .......... .......... .......... .......... 32% 18.0M 1s
  3550K .......... .......... .......... .......... .......... 32% 23.0M 1s
  3600K .......... .......... .......... .......... .......... 33% 13.0M 1s
  3650K .......... .......... .......... .......... .......... 33% 13.3M 1s
  3700K .......... .......... .......... .......... .......... 34% 18.1M 1s
  3750K .......... .......... .......... .......... .......... 34% 17.9M 1s
  3800K .......... .......... .......... .......... .......... 34% 15.6M 1s
  3850K .......... .......... .......... .......... .......... 35% 16.4M 1s
  3900K .......... .......... .......... .......... .......... 35% 20.1M 1s
  3950K .......... .......... .......... .......... .......... 36% 15.5M 1s
  4000K .......... .......... .......... .......... .......... 36% 15.1M 1s
  4050K .......... .......... .......... .......... .......... 37% 14.9M 1s
  4100K .......... .......... .......... .......... .......... 37% 19.0M 1s
  4150K .......... .......... .......... .......... .......... 38% 21.7M 1s
  4200K .......... .......... .......... .......... .......... 38% 17.8M 1s
  4250K .......... .......... .......... .......... .......... 39% 13.2M 1s
  4300K .......... .......... .......... .......... .......... 39% 30.9M 1s
  4350K .......... .......... .......... .......... .......... 39% 13.2M 1s
  4400K .......... .......... .......... .......... .......... 40% 20.4M 1s
  4450K .......... .......... .......... .......... .......... 40% 6.91M 1s
  4500K .......... .......... .......... .......... .......... 41% 39.7M 1s
  4550K .......... .......... .......... .......... .......... 41% 6.79M 1s
  4600K .......... .......... .......... .......... .......... 42% 41.2M 1s
  4650K .......... .......... .......... .......... .......... 42% 14.7M 1s
  4700K .......... .......... .......... .......... .......... 43% 18.2M 0s
  4750K .......... .......... .......... .......... .......... 43% 19.5M 0s
  4800K .......... .......... .......... .......... .......... 44% 11.6M 0s
  4850K .......... .......... .......... .......... .......... 44% 19.8M 0s
  4900K .......... .......... .......... .......... .......... 44% 37.2M 0s
  4950K .......... .......... .......... .......... .......... 45% 11.5M 0s
  5000K .......... .......... .......... .......... .......... 45% 23.6M 0s
  5050K .......... .......... .......... .......... .......... 46% 19.2M 0s
  5100K .......... .......... .......... .......... .......... 46% 16.6M 0s
  5150K .......... .......... .......... .......... .......... 47% 15.4M 0s
  5200K .......... .......... .......... .......... .......... 47% 18.9M 0s
  5250K .......... .......... .......... .......... .......... 48% 29.2M 0s
  5300K .......... .......... .......... .......... .......... 48% 13.7M 0s
  5350K .......... .......... .......... .......... .......... 49% 16.4M 0s
  5400K .......... .......... .......... .......... .......... 49% 19.9M 0s
  5450K .......... .......... .......... .......... .......... 49% 6.74M 0s
  5500K .......... .......... .......... .......... .......... 50% 11.1M 0s
  5550K .......... .......... .......... .......... .......... 50% 19.2M 0s
  5600K .......... .......... .......... .......... .......... 51% 12.4M 0s
  5650K .......... .......... .......... .......... .......... 51% 34.1M 0s
  5700K .......... .......... .......... .......... .......... 52% 17.7M 0s
  5750K .......... .......... .......... .......... .......... 52% 20.0M 0s
  5800K .......... .......... .......... .......... .......... 53% 17.4M 0s
  5850K .......... .......... .......... .......... .......... 53% 16.9M 0s
  5900K .......... .......... .......... .......... .......... 54% 27.1M 0s
  5950K .......... .......... .......... .......... .......... 54% 2.57M 0s
  6000K .......... .......... .......... .......... .......... 54% 24.1M 0s
  6050K .......... .......... .......... .......... .......... 55% 17.1M 0s
  6100K .......... .......... .......... .......... .......... 55% 19.9M 0s
  6150K .......... .......... .......... .......... .......... 56% 26.0M 0s
  6200K .......... .......... .......... .......... .......... 56% 21.7M 0s
  6250K .......... .......... .......... .......... .......... 57% 22.0M 0s
  6300K .......... .......... .......... .......... .......... 57% 27.7M 0s
  6350K .......... .......... .......... .......... .......... 58% 16.3M 0s
  6400K .......... .......... .......... .......... .......... 58% 31.6M 0s
  6450K .......... .......... .......... .......... .......... 59% 24.2M 0s
  6500K .......... .......... .......... .......... .......... 59% 19.4M 0s
  6550K .......... .......... .......... .......... .......... 59% 11.9M 0s
  6600K .......... .......... .......... .......... .......... 60% 23.7M 0s
  6650K .......... .......... .......... .......... .......... 60% 26.2M 0s
  6700K .......... .......... .......... .......... .......... 61% 15.4M 0s
  6750K .......... .......... .......... .......... .......... 61% 26.6M 0s
  6800K .......... .......... .......... .......... .......... 62% 4.73M 0s
  6850K .......... .......... .......... .......... .......... 62% 40.2M 0s
  6900K .......... .......... .......... .......... .......... 63% 42.9M 0s
  6950K .......... .......... .......... .......... .......... 63% 41.9M 0s
  7000K .......... .......... .......... .......... .......... 64% 38.6M 0s
  7050K .......... .......... .......... .......... .......... 64% 13.9M 0s
  7100K .......... .......... .......... .......... .......... 64% 13.5M 0s
  7150K .......... .......... .......... .......... .......... 65% 9.00M 0s
  7200K .......... .......... .......... .......... .......... 65% 13.7M 0s
  7250K .......... .......... .......... .......... .......... 66% 7.92M 0s
  7300K .......... .......... .......... .......... .......... 66% 14.9M 0s
  7350K .......... .......... .......... .......... .......... 67% 13.4M 0s
  7400K .......... .......... .......... .......... .......... 67% 19.9M 0s
  7450K .......... .......... .......... .......... .......... 68% 6.74M 0s
  7500K .......... .......... .......... .......... .......... 68% 27.8M 0s
  7550K .......... .......... .......... .......... .......... 69% 9.80M 0s
  7600K .......... .......... .......... .......... .......... 69% 8.46M 0s
  7650K .......... .......... .......... .......... .......... 69% 16.5M 0s
  7700K .......... .......... .......... .......... .......... 70% 13.0M 0s
  7750K .......... .......... .......... .......... .......... 70% 11.8M 0s
  7800K .......... .......... .......... .......... .......... 71% 16.0M 0s
  7850K .......... .......... .......... .......... .......... 71% 29.7M 0s
  7900K .......... .......... .......... .......... .......... 72% 17.1M 0s
  7950K .......... .......... .......... .......... .......... 72% 12.9M 0s
  8000K .......... .......... .......... .......... .......... 73% 18.7M 0s
  8050K .......... .......... .......... .......... .......... 73% 17.1M 0s
  8100K .......... .......... .......... .......... .......... 74% 17.5M 0s
  8150K .......... .......... .......... .......... .......... 74% 11.4M 0s
  8200K .......... .......... .......... .......... .......... 74% 38.5M 0s
  8250K .......... .......... .......... .......... .......... 75% 13.8M 0s
  8300K .......... .......... .......... .......... .......... 75% 22.3M 0s
  8350K .......... .......... .......... .......... .......... 76% 21.3M 0s
  8400K .......... .......... .......... .......... .......... 76% 14.4M 0s
  8450K .......... .......... .......... .......... .......... 77% 26.1M 0s
  8500K .......... .......... .......... .......... .......... 77% 10.6M 0s
  8550K .......... .......... .......... .......... .......... 78% 21.5M 0s
  8600K .......... .......... .......... .......... .......... 78% 22.4M 0s
  8650K .......... .......... .......... .......... .......... 79% 17.5M 0s
  8700K .......... .......... .......... .......... .......... 79% 13.5M 0s
  8750K .......... .......... .......... .......... .......... 79% 11.2M 0s
  8800K .......... .......... .......... .......... .......... 80% 19.5M 0s
  8850K .......... .......... .......... .......... .......... 80%  122M 0s
  8900K .......... .......... .......... .......... .......... 81% 19.5M 0s
  8950K .......... .......... .......... .......... .......... 81% 19.7M 0s
  9000K .......... .......... .......... .......... .......... 82% 17.1M 0s
  9050K .......... .......... .......... .......... .......... 82% 18.0M 0s
  9100K .......... .......... .......... .......... .......... 83% 32.8M 0s
  9150K .......... .......... .......... .......... .......... 83% 15.4M 0s
  9200K .......... .......... .......... .......... .......... 84% 15.6M 0s
  9250K .......... .......... .......... .......... .......... 84% 14.1M 0s
  9300K .......... .......... .......... .......... .......... 84% 16.1M 0s
  9350K .......... .......... .......... .......... .......... 85% 8.51M 0s
  9400K .......... .......... .......... .......... .......... 85% 10.6M 0s
  9450K .......... .......... .......... .......... .......... 86% 24.2M 0s
  9500K .......... .......... .......... .......... .......... 86% 28.4M 0s
  9550K .......... .......... .......... .......... .......... 87% 21.7M 0s
  9600K .......... .......... .......... .......... .......... 87% 12.2M 0s
  9650K .......... .......... .......... .......... .......... 88% 16.5M 0s
  9700K .......... .......... .......... .......... .......... 88% 38.8M 0s
  9750K .......... .......... .......... .......... .......... 89% 11.7M 0s
  9800K .......... .......... .......... .......... .......... 89% 2.29M 0s
  9850K .......... .......... .......... .......... .......... 89% 9.38M 0s
  9900K .......... .......... .......... .......... .......... 90% 8.70M 0s
  9950K .......... .......... .......... .......... .......... 90% 14.0M 0s
 10000K .......... .......... .......... .......... .......... 91% 14.3M 0s
 10050K .......... .......... .......... .......... .......... 91% 23.5M 0s
 10100K .......... .......... .......... .......... .......... 92% 61.5M 0s
 10150K .......... .......... .......... .......... .......... 92% 11.2M 0s
 10200K .......... .......... .......... .......... .......... 93% 33.8M 0s
 10250K .......... .......... .......... .......... .......... 93% 17.5M 0s
 10300K .......... .......... .......... .......... .......... 94% 31.8M 0s
 10350K .......... .......... .......... .......... .......... 94% 19.1M 0s
 10400K .......... .......... .......... .......... .......... 94% 25.3M 0s
 10450K .......... .......... .......... .......... .......... 95% 15.0M 0s
 10500K .......... .......... .......... .......... .......... 95% 20.5M 0s
 10550K .......... .......... .......... .......... .......... 96% 21.0M 0s
 10600K .......... .......... .......... .......... .......... 96% 13.1M 0s
 10650K .......... .......... .......... .......... .......... 97% 16.1M 0s
 10700K .......... .......... .......... .......... .......... 97% 25.0M 0s
 10750K .......... .......... .......... .......... .......... 98% 18.1M 0s
 10800K .......... .......... .......... .......... .......... 98% 22.5M 0s
 10850K .......... .......... .......... .......... .......... 99% 7.48M 0s
 10900K .......... .......... .......... .......... .......... 99% 41.1M 0s
 10950K .......... .......... .......... .......... .......... 99% 9.12M 0s
 11000K ..                                                    100% 3.47M=0.8sgodel version 0.27.0
```

The godelinit program can be used to install or update gödel. Although it is technically possible to use the
`./godelw update` mechanism to get the new binary, pre-2.0 versions of gödel do not know about configuration upgrades,
so running `./godelw update` will only update the gödel version. Using godelinit has the advantage that it knows to
invoke the commands to upgrade the configuration after gödel itself has been updated.

First, ensure that the latest version of godelinit is installed:

```
➜ go get -u github.com/palantir/godel/godelinit
```

Then, run `godelinit`. Running the program with no arguments determines the latest released version of gödel and either
installs it or upgrades the version of gödel in the current project to that version. In this case, because the newest
version of gödel is >=2 and the version of gödel in the current project is <2, the legacy configuration will be
upgraded:

```
➜ godelinit
Upgraded configuration for godel.yml
Upgraded configuration for dist-plugin.yml
Upgraded configuration for license-plugin.yml
Upgraded configuration for format-plugin.yml
Upgraded configuration for test-plugin.yml
Upgraded configuration for check-plugin.yml
WARNING: The following configuration file(s) were non-empty and had no known upgraders for legacy configuration: [generate.yml]
         If these configuration file(s) are for plugins, add the plugins to the configuration and rerun the upgrade task.
```

As indicated by the output, most of the configuration is automatically upgraded. Verify that this is the case:

```
➜ git status
On branch master
Your branch is up-to-date with 'origin/master'.
Changes not staged for commit:
  (use "git add/rm <file>..." to update what will be committed)
  (use "git checkout -- <file>..." to discard changes in working directory)

	deleted:    godel/config/check.yml
	deleted:    godel/config/dist.yml
	deleted:    godel/config/exclude.yml
	deleted:    godel/config/format.yml
	modified:   godel/config/godel.properties
	deleted:    godel/config/imports.yml
	deleted:    godel/config/license.yml
	deleted:    godel/config/test.yml
	modified:   godelw

Untracked files:
  (use "git add <file>..." to include in what will be committed)

	godel/config/check-plugin.yml
	godel/config/dist-plugin.yml
	godel/config/format-plugin.yml
	godel/config/godel.yml
	godel/config/license-plugin.yml
	godel/config/test-plugin.yml

no changes added to commit (use "git add" and/or "git commit -a")
```

However, note that the output said that `generate.yml` was non-empty but there were no known upgraders for it. This is
because, although `generate` was a builtin task prior to gödel 2.0, starting with gödel 2.0 this task is provided by a
plugin. We can upgrade the configuration by adding the plugin and explicitly running the `upgrade-config` task with the
`--legacy` flag set. Start by updating the `godel.yml` configuration to add the plugin:

```
➜ echo 'plugins:
  resolvers:
    - "https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"
  plugins:
    - locator:
        id: "com.palantir.godel-generate-plugin:generate-plugin:1.0.0"
exclude:
  names:
    - "\\\\..+"
    - "vendor"
    - ".+_string.go"
  paths:
    - "godel"' > godel/config/godel.yml
```

Now that the plugin has been added to the configuration, run the `upgrade-config` task in legacy mode:

```
➜ ./godelw upgrade-config --legacy
Upgraded configuration for generate-plugin.yml
```

The message indicates that the configuration was upgraded. You can verify this by observing that a `generate-plugin.yml`
file is now present in the configuration directory:

```
➜ cat godel/config/generate-plugin.yml
generators:
  stringer:
    go-generate-dir: generator
    gen-paths:
      paths:
      - echo/type_string.go
```

Run the verify task to verify that the upgraded configuration is working correctly:

```
➜ ./godelw verify --apply=false
Running format...
Running generate...
Running license...
Running check...
[errcheck]      Running errcheck...
[compiles]      Running compiles...
[deadcode]      Running deadcode...
[extimport]     Running extimport...
[extimport]     Finished extimport
[golint]        Running golint...
[golint]        Finished golint
[govet]         Running govet...
[govet]         Finished govet
[importalias]   Running importalias...
[importalias]   Finished importalias
[ineffassign]   Running ineffassign...
[ineffassign]   Finished ineffassign
[novendor]      Running novendor...
[novendor]      Finished novendor
[outparamcheck] Running outparamcheck...
[deadcode]      Finished deadcode
[unconvert]     Running unconvert...
[compiles]      Finished compiles
[varcheck]      Running varcheck...
[errcheck]      Finished errcheck
[outparamcheck] Finished outparamcheck
[unconvert]     Finished unconvert
[varcheck]      Finished varcheck
Running test...
?   	github.com/nmiyake/echgo                 	[no test files]
ok  	github.com/nmiyake/echgo/echo            	0.002s
?   	github.com/nmiyake/echgo/generator       	[no test files]
--- FAIL: TestInvalidType (0.73s)
panic: command [/go/src/github.com/nmiyake/echgo/godelw artifacts build --absolute --os-arch=linux-amd64 --requires-build echgo] failed with output:
Error: unknown flag: --os-arch
Usage:
  distgo artifacts build [flags] [product-build-ids]

Flags:
      --absolute         print the absolute path for artifacts
  -h, --help             help for build
      --requires-build   only prints the artifacts that require building (omits artifacts that are already built and are up-to-date)

Global Flags:
      --assets stringSlice    path(s) to the plugin asset(s)
      --config string         path to the plugin configuration file
      --debug                 run in debug mode
      --godel-config string   path to the godel.yml configuration file
      --project-dir string    path to project directory
Error: exit status 1 [recovered]
	panic: command [/go/src/github.com/nmiyake/echgo/godelw artifacts build --absolute --os-arch=linux-amd64 --requires-build echgo] failed with output:
Error: unknown flag: --os-arch
Usage:
  distgo artifacts build [flags] [product-build-ids]

Flags:
      --absolute         print the absolute path for artifacts
  -h, --help             help for build
      --requires-build   only prints the artifacts that require building (omits artifacts that are already built and are up-to-date)

Global Flags:
      --assets stringSlice    path(s) to the plugin asset(s)
      --config string         path to the plugin configuration file
      --debug                 run in debug mode
      --godel-config string   path to the godel.yml configuration file
      --project-dir string    path to project directory
Error: exit status 1

goroutine 5 [running]:
testing.tRunner.func1(0xc42000e2d0)
	/usr/local/go/src/testing/testing.go:742 +0x29d
panic(0x516640, 0xc4200de080)
	/usr/local/go/src/runtime/panic.go:505 +0x229
github.com/nmiyake/echgo/integration_test_test.TestInvalidType(0xc42000e2d0)
	/go/src/github.com/nmiyake/echgo/integration_test/integration_test.go:27 +0x49e
testing.tRunner(0xc42000e2d0, 0x548608)
	/usr/local/go/src/testing/testing.go:777 +0xd0
created by testing.(*T).Run
	/usr/local/go/src/testing/testing.go:824 +0x2e0
FAIL	github.com/nmiyake/echgo/integration_test	0.732s
Error: 1 package(s) had failing tests:
	github.com/nmiyake/echgo/integration_test
Failed tasks:
	test
```

Although all of the checks pass, one of the tests now fails. This is because the integration test uses the
`github.com/palantir/godel/pkg/products`. This package had an API-breaking change: projects that use version 2 or later
of gödel that use this package must use the v2 version of the package instead. In order to fix this, vendor the v2
package and update the test to use it:

```
➜ cp -r ${GOPATH}/src/github.com/palantir/godel/pkg/products/v2/ vendor/github.com/palantir/godel/pkg/products/
➜ sed -i 's:github.com/palantir/godel/pkg/products:github.com/palantir/godel/pkg/products/v2/products:g' integration_test/integration_test.go
```

Confirm that the "verify" task now succeeds:

```
➜ ./godelw verify --apply=false
Running format...
Running generate...
Running license...
Running check...
[extimport]     Running extimport...
[compiles]      Running compiles...
[deadcode]      Running deadcode...
[errcheck]      Running errcheck...
[extimport]     Finished extimport
[golint]        Running golint...
[golint]        Finished golint
[govet]         Running govet...
[govet]         Finished govet
[importalias]   Running importalias...
[importalias]   Finished importalias
[ineffassign]   Running ineffassign...
[ineffassign]   Finished ineffassign
[novendor]      Running novendor...
[novendor]      Finished novendor
[outparamcheck] Running outparamcheck...
[deadcode]      Finished deadcode
[unconvert]     Running unconvert...
[compiles]      Finished compiles
[varcheck]      Running varcheck...
[errcheck]      Finished errcheck
[outparamcheck] Finished outparamcheck
[unconvert]     Finished unconvert
[varcheck]      Finished varcheck
Running test...
?   	github.com/nmiyake/echgo                 	[no test files]
ok  	github.com/nmiyake/echgo/echo            	(cached)
?   	github.com/nmiyake/echgo/generator       	[no test files]
ok  	github.com/nmiyake/echgo/integration_test	2.287s
```

Finally, clean up the state:

```
➜ cd ${GOPATH}/src/${PROJECT_PATH}
➜ rm -rf ${GOPATH}/src/github.com/nmiyake/echgo
```

Tutorial end state
------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist-plugin.yml` is configured to build `echgo2`
* Project is tagged as 0.0.1
* `godel/config/dist-plugin.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* `dockerctx` directory exists and `godel/config/dist-plugin.yml` is configured to build Docker images for the product
* Go files have license headers
* `godel/config/godel.yml` is configured to add the go-generate plugin
* `godel/config/generate-plugin.yml` is configured to generate string function
* `godel/config/godel.yml` is configured to ignore all `.+_string.go` files
* `integration_test` contains integration tests
* `godel/config/test-plugin.yml` is configured to specify the "integration" tag
* `docs` contains documentation
* `.circleci/config.yml` exists
* Project is tagged as 1.0.0
* `godelw` is updated to the latest version

Tutorial next step
------------------
[Other commands](https://github.com/palantir/godel/wiki/Other-commands)

More
----
### Update to a specific version
The `--version` flag can be used to upgrade gödel to a specific version.

### Provide a checksum
The `--checksum` flag can be used to specify the expected checksum for the update.

### Sync installation to contents of `godel/config/godel.properties`
The `--sync` flag can be used to specify that the update operation should update gödel to match the values specified
in the `godel/config/godel.properties` file. The file should be edited to match the desired state first.

Updating gödel in this manner requires knowing the distribution URL for the target version. Although it is optional, it
is recommended to have the SHA-256 checksum of the distribution as well to ensure the integrity of the update.
