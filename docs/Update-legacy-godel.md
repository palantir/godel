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
* Go files have license headers
* `godel/config/godel.yml` is configured to add the go-generate plugin
* `godel/config/generate-plugin.yml` is configured to generate string function
* `godel/config/godel.yml` is configured to ignore all `.+_string.go` files
* `integration_test` contains integration tests
* `godel/config/test-plugin.yml` is configured to specify the "integration" tag
* `docs` contains documentation
* `.circleci/config.yml` exists
* Project is tagged as 1.0.0

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

     0K .......... .......... .......... .......... ..........  0% 1.42M 8s
    50K .......... .......... .......... .......... ..........  0% 1.51M 7s
   100K .......... .......... .......... .......... ..........  1% 1.78M 7s
   150K .......... .......... .......... .......... ..........  1% 5.36M 6s
   200K .......... .......... .......... .......... ..........  2%  600K 8s
   250K .......... .......... .......... .......... ..........  2% 5.57M 7s
   300K .......... .......... .......... .......... ..........  3% 15.8M 6s
   350K .......... .......... .......... .......... ..........  3%  958K 7s
   400K .......... .......... .......... .......... ..........  4% 1.48M 7s
   450K .......... .......... .......... .......... ..........  4% 1.06M 7s
   500K .......... .......... .......... .......... ..........  4% 2.68M 7s
   550K .......... .......... .......... .......... ..........  5% 1.69M 7s
   600K .......... .......... .......... .......... ..........  5% 1.42M 7s
   650K .......... .......... .......... .......... ..........  6% 5.99M 6s
   700K .......... .......... .......... .......... ..........  6% 1.18M 6s
   750K .......... .......... .......... .......... ..........  7% 6.38M 6s
   800K .......... .......... .......... .......... ..........  7% 1.68M 6s
   850K .......... .......... .......... .......... ..........  8% 4.64M 6s
   900K .......... .......... .......... .......... ..........  8% 1.01M 6s
   950K .......... .......... .......... .......... ..........  9% 1.20M 6s
  1000K .......... .......... .......... .......... ..........  9% 3.44M 6s
  1050K .......... .......... .......... .......... ..........  9% 1.38M 6s
  1100K .......... .......... .......... .......... .......... 10% 7.06M 6s
  1150K .......... .......... .......... .......... .......... 10% 1.95M 6s
  1200K .......... .......... .......... .......... .......... 11% 4.25M 5s
  1250K .......... .......... .......... .......... .......... 11% 2.21M 5s
  1300K .......... .......... .......... .......... .......... 12% 1.55M 5s
  1350K .......... .......... .......... .......... .......... 12% 13.4M 5s
  1400K .......... .......... .......... .......... .......... 13% 1.85M 5s
  1450K .......... .......... .......... .......... .......... 13% 4.91M 5s
  1500K .......... .......... .......... .......... .......... 14% 1.90M 5s
  1550K .......... .......... .......... .......... .......... 14% 2.05M 5s
  1600K .......... .......... .......... .......... .......... 14% 8.85M 5s
  1650K .......... .......... .......... .......... .......... 15%  880K 5s
  1700K .......... .......... .......... .......... .......... 15% 5.42M 5s
  1750K .......... .......... .......... .......... .......... 16% 2.38M 5s
  1800K .......... .......... .......... .......... .......... 16% 4.62M 5s
  1850K .......... .......... .......... .......... .......... 17% 1.04M 5s
  1900K .......... .......... .......... .......... .......... 17% 35.0M 5s
  1950K .......... .......... .......... .......... .......... 18% 1.33M 5s
  2000K .......... .......... .......... .......... .......... 18% 3.19M 5s
  2050K .......... .......... .......... .......... .......... 19% 1.97M 4s
  2100K .......... .......... .......... .......... .......... 19% 4.13M 4s
  2150K .......... .......... .......... .......... .......... 19% 1.33M 4s
  2200K .......... .......... .......... .......... .......... 20% 9.06M 4s
  2250K .......... .......... .......... .......... .......... 20% 1.72M 4s
  2300K .......... .......... .......... .......... .......... 21% 3.73M 4s
  2350K .......... .......... .......... .......... .......... 21% 1.67M 4s
  2400K .......... .......... .......... .......... .......... 22% 2.06M 4s
  2450K .......... .......... .......... .......... .......... 22% 1.48M 4s
  2500K .......... .......... .......... .......... .......... 23% 3.89M 4s
  2550K .......... .......... .......... .......... .......... 23% 3.32M 4s
  2600K .......... .......... .......... .......... .......... 24% 2.12M 4s
  2650K .......... .......... .......... .......... .......... 24% 2.95M 4s
  2700K .......... .......... .......... .......... .......... 24% 1.99M 4s
  2750K .......... .......... .......... .......... .......... 25% 7.02M 4s
  2800K .......... .......... .......... .......... .......... 25% 2.21M 4s
  2850K .......... .......... .......... .......... .......... 26% 2.32M 4s
  2900K .......... .......... .......... .......... .......... 26% 5.67M 4s
  2950K .......... .......... .......... .......... .......... 27% 1.87M 4s
  3000K .......... .......... .......... .......... .......... 27% 3.59M 4s
  3050K .......... .......... .......... .......... .......... 28% 1.33M 4s
  3100K .......... .......... .......... .......... .......... 28% 2.29M 4s
  3150K .......... .......... .......... .......... .......... 29% 2.36M 4s
  3200K .......... .......... .......... .......... .......... 29% 3.60M 4s
  3250K .......... .......... .......... .......... .......... 29% 2.58M 4s
  3300K .......... .......... .......... .......... .......... 30% 1.75M 4s
  3350K .......... .......... .......... .......... .......... 30% 4.92M 4s
  3400K .......... .......... .......... .......... .......... 31% 2.06M 4s
  3450K .......... .......... .......... .......... .......... 31% 3.88M 3s
  3500K .......... .......... .......... .......... .......... 32% 2.95M 3s
  3550K .......... .......... .......... .......... .......... 32% 1.44M 3s
  3600K .......... .......... .......... .......... .......... 33% 11.3M 3s
  3650K .......... .......... .......... .......... .......... 33% 1.55M 3s
  3700K .......... .......... .......... .......... .......... 34% 5.14M 3s
  3750K .......... .......... .......... .......... .......... 34% 2.75M 3s
  3800K .......... .......... .......... .......... .......... 34% 2.84M 3s
  3850K .......... .......... .......... .......... .......... 35% 2.25M 3s
  3900K .......... .......... .......... .......... .......... 35% 1.74M 3s
  3950K .......... .......... .......... .......... .......... 36% 6.94M 3s
  4000K .......... .......... .......... .......... .......... 36% 1.19M 3s
  4050K .......... .......... .......... .......... .......... 37% 15.2M 3s
  4100K .......... .......... .......... .......... .......... 37% 1.58M 3s
  4150K .......... .......... .......... .......... .......... 38% 3.57M 3s
  4200K .......... .......... .......... .......... .......... 38% 5.04M 3s
  4250K .......... .......... .......... .......... .......... 39% 1.59M 3s
  4300K .......... .......... .......... .......... .......... 39% 2.48M 3s
  4350K .......... .......... .......... .......... .......... 39% 1.46M 3s
  4400K .......... .......... .......... .......... .......... 40% 6.25M 3s
  4450K .......... .......... .......... .......... .......... 40% 1.72M 3s
  4500K .......... .......... .......... .......... .......... 41% 4.38M 3s
  4550K .......... .......... .......... .......... .......... 41%  966K 3s
  4600K .......... .......... .......... .......... .......... 42% 1.72M 3s
  4650K .......... .......... .......... .......... .......... 42% 3.05M 3s
  4700K .......... .......... .......... .......... .......... 43% 1.85M 3s
  4750K .......... .......... .......... .......... .......... 43% 4.91M 3s
  4800K .......... .......... .......... .......... .......... 44% 1.97M 3s
  4850K .......... .......... .......... .......... .......... 44% 2.76M 3s
  4900K .......... .......... .......... .......... .......... 44% 3.11M 3s
  4950K .......... .......... .......... .......... .......... 45% 1.51M 3s
  5000K .......... .......... .......... .......... .......... 45% 5.03M 3s
  5050K .......... .......... .......... .......... .......... 46% 1.40M 3s
  5100K .......... .......... .......... .......... .......... 46% 2.33M 3s
  5150K .......... .......... .......... .......... .......... 47% 1.40M 3s
  5200K .......... .......... .......... .......... .......... 47% 1.47M 3s
  5250K .......... .......... .......... .......... .......... 48% 4.83M 3s
  5300K .......... .......... .......... .......... .......... 48% 1.53M 3s
  5350K .......... .......... .......... .......... .......... 49% 3.58M 3s
  5400K .......... .......... .......... .......... .......... 49% 1.08M 3s
  5450K .......... .......... .......... .......... .......... 49% 1.43M 3s
  5500K .......... .......... .......... .......... .......... 50% 27.8M 2s
  5550K .......... .......... .......... .......... .......... 50% 2.82M 2s
  5600K .......... .......... .......... .......... .......... 51% 2.60M 2s
  5650K .......... .......... .......... .......... .......... 51% 1.69M 2s
  5700K .......... .......... .......... .......... .......... 52% 10.9M 2s
  5750K .......... .......... .......... .......... .......... 52% 1.24M 2s
  5800K .......... .......... .......... .......... .......... 53% 3.02M 2s
  5850K .......... .......... .......... .......... .......... 53% 1.14M 2s
  5900K .......... .......... .......... .......... .......... 54% 2.05M 2s
  5950K .......... .......... .......... .......... .......... 54% 5.36M 2s
  6000K .......... .......... .......... .......... .......... 54% 1.49M 2s
  6050K .......... .......... .......... .......... .......... 55% 4.30M 2s
  6100K .......... .......... .......... .......... .......... 55% 1.30M 2s
  6150K .......... .......... .......... .......... .......... 56% 6.45M 2s
  6200K .......... .......... .......... .......... .......... 56% 1.42M 2s
  6250K .......... .......... .......... .......... .......... 57% 4.87M 2s
  6300K .......... .......... .......... .......... .......... 57% 2.60M 2s
  6350K .......... .......... .......... .......... .......... 58% 1.72M 2s
  6400K .......... .......... .......... .......... .......... 58% 7.07M 2s
  6450K .......... .......... .......... .......... .......... 59%  454K 2s
  6500K .......... .......... .......... .......... .......... 59% 24.5M 2s
  6550K .......... .......... .......... .......... .......... 59% 2.80M 2s
  6600K .......... .......... .......... .......... .......... 60% 17.4M 2s
  6650K .......... .......... .......... .......... .......... 60% 1.61M 2s
  6700K .......... .......... .......... .......... .......... 61% 4.23M 2s
  6750K .......... .......... .......... .......... .......... 61% 2.38M 2s
  6800K .......... .......... .......... .......... .......... 62% 1.76M 2s
  6850K .......... .......... .......... .......... .......... 62% 10.1M 2s
  6900K .......... .......... .......... .......... .......... 63% 1.75M 2s
  6950K .......... .......... .......... .......... .......... 63% 3.36M 2s
  7000K .......... .......... .......... .......... .......... 64% 2.12M 2s
  7050K .......... .......... .......... .......... .......... 64% 3.14M 2s
  7100K .......... .......... .......... .......... .......... 64% 1.21M 2s
  7150K .......... .......... .......... .......... .......... 65% 2.85M 2s
  7200K .......... .......... .......... .......... .......... 65% 1.51M 2s
  7250K .......... .......... .......... .......... .......... 66% 5.18M 2s
  7300K .......... .......... .......... .......... .......... 66% 1.16M 2s
  7350K .......... .......... .......... .......... .......... 67% 12.6M 2s
  7400K .......... .......... .......... .......... .......... 67% 1.25M 2s
  7450K .......... .......... .......... .......... .......... 68% 9.61M 2s
  7500K .......... .......... .......... .......... .......... 68% 1.41M 2s
  7550K .......... .......... .......... .......... .......... 69% 9.78M 2s
  7600K .......... .......... .......... .......... .......... 69% 2.50M 2s
  7650K .......... .......... .......... .......... .......... 69% 3.81M 1s
  7700K .......... .......... .......... .......... .......... 70% 2.60M 1s
  7750K .......... .......... .......... .......... .......... 70% 4.73M 1s
  7800K .......... .......... .......... .......... .......... 71% 1.85M 1s
  7850K .......... .......... .......... .......... .......... 71% 2.46M 1s
  7900K .......... .......... .......... .......... .......... 72% 1.43M 1s
  7950K .......... .......... .......... .......... .......... 72% 5.04M 1s
  8000K .......... .......... .......... .......... .......... 73% 1.76M 1s
  8050K .......... .......... .......... .......... .......... 73% 5.72M 1s
  8100K .......... .......... .......... .......... .......... 74% 2.42M 1s
  8150K .......... .......... .......... .......... .......... 74% 1.83M 1s
  8200K .......... .......... .......... .......... .......... 74% 16.8M 1s
  8250K .......... .......... .......... .......... .......... 75% 1.26M 1s
  8300K .......... .......... .......... .......... .......... 75% 5.32M 1s
  8350K .......... .......... .......... .......... .......... 76% 3.35M 1s
  8400K .......... .......... .......... .......... .......... 76% 1.63M 1s
  8450K .......... .......... .......... .......... .......... 77% 30.0M 1s
  8500K .......... .......... .......... .......... .......... 77% 1.04M 1s
  8550K .......... .......... .......... .......... .......... 78% 13.5M 1s
  8600K .......... .......... .......... .......... .......... 78% 24.4M 1s
  8650K .......... .......... .......... .......... .......... 79%  879K 1s
  8700K .......... .......... .......... .......... .......... 79% 2.63M 1s
  8750K .......... .......... .......... .......... .......... 79% 4.32M 1s
  8800K .......... .......... .......... .......... .......... 80% 1.05M 1s
  8850K .......... .......... .......... .......... .......... 80% 5.39M 1s
  8900K .......... .......... .......... .......... .......... 81% 19.8M 1s
  8950K .......... .......... .......... .......... .......... 81% 1.70M 1s
  9000K .......... .......... .......... .......... .......... 82% 3.44M 1s
  9050K .......... .......... .......... .......... .......... 82% 3.49M 1s
  9100K .......... .......... .......... .......... .......... 83% 2.44M 1s
  9150K .......... .......... .......... .......... .......... 83% 7.85M 1s
  9200K .......... .......... .......... .......... .......... 84% 1.63M 1s
  9250K .......... .......... .......... .......... .......... 84% 6.74M 1s
  9300K .......... .......... .......... .......... .......... 84% 7.38M 1s
  9350K .......... .......... .......... .......... .......... 85% 2.11M 1s
  9400K .......... .......... .......... .......... .......... 85%  849K 1s
  9450K .......... .......... .......... .......... .......... 86% 11.9M 1s
  9500K .......... .......... .......... .......... .......... 86% 2.64M 1s
  9550K .......... .......... .......... .......... .......... 87% 3.96M 1s
  9600K .......... .......... .......... .......... .......... 87% 3.09M 1s
  9650K .......... .......... .......... .......... .......... 88% 4.80M 1s
  9700K .......... .......... .......... .......... .......... 88% 1.41M 1s
  9750K .......... .......... .......... .......... .......... 89% 9.56M 1s
  9800K .......... .......... .......... .......... .......... 89% 3.19M 1s
  9850K .......... .......... .......... .......... .......... 89% 1.94M 0s
  9900K .......... .......... .......... .......... .......... 90% 9.54M 0s
  9950K .......... .......... .......... .......... .......... 90% 2.44M 0s
 10000K .......... .......... .......... .......... .......... 91% 2.60M 0s
 10050K .......... .......... .......... .......... .......... 91% 9.72M 0s
 10100K .......... .......... .......... .......... .......... 92% 2.36M 0s
 10150K .......... .......... .......... .......... .......... 92% 2.01M 0s
 10200K .......... .......... .......... .......... .......... 93% 8.96M 0s
 10250K .......... .......... .......... .......... .......... 93% 1.55M 0s
 10300K .......... .......... .......... .......... .......... 94% 1.67M 0s
 10350K .......... .......... .......... .......... .......... 94% 2.18M 0s
 10400K .......... .......... .......... .......... .......... 94% 9.56M 0s
 10450K .......... .......... .......... .......... .......... 95% 3.74M 0s
 10500K .......... .......... .......... .......... .......... 95% 1.63M 0s
 10550K .......... .......... .......... .......... .......... 96% 9.54M 0s
 10600K .......... .......... .......... .......... .......... 96% 1.78M 0s
 10650K .......... .......... .......... .......... .......... 97% 2.91M 0s
 10700K .......... .......... .......... .......... .......... 97% 5.92M 0s
 10750K .......... .......... .......... .......... .......... 98% 1.48M 0s
 10800K .......... .......... .......... .......... .......... 98% 4.93M 0s
 10850K .......... .......... .......... .......... .......... 99% 3.05M 0s
 10900K .......... .......... .......... .......... .......... 99% 6.06M 0s
 10950K .......... .......... .......... .......... .......... 99% 3.55M 0s
 11000K ..                                                    100% 34.4M=4.7sgodel version 0.27.0
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
        id: "com.palantir.godel-generate-plugin:generate-plugin:1.0.0-rc1"
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
[extimport]     Running extimport...
[deadcode]      Running deadcode...
[compiles]      Running compiles...
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
[deadcode]      Finished deadcode
[unconvert]     Running unconvert...
[varcheck]      Running varcheck...
[compiles]      Finished compiles
[outparamcheck] Finished outparamcheck
[unconvert]     Finished unconvert
[varcheck]      Finished varcheck
Running test...
?   	github.com/nmiyake/echgo                 	[no test files]
ok  	github.com/nmiyake/echgo/echo            	0.003s
?   	github.com/nmiyake/echgo/generator       	[no test files]
--- FAIL: TestInvalidType (0.10s)
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
FAIL	github.com/nmiyake/echgo/integration_test	0.107s
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
ok  	github.com/nmiyake/echgo/integration_test	0.542s
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
