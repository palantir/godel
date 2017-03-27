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

package dist

import (
	"sort"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
)

func Schedule(specsWithDeps []params.ProductBuildSpecWithDeps) ([]params.ProductBuildSpecWithDeps, error) {
	var schedule []params.ProductBuildSpecWithDeps
	graph := make(map[string]map[string]bool)
	specMap := make(map[string][]params.ProductBuildSpecWithDeps)
	// create a graph of dependencies
	for _, curSpec := range specsWithDeps {
		product := curSpec.Spec.ProductName
		if graph[product] == nil {
			graph[product] = make(map[string]bool)
		}
		if specMap[product] == nil {
			specMap[product] = make([]params.ProductBuildSpecWithDeps, 0)
		}
		specMap[product] = append(specMap[product], curSpec)
		for _, curDist := range curSpec.Spec.Dist {
			for _, dep := range curDist.Info.Deps() {
				if graph[dep] == nil {
					graph[dep] = make(map[string]bool)
				}
				graph[dep][product] = true
			}
		}
	}

	// get the topological ordering among the products
	order, err := topologicalOrdering(graph)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to generate ordering among the products. The dist dependencies between the products contains a cycle.")
	}

	// construct the final schedule
	for _, product := range order {
		for _, spec := range specMap[product] {
			schedule = append(schedule, spec)
		}
	}

	return schedule, nil
}

func topologicalOrdering(graph map[string]map[string]bool) ([]string, error) {
	var order []string
	// get all nodes in the graph and sort lexicographically for deterministic order
	var nodes []string
	indeg := make(map[string]int)
	for node := range graph {
		indeg[node] = 0
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)
	// compute the incoming edges on each vertex
	for _, v := range nodes {
		for neighbour := range graph[v] {
			indeg[neighbour]++
		}
	}
	// q contains all vertices with in-degree zero
	var q []string
	for _, v := range nodes {
		if indeg[v] == 0 {
			q = append(q, v)
		}
	}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		order = append(order, cur)
		var neighbours []string
		// sort all the neighbours to ensure deterministic order
		for neighbour := range graph[cur] {
			neighbours = append(neighbours, neighbour)
		}
		sort.Strings(neighbours)
		for _, neighbour := range neighbours {
			indeg[neighbour]--
			if indeg[neighbour] == 0 {
				q = append(q, neighbour)
			}
		}
	}
	if len(order) != len(graph) {
		return nil, errors.New("Please provide a valid DAG as input")
	}
	return order, nil
}
