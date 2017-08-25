// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config_test

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/apps/distgo/config"
)

func Example_sls() {
	yml := `
products:
  cache-service:
    build:
      main-pkg: ./main/cache
      os-archs:
        - os: linux
          arch: amd64
      version-var: main.Version
    dist:
      input-dir: cache/dist/sls
      output-dir: cache/build/distributions
      dist-type:
        type: sls
        info:
          service-args: "--config var/conf/cache.yml server"
          manifest-extensions:
            cache: true
          reloadable: true
group-id: com.palantir.cache`

	cfg := configFromYML(yml)
	fmt.Printf("%q", fmt.Sprintf("%+v", cfg))
	// Output: "{Products:map[cache-service:{Build:{Skip:false Script: MainPkg:./main/cache OutputDir: BuildArgsScript: VersionVar:main.Version Environment:map[] OSArchs:[linux-amd64]} Run:{Args:[]} Dist:[{OutputDir:cache/build/distributions InputDir:cache/dist/sls InputProducts:[] Script: DistType:{Type:sls Info:{InitShTemplateFile: ManifestTemplateFile: ServiceArgs:--config var/conf/cache.yml server ProductType: ManifestExtensions:map[cache:true] Reloadable:true YMLValidationExclude:{Names:[] Paths:[]}}} Publish:{GroupID: Almanac:{Metadata:map[] Tags:[]}}}] DockerImages:[] DefaultPublish:{GroupID: Almanac:{Metadata:map[] Tags:[]}}}] BuildOutputDir: DistOutputDir: DistScriptInclude: GroupID:com.palantir.cache Exclude:{Names:[] Paths:[]}}"
}

func Example_bin() {
	yml := `
products:
  godel:
    build:
      main-pkg: ./cmd/godel
      environment:
        CGO_ENABLED: "0"
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
      version-var: main.Version
    dist:
      dist-type:
        type: bin
      script: |
              function setup_wrapper {
                # logic for function (omitted for brevity)
              }

              # copy contents of resources directory
              mkdir -p "$DIST_DIR/wrapper"
              setup_wrapper "$DIST_DIR/wrapper"
group-id: com.palantir.godel`

	cfg := configFromYML(yml)
	fmt.Printf("%q", fmt.Sprintf("%+v", cfg))
	// Output: "{Products:map[godel:{Build:{Skip:false Script: MainPkg:./cmd/godel OutputDir: BuildArgsScript: VersionVar:main.Version Environment:map[CGO_ENABLED:0] OSArchs:[darwin-amd64 linux-amd64]} Run:{Args:[]} Dist:[{OutputDir: InputDir: InputProducts:[] Script:function setup_wrapper {\n  # logic for function (omitted for brevity)\n}\n\n# copy contents of resources directory\nmkdir -p \"$DIST_DIR/wrapper\"\nsetup_wrapper \"$DIST_DIR/wrapper\"\n DistType:{Type:bin Info:{OmitInitSh:<nil> InitShTemplateFile:}} Publish:{GroupID: Almanac:{Metadata:map[] Tags:[]}}}] DockerImages:[] DefaultPublish:{GroupID: Almanac:{Metadata:map[] Tags:[]}}}] BuildOutputDir: DistOutputDir: DistScriptInclude: GroupID:com.palantir.godel Exclude:{Names:[] Paths:[]}}"
}

func Example_manual() {
	yml := `
products:
  godel:
    build:
      skip: true
    dist:
      dist-type:
        type: manual
        info:
          extension: tgz
      script: |
              # filler/dummy example. Real operation should perform dist or packaging
              # operation. Important thing is that output must be at "$DIST_DIR/$PRODUCT-$VERSION.[extension]".
              echo "test-dist-contents" > "$DIST_DIR/$PRODUCT-$VERSION.tgz"
group-id: com.palantir.godel`

	cfg := configFromYML(yml)
	fmt.Printf("%q", fmt.Sprintf("%+v", cfg))
	// Output: "{Products:map[godel:{Build:{Skip:true Script: MainPkg: OutputDir: BuildArgsScript: VersionVar: Environment:map[] OSArchs:[]} Run:{Args:[]} Dist:[{OutputDir: InputDir: InputProducts:[] Script:# filler/dummy example. Real operation should perform dist or packaging\n# operation. Important thing is that output must be at \"$DIST_DIR/$PRODUCT-$VERSION.[extension]\".\necho \"test-dist-contents\" > \"$DIST_DIR/$PRODUCT-$VERSION.tgz\"\n DistType:{Type:manual Info:{Extension:tgz}} Publish:{GroupID: Almanac:{Metadata:map[] Tags:[]}}}] DockerImages:[] DefaultPublish:{GroupID: Almanac:{Metadata:map[] Tags:[]}}}] BuildOutputDir: DistOutputDir: DistScriptInclude: GroupID:com.palantir.godel Exclude:{Names:[] Paths:[]}}"
}

