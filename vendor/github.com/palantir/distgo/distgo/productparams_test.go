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

import (
	"testing"

	"github.com/palantir/godel/pkg/osarch"
	"github.com/stretchr/testify/assert"

	"github.com/palantir/distgo/dister"
	"github.com/palantir/distgo/distgo"
)

func TestProductParamsForProductArgs(t *testing.T) {
	for i, tc := range []struct {
		projectParam distgo.ProjectParam
		productIDs   []distgo.ProductID
		want         []distgo.ProductParam
		wantError    string
	}{
		{
			projectParam: distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"foo": {ID: "foo"},
					"bar": {ID: "bar"},
				},
			},
			productIDs: nil,
			want: []distgo.ProductParam{
				{ID: "bar"},
				{ID: "foo"},
			},
		},
		{
			projectParam: distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"foo": {ID: "foo"},
					"bar": {ID: "bar"},
				},
			},
			productIDs: []distgo.ProductID{
				"foo",
			},
			want: []distgo.ProductParam{
				{ID: "foo"},
			},
		},
		{
			projectParam: distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"foo": {ID: "foo"},
					"bar": {ID: "bar"},
				},
			},
			productIDs: []distgo.ProductID{
				"baz",
			},
			wantError: "product(s) [baz] not valid -- valid values are [bar foo]",
		},
	} {
		products, err := distgo.ProductParamsForProductArgs(tc.projectParam.Products, tc.productIDs...)
		if tc.wantError == "" {
			assert.Equal(t, tc.want, products, "Case %d", i)
		} else {
			assert.EqualError(t, err, tc.wantError, "Case %d", i)
		}
	}
}

func TestProductParamsForBuildProductArgs(t *testing.T) {
	mustOSArch := func(in string) osarch.OSArch {
		out, err := osarch.New(in)
		if err != nil {
			panic(err)
		}
		return out
	}

	for i, tc := range []struct {
		projectParam    distgo.ProjectParam
		productBuildIDs []distgo.ProductBuildID
		want            []distgo.ProductParam
		wantError       string
	}{
		{
			projectParam: distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"foo": {
						ID: "foo",
						Build: &distgo.BuildParam{
							OSArchs: []osarch.OSArch{
								mustOSArch("darwin-amd64"),
								mustOSArch("linux-amd64"),
							},
						},
					},
					"bar": {
						ID: "bar",
						Build: &distgo.BuildParam{
							OSArchs: []osarch.OSArch{
								mustOSArch("darwin-amd64"),
							},
						},
					},
				},
			},
			productBuildIDs: nil,
			want: []distgo.ProductParam{
				{
					ID: "bar",
					Build: &distgo.BuildParam{
						OSArchs: []osarch.OSArch{
							mustOSArch("darwin-amd64"),
						},
					},
				},
				{
					ID: "foo",
					Build: &distgo.BuildParam{
						OSArchs: []osarch.OSArch{
							mustOSArch("darwin-amd64"),
							mustOSArch("linux-amd64"),
						},
					},
				},
			},
		},
		{
			projectParam: distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"foo": {
						ID: "foo",
						Build: &distgo.BuildParam{
							OSArchs: []osarch.OSArch{
								mustOSArch("darwin-amd64"),
								mustOSArch("linux-amd64"),
							},
						},
					},
					"bar": {
						ID: "bar",
						Build: &distgo.BuildParam{
							OSArchs: []osarch.OSArch{
								mustOSArch("darwin-amd64"),
							},
						},
					},
				},
			},
			productBuildIDs: []distgo.ProductBuildID{
				"foo.darwin-amd64",
				"bar",
			},
			want: []distgo.ProductParam{
				{
					ID: "bar",
					Build: &distgo.BuildParam{
						OSArchs: []osarch.OSArch{
							mustOSArch("darwin-amd64"),
						},
					},
				},
				{
					ID: "foo",
					Build: &distgo.BuildParam{
						OSArchs: []osarch.OSArch{
							mustOSArch("darwin-amd64"),
						},
					},
				},
			},
		},
		{
			projectParam: distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"foo": {
						ID: "foo",
						Build: &distgo.BuildParam{
							OSArchs: []osarch.OSArch{
								mustOSArch("darwin-amd64"),
								mustOSArch("linux-amd64"),
							},
						},
					},
					"bar": {
						ID: "bar",
						Build: &distgo.BuildParam{
							OSArchs: []osarch.OSArch{
								mustOSArch("darwin-amd64"),
							},
						},
					},
				},
			},
			productBuildIDs: []distgo.ProductBuildID{
				"baz",
			},
			wantError: "build product(s) [baz] not valid -- valid values are [bar bar.darwin-amd64 foo foo.darwin-amd64 foo.linux-amd64]",
		},
		{
			projectParam: distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"foo": {
						ID: "foo",
						Build: &distgo.BuildParam{
							OSArchs: []osarch.OSArch{
								mustOSArch("darwin-amd64"),
								mustOSArch("linux-amd64"),
							},
						},
					},
					"bar": {
						ID: "bar",
						Build: &distgo.BuildParam{
							OSArchs: []osarch.OSArch{
								mustOSArch("darwin-amd64"),
							},
						},
					},
				},
			},
			productBuildIDs: []distgo.ProductBuildID{
				"bar.linux-amd64",
			},
			wantError: "build product(s) [bar.linux-amd64] not valid -- valid values are [bar bar.darwin-amd64 foo foo.darwin-amd64 foo.linux-amd64]",
		},
	} {
		products, err := distgo.ProductParamsForBuildProductArgs(tc.projectParam.Products, tc.productBuildIDs...)
		if tc.wantError == "" {
			assert.Equal(t, tc.want, products, "Case %d", i)
		} else {
			assert.EqualError(t, err, tc.wantError, "Case %d", i)
		}
	}
}

