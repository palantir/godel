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

package distgo_test

//func TestDistTopologicalSort(t *testing.T) {
//	for i, tc := range []struct {
//		name          string
//		projectCfgYML string
//		want          []string
//	}{
//		{
//			"1 -> 2",
//			`
//products:
//  test-1:
//    dist:
//      dist-type:
//        type: os-arch-bin
//    dependencies:
//      - test-2
//  test-2:
//    dist:
//      dist-type:
//        type: os-arch-bin
//`,
//			[]string{
//				"test-2",
//				"test-1",
//			},
//		},
//		{
//			"1 -> 2,3",
//			`
//products:
//  test-1:
//    dist:
//      dist-type:
//        type: os-arch-bin
//    dependencies:
//      - test-2
//      - test-3
//  test-2:
//    dist:
//      dist-type:
//        type: os-arch-bin
//  test-3:
//    dist:
//      dist-type:
//        type: os-arch-bin
//`,
//			[]string{
//				"test-2",
//				"test-3",
//				"test-1",
//			},
//		},
//		{
//			"1 -> 2 -> 3",
//			`
//products:
//  test-1:
//    dist:
//      dist-type:
//        type: os-arch-bin
//    dependencies:
//      - test-2
//  test-2:
//    dist:
//      dist-type:
//        type: os-arch-bin
//    dependencies:
//      - test-3
//  test-3:
//    dist:
//      dist-type:
//        type: os-arch-bin
//`,
//			[]string{
//				"test-3",
//				"test-2",
//				"test-1",
//			},
//		},
//	} {
//		disterFactory, err := dister.NewDisterFactory()
//		require.NoError(t, err, "Case %d: %s", i, tc.name)
//		defaultDistCfg, err := dister.DefaultConfig()
//		require.NoError(t, err, "Case %d: %s", i, tc.name)
//
//		cfg, err := distgo.LoadConfig([]byte(tc.projectCfgYML))
//		require.NoError(t, err, "Case %d: %s", i, tc.name)
//		projectParam, err := cfg.ToParam("", disterFactory, defaultDistCfg, dockerBuilderFactory)
//		require.NoError(t, err, "Case %d: %s", i, tc.name)
//
//		dependencyGraph, err := toDependencyGraph(projectParam.Products)
//		require.NoError(t, err, "Case %d: %s", i, tc.name)
//
//		got, err := topologicalOrdering(dependencyGraph)
//		require.NoError(t, err, "Case %d: %s", i, tc.name)
//
//		assert.Equal(t, tc.want, got, "Case %d: %s", i, tc.name)
//	}
//}