func Example_rpm() {
	yml := `
products:
  orchestrator:
    dist:
      input-dir: ./rpm
      dist-type:
        type: rpm
        info:
          config-files:
            - /usr/lib/systemd/system/orchestrator.service
          before-install-script: |
              /usr/bin/getent group orchestrator || /usr/sbin/groupadd \
                      -g 380 orchestrator
              /usr/bin/getent passwd orchestrator || /usr/sbin/useradd -r \
                      -d /var/lib/orchestrator -g orchestrator -u 380 -m \
                      -s /sbin/nologin orchestrator
          after-install-script: |
              systemctl daemon-reload
          after-remove-script: |
              systemctl daemon-reload
      script: |
          mkdir "$DIST_DIR"/usr/libexec/orchestrator
          cp build/linux-amd64/orchestrator "$DIST_DIR"/usr/libexec/orchestrator
group-id: com.palantir.pcloud`

	cfg := configFromYML(yml)
	fmt.Printf("%q", fmt.Sprintf("%+v", cfg))
	// Output: "{Products:map[orchestrator:{Build:{Skip:false Script: MainPkg: OutputDir: BuildArgsScript: VersionVar: Environment:map[] OSArchs:[]} Run:{Args:[]} Dist:[{OutputDir: InputDir:./rpm InputProducts:[] Script:mkdir \"$DIST_DIR\"/usr/libexec/orchestrator\ncp build/linux-amd64/orchestrator \"$DIST_DIR\"/usr/libexec/orchestrator\n DistType:{Type:rpm Info:{Release: ConfigFiles:[/usr/lib/systemd/system/orchestrator.service] BeforeInstallScript:/usr/bin/getent group orchestrator || /usr/sbin/groupadd \\\n        -g 380 orchestrator\n/usr/bin/getent passwd orchestrator || /usr/sbin/useradd -r \\\n        -d /var/lib/orchestrator -g orchestrator -u 380 -m \\\n        -s /sbin/nologin orchestrator\n AfterInstallScript:systemctl daemon-reload\n AfterRemoveScript:systemctl daemon-reload\n}} Publish:{GroupID: Almanac:{Metadata:map[] Tags:[]}}}] DockerImages:[] DefaultPublish:{GroupID: Almanac:{Metadata:map[] Tags:[]}}}] BuildOutputDir: DistOutputDir: DistScriptInclude: GroupID:com.palantir.pcloud Exclude:{Names:[] Paths:[]}}"
}

func Example_docker() {
	yml := `
products:
  godel:
    build:
      main-pkg: ./cmd/godel
      environment:
        CGO_ENABLED: "0"
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
      version-var: main.Version
    dist:
      dist-type:
        type: bin
      script: |
              function setup_wrapper {
                # logic for function (omitted for brevity)
              }

              # copy contents of resources directory
              mkdir -p "$DIST_DIR/wrapper"
              setup_wrapper "$DIST_DIR/wrapper"
    docker:
    -
      repository: palantir/godel
      tag: latest
      context-dir: path/to/context/dir
group-id: com.palantir.godel`

	cfg := configFromYML(yml)
	fmt.Printf("%q", fmt.Sprintf("%+v", cfg))
	// Output: "{Products:map[godel:{Build:{Skip:false Script: MainPkg:./cmd/godel OutputDir: BuildArgsScript: VersionVar:main.Version Environment:map[CGO_ENABLED:0] OSArchs:[darwin-amd64 linux-amd64]} Run:{Args:[]} Dist:[{OutputDir: InputDir: InputProducts:[] Script:function setup_wrapper {\n  # logic for function (omitted for brevity)\n}\n\n# copy contents of resources directory\nmkdir -p \"$DIST_DIR/wrapper\"\nsetup_wrapper \"$DIST_DIR/wrapper\"\n DistType:{Type:bin Info:{OmitInitSh:<nil> InitShTemplateFile:}} Publish:{GroupID: Almanac:{Metadata:map[] Tags:[]}}}] DockerImages:[{Repository:palantir/godel Tag:latest ContextDir:path/to/context/dir Deps:[] Info:{Type: Data:<nil>} BuildArgsScript:}] DefaultPublish:{GroupID: Almanac:{Metadata:map[] Tags:[]}}}] BuildOutputDir: DistOutputDir: DistScriptInclude: GroupID:com.palantir.godel Exclude:{Names:[] Paths:[]}}"
}

func configFromYML(yml string) config.Project {
	cfg := config.Project{}
	if err := yaml.Unmarshal([]byte(yml), &cfg); err != nil {
		panic(err)
	}
	if _, err := cfg.ToParams(); err != nil {
		panic(err)
	}
	return cfg
}