func TestProductParamsForDistProductArgs(t *testing.T) {
	for i, tc := range []struct {
		projectParam   distgo.ProjectParam
		productDistIDs []distgo.ProductDistID
		want           []distgo.ProductParam
		wantError      string
	}{
		{
			projectParam: distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"foo": {
						ID: "foo",
						Dist: &distgo.DistParam{
							DistParams: map[distgo.DistID]distgo.DisterParam{
								dister.OSArchBinDistTypeName: {
									Dister: dister.NewOSArchBinDister(),
								},
							},
						},
					},
					"bar": {
						ID: "bar",
						Dist: &distgo.DistParam{
							DistParams: map[distgo.DistID]distgo.DisterParam{
								dister.OSArchBinDistTypeName: {
									Dister: dister.NewOSArchBinDister(),
								},
							},
						},
					},
				},
			},
			productDistIDs: nil,
			want: []distgo.ProductParam{
				{
					ID: "bar",
					Dist: &distgo.DistParam{
						DistParams: map[distgo.DistID]distgo.DisterParam{
							dister.OSArchBinDistTypeName: {
								Dister: dister.NewOSArchBinDister(),
							},
						},
					},
				},
				{
					ID: "foo",
					Dist: &distgo.DistParam{
						DistParams: map[distgo.DistID]distgo.DisterParam{
							dister.OSArchBinDistTypeName: {
								Dister: dister.NewOSArchBinDister(),
							},
						},
					},
				},
			},
		},
		{
			projectParam: distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"foo": {
						ID: "foo",
						Dist: &distgo.DistParam{
							DistParams: map[distgo.DistID]distgo.DisterParam{
								"dister-1": {
									Dister: dister.NewOSArchBinDister(),
								},
								"dister-2": {
									Dister: dister.NewOSArchBinDister(),
								},
							},
						},
					},
					"bar": {
						ID: "bar",
						Dist: &distgo.DistParam{
							DistParams: map[distgo.DistID]distgo.DisterParam{
								dister.OSArchBinDistTypeName: {
									Dister: dister.NewOSArchBinDister(),
								},
							},
						},
					},
				},
			},
			productDistIDs: []distgo.ProductDistID{
				"foo.dister-2",
				"bar",
			},
			want: []distgo.ProductParam{
				{
					ID: "bar",
					Dist: &distgo.DistParam{
						DistParams: map[distgo.DistID]distgo.DisterParam{
							dister.OSArchBinDistTypeName: {
								Dister: dister.NewOSArchBinDister(),
							},
						},
					},
				},
				{
					ID: "foo",
					Dist: &distgo.DistParam{
						DistParams: map[distgo.DistID]distgo.DisterParam{
							"dister-2": {
								Dister: dister.NewOSArchBinDister(),
							},
						},
					},
				},
			},
		},
		{
			projectParam: distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"foo": {
						ID: "foo",
						Dist: &distgo.DistParam{
							DistParams: map[distgo.DistID]distgo.DisterParam{
								"dister-1": {
									Dister: dister.NewOSArchBinDister(),
								},
								"dister-2": {
									Dister: dister.NewOSArchBinDister(),
								},
							},
						},
					},
					"bar": {
						ID: "bar",
						Dist: &distgo.DistParam{
							DistParams: map[distgo.DistID]distgo.DisterParam{
								dister.OSArchBinDistTypeName: {
									Dister: dister.NewOSArchBinDister(),
								},
							},
						},
					},
				},
			},
			productDistIDs: []distgo.ProductDistID{
				"baz",
			},
			wantError: "dist product(s) [baz] not valid -- valid values are [bar bar.os-arch-bin foo foo.dister-1 foo.dister-2]",
		},
		{
			projectParam: distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"foo": {
						ID: "foo",
						Dist: &distgo.DistParam{
							DistParams: map[distgo.DistID]distgo.DisterParam{
								"dister-1": {
									Dister: dister.NewOSArchBinDister(),
								},
								"dister-2": {
									Dister: dister.NewOSArchBinDister(),
								},
							},
						},
					},
					"bar": {
						ID: "bar",
						Dist: &distgo.DistParam{
							DistParams: map[distgo.DistID]distgo.DisterParam{
								dister.OSArchBinDistTypeName: {
									Dister: dister.NewOSArchBinDister(),
								},
							},
						},
					},
				},
			},
			productDistIDs: []distgo.ProductDistID{
				"bar.bad-dister",
			},
			wantError: "dist product(s) [bar.bad-dister] not valid -- valid values are [bar bar.os-arch-bin foo foo.dister-1 foo.dister-2]",
		},
	} {
		products, err := distgo.ProductParamsForDistProductArgs(tc.projectParam.Products, tc.productDistIDs...)
		if tc.wantError == "" {
			assert.Equal(t, tc.want, products, "Case %d", i)
		} else {
			assert.EqualError(t, err, tc.wantError, "Case %d", i)
		}
	}
}
