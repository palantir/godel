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

package publisher_test

import (
	"fmt"
	"testing"

	"github.com/palantir/godel/pkg/osarch"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/dister"
	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/publisher"
)

func TestBintrayPublisher(t *testing.T) {
	runPublisherTests(t, publisher.NewBintrayPublisherCreator().Publisher(), true,
		publisherTestCase{
			name: "publishes artifact to Bintray",
			projectCfg: distgo.ProjectConfig{
				Products: map[distgo.ProductID]distgo.ProductConfig{
					"foo": {
						Build: &distgo.BuildConfig{
							MainPkg: stringPtr("foo"),
						},
						Dist: &distgo.DistConfig{
							Disters: &distgo.DistersConfig{
								dister.OSArchBinDistTypeName: {
									Type: stringPtr(dister.OSArchBinDistTypeName),
								},
							},
						},
						Publish: &distgo.PublishConfig{
							GroupID: stringPtr("com.test.group"),
							PublishInfo: &map[distgo.PublishID]yaml.MapSlice{
								publisher.BintrayPublishTypeName: *mustMapSlicePtr(publisher.BintrayPublishConfig{
									BasicConnectionInfo: publisher.BasicConnectionInfo{
										URL:      "http://artifactory.domain.com",
										Username: "testUsername",
										Password: "testPassword",
									},
									Subject:       "testSubject",
									Repository:    "testRepo",
									Product:       "testProduct",
									Publish:       true,
									DownloadsList: true,
								}),
							},
						},
					},
				},
			},
			wantOutput: func(projectDir string) string {
				return fmt.Sprintf(`[DRY RUN] Uploading %s/out/dist/foo/1.0.0/os-arch-bin/foo-1.0.0-%s.tgz to http://artifactory.domain.com/content/testSubject/testRepo/testProduct/1.0.0/com.test.group/foo/1.0.0/foo-1.0.0-%s.tgz
[DRY RUN] Running Bintray publish for uploaded artifacts...done
[DRY RUN] Adding artifact to Bintray downloads list for package...done
`, projectDir, osarch.Current().String(), osarch.Current().String())
			},
		},
		publisherTestCase{
			name: "omitting product configuration defaults to ProductID",
			projectCfg: distgo.ProjectConfig{
				Products: map[distgo.ProductID]distgo.ProductConfig{
					"foo": {
						Build: &distgo.BuildConfig{
							MainPkg: stringPtr("foo"),
						},
						Dist: &distgo.DistConfig{
							Disters: &distgo.DistersConfig{
								dister.OSArchBinDistTypeName: {
									Type: stringPtr(dister.OSArchBinDistTypeName),
								},
							},
						},
						Publish: &distgo.PublishConfig{
							GroupID: stringPtr("com.test.group"),
							PublishInfo: &map[distgo.PublishID]yaml.MapSlice{
								publisher.BintrayPublishTypeName: *mustMapSlicePtr(publisher.BintrayPublishConfig{
									BasicConnectionInfo: publisher.BasicConnectionInfo{
										URL:      "http://artifactory.domain.com",
										Username: "testUsername",
										Password: "testPassword",
									},
									Subject:       "testSubject",
									Repository:    "testRepo",
									Publish:       true,
									DownloadsList: true,
								}),
							},
						},
					},
				},
			},
			wantOutput: func(projectDir string) string {
				return fmt.Sprintf(`[DRY RUN] Uploading %s/out/dist/foo/1.0.0/os-arch-bin/foo-1.0.0-%s.tgz to http://artifactory.domain.com/content/testSubject/testRepo/foo/1.0.0/com.test.group/foo/1.0.0/foo-1.0.0-%s.tgz
[DRY RUN] Running Bintray publish for uploaded artifacts...done
[DRY RUN] Adding artifact to Bintray downloads list for package...done
`, projectDir, osarch.Current().String(), osarch.Current().String())
			},
		},
	)
}
