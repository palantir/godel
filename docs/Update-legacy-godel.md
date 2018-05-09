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

     0K .......... .......... .......... .......... ..........  0%  633K 17s
    50K .......... .......... .......... .......... ..........  0% 1.32M 13s
   100K .......... .......... .......... .......... ..........  1% 10.5M 9s
   150K .......... .......... .......... .......... ..........  1%  981K 9s
   200K .......... .......... .......... .......... ..........  2% 6.63M 8s
   250K .......... .......... .......... .......... ..........  2% 1.92M 7s
   300K .......... .......... .......... .......... ..........  3% 6.16M 6s
   350K .......... .......... .......... .......... ..........  3% 6.59M 6s
   400K .......... .......... .......... .......... ..........  4% 4.51M 5s
   450K .......... .......... .......... .......... ..........  4% 5.59M 5s
   500K .......... .......... .......... .......... ..........  4% 5.38M 5s
   550K .......... .......... .......... .......... ..........  5% 5.38M 4s
   600K .......... .......... .......... .......... ..........  5% 2.69M 4s
   650K .......... .......... .......... .......... ..........  6% 5.82M 4s
   700K .......... .......... .......... .......... ..........  6% 2.64M 4s
   750K .......... .......... .......... .......... ..........  7% 6.02M 4s
   800K .......... .......... .......... .......... ..........  7% 3.98M 4s
   850K .......... .......... .......... .......... ..........  8% 1.11M 4s
   900K .......... .......... .......... .......... ..........  8% 2.38M 4s
   950K .......... .......... .......... .......... ..........  9% 3.59M 4s
  1000K .......... .......... .......... .......... ..........  9% 1.39M 4s
  1050K .......... .......... .......... .......... ..........  9% 33.6M 4s
  1100K .......... .......... .......... .......... .......... 10% 2.72M 4s
  1150K .......... .......... .......... .......... .......... 10% 1.42M 4s
  1200K .......... .......... .......... .......... .......... 11% 11.7M 4s
  1250K .......... .......... .......... .......... .......... 11% 1.07M 4s
  1300K .......... .......... .......... .......... .......... 12% 4.89M 4s
  1350K .......... .......... .......... .......... .......... 12% 1.56M 4s
  1400K .......... .......... .......... .......... .......... 13% 9.41M 4s
  1450K .......... .......... .......... .......... .......... 13%  588K 4s
  1500K .......... .......... .......... .......... .......... 14% 3.80M 4s
  1550K .......... .......... .......... .......... .......... 14% 1.05M 4s
  1600K .......... .......... .......... .......... .......... 14% 3.20M 4s
  1650K .......... .......... .......... .......... .......... 15% 1.87M 4s
  1700K .......... .......... .......... .......... .......... 15% 2.40M 4s
  1750K .......... .......... .......... .......... .......... 16% 2.06M 4s
  1800K .......... .......... .......... .......... .......... 16% 1.46M 4s
  1850K .......... .......... .......... .......... .......... 17% 1.84M 4s
  1900K .......... .......... .......... .......... .......... 17% 2.06M 4s
  1950K .......... .......... .......... .......... .......... 18% 2.87M 4s
  2000K .......... .......... .......... .......... .......... 18% 1.27M 4s
  2050K .......... .......... .......... .......... .......... 19% 1.49M 4s
  2100K .......... .......... .......... .......... .......... 19% 1.68M 4s
  2150K .......... .......... .......... .......... .......... 19% 2.15M 4s
  2200K .......... .......... .......... .......... .......... 20% 2.10M 4s
  2250K .......... .......... .......... .......... .......... 20% 2.18M 4s
  2300K .......... .......... .......... .......... .......... 21% 2.44M 4s
  2350K .......... .......... .......... .......... .......... 21% 2.07M 4s
  2400K .......... .......... .......... .......... .......... 22% 2.90M 4s
  2450K .......... .......... .......... .......... .......... 22% 2.25M 4s
  2500K .......... .......... .......... .......... .......... 23% 2.94M 4s
  2550K .......... .......... .......... .......... .......... 23% 3.50M 4s
  2600K .......... .......... .......... .......... .......... 24% 1021K 4s
  2650K .......... .......... .......... .......... .......... 24% 4.72M 4s
  2700K .......... .......... .......... .......... .......... 24% 1.57M 4s
  2750K .......... .......... .......... .......... .......... 25% 2.37M 4s
  2800K .......... .......... .......... .......... .......... 25% 3.33M 4s
  2850K .......... .......... .......... .......... .......... 26% 2.40M 4s
  2900K .......... .......... .......... .......... .......... 26% 1.05M 4s
  2950K .......... .......... .......... .......... .......... 27% 3.27M 4s
  3000K .......... .......... .......... .......... .......... 27% 1.23M 4s
  3050K .......... .......... .......... .......... .......... 28% 2.06M 4s
  3100K .......... .......... .......... .......... .......... 28% 3.59M 4s
  3150K .......... .......... .......... .......... .......... 29% 1.12M 4s
  3200K .......... .......... .......... .......... .......... 29% 1.29M 4s
  3250K .......... .......... .......... .......... .......... 29% 1.95M 4s
  3300K .......... .......... .......... .......... .......... 30% 2.00M 4s
  3350K .......... .......... .......... .......... .......... 30% 4.64M 4s
  3400K .......... .......... .......... .......... .......... 31% 1.94M 4s
  3450K .......... .......... .......... .......... .......... 31% 5.74M 4s
  3500K .......... .......... .......... .......... .......... 32% 1.00M 4s
  3550K .......... .......... .......... .......... .......... 32% 14.2M 4s
  3600K .......... .......... .......... .......... .......... 33% 1.19M 4s
  3650K .......... .......... .......... .......... .......... 33% 2.25M 3s
  3700K .......... .......... .......... .......... .......... 34% 9.38M 3s
  3750K .......... .......... .......... .......... .......... 34% 1.71M 3s
  3800K .......... .......... .......... .......... .......... 34% 1.38M 3s
  3850K .......... .......... .......... .......... .......... 35% 2.19M 3s
  3900K .......... .......... .......... .......... .......... 35% 3.93M 3s
  3950K .......... .......... .......... .......... .......... 36% 1.74M 3s
  4000K .......... .......... .......... .......... .......... 36% 2.21M 3s
  4050K .......... .......... .......... .......... .......... 37% 6.86M 3s
  4100K .......... .......... .......... .......... .......... 37% 2.04M 3s
  4150K .......... .......... .......... .......... .......... 38% 4.03M 3s
  4200K .......... .......... .......... .......... .......... 38% 1.49M 3s
  4250K .......... .......... .......... .......... .......... 39% 9.93M 3s
  4300K .......... .......... .......... .......... .......... 39% 2.02M 3s
  4350K .......... .......... .......... .......... .......... 39% 2.14M 3s
  4400K .......... .......... .......... .......... .......... 40% 9.25M 3s
  4450K .......... .......... .......... .......... .......... 40% 2.02M 3s
  4500K .......... .......... .......... .......... .......... 41% 3.17M 3s
  4550K .......... .......... .......... .......... .......... 41% 1.12M 3s
  4600K .......... .......... .......... .......... .......... 42% 4.35M 3s
  4650K .......... .......... .......... .......... .......... 42% 2.94M 3s
  4700K .......... .......... .......... .......... .......... 43% 5.00M 3s
  4750K .......... .......... .......... .......... .......... 43% 2.18M 3s
  4800K .......... .......... .......... .......... .......... 44% 6.28M 3s
  4850K .......... .......... .......... .......... .......... 44% 4.12M 3s
  4900K .......... .......... .......... .......... .......... 44% 2.55M 3s
  4950K .......... .......... .......... .......... .......... 45% 9.26M 3s
  5000K .......... .......... .......... .......... .......... 45% 1.28M 3s
  5050K .......... .......... .......... .......... .......... 46% 2.91M 3s
  5100K .......... .......... .......... .......... .......... 46% 2.13M 3s
  5150K .......... .......... .......... .......... .......... 47% 2.92M 3s
  5200K .......... .......... .......... .......... .......... 47% 10.1M 3s
  5250K .......... .......... .......... .......... .......... 48% 1.85M 3s
  5300K .......... .......... .......... .......... .......... 48% 1.29M 3s
  5350K .......... .......... .......... .......... .......... 49% 5.29M 3s
  5400K .......... .......... .......... .......... .......... 49% 1.35M 3s
  5450K .......... .......... .......... .......... .......... 49% 3.37M 2s
  5500K .......... .......... .......... .......... .......... 50% 3.23M 2s
  5550K .......... .......... .......... .......... .......... 50% 5.46M 2s
  5600K .......... .......... .......... .......... .......... 51% 2.00M 2s
  5650K .......... .......... .......... .......... .......... 51% 6.25M 2s
  5700K .......... .......... .......... .......... .......... 52% 3.51M 2s
  5750K .......... .......... .......... .......... .......... 52% 4.32M 2s
  5800K .......... .......... .......... .......... .......... 53% 2.65M 2s
  5850K .......... .......... .......... .......... .......... 53% 3.62M 2s
  5900K .......... .......... .......... .......... .......... 54% 6.00M 2s
  5950K .......... .......... .......... .......... .......... 54% 3.29M 2s
  6000K .......... .......... .......... .......... .......... 54% 3.34M 2s
  6050K .......... .......... .......... .......... .......... 55% 3.15M 2s
  6100K .......... .......... .......... .......... .......... 55% 4.66M 2s
  6150K .......... .......... .......... .......... .......... 56% 6.48M 2s
  6200K .......... .......... .......... .......... .......... 56% 2.78M 2s
  6250K .......... .......... .......... .......... .......... 57% 4.93M 2s
  6300K .......... .......... .......... .......... .......... 57% 3.25M 2s
  6350K .......... .......... .......... .......... .......... 58% 9.13M 2s
  6400K .......... .......... .......... .......... .......... 58% 1.35M 2s
  6450K .......... .......... .......... .......... .......... 59% 2.82M 2s
  6500K .......... .......... .......... .......... .......... 59% 4.19M 2s
  6550K .......... .......... .......... .......... .......... 59% 9.07M 2s
  6600K .......... .......... .......... .......... .......... 60% 2.54M 2s
  6650K .......... .......... .......... .......... .......... 60% 6.38M 2s
  6700K .......... .......... .......... .......... .......... 61% 2.18M 2s
  6750K .......... .......... .......... .......... .......... 61% 10.2M 2s
  6800K .......... .......... .......... .......... .......... 62% 2.60M 2s
  6850K .......... .......... .......... .......... .......... 62% 3.86M 2s
  6900K .......... .......... .......... .......... .......... 63% 6.06M 2s
  6950K .......... .......... .......... .......... .......... 63% 3.35M 2s
  7000K .......... .......... .......... .......... .......... 64% 6.54M 2s
  7050K .......... .......... .......... .......... .......... 64% 6.01M 2s
  7100K .......... .......... .......... .......... .......... 64% 3.07M 2s
  7150K .......... .......... .......... .......... .......... 65% 4.67M 2s
  7200K .......... .......... .......... .......... .......... 65% 5.40M 2s
  7250K .......... .......... .......... .......... .......... 66% 2.62M 1s
  7300K .......... .......... .......... .......... .......... 66% 8.38M 1s
  7350K .......... .......... .......... .......... .......... 67% 3.47M 1s
  7400K .......... .......... .......... .......... .......... 67% 9.79M 1s
  7450K .......... .......... .......... .......... .......... 68% 4.62M 1s
  7500K .......... .......... .......... .......... .......... 68% 1.31M 1s
  7550K .......... .......... .......... .......... .......... 69% 1.80M 1s
  7600K .......... .......... .......... .......... .......... 69% 3.14M 1s
  7650K .......... .......... .......... .......... .......... 69% 10.2M 1s
  7700K .......... .......... .......... .......... .......... 70% 3.93M 1s
  7750K .......... .......... .......... .......... .......... 70% 3.06M 1s
  7800K .......... .......... .......... .......... .......... 71% 7.05M 1s
  7850K .......... .......... .......... .......... .......... 71% 2.53M 1s
  7900K .......... .......... .......... .......... .......... 72% 3.45M 1s
  7950K .......... .......... .......... .......... .......... 72% 8.98M 1s
  8000K .......... .......... .......... .......... .......... 73% 2.92M 1s
  8050K .......... .......... .......... .......... .......... 73% 3.92M 1s
  8100K .......... .......... .......... .......... .......... 74% 4.87M 1s
  8150K .......... .......... .......... .......... .......... 74% 4.88M 1s
  8200K .......... .......... .......... .......... .......... 74% 2.28M 1s
  8250K .......... .......... .......... .......... .......... 75% 5.93M 1s
  8300K .......... .......... .......... .......... .......... 75% 2.23M 1s
  8350K .......... .......... .......... .......... .......... 76% 4.66M 1s
  8400K .......... .......... .......... .......... .......... 76% 6.00M 1s
  8450K .......... .......... .......... .......... .......... 77% 1.60M 1s
  8500K .......... .......... .......... .......... .......... 77% 3.92M 1s
  8550K .......... .......... .......... .......... .......... 78% 4.92M 1s
  8600K .......... .......... .......... .......... .......... 78% 2.73M 1s
  8650K .......... .......... .......... .......... .......... 79% 5.33M 1s
  8700K .......... .......... .......... .......... .......... 79%  563K 1s
  8750K .......... .......... .......... .......... .......... 79% 6.23M 1s
  8800K .......... .......... .......... .......... .......... 80% 1.01M 1s
  8850K .......... .......... .......... .......... .......... 80% 2.55M 1s
  8900K .......... .......... .......... .......... .......... 81%  606K 1s
  8950K .......... .......... .......... .......... .......... 81%  489K 1s
  9000K .......... .......... .......... .......... .......... 82% 1.72M 1s
  9050K .......... .......... .......... .......... .......... 82%  662K 1s
  9100K .......... .......... .......... .......... .......... 83%  894K 1s
  9150K .......... .......... .......... .......... .......... 83%  642K 1s
  9200K .......... .......... .......... .......... .......... 84%  961K 1s
  9250K .......... .......... .......... .......... .......... 84% 1018K 1s
  9300K .......... .......... .......... .......... .......... 84% 1.33M 1s
  9350K .......... .......... .......... .......... .......... 85% 1.51M 1s
  9400K .......... .......... .......... .......... .......... 85% 2.19M 1s
  9450K .......... .......... .......... .......... .......... 86%  847K 1s
  9500K .......... .......... .......... .......... .......... 86% 2.94M 1s
  9550K .......... .......... .......... .......... .......... 87%  915K 1s
  9600K .......... .......... .......... .......... .......... 87%  831K 1s
  9650K .......... .......... .......... .......... .......... 88% 1.66M 1s
  9700K .......... .......... .......... .......... .......... 88%  634K 1s
  9750K .......... .......... .......... .......... .......... 89% 1.09M 1s
  9800K .......... .......... .......... .......... .......... 89% 1.27M 1s
  9850K .......... .......... .......... .......... .......... 89%  658K 1s
  9900K .......... .......... .......... .......... .......... 90% 1014K 0s
  9950K .......... .......... .......... .......... .......... 90% 1.45M 0s
 10000K .......... .......... .......... .......... .......... 91% 1.75M 0s
 10050K .......... .......... .......... .......... .......... 91%  695K 0s
 10100K .......... .......... .......... .......... .......... 92% 1.60M 0s
 10150K .......... .......... .......... .......... .......... 92% 1.33M 0s
 10200K .......... .......... .......... .......... .......... 93% 1.10M 0s
 10250K .......... .......... .......... .......... .......... 93% 2.80M 0s
 10300K .......... .......... .......... .......... .......... 94% 1.12M 0s
 10350K .......... .......... .......... .......... .......... 94% 16.6M 0s
 10400K .......... .......... .......... .......... .......... 94% 1.42M 0s
 10450K .......... .......... .......... .......... .......... 95% 1.96M 0s
 10500K .......... .......... .......... .......... .......... 95% 3.17M 0s
 10550K .......... .......... .......... .......... .......... 96% 2.27M 0s
 10600K .......... .......... .......... .......... .......... 96% 3.10M 0s
 10650K .......... .......... .......... .......... .......... 97% 2.42M 0s
 10700K .......... .......... .......... .......... .......... 97% 5.59M 0s
 10750K .......... .......... .......... .......... .......... 98% 1.58M 0s
 10800K .......... .......... .......... .......... .......... 98% 1.18M 0s
 10850K .......... .......... .......... .......... .......... 99% 14.3M 0s
 10900K .......... .......... .......... .......... .......... 99% 1.92M 0s
 10950K .......... .......... .......... .......... .......... 99% 4.14M 0s
 11000K ..                                                    100% 33.0M=5.2sgodel version 0.27.0
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
[compiles]      Running compiles...
[extimport]     Running extimport...
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
[errcheck]      Finished errcheck
[unconvert]     Running unconvert...
[deadcode]      Finished deadcode
[varcheck]      Running varcheck...
[compiles]      Finished compiles
[outparamcheck] Finished outparamcheck
[unconvert]     Finished unconvert
[varcheck]      Finished varcheck
Running test...
?   	github.com/nmiyake/echgo                 	[no test files]
ok  	github.com/nmiyake/echgo/echo            	0.002s
?   	github.com/nmiyake/echgo/generator       	[no test files]
--- FAIL: TestInvalidType (0.69s)
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

goroutine 18 [running]:
testing.tRunner.func1(0xc42007c1e0)
	/usr/local/go/src/testing/testing.go:742 +0x29d
panic(0x516640, 0xc42004c090)
	/usr/local/go/src/runtime/panic.go:505 +0x229
github.com/nmiyake/echgo/integration_test_test.TestInvalidType(0xc42007c1e0)
	/go/src/github.com/nmiyake/echgo/integration_test/integration_test.go:27 +0x49e
testing.tRunner(0xc42007c1e0, 0x548608)
	/usr/local/go/src/testing/testing.go:777 +0xd0
created by testing.(*T).Run
	/usr/local/go/src/testing/testing.go:824 +0x2e0
FAIL	github.com/nmiyake/echgo/integration_test	0.693s
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
[compiles]      Running compiles...
[extimport]     Running extimport...
[errcheck]      Running errcheck...
[deadcode]      Running deadcode...
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
ok  	github.com/nmiyake/echgo/integration_test	2.291s
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
